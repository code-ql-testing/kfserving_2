/*
Copyright 2019 The Knative Authors

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
	"context"

	"github.com/knative/pkg/apis"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/api/equality"
)

// Validate inspects and validates ClusterServerlessService object.
func (ci *ServerlessService) Validate(ctx context.Context) *apis.FieldError {
	return ci.Spec.Validate(apis.WithinSpec(ctx)).ViaField("spec")
}

// Validate inspects and validates ServerlessServiceSpec object.
func (spec *ServerlessServiceSpec) Validate(ctx context.Context) *apis.FieldError {
	// Spec must not be empty.
	if equality.Semantic.DeepEqual(spec, &ServerlessServiceSpec{}) {
		return apis.ErrMissingField(apis.CurrentField)
	}
	var all *apis.FieldError
	// Spec mode must be from the enum and
	switch spec.Mode {
	case SKSOperationModeProxy, SKSOperationModeServe:
		break
	case "":
		all = all.Also(apis.ErrMissingField("mode"))
	default:
		all = all.Also(apis.ErrInvalidValue(spec.Mode, "mode"))
	}

	all = all.Also(validateReference(spec.ObjectRef).ViaField("objectRef"))

	return all.Also(spec.ProtocolType.Validate(ctx).ViaField("protocolType"))
}

func validateReference(ref autoscalingv1.CrossVersionObjectReference) *apis.FieldError {
	if equality.Semantic.DeepEqual(ref, autoscalingv1.CrossVersionObjectReference{}) {
		return apis.ErrMissingField(apis.CurrentField)
	}
	var errs *apis.FieldError
	if ref.Kind == "" {
		errs = errs.Also(apis.ErrMissingField("kind"))
	}
	if ref.Name == "" {
		errs = errs.Also(apis.ErrMissingField("name"))
	}
	if ref.APIVersion == "" {
		errs = errs.Also(apis.ErrMissingField("apiVersion"))
	}
	return errs
}