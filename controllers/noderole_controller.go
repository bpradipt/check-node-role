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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	nodeattrv1alpha1 "github.com/check-node-role/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// NodeRoleReconciler reconciles a NodeRole object
type NodeRoleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nodeattr.power.io,resources=noderoles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nodeattr.power.io,resources=noderoles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch

func (r *NodeRoleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("noderole", req.NamespacedName)

	// Fetch the NodeRole instance

	noderole := &nodeattrv1alpha1.NodeRole{}

	err := r.Get(ctx, req.NamespacedName, noderole)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("NodeRole resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get NodeRole")
		return ctrl.Result{}, err
	}
	// NodeRole instance found

	controllerList := &corev1.NodeList{}
	controllerListOpts := []client.ListOption{
		// node-role.kubernetes.io/master=
		client.MatchingLabels(map[string]string{"node-role.kubernetes.io/master": ""}),
	}

	infraList := &corev1.NodeList{}
	infraListOpts := []client.ListOption{
		// node-role.kubernetes.io/infra=
		client.MatchingLabels(map[string]string{"node-role.kubernetes.io/infra": ""}),
	}

	workerList := &corev1.NodeList{}
	workerListOpts := []client.ListOption{
		// node-role.kubernetes.io/worker=
		client.MatchingLabels(map[string]string{"node-role.kubernetes.io/worker": ""}),
	}

	if noderole.Spec.Controller != nil && *noderole.Spec.Controller {
		if err = r.List(ctx, controllerList, controllerListOpts...); err != nil {
			log.Error(err, "Failed to get controller nodes")
		} else {
			log.Info("Got list of controller nodes")
			noderole.Status.ControllerNodes = getNodeNames(controllerList.Items)
		}
	}

	if noderole.Spec.Worker != nil && *noderole.Spec.Worker {
		if err = r.List(ctx, workerList, workerListOpts...); err != nil {
			log.Error(err, "Failed to get worker nodes")
		} else {
			log.Info("Got list of worker nodes")
			noderole.Status.WorkerNodes = getNodeNames(workerList.Items)
		}
	}

	if noderole.Spec.Infra != nil && *noderole.Spec.Infra {
		if err = r.List(ctx, infraList, infraListOpts...); err != nil {
			log.Error(err, "Failed to get infra nodes")
		} else {
			log.Info("Got list of infra nodes")
			noderole.Status.InfraNodes = getNodeNames(infraList.Items)
		}
	}

	err = r.Status().Update(ctx, noderole)
	if err != nil {
		log.Error(err, "Failed to update noderole status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *NodeRoleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nodeattrv1alpha1.NodeRole{}).
		Watches(&source.Kind{Type: &corev1.Node{}}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(nodeRoleObj handler.MapObject) []reconcile.Request {
				nodeRoleList := &nodeattrv1alpha1.NodeRoleList{}
				err := r.Client.List(context.TODO(), nodeRoleList)
				if err != nil {
					return []reconcile.Request{}
				}
				reconcileRequests := make([]reconcile.Request, len(nodeRoleList.Items))
				for _, noderole := range nodeRoleList.Items {
					reconcileRequests = append(reconcileRequests, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Name:      noderole.Name,
							Namespace: noderole.Namespace,
						},
					})
				}
				return reconcileRequests
			}),
		}).
		Complete(r)
}

// Get NodeNames list
func getNodeNames(nodes []corev1.Node) []string {
	var nodeNames []string
	for _, node := range nodes {
		nodeNames = append(nodeNames, node.Name)
	}
	return nodeNames
}
