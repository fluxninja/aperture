/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/fluxninja/aperture/v2/operator/controllers"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/fluxninja/aperture/v2/operator/api"
	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/net/tlsconfig"
	"github.com/go-logr/logr"
)

// ControllerReconciler reconciles a Controller object.
type ControllerReconciler struct {
	client.Client
	DynamicClient       dynamic.Interface
	Scheme              *runtime.Scheme
	Recorder            record.EventRecorder
	MultipleControllers bool
	ResourcesDeleted    map[string]bool
	defaultsExecuted    bool
	mutex               sync.RWMutex
}

var (
	controllerCert       *bytes.Buffer
	controllerKey        *bytes.Buffer
	controllerClientCert *bytes.Buffer
)

func (r *ControllerReconciler) verifySingletonControllerExists(ctx context.Context, instance *controllerv1alpha1.Controller, instances *controllerv1alpha1.ControllerList) (bool, error) {
	// Check if this is an update request. If not, skip the resource creation as only single Controller is required in a cluster.
	if instances.Items != nil && len(instances.Items) != 0 {
		for _, ins := range instances.Items {
			if ins.GetNamespace() == instance.GetNamespace() && ins.GetName() == instance.GetName() {
				continue
			}
			if ins.GetDeletionTimestamp() == nil && (ins.Status.Resources == "creating" || ins.Status.Resources == "created") {
				r.Recorder.Event(instance, corev1.EventTypeWarning, "ResourcesExist",
					"The required resources are already deployed. Skipping resource creation as currently, the Controller does not support multiple replicas.")

				instance.Status.Resources = "skipped"
				if err := r.updateStatus(ctx, instance.DeepCopy()); err != nil {
					return false, err
				}

				return false, nil
			}
		}
	}
	return true, nil
}

