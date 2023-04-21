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

package namespace

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"time"

	"github.com/fluxninja/aperture/operator/controllers"
	"github.com/fluxninja/aperture/operator/controllers/agent"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// NamespaceReconciler reconciles a Namespace object.
type NamespaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Namespace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	namespace := req.Namespace
	if namespace == "" {
		namespace = req.Name
	}

	ns := &corev1.Namespace{}
	err := r.Get(ctx, types.NamespacedName{Name: namespace}, ns)
	if (err != nil && errors.IsNotFound(err)) || ns.DeletionTimestamp != nil {
		// Namespace is deleted so no need to reconcile.
		return ctrl.Result{}, nil
	}

	// Checking if the injection is enabled and time check for newly create namespaces
	if (ns.Labels == nil || ns.Labels[controllers.SidecarLabelKey] != controllers.Enabled) &&
		time.Now().UTC().Sub(ns.GetCreationTimestamp().Time.UTC()) > 5*time.Second {
		return ctrl.Result{}, nil
	}

	instance := &agentv1alpha1.Agent{}
	instances := &agentv1alpha1.AgentList{}
	err = r.List(ctx, instances)
	if err != nil {
		logger.Error(err, "failed to list Agent")
		return ctrl.Result{}, err
	}

	// Checking if Aperture CR is created and sidecar mode is enabled
	createResource := false
	if instances.Items != nil && len(instances.Items) != 0 {
		for index := range instances.Items {
			ins := instances.Items[index]
			if ins.GetDeletionTimestamp() == nil && ins.Spec.Sidecar.Enabled && (ins.Status.Resources == "creating" || ins.Status.Resources == "created") {
				instance = &ins
				createResource = true
				break
			}
		}
	}

	if createResource {
		if ns.Labels == nil || ns.Labels[controllers.SidecarLabelKey] != controllers.Enabled {
			if slices.Contains(instance.Spec.Sidecar.EnableNamespaceByDefault, ns.GetName()) {
				if ns.Labels == nil {
					ns.Labels = map[string]string{}
				}
				ns.Labels[controllers.SidecarLabelKey] = controllers.Enabled
				if err := controllers.UpdateResource(r.Client, ctx, ns); err != nil {
					return ctrl.Result{}, err
				}
			} else {
				// No need to reconcile as neither the label exists nor the default injection is enabled
				return ctrl.Result{}, nil
			}
		}
		if err := r.manageResources(ctx, logger, instance, namespace); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// manageResources creates/updates required resources.
func (r *NamespaceReconciler) manageResources(ctx context.Context, log logr.Logger, instance *agentv1alpha1.Agent, namespace string) error {
	if err := r.reconcileControllerCertConfigMap(ctx, log, instance, namespace); err != nil {
		return err
	}

	if err := r.reconcileConfigMap(ctx, log, instance.DeepCopy(), namespace); err != nil {
		return err
	}

	if err := r.reconcileSecret(ctx, log, instance.DeepCopy(), namespace); err != nil {
		return err
	}

	return nil
}

// reconcileControllerCertConfigMap prepares the desired states for Agent configmaps containing Controller cert and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *NamespaceReconciler) reconcileControllerCertConfigMap(ctx context.Context, log logr.Logger, instance *agentv1alpha1.Agent, namespace string) error {
	if instance.Spec.ControllerClientCertConfig.ConfigMapName != "" &&
		instance.Spec.ControllerClientCertConfig.ConfigMapName != controllers.AgentControllerClientCertCMName {
		return nil
	}

	if instance.Spec.ControllerClientCertConfig.ClientCertKeyName == "" {
		instance.Spec.ControllerClientCertConfig.ClientCertKeyName = controllers.ControllerClientCertKey
	}

	configMap := agent.CreateAgentControllerClientCertConfigMapInNamespace(ctx, r.Client, instance, namespace)

	if configMap == nil {
		return nil
	}

	instance.Spec.ControllerClientCertConfig.ConfigMapName = controllers.AgentControllerClientCertCMName
	res, err := agent.CreateConfigMapForAgent(r.Client, nil, configMap, ctx, instance)
	if err != nil {
		msg := fmt.Sprintf("failed to create/update ConfigMap in namespace '%s' for Instance '%s' in Namespace '%s'. Error='%s'",
			namespace, instance.GetName(), instance.GetNamespace(), err.Error())
		log.Error(err, msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		log.Info(fmt.Sprintf("Created ConfigMap '%s' in Namespace '%s'", configMap.GetName(), configMap.GetNamespace()))
	case controllerutil.OperationResultUpdated:
		log.Info(fmt.Sprintf("Updated ConfigMap '%s' in Namespace '%s'", configMap.GetName(), configMap.GetNamespace()))
	case controllerutil.OperationResultNone:
	default:
	}
	return nil
}

// reconcileConfigMap prepares the desired states for Agent configmap and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *NamespaceReconciler) reconcileConfigMap(ctx context.Context, log logr.Logger, instance *agentv1alpha1.Agent, namespace string) error {
	if len(instance.Spec.ConfigSpec.AgentFunctions.Endpoints) > 0 &&
		instance.Spec.ControllerClientCertConfig.ConfigMapName != "" {
		instance.Spec.ConfigSpec.AgentFunctions.ClientConfig.GRPCClient.ClientTLSConfig.CAFile = path.Join(controllers.AgentControllerClientCertPath, instance.Spec.ControllerClientCertConfig.ClientCertKeyName)
	}

	configMap := agent.CreateAgentConfigMapInNamespace(ctx, r.Client, instance, namespace)

	res, err := agent.CreateConfigMapForAgent(r.Client, nil, configMap, ctx, instance)
	if err != nil {
		msg := fmt.Sprintf("failed to create/update ConfigMap in namespace '%s' for Instance '%s' in Namespace '%s'. Error='%s'",
			namespace, instance.GetName(), instance.GetNamespace(), err.Error())
		log.Error(err, msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		log.Info(fmt.Sprintf("Created ConfigMap '%s' in Namespace '%s'", configMap.GetName(), configMap.GetNamespace()))
	case controllerutil.OperationResultUpdated:
		log.Info(fmt.Sprintf("Updated ConfigMap '%s' in Namespace '%s'", configMap.GetName(), configMap.GetNamespace()))
	case controllerutil.OperationResultNone:
	default:
	}
	return nil
}

// reconcileSecret prepares the desired states for Agent ApiKey secret and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *NamespaceReconciler) reconcileSecret(ctx context.Context, log logr.Logger, instance *agentv1alpha1.Agent, namespace string) error {
	if !instance.Spec.Secrets.FluxNinjaExtension.Create {
		return nil
	}

	secret, err := agent.CreateAgentSecretInNamespace(instance, namespace)
	if err != nil {
		msg := fmt.Sprintf("failed to create Secret instance in namespace '%s' for Instance '%s' in Namespace '%s'. Error='%s'",
			namespace, instance.GetName(), instance.GetNamespace(), err.Error())
		log.Error(err, msg)
		return fmt.Errorf(msg)
	}

	res, err := agent.CreateSecretForAgent(r.Client, nil, secret, ctx, instance)
	if err != nil {
		msg := fmt.Sprintf("failed to create/update Secret in namespace '%s' for Instance '%s' in Namespace '%s'. Error='%s'",
			namespace, instance.GetName(), instance.GetNamespace(), err.Error())
		log.Error(err, msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		log.Info(fmt.Sprintf("Created Secret '%s' in Namespace '%s'", secret.GetName(), secret.GetNamespace()))
	case controllerutil.OperationResultUpdated:
		log.Info(fmt.Sprintf("Updated Secret '%s' in Namespace '%s'", secret.GetName(), secret.GetNamespace()))
	case controllerutil.OperationResultNone:
	default:
	}
	return nil
}

// eventFilters sets up a Predicate filter for the received events.
func namespaceEventFilters() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(create event.CreateEvent) bool {
			_, ok := create.Object.(*corev1.Namespace)

			return ok
		},
		UpdateFunc: func(update event.UpdateEvent) bool {
			nsOld, ok1 := update.ObjectOld.(*corev1.Namespace)
			nsNew, ok2 := update.ObjectNew.(*corev1.Namespace)
			if ok1 && ok2 {
				if (nsOld.Labels == nil || nsOld.Labels[controllers.SidecarLabelKey] != controllers.Enabled) &&
					(nsNew.Labels != nil && nsNew.Labels[controllers.SidecarLabelKey] == controllers.Enabled) {
					return true
				}

				return false
			}

			cm, ok1 := update.ObjectOld.(*corev1.ConfigMap)
			secret, ok2 := update.ObjectOld.(*corev1.Secret)

			if ok1 && cm.Labels != nil && cm.Labels["app.kubernetes.io/component"] == controllers.AgentServiceName {
				newCm, _ := update.ObjectNew.(*corev1.ConfigMap)
				return !reflect.DeepEqual(cm.Data, newCm.Data) || newCm.GetDeletionTimestamp() != nil
			} else if ok2 && secret.Labels != nil && secret.Labels["app.kubernetes.io/component"] == controllers.AgentServiceName {
				newSecret, _ := update.ObjectNew.(*corev1.Secret)
				return !reflect.DeepEqual(secret.Data, newSecret.Data) || newSecret.GetDeletionTimestamp() != nil
			} else {
				return false
			}
		},
		DeleteFunc: func(delete event.DeleteEvent) bool {
			cm, ok1 := delete.Object.(*corev1.ConfigMap)
			secret, ok2 := delete.Object.(*corev1.Secret)

			if ok1 && cm.Labels != nil && cm.Labels["app.kubernetes.io/component"] == controllers.AgentServiceName {
				return true
			} else if ok2 && secret.Labels != nil && secret.Labels["app.kubernetes.io/component"] == controllers.AgentServiceName {
				return true
			} else {
				return false
			}
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Watches(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForObject{}).
		WithEventFilter(namespaceEventFilters()).
		Complete(r)
}
