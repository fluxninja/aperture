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
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/yaml"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	policy "github.com/fluxninja/aperture/pkg/policies/controlplane"
)

// PolicyReconciler reconciles a Policy object.
type PolicyReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	Recorder         record.EventRecorder
	resourcesDeleted map[types.NamespacedName]bool
}

//+kubebuilder:rbac:groups=fluxninja.com,resources=policies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fluxninja.com,resources=policies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fluxninja.com,resources=policies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Policy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *PolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	if r.resourcesDeleted == nil {
		r.resourcesDeleted = make(map[types.NamespacedName]bool)
	}

	instance := &v1alpha1.Policy{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil && errors.IsNotFound(err) {
		if !r.resourcesDeleted[req.NamespacedName] {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info(fmt.Sprintf("Handling deletion of resources for Instance '%s' in Namespace '%s'", req.Name, req.Namespace))
			instance.Name = req.Name
			instance.Namespace = req.Namespace
			r.deleteResources(ctx, instance.DeepCopy())
			r.resourcesDeleted[req.NamespacedName] = true
			logger.Info("Policy resource not found. Ignoring since object must be deleted")
		}
		return ctrl.Result{}, nil
	} else if err != nil {
		// Error reading the object - requeue the request.
		logger.Error(err, "failed to get Policy")
		return ctrl.Result{}, err
	}

	// Handing delete operation
	if instance.GetDeletionTimestamp() != nil {
		logger.Info(fmt.Sprintf("Handling deletion of resources for Instance '%s' in Namespace '%s'", instance.GetName(), instance.GetNamespace()))
		r.deleteResources(ctx, instance.DeepCopy())
		r.resourcesDeleted[req.NamespacedName] = true
		return ctrl.Result{}, nil
	}

	if instance.Annotations == nil || instance.Annotations[defaulterAnnotationKey] != "true" {
		err = r.validatePolicy(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		} else if instance.Status.Status == failedStatus {
			return ctrl.Result{}, nil
		}

		if err := r.updatePolicy(ctx, instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	instance.Status.Status = "uploading"
	if err := r.updateStatus(ctx, instance.DeepCopy()); err != nil {
		return ctrl.Result{}, err
	}
	r.resourcesDeleted[req.NamespacedName] = false

	if err := r.reconcilePolicy(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	instance.Status.Status = "uploaded"
	if err := r.updateStatus(ctx, instance.DeepCopy()); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PolicyReconciler) deleteResources(ctx context.Context, instance *v1alpha1.Policy) {
	filename := filepath.Join(policyFilePath, getPolicyFileName(instance.GetName(), instance.GetNamespace()))
	err := os.Remove(filename)
	if err != nil {
		log.FromContext(ctx).Info(fmt.Sprintf("Failed to write Policy to file '%s'. Error: '%s'", filename, err.Error()))
	}
}

func (r *PolicyReconciler) updatePolicy(ctx context.Context, instance *v1alpha1.Policy) error {
	attempt := 5
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
				instance.Annotations = annotations
				continue
			}
			return err
		}
	}

	return nil
}

// updateResource updates the Aperture resource in Kubernetes.
func (r *PolicyReconciler) updateStatus(ctx context.Context, instance *v1alpha1.Policy) error {
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

func (r *PolicyReconciler) validatePolicy(ctx context.Context, instance *v1alpha1.Policy) error {
	_, valid, msg, err := policy.ValidateAndCompile(ctx, "", instance.Spec.Raw)
	if err != nil || !valid {
		instance.Status.Status = failedStatus
		if msg == "" {
			msg = err.Error()
		}
		r.Recorder.Eventf(instance, corev1.EventTypeWarning, "ValidationFailed", "Failed to validate Policy. Error: '%s'", msg)
		errUpdate := r.updateStatus(ctx, instance)
		if errUpdate != nil {
			return errUpdate
		}
		return nil
	}

	if instance.Status.Status == failedStatus {
		instance.Status.Status = ""
	}
	return nil
}

func (r *PolicyReconciler) reconcilePolicy(ctx context.Context, instance *v1alpha1.Policy) error {
	filename := filepath.Join(policyFilePath, getPolicyFileName(instance.GetName(), instance.GetNamespace()))
	yamlContent, err := yaml.JSONToYAML(instance.Spec.Raw)
	if err != nil {
		r.Recorder.Eventf(instance, corev1.EventTypeWarning, "UploadFailed", "Failed to write Policy to file '%s'. Error: '%s'", filename, err.Error())
		return err
	}
	err = os.WriteFile(filename, yamlContent, 0o600)
	if err != nil {
		r.Recorder.Eventf(instance, corev1.EventTypeWarning, "UploadFailed", "Failed to write Policy to file '%s'. Error: '%s'", filename, err.Error())
		return err
	}

	r.Recorder.Eventf(instance, corev1.EventTypeWarning, "UploadSuccessful", "Wrote Policy to file '%s'.", filename)
	return nil
}

// eventFilters sets up a Predicate filter for the received events.
func eventFiltersForPolicy() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(create event.CreateEvent) bool {
			_, ok := create.Object.(*v1alpha1.Policy)
			return ok
		},
		UpdateFunc: func(update event.UpdateEvent) bool {
			new, ok1 := update.ObjectNew.(*v1alpha1.Policy)
			old, ok2 := update.ObjectOld.(*v1alpha1.Policy)
			if !ok1 || !ok2 {
				return false
			}

			return !bytes.Equal(old.Spec.Raw, new.Spec.Raw) || new.GetDeletionTimestamp() != nil
		},
		DeleteFunc: func(delete event.DeleteEvent) bool {
			return !delete.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Policy{}).
		WithEventFilter(eventFiltersForPolicy()).
		Complete(r)
}
