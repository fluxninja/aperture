// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package flowcontrolv1

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto supports using AuthzResponse within kubernetes types, where deepcopy-gen is used.
func (in *AuthzResponse) DeepCopyInto(out *AuthzResponse) {
	p := proto.Clone(in).(*AuthzResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthzResponse. Required by controller-gen.
func (in *AuthzResponse) DeepCopy() *AuthzResponse {
	if in == nil {
		return nil
	}
	out := new(AuthzResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AuthzResponse. Required by controller-gen.
func (in *AuthzResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
