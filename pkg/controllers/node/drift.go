/*
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

package node

import (
	"context"
	"github.com/aws/karpenter-core/pkg/cloudprovider"

	"k8s.io/utils/clock"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/aws/karpenter-core/pkg/apis/provisioning/v1alpha5"
	"github.com/aws/karpenter-core/pkg/controllers/state"
)

// Emptiness is a subreconciler that deletes nodes that are empty after a ttl
type Drift struct {
	kubeClient    client.Client
	clock         clock.Clock
	cluster       *state.Cluster
	cloudProvider cloudprovider.CloudProvider
}

// Reconcile reconciles the node
func (r *Drift) Reconcile(ctx context.Context, provisioner *v1alpha5.Provisioner, n *v1.Node) (reconcile.Result, error) {
	if n.Annotations[v1alpha5.AmiAnnotationKey] == "" {
		//Skip because we dont have the ami-id, but who adds the ami label, do we need to put an annotation while provisioning ?
		return reconcile.Result{}, nil
	}
	isDrifted := r.cloudProvider.IsNodeDrifted(ctx, n, provisioner)
	if isDrifted {
		n.Annotations[v1alpha5.DriftedAnnotationKey] = "true"
	}
	return reconcile.Result{}, nil
}
