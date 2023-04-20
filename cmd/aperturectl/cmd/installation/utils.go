package installation

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/imdario/mergo"
	"golang.org/x/exp/slices"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/fluxninja/aperture/operator/controllers"
	"github.com/fluxninja/aperture/pkg/log"
)

var (
	valuesFile     string
	kubeConfig     string
	kubeRestConfig *rest.Config
	version        string
	latestVersion  string
	namespace      string
	kubeClient     client.Client
	timeout        int
	generateCert   bool
)

const (
	apertureLatestVersion  = "latest"
	defaultNS              = "default"
	controller             = "controller"
	agent                  = "agent"
	istioConfig            = "istioconfig"
	istioConfigReleaseName = "aperture-envoy-filter"
	apertureAgent          = "aperture-agent"
	apertureController     = "aperture-controller"
	aperturectl            = "aperturectl"
)

// getTemplets loads CRDs, hooks and manifests from the Helm chart.
func getTemplets(chartName, releaseName string, order releaseutil.KindSortOrder) ([]chart.CRD, []*release.Hook, []releaseutil.Manifest, error) {
	chartURL := fmt.Sprintf("https://fluxninja.github.io/aperture/%s-%s.tgz", chartName, version)

	resp, err := http.Get(chartURL) //nolint
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to download chart: %s", err)
	}
	defer resp.Body.Close()

	ch, err := loader.LoadArchive(resp.Body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load chart: %s", err)
	}

	values := ch.Values
	if valuesFile != "" {
		var customValues chartutil.Values
		customValues, err = chartutil.ReadValuesFile(valuesFile)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read values: %s", err)
		}

		if err = mergo.Merge(&values, customValues.AsMap(), mergo.WithOverride); err != nil {
			values = customValues
		}
	} else if releaseName == agent && slices.Compare(order, releaseutil.UninstallOrder) == 0 {
		values = map[string]interface{}{
			"agent": map[string]interface{}{
				"config": map[string]interface{}{
					"etcd": map[string]interface{}{
						"endpoints": []string{"dummy"},
					},
					"prometheus": map[string]interface{}{
						"address": "dummy",
					},
				},
			},
		}
	}

	componentValues, ok := values[releaseName].(map[string]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to get %s values", releaseName)
	}

	isNamespaceScoped, ok := componentValues["namespaceScoped"].(bool)
	if !ok {
		isNamespaceScoped = false
	}

	if releaseName == controller && isNamespaceScoped {
		values, err = manageControllerCertificateSecret(values, fmt.Sprintf("%s-%s", controller, apertureController), namespace, order)
		if err != nil {
			return nil, nil, nil, err
		}
	} else if releaseName == agent && isNamespaceScoped {
		values, err = manageAgentControllerClientCertConfigMap(values, fmt.Sprintf("%s-%s", agent, apertureAgent), namespace, order)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	renderedValues, err := chartutil.ToRenderValues(ch, values, chartutil.ReleaseOptions{Name: releaseName, Namespace: namespace}, chartutil.DefaultCapabilities)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read values: %s", err)
	}

	if err = chartutil.ProcessDependencies(ch, renderedValues); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to process dependencies: %s", err)
	}

	files, err := engine.RenderWithClient(ch, renderedValues, kubeRestConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to render chart: %s", err)
	}

	for k := range files {
		if strings.HasSuffix(k, "NOTES.txt") {
			delete(files, k)
		}
	}

	hooks, manifests, err := releaseutil.SortManifests(files, chartutil.DefaultVersionSet, order)
	crds := ch.CRDObjects()

	if isNamespaceScoped {
		crds = []chart.CRD{}
		hooks = []*release.Hook{}
	}

	return crds, hooks, manifests, err
}

// applyManifest creates/updates the generated manifest to Kubernetes.
func applyManifest(manifest string) error {
	unstructuredObject, err := prepareUnstructuredObject(manifest)
	if err != nil {
		return err
	}

	err = applyObjectToKubernetesWithRetry(unstructuredObject)
	if err != nil {
		return fmt.Errorf("failed to apply - %s/%s, Error - '%s'", unstructuredObject.GetKind(), unstructuredObject.GetName(), err)
	}
	return nil
}

