package v1alpha1

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="action",type=string,JSONPath=`.spec.action`
// +kubebuilder:printcolumn:name="duration",type=string,JSONPath=`.spec.duration`
// +chaos-mesh:experiment
// +chaos-mesh:oneshot=in.Spec.Action==VmRestart

// OpenStackChaos is the Schema for the openstackchaos API
type OpenStackChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackChaosSpec   `json:"spec"`
	Status OpenStackChaosStatus `json:"status,omitempty"`
}

var _ InnerObjectWithSelector = (*OpenStackChaos)(nil)
var _ InnerObject = (*OpenStackChaos)(nil)

const (
	// VmStop represents the chaos action of stopping an OpenStack VM.
	VmStop OpenStackChaosAction = "vm-stop"
	// VmRestart represents the chaos action of restarting an OpenStack VM.
	VmRestart OpenStackChaosAction = "vm-restart"
)

// OpenStackChaosSpec is the content of the specification for an OpenStackChaos
type OpenStackChaosSpec struct {
	// Action defines the specific OpenStack chaos action.
	// Supported action: vm-stop / vm-restart / volume-detach
	// Default action: vm-stop
	// +kubebuilder:validation:Enum=vm-stop;vm-restart;volume-detach
	Action OpenStackChaosAction `json:"action"`

	// Duration represents the duration of the chaos action.
	// +optional
	Duration *string `json:"duration,omitempty" webhook:"Duration"`

	// SecretName defines the name of kubernetes secret. It is used for OpenStack credentials.
	// +optional
	SecretName *string `json:"secretName,omitempty"`

	OpenStackSelector `json:",inline"`
}

// OpenStackChaosStatus represents the status of an OpenStackChaos
type OpenStackChaosStatus struct {
	ChaosStatus `json:",inline"`
}

type OpenStackSelector struct {
	// TenantID defines the ID of the OpenStack tenant.
	TenantID string `json:"tenantID"`

	// AuthURL defines the URL for authentication in OpenStack.
	AuthURL string `json:"authURL"`

	// VMID defines the ID of the Virtual Machine in OpenStack to target.
	VMID string `json:"vmID"`

	// VolumeID indicates the ID of the volume.
	// Needed in volume-detach.
	// +optional
	VolumeID *string `json:"volumeID,omitempty" webhook:"VolumeID,nilable"`

	// SecretName defines the name of kubernetes secret. It is used for OpenStack credentials.
	// +optional
	SecretName *string `json:"secretName,omitempty"`

	// RemoteCluster represents the remote cluster where the chaos will be deployed
	// +optional
	RemoteCluster string `json:"remoteCluster,omitempty"`
}

func (obj *OpenStackChaos) GetSelectorSpecs() map[string]interface{} {
	return map[string]interface{}{
		".": &obj.Spec.OpenStackSelector,
	}
}

func (selector *OpenStackSelector) Id() string {
	json, _ := json.Marshal(selector)
	return string(json)
}
