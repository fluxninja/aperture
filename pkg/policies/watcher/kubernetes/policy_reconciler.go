package kubernetes

import (
	"bytes"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

// PolicyReconciler reconciles a Policy Custom Resource.
type PolicyReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	Recorder         record.EventRecorder
	Trackers         notifiers.Trackers
	resourcesDeleted map[types.NamespacedName]bool
}

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
	if r.resourcesDeleted == nil {
		r.resourcesDeleted = make(map[types.NamespacedName]bool)
	}

	instance := &policyv1alpha1.Policy{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil && errors.IsNotFound(err) {
		if !r.resourcesDeleted[req.NamespacedName] {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Debug().Msg(fmt.Sprintf("Handling deletion of resources for Instance '%s' in Namespace '%s'", req.Name, req.Namespace))
			instance.Name = req.Name
			instance.Namespace = req.Namespace
			r.deleteResources(ctx, instance.DeepCopy())
			r.resourcesDeleted[req.NamespacedName] = true
			log.Debug().Msg("Policy resource not found. Ignoring since object must be deleted")
		}
		return ctrl.Result{}, nil
	} else if err != nil {
		// Error reading the object - requeue the request.
		log.Error().Err(err).Msg("failed to get Policy")
		return ctrl.Result{}, err
	}

	// Handing delete operation
	if instance.GetDeletionTimestamp() != nil {
		log.Debug().Msg(fmt.Sprintf("Handling deletion of resources for Instance '%s' in Namespace '%s'", instance.GetName(), instance.GetNamespace()))
		r.deleteResources(ctx, instance.DeepCopy())
		r.resourcesDeleted[req.NamespacedName] = true
		return ctrl.Result{}, nil
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

func (r *PolicyReconciler) deleteResources(ctx context.Context, instance *policyv1alpha1.Policy) {
	r.Trackers.RemoveEvent(notifiers.Key(instance.GetName()))
}

// updateResource updates the Aperture resource in Kubernetes.
func (r *PolicyReconciler) updateStatus(ctx context.Context, instance *policyv1alpha1.Policy) error {
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

// reconcilePolicy sends a write event to notifier to get it uploaded on the Etcd.
func (r *PolicyReconciler) reconcilePolicy(ctx context.Context, instance *policyv1alpha1.Policy) error {
	r.Trackers.WriteEvent(notifiers.Key(instance.GetName()), instance.Spec.Raw)

	r.Recorder.Eventf(instance, corev1.EventTypeWarning, "UploadSuccessful", "Uploaded policy to Etcd.")
	return nil
}

// eventFilters sets up a Predicate filter for the received events.
func eventFiltersForPolicy() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(create event.CreateEvent) bool {
			_, ok := create.Object.(*policyv1alpha1.Policy)
			return ok
		},
		UpdateFunc: func(update event.UpdateEvent) bool {
			new, ok1 := update.ObjectNew.(*policyv1alpha1.Policy)
			old, ok2 := update.ObjectOld.(*policyv1alpha1.Policy)
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
		For(&policyv1alpha1.Policy{}).
		WithEventFilter(eventFiltersForPolicy()).
		Complete(r)
}
