// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package wrappersv1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using RateLimiterDynamicConfigWrapper within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiterDynamicConfigWrapper) DeepCopyInto(out *RateLimiterDynamicConfigWrapper) {
	p := proto.Clone(in).(*RateLimiterDynamicConfigWrapper)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterDynamicConfigWrapper. Required by controller-gen.
func (in *RateLimiterDynamicConfigWrapper) DeepCopy() *RateLimiterDynamicConfigWrapper {
	if in == nil {
		return nil
	}
	out := new(RateLimiterDynamicConfigWrapper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterDynamicConfigWrapper. Required by controller-gen.
func (in *RateLimiterDynamicConfigWrapper) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
