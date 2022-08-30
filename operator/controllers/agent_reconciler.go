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

package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// AgentReconciler reconciles a Agent object.
type AgentReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	Recorder         record.EventRecorder
	ApertureInjector *ApertureInjector
	resourcesDeleted bool
}

//+kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fluxninja.com,resources=agents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fluxninja.com,resources=agents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fluxninja.com,resources=agents/finalizers,verbs=update
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=componentstatuses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
//+kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=namespaces/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=nodes/metrics,verbs=get
//+kubebuilder:rbac:groups=core,resources=nodes/spec,verbs=get
//+kubebuilder:rbac:groups=core,resources=nodes/proxy,verbs=get
//+kubebuilder:rbac:groups=core,resources=nodes/stats,verbs=get
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=policy,resources=podsecuritypolicies,verbs=use
//+kubebuilder:rbac:groups=quota.openshift.io,resources=clusterresourcequotas,verbs=get;list
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=security.openshift.io,resources=securitycontextconstraints,verbs=use
//+kubebuilder:rbac:urls=/version;/healthz;/metrics,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Agent object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *AgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	instance := &v1alpha1.Agent{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil && errors.IsNotFound(err) {
		if !r.resourcesDeleted {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Agent resource not found. Ignoring since object must be deleted")
		}
		return ctrl.Result{}, nil
	} else if err != nil {
		// Error reading the object - requeue the request.
		logger.Error(err, "failed to get Aperture")
		return ctrl.Result{}, err
	}

	// Handing delete operation
	if instance.GetDeletionTimestamp() != nil {
		logger.Info(fmt.Sprintf("Handling deletion of resources for Instance '%s' in Namespace '%s'", instance.GetName(), instance.GetNamespace()))
		if controllerutil.ContainsFinalizer(instance, finalizerName) {
			if err = r.deleteResources(ctx, logger, instance.DeepCopy()); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(instance, finalizerName)
			if err = r.updateAgent(ctx, instance); err != nil && !errors.IsNotFound(err) {
				return ctrl.Result{}, err
			}
		}

		r.resourcesDeleted = true
		return ctrl.Result{}, nil
	}

	instances := &v1alpha1.AgentList{}
	err = r.List(ctx, instances)
	if err != nil {
		logger.Error(err, "failed to list Agent")
		return ctrl.Result{}, err
	}

	// Check if this is an update request. If not, skip the resource creation as only single Controller is required in a cluster.
	if instances.Items != nil && len(instances.Items) != 0 {
		for _, ins := range instances.Items {
			if ins.GetNamespace() == instance.GetNamespace() && ins.GetName() == instance.GetName() {
				continue
			}
			if ins.GetDeletionTimestamp() == nil && (ins.Status.Resources == "creating" || ins.Status.Resources == "created") {
				r.Recorder.Event(instance, corev1.EventTypeWarning, "ResourcesExist",
					"The required resources are already deployed. Skipping resource creation as currently, the Agent doesn't require multiple replicas.")

				instance.Status.Resources = "skipped"
				if err := r.updateStatus(ctx, instance.DeepCopy()); err != nil {
					return ctrl.Result{}, err
				}

				return ctrl.Result{}, nil
			}
		}
	}

	instance.Status.Resources = "creating"
	if err := r.updateStatus(ctx, instance.DeepCopy()); err != nil {
		return ctrl.Result{}, err
	}
	r.resourcesDeleted = false

	if err := r.manageResources(ctx, logger, instance); err != nil {
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(instance, finalizerName) {
		controllerutil.AddFinalizer(instance, finalizerName)
	}

	if err := r.updateAgent(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	instance.Status.Resources = "created"
	if err := r.updateStatus(ctx, instance.DeepCopy()); err != nil {
		return ctrl.Result{}, err
	}

	r.ApertureInjector.Instance = instance
	return ctrl.Result{}, nil
}

// updateResource updates the Aperture resource in Kubernetes.
func (r *AgentReconciler) updateStatus(ctx context.Context, instance *v1alpha1.Agent) error {
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

// deleteResources deletes cluster-scoped resources for which owner-reference is not added.
func (r *AgentReconciler) deleteResources(ctx context.Context, log logr.Logger, instance *v1alpha1.Agent) error {
	deleteClusterRole := true
	instances := &v1alpha1.ControllerList{}
	err := r.List(ctx, instances)
	if err != nil {
		log.Error(err, "failed to list Controller")
		return err
	}
	if instances.Items != nil && len(instances.Items) != 0 {
		for _, ins := range instances.Items {
			if ins.Status.Resources == "created" {
				deleteClusterRole = false
			}
		}
	}
	if deleteClusterRole {
		if err := r.Delete(ctx, clusterRoleForAgent(instance)); err != nil && !errors.IsNotFound(err) {
			log.Error(err, "failed to delete object of ClusterRole")
			return err
		}
	}

	if err := r.Delete(ctx, clusterRoleBindingForAgent(instance)); err != nil && !errors.IsNotFound(err) {
		log.Error(err, "failed to delete object of ClusterRoleBinding")
		return err
	}

	if instance.Spec.Sidecar.Enabled {
		mwc, err := mutatingWebhookConfiguration(instance)
		if err != nil {
			return err
		}
		if err = r.Delete(ctx, mwc); err != nil {
			log.Error(err, "failed to delete object of MutatingWebhookConfiguration")
			return err
		}

		nsList := &corev1.NamespaceList{}
		err = r.List(ctx, nsList)
		if err != nil {
			return fmt.Errorf("failed to list Namespaces. Error: %+v", err)
		}

		for _, ns := range nsList.Items {
			if ns.Labels == nil || ns.Labels[sidecarLabelKey] != enabled {
				continue
			}

			configMap, err := configMapForAgentConfig(instance, nil)
			if err != nil && !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("failed to fetch object of ConfigMap '%s' in namespace %s", configMap.GetName(), ns.GetName()))
			}

			configMap.Namespace = ns.GetName()
			configMap.Annotations = getAgentAnnotationsWithOwnerRef(instance)
			if err := r.Delete(ctx, configMap); err != nil && !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("failed to delete object of ConfigMap '%s' in namespace %s", configMap.GetName(), ns.GetName()))
			}

			if instance.Spec.FluxNinjaPlugin.APIKeySecret.Create && instance.Spec.FluxNinjaPlugin.Enabled {
				secret, err := secretForAgentAPIKey(instance, nil)
				if err != nil && !errors.IsNotFound(err) {
					log.Error(err, fmt.Sprintf("failed to delete object of Secret '%s' in namespace %s", configMap.GetName(), ns.GetName()))
				}

				secret.Namespace = ns.GetName()
				secret.Annotations = getAgentAnnotationsWithOwnerRef(instance)
				if err := r.Delete(ctx, secret); err != nil && !errors.IsNotFound(err) {
					log.Error(err, fmt.Sprintf("failed to delete object of Secret '%s' in namespace %s", configMap.GetName(), ns.GetName()))
				}
			}
		}
	}
	return nil
}

