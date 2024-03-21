package vmstop

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	impltypes "github.com/chaos-mesh/chaos-mesh/controllers/chaosimpl/types"
)

// Ensures Impl implements the ChaosImpl interface.
var _ impltypes.ChaosImpl = (*Impl)(nil)

// Impl represents the implementation of the OpenStack chaos action to stop a VM.
type Impl struct {
	client.Client

	Log logr.Logger
}

// Apply attempts to stop a specified OpenStack VM using credentials stored in a Kubernetes secret.
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

	// Stop the VM
	err = servers.Stop(computeClient, records[index].Id).ExtractErr()
	if err != nil {
		impl.Log.Error(err, "Failed to stop the OpenStack VM")
		return v1alpha1.NotInjected, err
	}

	return v1alpha1.Injected, nil
}

func (impl *Impl) Recover(ctx context.Context, index int, records []*v1alpha1.Record, obj v1alpha1.InnerObject) (v1alpha1.Phase, error) {
	openstackChaos := obj.(*v1alpha1.OpenStackChaos)

	// Assuming records[index].Id contains the VM ID
	vmID := records[index].Id

	// Retrieve the secret containing the credentials
	secret := &v1.Secret{}
	err := impl.Client.Get(ctx, types.NamespacedName{
		Name:      *openstackChaos.Spec.SecretName,
		Namespace: openstackChaos.Namespace,
	}, secret)
	if err != nil {
		impl.Log.Error(err, "Failed to get cloud secret for OpenStack")
		return v1alpha1.Injected, err
	}

	// Authenticate with OpenStack using credentials from the secret
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: string(secret.Data["identity_endpoint"]),
		Username:         string(secret.Data["username"]),
		Password:         string(secret.Data["password"]),
		TenantID:         string(secret.Data["tenant_id"]),
		DomainName:       string(secret.Data["domain_name"]),
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		impl.Log.Error(err, "Failed to authenticate with OpenStack")
		return v1alpha1.Injected, err
	}

	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: string(secret.Data["region"]),
	})
	if err != nil {
		impl.Log.Error(err, "Failed to create compute client")
		return v1alpha1.Injected, err
	}

	// Start the VM
	err = servers.Start(computeClient, vmID).ExtractErr()
	if err != nil {
		impl.Log.Error(err, "Failed to start the OpenStack VM")
		return v1alpha1.Injected, err
	}

	return v1alpha1.NotInjected, nil
}

// NewImpl creates a new instance of the OpenStack chaos implementation for stopping VMs.
// c: The Kubernetes client.
// log: A logger instance.
// Returns a new Impl instance configured with the provided client and logger.
func NewImpl(c client.Client, log logr.Logger) *Impl {
	return &Impl{
		Client: c,
		Log:    log.WithName("vmstop"),
	}
}