//+kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=validatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fluxninja.com,resources=controllers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fluxninja.com,resources=controllers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fluxninja.com,resources=controllers/finalizers,verbs=update
//+kubebuilder:rbac:groups=fluxninja.com,resources=policies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fluxninja.com,resources=policies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fluxninja.com,resources=policies/finalizers,verbs=update
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Controller object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ControllerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	instance := &controllerv1alpha1.Controller{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil && errors.IsNotFound(err) {
		r.mutex.RLock()
		defer r.mutex.RUnlock()
		if deleted, ok := r.ResourcesDeleted[req.NamespacedName.String()]; !ok || !deleted {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and do not requeue
			logger.Info("Controller resource not found. Ignoring since object must be deleted")
		}
		return ctrl.Result{}, nil
	} else if err != nil {
		// Error reading the object - requeue the request.
		logger.Error(err, "failed to get Controller")
		return ctrl.Result{}, err
	}

	// Checking if the Minimum kubernetes version condition is satisfied.
	if len(instance.Status.Resources) == 0 && !controllers.MinimumKubernetesVersionBool {
		r.Recorder.Event(instance, corev1.EventTypeWarning, "MinimumKubernetesVersionFail",
			fmt.Sprintf("Kubernetes version %v is not supported. Please use a kubernetes cluster with version %v or above.", controllers.CurrentKubernetesVersion.String(), controllers.MinimumKubernetesVersion))
		instance.Status.Resources = "failed"
		err = r.updateStatus(ctx, instance.DeepCopy())
		return ctrl.Result{}, err
	}

	// Handing delete operation
	if instance.GetDeletionTimestamp() != nil {
		logger.Info(fmt.Sprintf("Handling deletion of resources for Instance '%s' in Namespace '%s'", instance.GetName(), instance.GetNamespace()))
		if controllerutil.ContainsFinalizer(instance, controllers.FinalizerName) {
			if !r.MultipleControllers {
				r.deleteResources(ctx, logger, instance.DeepCopy())
			}

			controllerutil.RemoveFinalizer(instance, controllers.FinalizerName)
			if err = r.updateController(ctx, instance); err != nil && !errors.IsNotFound(err) {
				return ctrl.Result{}, err
			}
		}

		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.ResourcesDeleted[req.NamespacedName.String()] = true
		return ctrl.Result{}, nil
	}

	instances := &controllerv1alpha1.ControllerList{}
	err = r.List(ctx, instances)
	if err != nil {
		logger.Error(err, "failed to list Controller")
		return ctrl.Result{}, err
	}

	if !r.MultipleControllers {
		passed, innerErr := r.verifySingletonControllerExists(ctx, instance, instances)
		if innerErr != nil {
			return ctrl.Result{}, innerErr
		}
		if !passed {
			return ctrl.Result{}, nil
		}
	}

	instance.Status.Resources = "creating"
	if err = r.updateStatus(ctx, instance.DeepCopy()); err != nil {
		return ctrl.Result{}, err
	}

	if instance.Annotations == nil || instance.Annotations[controllers.DefaulterAnnotationKey] != "true" || !r.defaultsExecuted {
		err = r.checkDefaults(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		} else if instance.Status.Resources == controllers.FailedStatus {
			return ctrl.Result{}, nil
		}
		r.defaultsExecuted = true
	}

	r.mutex.Lock()
	r.ResourcesDeleted[req.NamespacedName.String()] = false
	r.mutex.Unlock()
	if err = r.manageResources(ctx, logger, instance); err != nil {
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(instance, controllers.FinalizerName) {
		controllerutil.AddFinalizer(instance, controllers.FinalizerName)
	}

	if !instance.Status.IsMigrationCompleted {
		// Reloading instances for the case when the migration is completed by another controller.
		err = r.List(ctx, instances)
		if err != nil {
			logger.Error(err, "failed to list Controller")
			return ctrl.Result{}, err
		}

		doMigration := true
		isCRNamedControllerPresent := false
		for _, ins := range instances.Items {
			if ins.GetName() == controllers.ControllerName {
				isCRNamedControllerPresent = true
			}

			if ins.GetUID() == instance.GetUID() {
				continue
			}

			if ins.Status.IsMigrationCompleted {
				doMigration = false
				break
			}
		}

		if doMigration {
			if err := r.deleteOlderInstances(ctx, logger, instance, isCRNamedControllerPresent); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if err := r.updateController(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	instance.Status.IsMigrationCompleted = true
	instance.Status.Resources = "created"
	if err := r.updateStatus(ctx, instance.DeepCopy()); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// updateController updates the Controller resource in Kubernetes.
func (r *ControllerReconciler) updateController(ctx context.Context, instance *controllerv1alpha1.Controller) error {
	attempt := 5
	finalizers := instance.DeepCopy().Finalizers
	spec := instance.DeepCopy().Spec
	annotations := instance.DeepCopy().Annotations
	for attempt > 0 {
		attempt -= 1
		if err := r.Update(ctx, instance); err != nil {
			if errors.IsConflict(err) {
				namespacesName := types.NamespacedName{
					Namespace: instance.GetNamespace(),
					Name:      instance.GetName(),
				}
				if err = r.Get(ctx, namespacesName, instance); err != nil {
					return err
				}
				instance.Finalizers = finalizers
				instance.Spec = spec
				instance.Annotations = annotations
				continue
			}
			return err
		}
	}

	return nil
}

// updateResource updates the Controller resource in Kubernetes.
func (r *ControllerReconciler) updateStatus(ctx context.Context, instance *controllerv1alpha1.Controller) error {
	attempt := 5
	status := instance.DeepCopy().Status
	for attempt > 0 {
		attempt -= 1
		if err := r.Status().Update(ctx, instance); err != nil {
			if errors.IsConflict(err) {
				namespacesName := types.NamespacedName{
					Namespace: instance.GetNamespace(),
					Name:      instance.GetName(),
				}
				if err = r.Get(ctx, namespacesName, instance); err != nil {
					return err
				}
				instance.Status = status
				continue
			}
			return err
		}
	}
	return nil
}

func (r *ControllerReconciler) deleteOlderInstances(ctx context.Context, log logr.Logger, instance *controllerv1alpha1.Controller, isCRNamedControllerPresent bool) error {
	singletonInstance := instance.DeepCopy()
	singletonInstance.Name = controllers.ControllerName

	if !isCRNamedControllerPresent {
		deployment, err := deploymentForController(singletonInstance.DeepCopy(), log, r.Scheme)
		if err != nil || deployment == nil {
			log.Error(err, "Failed to create object for old Deployment during Migration")
		} else {
			if err = r.Delete(ctx, deployment); err != nil && !errors.IsNotFound(err) {
				log.Error(err, "Failed to delete old Deployment during Migration")
			}
		}

		if singletonInstance.Spec.ServiceAccountSpec.Create {
			var serviceAccount *corev1.ServiceAccount
			serviceAccount, err = serviceAccountForController(singletonInstance.DeepCopy(), r.Scheme)
			if err != nil || serviceAccount == nil {
				log.Error(err, "Failed to create object for old ServiceAccount during Migration")
			} else {
				if err = r.Delete(ctx, serviceAccount); err != nil && !errors.IsNotFound(err) {
					log.Error(err, "Failed to delete old ServiceAccount during Migration")
				}
			}
		}

		configMap, err := configMapForControllerConfig(singletonInstance.DeepCopy(), r.Scheme)
		if err != nil || configMap == nil {
			log.Error(err, "Failed to create object for old ConfigMap during Migration")
		} else {
			if err = r.Delete(ctx, configMap); err != nil && !errors.IsNotFound(err) {
				log.Error(err, "Failed to delete old ConfigMap during Migration")
			}
		}

		service, err := serviceForController(singletonInstance.DeepCopy(), log, r.Scheme)
		if err != nil || service == nil {
			log.Error(err, "Failed to create object for old Service during Migration")
		} else {
			if err = r.Delete(ctx, service); err != nil && !errors.IsNotFound(err) {
				log.Error(err, "Failed to delete old Service during Migration")
			}
		}
	}

	var certBytes []byte
	if controllerClientCert == nil {
		certBytes = []byte{}
	} else {
		certBytes = controllerClientCert.Bytes()
	}

	vwc := validatingWebhookConfiguration(singletonInstance.DeepCopy(), certBytes)
	vwc.Name = controllers.ControllerServiceName
	if err := r.Delete(ctx, vwc); err != nil && !errors.IsNotFound(err) {
		log.Error(err, "Failed to delete old ValidatingWebhookConfiguration during Migration")
	}

	crb := clusterRoleBindingForController(singletonInstance.DeepCopy())
	crb.Name = controllers.ControllerServiceName
	if err := r.Delete(ctx, crb); err != nil && !errors.IsNotFound(err) {
		log.Error(err, "Failed to delete old ClusterRoleBinding during Migration")
	}

	cr := clusterRoleForController(singletonInstance.DeepCopy())
	cr.Name = controllers.ControllerServiceName
	if err := r.Delete(ctx, cr); err != nil && !errors.IsNotFound(err) {
		log.Error(err, "Failed to delete old ClusterRole during Migration")
	}

	return nil
}

// deleteResources deletes cluster-scoped resources for which owner-reference is not added.
func (r *ControllerReconciler) deleteResources(ctx context.Context, log logr.Logger, instance *controllerv1alpha1.Controller) {
	if err := r.Delete(ctx, clusterRoleForController(instance)); err != nil {
		log.Error(err, "failed to delete object of ClusterRole")
	}

	if err := r.Delete(ctx, clusterRoleBindingForController(instance)); err != nil {
		log.Error(err, "failed to delete object of ClusterRoleBinding")
	}

	if err := r.Delete(ctx, validatingWebhookConfiguration(instance, nil)); err != nil {
		log.Error(err, "failed to delete object of ValidatingWebhookConfiguration")
	}
}

// checkDefaults checks and sets defaults when the Defaulter webhook is not triggered.
func (r *ControllerReconciler) checkDefaults(ctx context.Context, instance *controllerv1alpha1.Controller) error {
	resource, err := r.DynamicClient.Resource(api.GroupVersion.WithResource("controllers")).Namespace(instance.GetNamespace()).Get(ctx, instance.GetName(), v1.GetOptions{})
	if err != nil {
		instance.Status.Resources = controllers.FailedStatus
		r.Recorder.Eventf(instance, corev1.EventTypeWarning, "FailedToFetch", "Failed to fetch Resource. Error: '%s'", err.Error())
		errUpdate := r.updateStatus(ctx, instance)
		if errUpdate != nil {
			return errUpdate
		}
		return nil
	}

	resourceBytes, err := resource.MarshalJSON()
	if err != nil {
		instance.Status.Resources = controllers.FailedStatus
		r.Recorder.Eventf(instance, corev1.EventTypeWarning, "FailedToMarshal", "Failed to marshal Resource into JSON. Error: '%s'", err.Error())
		errUpdate := r.updateStatus(ctx, instance)
		if errUpdate != nil {
			return errUpdate
		}
		return nil
	}

	err = config.UnmarshalYAML(resourceBytes, instance)
	if err != nil {
		instance.Status.Resources = controllers.FailedStatus
		r.Recorder.Eventf(instance, corev1.EventTypeWarning, "FailedToSetDefaults", "Failed to set defaults. Error: '%s'", err.Error())
		errUpdate := r.updateStatus(ctx, instance)
		if errUpdate != nil {
			return errUpdate
		}
		return nil
	}

	if instance.Spec.Secrets.FluxNinjaExtension.Create && instance.Spec.Secrets.FluxNinjaExtension.Value == "" {
		instance.Status.Resources = controllers.FailedStatus
		r.Recorder.Eventf(instance, corev1.EventTypeWarning, "ValidationFailed", "The value for 'spec.secrets.fluxNinjaExtension.value' can not be empty when 'spec.secrets.fluxNinjaExtension.create' is set to true")
		errUpdate := r.updateStatus(ctx, instance)
		if errUpdate != nil {
			return errUpdate
		}
		return nil
	}

	if (instance.Spec.Image.Digest == "" && instance.Spec.Image.Tag == "") || (instance.Spec.Image.Digest != "" && instance.Spec.Image.Tag != "") {
		instance.Status.Resources = controllers.FailedStatus
		r.Recorder.Eventf(instance, corev1.EventTypeWarning, "ValidationFailed", "Either 'spec.image.digest' or 'spec.image.tag' should be provided.")
		errUpdate := r.updateStatus(ctx, instance)
		if errUpdate != nil {
			return errUpdate
		}
		return nil
	}

	if instance.Status.Resources == controllers.FailedStatus {
		instance.Status.Resources = ""
	}
	return nil
}

// manageResources creates/updates required resources.
func (r *ControllerReconciler) manageResources(ctx context.Context, log logr.Logger, instance *controllerv1alpha1.Controller) error {
	certFilePath := instance.Spec.ConfigSpec.Server.TLS.CertFile
	if certFilePath == "" {
		certFilePath = path.Join(controllers.ControllerCertPath, controllers.ControllerCertName)
	}

	keyFilePath := instance.Spec.ConfigSpec.Server.TLS.KeyFile
	if keyFilePath == "" {
		keyFilePath = path.Join(controllers.ControllerCertPath, controllers.ControllerCertKeyName)
	}

	instance.Spec.ConfigSpec.Server.TLS = tlsconfig.ServerTLSConfig{
		CertFile: certFilePath,
		KeyFile:  keyFilePath,
		Enabled:  true,
	}

	if !r.MultipleControllers {
		instance.Spec.ConfigSpec.Policies.CRWatcher.Enabled = true
		if err := r.reconcileClusterRole(ctx, instance); err != nil {
			return err
		}

		if err := r.reconcileClusterRoleBinding(ctx, instance); err != nil {
			return err
		}

		if err := r.reconcileValidatingWebhookConfigurationAndCertSecret(ctx, instance); err != nil {
			return err
		}
	}

	if err := r.reconcileService(ctx, log, instance); err != nil {
		return err
	}

	if err := r.reconcileConfigMap(ctx, instance); err != nil {
		return err
	}

	if err := r.reconcileServiceAccount(ctx, log, instance); err != nil {
		return err
	}

	if err := r.reconcileDeployment(ctx, log, instance); err != nil {
		return err
	}

	if err := r.reconcileSecret(ctx, instance); err != nil {
		return err
	}

	return nil
}

// reconcileConfigMap prepares the desired states for Controller configmaps and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *ControllerReconciler) reconcileConfigMap(ctx context.Context, instance *controllerv1alpha1.Controller) error {
	configMap, err := configMapForControllerConfig(instance.DeepCopy(), r.Scheme)
	if err != nil {
		return err
	}

	if err = ctrl.SetControllerReference(instance, configMap, r.Scheme); err != nil {
		return err
	}

	if _, err = createConfigMapForController(r.Client, r.Recorder, configMap, ctx, instance); err != nil {
		return err
	}

	return nil
}

// reconcileService prepares the desired states for Controller services and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *ControllerReconciler) reconcileService(ctx context.Context, log logr.Logger, instance *controllerv1alpha1.Controller) error {
	service, err := serviceForController(instance.DeepCopy(), log, r.Scheme)
	if err != nil {
		return err
	}
	if err = r.createService(service, ctx, instance); err != nil {
		return err
	}

	return nil
}

// reconcileClusterRole prepares the desired states for Controller ClusterRole and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *ControllerReconciler) reconcileClusterRole(ctx context.Context, instance *controllerv1alpha1.Controller) error {
	cr := clusterRoleForController(instance.DeepCopy())
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, cr, controllers.ClusterRoleMutate(cr, cr.Rules))
	if err != nil {
		if errors.IsConflict(err) {
			return r.reconcileClusterRole(ctx, instance)
		}

		msg := fmt.Sprintf("failed to create ClusterRole '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			cr.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		r.Recorder.Event(instance, corev1.EventTypeNormal, "ClusterRoleCreationFailed", msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ClusterRoleCreationSuccessful", "Created ClusterRole '%s'", cr.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ClusterRoleUpdationSuccessful", "Updated ClusterRole '%s'", cr.GetName())
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}

// reconcileClusterRoleBinding prepares the desired states for Controller ClusterRoleBinding and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *ControllerReconciler) reconcileClusterRoleBinding(ctx context.Context, instance *controllerv1alpha1.Controller) error {
	log := log.FromContext(ctx)
	crb := clusterRoleBindingForController(instance.DeepCopy())
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, crb, controllers.ClusterRoleBindingMutate(crb, crb.RoleRef, crb.Subjects))
	if err != nil {
		if errors.IsConflict(err) {
			return r.reconcileClusterRoleBinding(ctx, instance)
		}

		// Checking invalid as Kubernetes does not allow updating RoleRef
		if errors.IsInvalid(err) && strings.Contains(err.Error(), "cannot change roleRef") {
			if err = r.Delete(ctx, clusterRoleBindingForController(instance)); err != nil {
				log.Error(err, "failed to delete object of ClusterRoleBinding")
			}
			return r.reconcileClusterRoleBinding(ctx, instance)
		}

		msg := fmt.Sprintf("failed to create ClusterRoleBinding '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			crb.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		r.Recorder.Event(instance, corev1.EventTypeNormal, "ClusterRoleBindingCreationFailed", msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ClusterRoleBindingCreationSuccessful", "Created ClusterRoleBinding '%s'", crb.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ClusterRoleBindingUpdationSuccessful", "Updated ClusterRoleBinding '%s'", crb.GetName())
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}

// reconcileServiceAccount prepares the desired states for Controller ServiceAccount and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *ControllerReconciler) reconcileServiceAccount(ctx context.Context, log logr.Logger, instance *controllerv1alpha1.Controller) error {
	if !instance.Spec.ServiceAccountSpec.Create {
		return nil
	}

	sa, err := serviceAccountForController(instance.DeepCopy(), r.Scheme)
	if err != nil {
		return err
	}

	if err = ctrl.SetControllerReference(instance, sa, r.Scheme); err != nil {
		return err
	}

	if err = r.createServiceAccount(sa, ctx, instance); err != nil {
		return err
	}

	return nil
}

// reconcileDeployment prepares the desired states for Controller Deployment and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *ControllerReconciler) reconcileDeployment(ctx context.Context, log logr.Logger, instance *controllerv1alpha1.Controller) error {
	dep, err := deploymentForController(instance.DeepCopy(), log, r.Scheme)
	if err != nil {
		return err
	}

	if err = ctrl.SetControllerReference(instance, dep, r.Scheme); err != nil {
		return err
	}

	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, dep, deploymentMutate(dep, dep.Spec))
	if err != nil {
		if errors.IsConflict(err) {
			return r.reconcileDeployment(ctx, log, instance)
		}

		msg := fmt.Sprintf("failed to create Deployment '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			dep.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		r.Recorder.Event(instance, corev1.EventTypeNormal, "DeploymentCreationFailed", msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "DeploymentCreationSuccessful", "Created Deployment '%s'", dep.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "DeploymentUpdationSuccessful", "Updated Deployment '%s'", dep.GetName())
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}

// reconcileValidatingWebhookConfigurationAndCertSecret prepares the desired states for ValidatingWebhookConfiguration and
// secret for its certificate and sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *ControllerReconciler) reconcileValidatingWebhookConfigurationAndCertSecret(ctx context.Context, instance *controllerv1alpha1.Controller) error {
	var err error
	if controllerCert == nil || controllerKey == nil || controllerClientCert == nil {
		controllerCert, controllerKey, controllerClientCert, err = controllers.GetOrGenerateCertificate(r.Client, instance.DeepCopy())
		if err != nil {
			return err
		}
	}

	secret, err := secretForControllerCert(instance.DeepCopy(), r.Scheme, controllerCert, controllerKey)
	if err != nil {
		return err
	}
	if _, err = createSecretForController(r.Client, r.Recorder, secret, ctx, instance); err != nil {
		return err
	}

	cm, err := configMapForControllerClientCert(instance.DeepCopy(), r.Scheme, controllerClientCert)
	if err != nil {
		return err
	}
	if _, err = createConfigMapForController(r.Client, r.Recorder, cm, ctx, instance); err != nil {
		return err
	}

	vwc := validatingWebhookConfiguration(instance.DeepCopy(), controllerClientCert.Bytes())
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, vwc, controllers.ValidatingWebhookConfigurationMutate(vwc, vwc.Webhooks))
	if err != nil {
		if errors.IsConflict(err) {
			return r.reconcileValidatingWebhookConfigurationAndCertSecret(ctx, instance)
		}

		msg := fmt.Sprintf("failed to create ValidatingWebhookConfiguration '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			vwc.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		r.Recorder.Event(instance, corev1.EventTypeNormal, "ValidatingWebhookConfigurationCreationFailed", msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal,
			"ValidatingWebhookConfigurationCreationSuccessful", "Created ValidatingWebhookConfiguration '%s'", vwc.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal,
			"ValidatingWebhookConfigurationUpdationSuccessful", "Updated ValidatingWebhookConfiguration '%s'", vwc.GetName())
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}

// reconcileSecret prepares the desired states for Controller ApiKey secret and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *ControllerReconciler) reconcileSecret(ctx context.Context, instance *controllerv1alpha1.Controller) error {
	if !instance.Spec.Secrets.FluxNinjaExtension.Create {
		return nil
	}
	secret, err := secretForControllerAPIKey(instance.DeepCopy(), r.Scheme)
	if err != nil {
		return err
	}
	if _, err = createSecretForController(r.Client, r.Recorder, secret, ctx, instance); err != nil {
		return err
	}

	instance.Spec.Secrets.FluxNinjaExtension.Create = false
	instance.Spec.Secrets.FluxNinjaExtension.Value = ""
	instance.Spec.Secrets.FluxNinjaExtension.SecretKeyRef.Name = controllers.SecretName(
		instance.GetName(), "controller", &instance.Spec.Secrets.FluxNinjaExtension)

	return nil
}

// eventFiltersForController sets up a Predicate filter for the received events.
func eventFiltersForController() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(create event.CreateEvent) bool {
			_, ok := create.Object.(*controllerv1alpha1.Controller)
			return ok
		},
		UpdateFunc: func(update event.UpdateEvent) bool {
			new, ok1 := update.ObjectNew.(*controllerv1alpha1.Controller)
			old, ok2 := update.ObjectOld.(*controllerv1alpha1.Controller)
			if !ok1 || !ok2 {
				return true
			}

			if new.GetDeletionTimestamp() != nil {
				return true
			}

			diffObjects := !reflect.DeepEqual(old.Spec, new.Spec)
			// Skipping update events for Secret updates
			if diffObjects && old.Spec.Secrets.FluxNinjaExtension.Value != "" &&
				new.Spec.Secrets.FluxNinjaExtension.Value == "" {
				return false
			}

			return diffObjects
		},
		DeleteFunc: func(delete event.DeleteEvent) bool {
			return !delete.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&controllerv1alpha1.Controller{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.ServiceAccount{}).
		WithEventFilter(eventFiltersForController()).
		Complete(r)
}
