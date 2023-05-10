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

package mutatingwebhook

import (
	"context"
	"fmt"

	"github.com/fluxninja/aperture/v2/operator/controllers"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// MutatingWebhookReconciler reconciles a Namespace object.
type MutatingWebhookReconciler struct {
	client.Client
	Scheme            *runtime.Scheme
	AgentManager      bool
	ControllerManager bool
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
func (r *MutatingWebhookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	mwc := &admissionregistrationv1.MutatingWebhookConfiguration{}
	err := r.Get(ctx, req.NamespacedName, mwc)
	if (err != nil && errors.IsNotFound(err)) || mwc.DeletionTimestamp != nil {
		// MutatingWebhookConfiguration is deleted so no need to reconcile.
		return ctrl.Result{}, nil
	}

	if r.AgentManager && !r.ControllerManager && mwc.GetName() == controllers.ControllerMutatingWebhookName ||
		!r.AgentManager && r.ControllerManager && mwc.GetName() == controllers.AgentMutatingWebhookName {
		return ctrl.Result{}, nil
	}

	if err = r.reconcileMutatingWebhookConfiguration(ctx, mwc); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileMutatingWebhookConfiguration prepares the desired states for MutatingWebhookConfiguration and
// sends an request to Kubernetes API to move the actual state to the prepared desired state.
func (r *MutatingWebhookReconciler) reconcileMutatingWebhookConfiguration(ctx context.Context, instance *admissionregistrationv1.MutatingWebhookConfiguration) error {
	log := log.FromContext(ctx)
	var mwc *admissionregistrationv1.MutatingWebhookConfiguration
	var err error

	if instance.GetName() == controllers.AgentMutatingWebhookName {
		mwc, err = agentMutatingWebhookConfiguration()
		if err != nil {
			return err
		}
	} else {
		mwc, err = controllerMutatingWebhookConfiguration()
		if err != nil {
			return err
		}
	}

	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, mwc, controllers.MutatingWebhookConfigurationMutate(mwc, mwc.Webhooks))
	if err != nil {
		if errors.IsConflict(err) {
			return r.reconcileMutatingWebhookConfiguration(ctx, instance)
		}

		msg := fmt.Sprintf("failed to create MutatingWebhookConfiguration '%s'. Response='%v', Error='%s'",
			mwc.GetName(), res, err.Error())
		log.Error(err, msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		log.Info(fmt.Sprintf("Created MutatingWebhookConfiguration '%s'", mwc.GetName()))
	case controllerutil.OperationResultUpdated:
		log.Info(fmt.Sprintf("Updated MutatingWebhookConfiguration '%s'", mwc.GetName()))
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}

// eventFilters sets up a Predicate filter for the received events.
func mwcEventFilters() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(create event.CreateEvent) bool {
			new, ok := create.Object.(*admissionregistrationv1.MutatingWebhookConfiguration)

			return ok && (new.GetName() == controllers.AgentMutatingWebhookName || new.GetName() == controllers.ControllerMutatingWebhookName)
		},
		UpdateFunc: func(update event.UpdateEvent) bool {
			new, ok1 := update.ObjectNew.(*admissionregistrationv1.MutatingWebhookConfiguration)
			old, ok2 := update.ObjectOld.(*admissionregistrationv1.MutatingWebhookConfiguration)
			if !ok1 || !ok2 {
				return false
			} else if new.GetName() != controllers.AgentMutatingWebhookName && new.GetName() != controllers.ControllerMutatingWebhookName {
				return false
			} else if old.GetName() != controllers.AgentMutatingWebhookName && old.GetName() != controllers.ControllerMutatingWebhookName {
				return false
			}

			return !equality.Semantic.DeepEqual(old.Webhooks, new.Webhooks)
		},
		DeleteFunc: func(delete event.DeleteEvent) bool {
			return !delete.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *MutatingWebhookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&admissionregistrationv1.MutatingWebhookConfiguration{}).
		WithEventFilter(mwcEventFilters()).
		Complete(r)
}
