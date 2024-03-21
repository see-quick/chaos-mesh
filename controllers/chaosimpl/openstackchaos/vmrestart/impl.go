// Copyright Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package vmrestart

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "k8s.io/api/core/v1"         // Correct import for v1
	"k8s.io/apimachinery/pkg/types" // Correct import for types

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	impltypes "github.com/chaos-mesh/chaos-mesh/controllers/chaosimpl/types"
)

// Ensures Impl implements the ChaosImpl interface.
var _ impltypes.ChaosImpl = (*Impl)(nil)

// Impl represents the implementation of the OpenStack chaos action to restart a VM.
type Impl struct {
	client.Client

	Log logr.Logger
}

// Apply attempts to restart a specified OpenStack VM using credentials stored in a Kubernetes secret.
func (impl *Impl) Apply(ctx context.Context, index int, records []*v1alpha1.Record, obj v1alpha1.InnerObject) (v1alpha1.Phase, error) {
	openstackChaos := obj.(*v1alpha1.OpenStackChaos)

	if openstackChaos.Spec.SecretName == nil {
		impl.Log.Error(nil, "Secret name not provided")
		return v1alpha1.NotInjected, fmt.Errorf("secret name not provided")
	}

	// Retrieve the secret containing the credentials
	secret := &v1.Secret{}
	err := impl.Client.Get(ctx, types.NamespacedName{
		Name:      *openstackChaos.Spec.SecretName,
		Namespace: openstackChaos.Namespace,
	}, secret)
	if err != nil {
		impl.Log.Error(err, "Failed to get cloud secret")
		return v1alpha1.NotInjected, err
	}

	// Use the credentials from the secret to authenticate with OpenStack
	opts, err := openstack.AuthOptionsFromEnv()
	opts.IdentityEndpoint = string(secret.Data["identity_endpoint"])
	opts.Username = string(secret.Data["username"])
	opts.Password = string(secret.Data["password"])
	opts.TenantID = string(secret.Data["tenant_id"])
	opts.DomainName = string(secret.Data["domain_name"])

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		impl.Log.Error(err, "Failed to authenticate with OpenStack")
		return v1alpha1.NotInjected, err
	}

	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: string(secret.Data["region"]),
	})
	if err != nil {
		impl.Log.Error(err, "Failed to create compute client")
		return v1alpha1.NotInjected, err
	}

	// Restart the VM using the correct method
	rebootOpts := servers.RebootOpts{
		Type: servers.SoftReboot, // Correctly use RebootOpts with SoftReboot
	}
	err = servers.Reboot(computeClient, records[index].Id, rebootOpts).ExtractErr()
	if err != nil {
		impl.Log.Error(err, "Failed to restart the OpenStack VM")
		return v1alpha1.NotInjected, err
	}

	return v1alpha1.Injected, nil
}

// Recover performs any cleanup or recovery actions after the chaos experiment.
// In the context of a VM restart, this might be a no-op as the state is managed by OpenStack.
func (impl *Impl) Recover(ctx context.Context, index int, records []*v1alpha1.Record, obj v1alpha1.InnerObject) (v1alpha1.Phase, error) {
	return v1alpha1.NotInjected, nil
}

// NewImpl creates a new instance of the OpenStack chaos implementation.
// c: The Kubernetes client.
// log: A logger instance.
// Returns a new Impl instance configured with the provided client and logger.
func NewImpl(c client.Client, log logr.Logger) *Impl {
	return &Impl{
		Client: c,
		Log:    log.WithName("vmrestart"),
	}
}