// updateAgent updates the Agent resource in Kubernetes.
func (r *AgentReconciler) updateAgent(ctx context.Context, instance *v1alpha1.Agent) error {
	attempt := 5
	finalizers := instance.DeepCopy().Finalizers
	spec := instance.DeepCopy().Spec
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
				continue
			}
			return err
		}
	}

	return nil
}

// manageResources creates/updates required resources.
func (r *AgentReconciler) manageResources(ctx context.Context, log logr.Logger, instance *v1alpha1.Agent) error {
	if err := r.reconcileConfigMap(ctx, instance); err != nil {
		return err
	}

	if err := r.reconcileService(ctx, log, instance); err != nil {
		return err
	}

	if err := r.reconcileClusterRole(ctx, instance); err != nil {
		return err
	}

	if err := r.reconcileClusterRoleBinding(ctx, instance); err != nil {
		return err
	}

	if err := r.reconcileServiceAccount(ctx, log, instance); err != nil {
		return err
	}

	if err := r.reconcileDaemonset(ctx, log, instance); err != nil {
		return err
	}

	if err := r.reconcileSecret(ctx, instance); err != nil {
		return err
	}

	if err := r.reconcileMutatingWebhookConfiguration(ctx, instance); err != nil {
		return err
	}

	if err := r.reconcileNamespacedResources(ctx, log, instance); err != nil {
		return err
	}
	return nil
}

