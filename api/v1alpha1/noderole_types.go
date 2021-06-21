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

package v1alpha1

import (
	conditions "github.com/openshift/custom-resource-status/conditions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodeRoleSpec defines the desired state of NodeRole
type NodeRoleSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Specifies which NodeRole type of nodes to find in the cluster
	Controller string `json:"controller,omitempty"`
	Worker     string `json:"worker,omitempty"`
	Infra      string `json:"infra,omitempty"`
}

// NodeRoleStatus defines the observed state of NodeRole
type NodeRoleStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Update the variables with list of specific nodes
	ControllerNodes []string `json:"controllerNodes"`
	WorkerNodes     []string `json:"workerNodes"`
	InfraNodes      []string `json:"infraNodes"`

	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	// Conditions is a list of conditions related to operator reconciliation
	Conditions []conditions.Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NodeRole is the Schema for the noderoles API
type NodeRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeRoleSpec   `json:"spec,omitempty"`
	Status NodeRoleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NodeRoleList contains a list of NodeRole
type NodeRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeRole{}, &NodeRoleList{})
}
