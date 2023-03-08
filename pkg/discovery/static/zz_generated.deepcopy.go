//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package static

import (
	v1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/discovery/entities/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StaticDiscoveryConfig) DeepCopyInto(out *StaticDiscoveryConfig) {
	*out = *in
	if in.Entities != nil {
		in, out := &in.Entities, &out.Entities
		*out = make([]v1.Entity, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StaticDiscoveryConfig.
func (in *StaticDiscoveryConfig) DeepCopy() *StaticDiscoveryConfig {
	if in == nil {
		return nil
	}
	out := new(StaticDiscoveryConfig)
	in.DeepCopyInto(out)
	return out
}
