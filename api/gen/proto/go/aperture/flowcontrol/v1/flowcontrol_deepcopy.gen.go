// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package flowcontrolv1

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto supports using CheckRequest within kubernetes types, where deepcopy-gen is used.
func (in *CheckRequest) DeepCopyInto(out *CheckRequest) {
	p := proto.Clone(in).(*CheckRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckRequest. Required by controller-gen.
func (in *CheckRequest) DeepCopy() *CheckRequest {
	if in == nil {
		return nil
	}
	out := new(CheckRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new CheckRequest. Required by controller-gen.
func (in *CheckRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using CheckResponse within kubernetes types, where deepcopy-gen is used.
func (in *CheckResponse) DeepCopyInto(out *CheckResponse) {
	p := proto.Clone(in).(*CheckResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckResponse. Required by controller-gen.
func (in *CheckResponse) DeepCopy() *CheckResponse {
	if in == nil {
		return nil
	}
	out := new(CheckResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new CheckResponse. Required by controller-gen.
func (in *CheckResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DecisionReason within kubernetes types, where deepcopy-gen is used.
func (in *DecisionReason) DeepCopyInto(out *DecisionReason) {
	p := proto.Clone(in).(*DecisionReason)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DecisionReason. Required by controller-gen.
func (in *DecisionReason) DeepCopy() *DecisionReason {
	if in == nil {
		return nil
	}
	out := new(DecisionReason)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DecisionReason. Required by controller-gen.
func (in *DecisionReason) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LimiterDecision within kubernetes types, where deepcopy-gen is used.
func (in *LimiterDecision) DeepCopyInto(out *LimiterDecision) {
	p := proto.Clone(in).(*LimiterDecision)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LimiterDecision. Required by controller-gen.
func (in *LimiterDecision) DeepCopy() *LimiterDecision {
	if in == nil {
		return nil
	}
	out := new(LimiterDecision)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LimiterDecision. Required by controller-gen.
func (in *LimiterDecision) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LimiterDecision_RateLimiter within kubernetes types, where deepcopy-gen is used.
func (in *LimiterDecision_RateLimiter) DeepCopyInto(out *LimiterDecision_RateLimiter) {
	p := proto.Clone(in).(*LimiterDecision_RateLimiter)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LimiterDecision_RateLimiter. Required by controller-gen.
func (in *LimiterDecision_RateLimiter) DeepCopy() *LimiterDecision_RateLimiter {
	if in == nil {
		return nil
	}
	out := new(LimiterDecision_RateLimiter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LimiterDecision_RateLimiter. Required by controller-gen.
func (in *LimiterDecision_RateLimiter) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LimiterDecision_ConcurrencyLimiter within kubernetes types, where deepcopy-gen is used.
func (in *LimiterDecision_ConcurrencyLimiter) DeepCopyInto(out *LimiterDecision_ConcurrencyLimiter) {
	p := proto.Clone(in).(*LimiterDecision_ConcurrencyLimiter)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LimiterDecision_ConcurrencyLimiter. Required by controller-gen.
func (in *LimiterDecision_ConcurrencyLimiter) DeepCopy() *LimiterDecision_ConcurrencyLimiter {
	if in == nil {
		return nil
	}
	out := new(LimiterDecision_ConcurrencyLimiter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LimiterDecision_ConcurrencyLimiter. Required by controller-gen.
func (in *LimiterDecision_ConcurrencyLimiter) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FluxMeter within kubernetes types, where deepcopy-gen is used.
func (in *FluxMeter) DeepCopyInto(out *FluxMeter) {
	p := proto.Clone(in).(*FluxMeter)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter. Required by controller-gen.
func (in *FluxMeter) DeepCopy() *FluxMeter {
	if in == nil {
		return nil
	}
	out := new(FluxMeter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter. Required by controller-gen.
func (in *FluxMeter) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
