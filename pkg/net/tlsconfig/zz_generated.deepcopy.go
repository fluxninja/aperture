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

package tlsconfig

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClientTLSConfig) DeepCopyInto(out *ClientTLSConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClientTLSConfig.
func (in *ClientTLSConfig) DeepCopy() *ClientTLSConfig {
	if in == nil {
		return nil
	}
	out := new(ClientTLSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerTLSConfig) DeepCopyInto(out *ServerTLSConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerTLSConfig.
func (in *ServerTLSConfig) DeepCopy() *ServerTLSConfig {
	if in == nil {
		return nil
	}
	out := new(ServerTLSConfig)
	in.DeepCopyInto(out)
	return out
}
