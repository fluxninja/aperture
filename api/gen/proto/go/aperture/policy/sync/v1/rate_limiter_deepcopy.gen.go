// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package syncv1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using RateLimiterWrapper within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiterWrapper) DeepCopyInto(out *RateLimiterWrapper) {
	p := proto.Clone(in).(*RateLimiterWrapper)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterWrapper. Required by controller-gen.
func (in *RateLimiterWrapper) DeepCopy() *RateLimiterWrapper {
	if in == nil {
		return nil
	}
	out := new(RateLimiterWrapper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterWrapper. Required by controller-gen.
func (in *RateLimiterWrapper) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using RateLimiterDecisionWrapper within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiterDecisionWrapper) DeepCopyInto(out *RateLimiterDecisionWrapper) {
	p := proto.Clone(in).(*RateLimiterDecisionWrapper)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterDecisionWrapper. Required by controller-gen.
func (in *RateLimiterDecisionWrapper) DeepCopy() *RateLimiterDecisionWrapper {
	if in == nil {
		return nil
	}
	out := new(RateLimiterDecisionWrapper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterDecisionWrapper. Required by controller-gen.
func (in *RateLimiterDecisionWrapper) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using RateLimiterDecision within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiterDecision) DeepCopyInto(out *RateLimiterDecision) {
	p := proto.Clone(in).(*RateLimiterDecision)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterDecision. Required by controller-gen.
func (in *RateLimiterDecision) DeepCopy() *RateLimiterDecision {
	if in == nil {
		return nil
	}
	out := new(RateLimiterDecision)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterDecision. Required by controller-gen.
func (in *RateLimiterDecision) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
