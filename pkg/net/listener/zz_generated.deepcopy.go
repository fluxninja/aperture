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

package listener

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ListenerConfig) DeepCopyInto(out *ListenerConfig) {
	*out = *in
	in.KeepAlive.DeepCopyInto(&out.KeepAlive)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListenerConfig.
func (in *ListenerConfig) DeepCopy() *ListenerConfig {
	if in == nil {
		return nil
	}
	out := new(ListenerConfig)
	in.DeepCopyInto(out)
	return out
}
