package installation

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fluxninja/aperture/pkg/log"
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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"
)

var (
	valuesFile     string
	kubeConfig     string
	kubeRestConfig *rest.Config
	version        string
	latestVersion  string
	namespace      string
	kubeClient     client.Client
)

const (
	apertureLatestVersion = "latest"
	apertureControllerNS  = "aperture-controller"
	apertureAgentNS       = "aperture-agent"
	controller            = "controller"
	agent                 = "agent"
)

// getTemplets loads CRDs, hooks and manifests from the Helm chart.
func getTemplets(chartName string, order releaseutil.KindSortOrder) ([]chart.CRD, []*release.Hook, []releaseutil.Manifest, error) {
	chartURL := fmt.Sprintf("https://fluxninja.github.io/aperture/aperture-%s-%s.tgz", chartName, version)

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
		values, err = chartutil.ReadValuesFile(valuesFile)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read values: %s", err)
		}
	} else if chartName == agent && slices.Compare(order, releaseutil.UninstallOrder) == 0 {
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

	renderedValues, err := chartutil.ToRenderValues(ch, values, chartutil.ReleaseOptions{Name: chartName, Namespace: namespace}, chartutil.DefaultCapabilities)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read values: %s", err)
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

	return ch.CRDObjects(), hooks, manifests, err
}

// applyManifest creates/updates the generated manifest to Kubernetes.
func applyManifest(manifest string) error {
	content := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(manifest), &content)
	if err != nil {
		return err
	}

	unstructuredObject := &unstructured.Unstructured{
		Object: content,
	}

	log.Info().Msgf("Installing - %s/%s", unstructuredObject.GetKind(), unstructuredObject.GetName())
	spec := unstructuredObject.Object["spec"]
	_, err = controllerutil.CreateOrUpdate(context.Background(), kubeClient, unstructuredObject, func() error {
		unstructuredObject.Object["spec"] = spec
		return nil
	})
	if err != nil && strings.Contains(err.Error(), "no matches for kind") {
		attempt := 0
		for attempt < 5 {
			time.Sleep(time.Second * time.Duration(attempt))
			if _, err = controllerutil.CreateOrUpdate(context.Background(), kubeClient, unstructuredObject, func() error {
				unstructuredObject.Object["spec"] = spec
				return nil
			}); err == nil {
				return nil
			}
		}
	}
	return err
}

// deleteManifest deletes the generated manifest from Kubernetes.
func deleteManifest(manifest string) error {
	content := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(manifest), &content)
	if err != nil {
		return err
	}

	unstructuredObject := &unstructured.Unstructured{
		Object: content,
	}
	log.Info().Msgf("Deleting - %s/%s", unstructuredObject.GetKind(), unstructuredObject.GetName())

	err = kubeClient.Delete(context.Background(), unstructuredObject)
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
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