// applyObjectToKubernetesWithRetry applies the given object to Kubernetes with retry.
func applyObjectToKubernetesWithRetry(unstructuredObject *unstructured.Unstructured) error {
	log.Info().Msgf("Applying - %s/%s", unstructuredObject.GetKind(), unstructuredObject.GetName())
	attempt := 0
	for attempt < 5 {
		attempt++
		err := applyObjectToKubernetes(unstructuredObject)
		if err == nil || (!strings.Contains(err.Error(), "no matches for kind") && !apierrors.IsConflict(err)) {
			return err
		}
		time.Sleep(time.Second * time.Duration(attempt))
	}
	return fmt.Errorf("failed to apply object after %d attempts", attempt)
}

// applyObjectToKubernetes applies the given object to Kubernetes.
func applyObjectToKubernetes(unstructuredObject *unstructured.Unstructured) error {
	key := types.NamespacedName{
		Name:      unstructuredObject.GetName(),
		Namespace: unstructuredObject.GetNamespace(),
	}
	existing := unstructuredObject.DeepCopy()
	err := kubeClient.Get(context.Background(), key, existing)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if !apierrors.IsNotFound(err) {
		opts := cmpopts.IgnoreFields(
			metav1.ObjectMeta{},
			"Generation", "ResourceVersion", "SelfLink", "UID",
			"CreationTimestamp", "DeletionTimestamp", "DeletionGracePeriodSeconds",
			"OwnerReferences", "Finalizers",
		)

		// Check if there are any differences between the existing and the new object
		if !cmp.Equal(unstructuredObject, existing, opts) {
			err = mergo.Map(&unstructuredObject.Object, &existing.Object)
			if err != nil {
				return err
			}
			err = kubeClient.Patch(context.Background(), unstructuredObject, client.MergeFrom(existing))
		}
	} else {
		err = kubeClient.Create(context.Background(), unstructuredObject)
	}

	return err
}

// deleteManifest deletes the generated manifest from Kubernetes.
func deleteManifest(manifest string) error {
	unstructuredObject, err := prepareUnstructuredObject(manifest)
	if err != nil {
		return err
	}

	log.Info().Msgf("Deleting - %s/%s", unstructuredObject.GetKind(), unstructuredObject.GetName())

	err = kubeClient.Delete(context.Background(), unstructuredObject)
	if apierrors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to delete - %s/%s, Error - '%s'", unstructuredObject.GetKind(), unstructuredObject.GetName(), err)
	}
	return nil
}

// manageNamespace creates namespace if not present.
func manageNamespace() error {
	ns := &corev1.Namespace{}
	err := kubeClient.Get(context.Background(), types.NamespacedName{Name: namespace}, ns)
	if apierrors.IsNotFound(err) {
		ns.Name = namespace
		if err = kubeClient.Create(context.Background(), ns); err != nil {
			return err
		}
	}

	return nil
}

// waitForHook waits for 1 successful execution of the hook.
func waitForHook(name string, ctx context.Context) error {
	job := &batchv1.Job{}
	for {
		err := kubeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, job)
		if err != nil || job.Status.Succeeded != 1 {
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	return nil
}

// prepareUnstructuredObject prepares unstructured.Unstructured from given YAML string.
func prepareUnstructuredObject(manifest string) (*unstructured.Unstructured, error) {
	content := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(manifest), &content)
	if err != nil {
		return nil, err
	}

	unstructuredObject := &unstructured.Unstructured{
		Object: content,
	}

	if unstructuredObject.GetNamespace() == "" {
		unstructuredObject.SetNamespace(namespace)
	}

	return unstructuredObject, nil
}

// handleInstall handles installation for given chart using given release name.
func handleInstall(chartName, releaseName string) error {
	crds, _, manifests, err := getTemplets(chartName, releaseName, releaseutil.InstallOrder)
	if err != nil {
		return err
	}

	errs := []error{}
	for _, crd := range crds {
		if err = applyManifest(string(crd.File.Data)); err != nil {
			errs = append(errs, err)
		}
	}

	for _, manifest := range manifests {
		if err = applyManifest(manifest.Content); err != nil {
			errs = append(errs, err)
		}
	}

	for _, err := range errs {
		log.Error().Msg(err.Error())
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to complete install successfully")
	}
	return nil
}

