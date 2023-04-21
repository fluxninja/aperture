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
	"golang.org/x/exp/slices"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ClusterRoleMutate returns a mutate function that can be used to update the ClusterRole's spec.
func ClusterRoleMutate(cr *rbacv1.ClusterRole, rules []rbacv1.PolicyRule) controllerutil.MutateFn {
	return func() error {
		cr.Rules = rules
		return nil
	}
}

// ClusterRoleBindingMutate returns a mutate function that can be used to update the Service's service spec.
func ClusterRoleBindingMutate(crb *rbacv1.ClusterRoleBinding, roleRef rbacv1.RoleRef, subjects []rbacv1.Subject) controllerutil.MutateFn {
	return func() error {
		crb.RoleRef = roleRef
		return nil
	}
}

// ConfigMapMutate returns a mutate function that can be used to update the ConfigMap's configuration data.
func ConfigMapMutate(cm *corev1.ConfigMap, files map[string]string) controllerutil.MutateFn {
	return func() error {
		cm.Data = files
		return nil
	}
}

// MutatingWebhookConfigurationMutate returns a mutate function that can be used to update the MutatingWebhookConfiguration's spec.
func MutatingWebhookConfigurationMutate(mwc *admissionregistrationv1.MutatingWebhookConfiguration, webhooks []admissionregistrationv1.MutatingWebhook) controllerutil.MutateFn {
	return func() error {
		mwc.Webhooks[0].AdmissionReviewVersions = webhooks[0].AdmissionReviewVersions
		mwc.Webhooks[0].ClientConfig = webhooks[0].ClientConfig
		mwc.Webhooks[0].NamespaceSelector = webhooks[0].NamespaceSelector
		mwc.Webhooks[0].Rules = webhooks[0].Rules
		mwc.Webhooks[0].FailurePolicy = webhooks[0].FailurePolicy
		mwc.Webhooks[0].SideEffects = webhooks[0].SideEffects
		mwc.Webhooks[0].TimeoutSeconds = webhooks[0].TimeoutSeconds
		return nil
	}
}

// SecretMutate returns a mutate function that can be used to update the Secret's data.
func SecretMutate(secret *corev1.Secret, data map[string][]byte, ownerReferences []v1.OwnerReference) controllerutil.MutateFn {
	return func() error {
		secret.Data = data
		if !slices.Equal(secret.OwnerReferences, ownerReferences) {
			secret.OwnerReferences = ownerReferences
		}
		return nil
	}
}

// ServiceAccountMutate returns a mutate function that can be used to update the ClusterRole's spec.
func ServiceAccountMutate(sa *corev1.ServiceAccount, automountServiceAccountToken *bool) controllerutil.MutateFn {
	return func() error {
		sa.AutomountServiceAccountToken = automountServiceAccountToken
		return nil
	}
}

// ServiceMutate returns a mutate function that can be used to update the Service's spec.
func ServiceMutate(svc *corev1.Service, spec corev1.ServiceSpec) controllerutil.MutateFn {
	return func() error {
		svc.Spec.Ports = spec.Ports
		svc.Spec.Selector = spec.Selector
		return nil
	}
}

// ValidatingWebhookConfigurationMutate returns a mutate function that can be used to update the ValidatingWebhookConfiguration's spec.
func ValidatingWebhookConfigurationMutate(vwc *admissionregistrationv1.ValidatingWebhookConfiguration, webhooks []admissionregistrationv1.ValidatingWebhook) controllerutil.MutateFn {
	return func() error {
		vwc.Webhooks[0].AdmissionReviewVersions = webhooks[0].AdmissionReviewVersions
		vwc.Webhooks[0].ClientConfig = webhooks[0].ClientConfig
		vwc.Webhooks[0].ObjectSelector = webhooks[0].ObjectSelector
		vwc.Webhooks[0].Rules = webhooks[0].Rules
		vwc.Webhooks[0].FailurePolicy = webhooks[0].FailurePolicy
		vwc.Webhooks[0].SideEffects = webhooks[0].SideEffects
		vwc.Webhooks[0].TimeoutSeconds = webhooks[0].TimeoutSeconds
		return nil
	}
}
