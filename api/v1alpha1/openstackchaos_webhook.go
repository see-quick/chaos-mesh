package v1alpha1

import (
	"reflect"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/chaos-mesh/chaos-mesh/api/genericwebhook"
)

// Custom type definitions for validation
type VMID string
type OpenStackChaosAction string

// Validate function for VMID type
func (in *VMID) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	openStackChaos := root.(*OpenStackChaos)
	if openStackChaos.Spec.Action == VmRestart || openStackChaos.Spec.Action == VmStop {
		if in == nil || *in == "" {
			err := errors.Wrapf(errInvalidValue, "the ID of VM is required for %s action", openStackChaos.Spec.Action)
			allErrs = append(allErrs, field.Invalid(path, in, err.Error()))
		}
	}

	return allErrs
}

// Validate function for OpenStackAction type
func (in *OpenStackChaosAction) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// in cannot be nil
	switch *in {
	case VmRestart, VmStop: // Assume these are defined OpenStackChaos actions
	default:
		err := errors.WithStack(errUnknownAction)
		log.Error(err, "Wrong OpenStackChaos Action type")

		allErrs = append(allErrs, field.Invalid(path, in, err.Error()))
	}
	return allErrs
}

func init() {
	// Register custom types with the webhook system for validation
	genericwebhook.Register("VMID", reflect.PtrTo(reflect.TypeOf(VMID(""))))
	// TODO: If you have other custom types to validate, register them here.. <----
}