// handleUnInstall handles uninstallation for given chart using given release name.
func handleUnInstall(chartName, releaseName string) error {
	crds, hooks, manifests, err := getTemplets(chartName, releaseName, releaseutil.UninstallOrder)
	if err != nil {
		return err
	}

	errs := []error{}
	for _, hook := range hooks {
		log.Info().Msgf("Executing hook - %s", hook.Name)
		if err = applyManifest(hook.Manifest); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		if err = waitForHook(hook.Name, ctx); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("timed out waiting for pre-delete hook completion")
			}
			return err
		}

		if err = deleteManifest(hook.Manifest); err != nil {
			errs = append(errs, err)
		}

		if err = kubeClient.DeleteAllOf(
			context.Background(), &corev1.Pod{}, client.InNamespace(namespace), client.MatchingLabels{"job-name": hook.Name}); err != nil {
			errs = append(errs, err)
		}
	}

	for _, manifest := range manifests {
		if err = deleteManifest(manifest.Content); err != nil {
			errs = append(errs, err)
		}
	}

	for _, crd := range crds {
		if err = deleteManifest(string(crd.File.Data)); err != nil {
			errs = append(errs, err)
		}
	}

	for _, err := range errs {
		log.Error().Msg(err.Error())
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to complete uninstall successfully")
	}
	return nil
}

// manageControllerCertificateSecret manages secret containing the Aperture Controller certificate.
func manageControllerCertificateSecret(values map[string]interface{}, releaseName, namespace string, order releaseutil.KindSortOrder) (map[string]interface{}, error) {
	controller, ok := values["controller"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to get controller values")
	}

	serverCert, ok := controller["serverCert"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to get serverCert values")
	}

	secretName, ok := serverCert["secretName"].(string)
	if !ok {
		if !generateCert && slices.Equal(order, releaseutil.InstallOrder) {
			return nil, fmt.Errorf(".Values.controller.serverCert.secretName must be set when .Values.controller.namespaceScoped is true and --generate-cert is not provided")
		}
		secretName = fmt.Sprintf("%s-cert", releaseName)
	}

	keyFileName, ok := serverCert["keyFileName"].(string)
	if !ok {
		keyFileName = controllers.ControllerCertKeyName
	}

	certFileName, ok := serverCert["certFileName"].(string)
	if !ok {
		certFileName = controllers.ControllerCertName
	}

	cert, key, clientCert, err := controllers.GenerateCertificate(releaseName, namespace)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretName,
			Namespace:   namespace,
			Labels:      controllers.CommonLabels(map[string]string{}, releaseName, controllers.ControllerServiceName),
			Annotations: map[string]string{},
		},
		Data: map[string][]byte{
			keyFileName:  key.Bytes(),
			certFileName: cert.Bytes(),
		},
	}
	secret.Labels["app.kubernetes.io/managed-by"] = aperturectl

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-client-cert", releaseName),
			Namespace:   namespace,
			Labels:      controllers.CommonLabels(map[string]string{}, releaseName, controllers.ControllerServiceName),
			Annotations: map[string]string{},
		},
		Data: map[string]string{
			controllers.ControllerClientCertKey: clientCert.String(),
		},
	}
	cm.Labels["app.kubernetes.io/managed-by"] = aperturectl

	if generateCert {
		createNew, err := CheckCertificate(secret, cm, certFileName, keyFileName)
		if err != nil {
			return nil, err
		}

		if createNew {
			unstructuredSecret, err := runtime.DefaultUnstructuredConverter.ToUnstructured(secret)
			if err != nil {
				return nil, err
			}

			gvk := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}
			unstructuredObject := &unstructured.Unstructured{Object: unstructuredSecret}
			unstructuredObject.SetGroupVersionKind(gvk)
			err = applyObjectToKubernetesWithRetry(unstructuredObject)
			if err != nil {
				return nil, err
			}

			var unstructuredCM map[string]interface{}
			unstructuredCM, err = runtime.DefaultUnstructuredConverter.ToUnstructured(cm)
			if err != nil {
				return nil, err
			}

			gvk = schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}
			unstructuredObject = &unstructured.Unstructured{Object: unstructuredCM}
			unstructuredObject.SetGroupVersionKind(gvk)
			err = applyObjectToKubernetesWithRetry(unstructuredObject)
			if err != nil {
				return nil, err
			}
		}
	}

	serverCert["secretName"] = secretName
	controller["serverCert"] = serverCert
	values["controller"] = controller
	return values, nil
}

