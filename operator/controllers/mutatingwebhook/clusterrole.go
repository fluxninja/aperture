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

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	agentv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
)

// updateClusterRoleBinding appends the Serviaccount in the ClusterRoleBinding if not exists.
func updateClusterRoleBinding(client client.Client, subject rbacv1.Subject, ctx context.Context, instance *agentv1alpha1.Agent) error {
	crb := &rbacv1.ClusterRoleBinding{}
	err := client.Get(ctx, types.NamespacedName{Name: controllers.AgentResourceName(instance), Namespace: instance.GetNamespace()}, crb)
	if err != nil {
		return fmt.Errorf("failed to Get the ClusterRoleBinding. Error: %+v", err)
	}

	for _, sub := range crb.Subjects {
		if sub.Name == subject.Name && sub.Namespace == subject.Namespace {
			return nil
		}
	}

	crb.Subjects = append(crb.Subjects, subject)
	err = client.Update(ctx, crb)
	if err != nil {
		if errors.IsConflict(err) {
			return updateClusterRoleBinding(client, subject, ctx, instance)
		}

		return fmt.Errorf("failed to Update the ClusterRoleBinding. Error: %+v", err)
	}

	return nil
}