// reconcileConfigMap prepares the desired states for Agent configmaps and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileConfigMap(ctx context.Context, instance *v1alpha1.Agent) error {
	if !instance.Spec.Sidecar.Enabled {
		configMap, err := configMapForAgentConfig(instance.DeepCopy(), r.Scheme)
		if err != nil {
			return err
		}
		if _, err = createConfigMapForAgent(r.Client, r.Recorder, configMap, ctx, instance); err != nil {
			return err
		}

	}

	return nil
}

// reconcileService prepares the desired states for Agent services and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileService(ctx context.Context, log logr.Logger, instance *v1alpha1.Agent) error {
	if !instance.Spec.Sidecar.Enabled {
		service, err := serviceForAgent(instance.DeepCopy(), log, r.Scheme)
		if err != nil {
			return err
		}
		if err = r.createService(service, ctx, instance); err != nil {
			return err
		}
	}

	return nil
}

// reconcileClusterRole prepares the desired states for Agent ClusterRole and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileClusterRole(ctx context.Context, instance *v1alpha1.Agent) error {
	cr := clusterRoleForAgent(instance.DeepCopy())
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, cr, clusterRoleMutate(cr, cr.Rules))
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
	case controllerutil.OperationResultNone:
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ClusterRoleCreationSuccessful", "Created ClusterRole '%s'", cr.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ClusterRoleUpdationSuccessful", "Updated ClusterRole '%s'", cr.GetName())
	default:
	}

	return nil
}

// reconcileClusterRoleBinding prepares the desired states for Agent ClusterRoleBinding and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileClusterRoleBinding(ctx context.Context, instance *v1alpha1.Agent) error {
	crb := clusterRoleBindingForAgent(instance.DeepCopy())
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, crb, clusterRoleBindingMutate(crb, crb.RoleRef, crb.Subjects))
	if err != nil {
		if errors.IsConflict(err) {
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

// reconcileServiceAccount prepares the desired states for Agent ServiceAccount and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileServiceAccount(ctx context.Context, log logr.Logger, instance *v1alpha1.Agent) error {
	if instance.Spec.ServiceAccountSpec.Create && !instance.Spec.Sidecar.Enabled {
		sa, err := serviceAccountForAgent(instance.DeepCopy(), r.Scheme)
		if err != nil {
			return err
		}
		if err = r.createServiceAccount(sa, ctx, instance); err != nil {
			return err
		}
	}

	return nil
}

// reconcileDaemonset prepares the desired states for Agent Daemonset and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileDaemonset(ctx context.Context, log logr.Logger, instance *v1alpha1.Agent) error {
	if instance.Spec.Sidecar.Enabled {
		return nil
	}

	dms, err := daemonsetForAgent(instance.DeepCopy(), log, r.Scheme)
	if err != nil {
		return err
	}

	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, dms, daemonsetMutate(dms, dms.Spec))
	if err != nil {
		if errors.IsConflict(err) {
			return r.reconcileDaemonset(ctx, log, instance)
		}

		msg := fmt.Sprintf("failed to create Daemonset '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			dms.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		r.Recorder.Event(instance, corev1.EventTypeNormal, "DaemonsetCreationFailed", msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "DaemonsetCreationSuccessful", "Created Daemonset '%s'", dms.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "DaemonsetUpdationSuccessful", "Updated Daemonset '%s'", dms.GetName())
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}

// reconcileMutatingWebhookConfiguration prepares the desired states for MutatingWebhookConfiguration and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileMutatingWebhookConfiguration(ctx context.Context, instance *v1alpha1.Agent) error {
	if !instance.Spec.Sidecar.Enabled {
		return nil
	}

	mwc, err := mutatingWebhookConfiguration(instance.DeepCopy())
	if err != nil {
		return err
	}

	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, mwc, mutatingWebhookConfigurationMutate(mwc, mwc.Webhooks))
	if err != nil {
		if errors.IsConflict(err) {
			return r.reconcileMutatingWebhookConfiguration(ctx, instance)
		}

		msg := fmt.Sprintf("failed to create MutatingWebhookConfiguration '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			mwc.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		r.Recorder.Event(instance, corev1.EventTypeNormal, "MutatingWebhookConfigurationCreationFailed", msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal,
			"MutatingWebhookConfigurationCreationSuccessful", "Created MutatingWebhookConfiguration '%s'", mwc.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal,
			"MutatingWebhookConfigurationUpdationSuccessful", "Updated MutatingWebhookConfiguration '%s'", mwc.GetName())
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}