// CheckCertificate checks if the certificate in the secret is valid.
func CheckCertificate(secret *corev1.Secret, cm *corev1.ConfigMap, certFileName, keyFileName string) (bool, error) {
	existingSecret := &corev1.Secret{}
	name := secret.GetName()
	err := kubeClient.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, existingSecret)
	if apierrors.IsNotFound(err) {
		return true, nil
	} else {
		existingCert, ok := existingSecret.Data[certFileName]
		if !ok {
			return false, fmt.Errorf(
				"failed to create Aperture Controller Certificate secret as a secret named '%s' already exists without '%s' certificate key", certFileName, name)
		}

		block, _ := pem.Decode(existingCert)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return true, nil
		}

		if time.Now().After(cert.NotAfter) {
			return true, nil
		}

		existingCM := &corev1.ConfigMap{}
		name = cm.GetName()
		err = kubeClient.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, existingCM)
		if apierrors.IsNotFound(err) {
			secret.Data[certFileName] = existingCert
			secret.Data[keyFileName] = existingSecret.Data[keyFileName]
			return true, nil
		}
	}

	return false, nil
}

// manageAgentControllerClientCertConfigMap manages the agent controller client certificate ConfigMap.
func manageAgentControllerClientCertConfigMap(values map[string]interface{}, releaseName, namespace string, order releaseutil.KindSortOrder) (map[string]interface{}, error) {
	agent, ok := values["agent"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to get agent values")
	}

	config, ok := agent["config"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to get agent config values")
	}

	controllerCert := agent["controllerCert"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to get controllerCert values")
	}

	cmName, ok := controllerCert["cmName"].(string)
	if !ok {
		cmName = fmt.Sprintf("%s-client-cert", controllers.AgentServiceName)
	}

	certFileName, ok := controllerCert["certFileName"].(string)
	if !ok {
		certFileName = controllers.ControllerClientCertKey
	}

	agentFunctions, ok := config["agent_functions"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	endpoints, ok := agentFunctions["endpoints"].([]interface{})
	if !ok {
		return nil, nil
	}

	// Convert the endpoints to []string
	var endpointsStr []string
	for _, endpoint := range endpoints {
		endpointStr, ok := endpoint.(string)
		if !ok {
			continue
		}
		endpointsStr = append(endpointsStr, endpointStr)
	}

	controllerClientCert := controllers.GetControllerClientCert(endpointsStr, kubeClient, context.Background())
	if controllerClientCert == nil {
		return nil, nil
	}

	if slices.Equal(order, releaseutil.InstallOrder) {
		existingCM := &corev1.ConfigMap{}
		err := kubeClient.Get(context.Background(), types.NamespacedName{Name: cmName, Namespace: namespace}, existingCM)
		if err == nil ||
			apierrors.IsNotFound(err) || (existingCM.Labels != nil && existingCM.Labels["app.kubernetes.io/managed-by"] == "aperturectl") {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:        cmName,
					Namespace:   namespace,
					Labels:      controllers.CommonLabels(map[string]string{}, releaseName, controllers.AgentServiceName),
					Annotations: map[string]string{},
				},
				Data: map[string]string{
					certFileName: string(controllerClientCert),
				},
			}

			cm.Labels["app.kubernetes.io/managed-by"] = aperturectl
			var unstructuredCM map[string]interface{}
			unstructuredCM, err = runtime.DefaultUnstructuredConverter.ToUnstructured(cm)
			if err != nil {
				return nil, err
			}

			gvk := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}
			unstructured := &unstructured.Unstructured{Object: unstructuredCM}
			unstructured.SetGroupVersionKind(gvk)
			err = applyObjectToKubernetesWithRetry(unstructured)
			if err != nil {
				return nil, err
			}
		}
	}

	controllerCert["cmName"] = cmName
	controllerCert["certFileName"] = certFileName
	agent["controllerCert"] = controllerCert
	values["agent"] = agent
	return values, nil
}