// reconcileSecret prepares the desired states for Agent ApiKey secret and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileSecret(ctx context.Context, instance *v1alpha1.Agent) error {
	if instance.Spec.FluxNinjaPlugin.APIKeySecret.Create && instance.Spec.FluxNinjaPlugin.Enabled {
		if !instance.Spec.Sidecar.Enabled {
			secret, err := secretForAgentAPIKey(instance.DeepCopy(), r.Scheme)
			if err != nil {
				return err
			}
			if _, err = createSecretForAgent(r.Client, r.Recorder, secret, ctx, instance); err != nil {
				return err
			}
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Create = false
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Value = ""
			instance.Spec.FluxNinjaPlugin.APIKeySecret.SecretKeyRef.Name = secretName(
				instance.GetName(), "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret)
		} else {
			value := instance.Spec.FluxNinjaPlugin.APIKeySecret.Value
			if !strings.HasPrefix(value, "enc::") && !strings.HasSuffix(value, "::enc") {
				instance.Spec.FluxNinjaPlugin.APIKeySecret.Value = fmt.Sprintf(
					"enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(instance.Spec.FluxNinjaPlugin.APIKeySecret.Value)))
			}
		}
	}

	return nil
}

// reconcileNamespacedResources prepares the desired states for Agent ConfigMap and Secret in existing namespace and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *AgentReconciler) reconcileNamespacedResources(ctx context.Context, log logr.Logger, instance *v1alpha1.Agent) error {
	if !instance.Spec.Sidecar.Enabled {
		return nil
	}

	nsList := &corev1.NamespaceList{}
	err := r.List(ctx, nsList)
	if err != nil {
		return fmt.Errorf("failed to list Namespaces. Error: %+v", err)
	}

	for index := range nsList.Items {
		ns := nsList.Items[index]
		if ns.GetDeletionTimestamp() != nil {
			continue
		}

		update := false
		if slices.Contains(instance.Spec.Sidecar.EnableNamespaceByDefault, ns.GetName()) &&
			(ns.Labels == nil || ns.Labels[sidecarLabelKey] != enabled) {
			if ns.Labels == nil {
				ns.Labels = map[string]string{}
			}
			ns.Labels[sidecarLabelKey] = enabled
			update = true
		}

		if ns.Labels == nil || ns.Labels[sidecarLabelKey] != enabled {
			continue
		}

		configMap := createAgentConfigMapInNamespace(instance.DeepCopy(), ns.GetName())
		if _, err = createConfigMapForAgent(r.Client, r.Recorder, configMap, ctx, instance); err != nil {
			return err
		}

		if instance.Spec.FluxNinjaPlugin.APIKeySecret.Create && instance.Spec.FluxNinjaPlugin.Enabled {
			secret, err := createAgentSecretInNamespace(instance.DeepCopy(), ns.GetName())
			if err != nil {
				return err
			}
			if _, err = createSecretForAgent(r.Client, r.Recorder, secret, ctx, instance); err != nil {
				return err
			}
		}

		if update {
			if err := updateResource(r.Client, ctx, &ns); err != nil {
				return err
			}
		}
	}

	return nil
}

// eventFilters sets up a Predicate filter for the received events.
func eventFiltersForAgent() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(create event.CreateEvent) bool {
			_, ok := create.Object.(*v1alpha1.Agent)
			return ok
		},
		UpdateFunc: func(update event.UpdateEvent) bool {
			new, ok1 := update.ObjectNew.(*v1alpha1.Agent)
			old, ok2 := update.ObjectOld.(*v1alpha1.Agent)
			if !ok1 || !ok2 {
				return true
			}

			if new.GetDeletionTimestamp() != nil {
				return true
			}

			diffObjects := !reflect.DeepEqual(old.Spec, new.Spec)
			// Skipping update events for Secret updates
			if diffObjects && old.Spec.FluxNinjaPlugin.APIKeySecret.Value != "" &&
				(new.Spec.FluxNinjaPlugin.APIKeySecret.Value == "" || strings.HasPrefix(new.Spec.FluxNinjaPlugin.APIKeySecret.Value, "enc::")) {
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
func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Agent{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.ServiceAccount{}).
		WithEventFilter(eventFiltersForAgent()).
		Complete(r)
}
