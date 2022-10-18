// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package languagev1

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto supports using GetPoliciesResponse within kubernetes types, where deepcopy-gen is used.
func (in *GetPoliciesResponse) DeepCopyInto(out *GetPoliciesResponse) {
	p := proto.Clone(in).(*GetPoliciesResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GetPoliciesResponse. Required by controller-gen.
func (in *GetPoliciesResponse) DeepCopy() *GetPoliciesResponse {
	if in == nil {
		return nil
	}
	out := new(GetPoliciesResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new GetPoliciesResponse. Required by controller-gen.
func (in *GetPoliciesResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Policies within kubernetes types, where deepcopy-gen is used.
func (in *Policies) DeepCopyInto(out *Policies) {
	p := proto.Clone(in).(*Policies)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Policies. Required by controller-gen.
func (in *Policies) DeepCopy() *Policies {
	if in == nil {
		return nil
	}
	out := new(Policies)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Policies. Required by controller-gen.
func (in *Policies) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Policy within kubernetes types, where deepcopy-gen is used.
func (in *Policy) DeepCopyInto(out *Policy) {
	p := proto.Clone(in).(*Policy)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Policy. Required by controller-gen.
func (in *Policy) DeepCopy() *Policy {
	if in == nil {
		return nil
	}
	out := new(Policy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Policy. Required by controller-gen.
func (in *Policy) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Circuit within kubernetes types, where deepcopy-gen is used.
func (in *Circuit) DeepCopyInto(out *Circuit) {
	p := proto.Clone(in).(*Circuit)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Circuit. Required by controller-gen.
func (in *Circuit) DeepCopy() *Circuit {
	if in == nil {
		return nil
	}
	out := new(Circuit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Circuit. Required by controller-gen.
func (in *Circuit) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Resources within kubernetes types, where deepcopy-gen is used.
func (in *Resources) DeepCopyInto(out *Resources) {
	p := proto.Clone(in).(*Resources)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resources. Required by controller-gen.
func (in *Resources) DeepCopy() *Resources {
	if in == nil {
		return nil
	}
	out := new(Resources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Resources. Required by controller-gen.
func (in *Resources) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Component within kubernetes types, where deepcopy-gen is used.
func (in *Component) DeepCopyInto(out *Component) {
	p := proto.Clone(in).(*Component)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Component. Required by controller-gen.
func (in *Component) DeepCopy() *Component {
	if in == nil {
		return nil
	}
	out := new(Component)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Component. Required by controller-gen.
func (in *Component) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using InPort within kubernetes types, where deepcopy-gen is used.
func (in *InPort) DeepCopyInto(out *InPort) {
	p := proto.Clone(in).(*InPort)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InPort. Required by controller-gen.
func (in *InPort) DeepCopy() *InPort {
	if in == nil {
		return nil
	}
	out := new(InPort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new InPort. Required by controller-gen.
func (in *InPort) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using OutPort within kubernetes types, where deepcopy-gen is used.
func (in *OutPort) DeepCopyInto(out *OutPort) {
	p := proto.Clone(in).(*OutPort)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutPort. Required by controller-gen.
func (in *OutPort) DeepCopy() *OutPort {
	if in == nil {
		return nil
	}
	out := new(OutPort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new OutPort. Required by controller-gen.
func (in *OutPort) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using GradientController within kubernetes types, where deepcopy-gen is used.
func (in *GradientController) DeepCopyInto(out *GradientController) {
	p := proto.Clone(in).(*GradientController)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GradientController. Required by controller-gen.
func (in *GradientController) DeepCopy() *GradientController {
	if in == nil {
		return nil
	}
	out := new(GradientController)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new GradientController. Required by controller-gen.
func (in *GradientController) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using GradientController_Ins within kubernetes types, where deepcopy-gen is used.
func (in *GradientController_Ins) DeepCopyInto(out *GradientController_Ins) {
	p := proto.Clone(in).(*GradientController_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GradientController_Ins. Required by controller-gen.
func (in *GradientController_Ins) DeepCopy() *GradientController_Ins {
	if in == nil {
		return nil
	}
	out := new(GradientController_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new GradientController_Ins. Required by controller-gen.
func (in *GradientController_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using GradientController_Outs within kubernetes types, where deepcopy-gen is used.
func (in *GradientController_Outs) DeepCopyInto(out *GradientController_Outs) {
	p := proto.Clone(in).(*GradientController_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GradientController_Outs. Required by controller-gen.
func (in *GradientController_Outs) DeepCopy() *GradientController_Outs {
	if in == nil {
		return nil
	}
	out := new(GradientController_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new GradientController_Outs. Required by controller-gen.
func (in *GradientController_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ControllerDynamicConfig within kubernetes types, where deepcopy-gen is used.
func (in *ControllerDynamicConfig) DeepCopyInto(out *ControllerDynamicConfig) {
	p := proto.Clone(in).(*ControllerDynamicConfig)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerDynamicConfig. Required by controller-gen.
func (in *ControllerDynamicConfig) DeepCopy() *ControllerDynamicConfig {
	if in == nil {
		return nil
	}
	out := new(ControllerDynamicConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ControllerDynamicConfig. Required by controller-gen.
func (in *ControllerDynamicConfig) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using EMA within kubernetes types, where deepcopy-gen is used.
func (in *EMA) DeepCopyInto(out *EMA) {
	p := proto.Clone(in).(*EMA)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EMA. Required by controller-gen.
func (in *EMA) DeepCopy() *EMA {
	if in == nil {
		return nil
	}
	out := new(EMA)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new EMA. Required by controller-gen.
func (in *EMA) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using EMA_Ins within kubernetes types, where deepcopy-gen is used.
func (in *EMA_Ins) DeepCopyInto(out *EMA_Ins) {
	p := proto.Clone(in).(*EMA_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EMA_Ins. Required by controller-gen.
func (in *EMA_Ins) DeepCopy() *EMA_Ins {
	if in == nil {
		return nil
	}
	out := new(EMA_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new EMA_Ins. Required by controller-gen.
func (in *EMA_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using EMA_Outs within kubernetes types, where deepcopy-gen is used.
func (in *EMA_Outs) DeepCopyInto(out *EMA_Outs) {
	p := proto.Clone(in).(*EMA_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EMA_Outs. Required by controller-gen.
func (in *EMA_Outs) DeepCopy() *EMA_Outs {
	if in == nil {
		return nil
	}
	out := new(EMA_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new EMA_Outs. Required by controller-gen.
func (in *EMA_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ArithmeticCombinator within kubernetes types, where deepcopy-gen is used.
func (in *ArithmeticCombinator) DeepCopyInto(out *ArithmeticCombinator) {
	p := proto.Clone(in).(*ArithmeticCombinator)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArithmeticCombinator. Required by controller-gen.
func (in *ArithmeticCombinator) DeepCopy() *ArithmeticCombinator {
	if in == nil {
		return nil
	}
	out := new(ArithmeticCombinator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ArithmeticCombinator. Required by controller-gen.
func (in *ArithmeticCombinator) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ArithmeticCombinator_Ins within kubernetes types, where deepcopy-gen is used.
func (in *ArithmeticCombinator_Ins) DeepCopyInto(out *ArithmeticCombinator_Ins) {
	p := proto.Clone(in).(*ArithmeticCombinator_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArithmeticCombinator_Ins. Required by controller-gen.
func (in *ArithmeticCombinator_Ins) DeepCopy() *ArithmeticCombinator_Ins {
	if in == nil {
		return nil
	}
	out := new(ArithmeticCombinator_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ArithmeticCombinator_Ins. Required by controller-gen.
func (in *ArithmeticCombinator_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ArithmeticCombinator_Outs within kubernetes types, where deepcopy-gen is used.
func (in *ArithmeticCombinator_Outs) DeepCopyInto(out *ArithmeticCombinator_Outs) {
	p := proto.Clone(in).(*ArithmeticCombinator_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArithmeticCombinator_Outs. Required by controller-gen.
func (in *ArithmeticCombinator_Outs) DeepCopy() *ArithmeticCombinator_Outs {
	if in == nil {
		return nil
	}
	out := new(ArithmeticCombinator_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ArithmeticCombinator_Outs. Required by controller-gen.
func (in *ArithmeticCombinator_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Decider within kubernetes types, where deepcopy-gen is used.
func (in *Decider) DeepCopyInto(out *Decider) {
	p := proto.Clone(in).(*Decider)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Decider. Required by controller-gen.
func (in *Decider) DeepCopy() *Decider {
	if in == nil {
		return nil
	}
	out := new(Decider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Decider. Required by controller-gen.
func (in *Decider) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Decider_Ins within kubernetes types, where deepcopy-gen is used.
func (in *Decider_Ins) DeepCopyInto(out *Decider_Ins) {
	p := proto.Clone(in).(*Decider_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Decider_Ins. Required by controller-gen.
func (in *Decider_Ins) DeepCopy() *Decider_Ins {
	if in == nil {
		return nil
	}
	out := new(Decider_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Decider_Ins. Required by controller-gen.
func (in *Decider_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Decider_Outs within kubernetes types, where deepcopy-gen is used.
func (in *Decider_Outs) DeepCopyInto(out *Decider_Outs) {
	p := proto.Clone(in).(*Decider_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Decider_Outs. Required by controller-gen.
func (in *Decider_Outs) DeepCopy() *Decider_Outs {
	if in == nil {
		return nil
	}
	out := new(Decider_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Decider_Outs. Required by controller-gen.
func (in *Decider_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Switcher within kubernetes types, where deepcopy-gen is used.
func (in *Switcher) DeepCopyInto(out *Switcher) {
	p := proto.Clone(in).(*Switcher)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Switcher. Required by controller-gen.
func (in *Switcher) DeepCopy() *Switcher {
	if in == nil {
		return nil
	}
	out := new(Switcher)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Switcher. Required by controller-gen.
func (in *Switcher) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Switcher_Ins within kubernetes types, where deepcopy-gen is used.
func (in *Switcher_Ins) DeepCopyInto(out *Switcher_Ins) {
	p := proto.Clone(in).(*Switcher_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Switcher_Ins. Required by controller-gen.
func (in *Switcher_Ins) DeepCopy() *Switcher_Ins {
	if in == nil {
		return nil
	}
	out := new(Switcher_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Switcher_Ins. Required by controller-gen.
func (in *Switcher_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Switcher_Outs within kubernetes types, where deepcopy-gen is used.
func (in *Switcher_Outs) DeepCopyInto(out *Switcher_Outs) {
	p := proto.Clone(in).(*Switcher_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Switcher_Outs. Required by controller-gen.
func (in *Switcher_Outs) DeepCopy() *Switcher_Outs {
	if in == nil {
		return nil
	}
	out := new(Switcher_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Switcher_Outs. Required by controller-gen.
func (in *Switcher_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using RateLimiter within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiter) DeepCopyInto(out *RateLimiter) {
	p := proto.Clone(in).(*RateLimiter)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter. Required by controller-gen.
func (in *RateLimiter) DeepCopy() *RateLimiter {
	if in == nil {
		return nil
	}
	out := new(RateLimiter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter. Required by controller-gen.
func (in *RateLimiter) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using RateLimiter_LazySync within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiter_LazySync) DeepCopyInto(out *RateLimiter_LazySync) {
	p := proto.Clone(in).(*RateLimiter_LazySync)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_LazySync. Required by controller-gen.
func (in *RateLimiter_LazySync) DeepCopy() *RateLimiter_LazySync {
	if in == nil {
		return nil
	}
	out := new(RateLimiter_LazySync)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_LazySync. Required by controller-gen.
func (in *RateLimiter_LazySync) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using RateLimiter_DynamicConfig within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiter_DynamicConfig) DeepCopyInto(out *RateLimiter_DynamicConfig) {
	p := proto.Clone(in).(*RateLimiter_DynamicConfig)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_DynamicConfig. Required by controller-gen.
func (in *RateLimiter_DynamicConfig) DeepCopy() *RateLimiter_DynamicConfig {
	if in == nil {
		return nil
	}
	out := new(RateLimiter_DynamicConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_DynamicConfig. Required by controller-gen.
func (in *RateLimiter_DynamicConfig) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using RateLimiter_Override within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiter_Override) DeepCopyInto(out *RateLimiter_Override) {
	p := proto.Clone(in).(*RateLimiter_Override)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_Override. Required by controller-gen.
func (in *RateLimiter_Override) DeepCopy() *RateLimiter_Override {
	if in == nil {
		return nil
	}
	out := new(RateLimiter_Override)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_Override. Required by controller-gen.
func (in *RateLimiter_Override) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using RateLimiter_Ins within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiter_Ins) DeepCopyInto(out *RateLimiter_Ins) {
	p := proto.Clone(in).(*RateLimiter_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_Ins. Required by controller-gen.
func (in *RateLimiter_Ins) DeepCopy() *RateLimiter_Ins {
	if in == nil {
		return nil
	}
	out := new(RateLimiter_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_Ins. Required by controller-gen.
func (in *RateLimiter_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ConcurrencyLimiter within kubernetes types, where deepcopy-gen is used.
func (in *ConcurrencyLimiter) DeepCopyInto(out *ConcurrencyLimiter) {
	p := proto.Clone(in).(*ConcurrencyLimiter)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConcurrencyLimiter. Required by controller-gen.
func (in *ConcurrencyLimiter) DeepCopy() *ConcurrencyLimiter {
	if in == nil {
		return nil
	}
	out := new(ConcurrencyLimiter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ConcurrencyLimiter. Required by controller-gen.
func (in *ConcurrencyLimiter) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Scheduler within kubernetes types, where deepcopy-gen is used.
func (in *Scheduler) DeepCopyInto(out *Scheduler) {
	p := proto.Clone(in).(*Scheduler)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler. Required by controller-gen.
func (in *Scheduler) DeepCopy() *Scheduler {
	if in == nil {
		return nil
	}
	out := new(Scheduler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler. Required by controller-gen.
func (in *Scheduler) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Scheduler_WorkloadParameters within kubernetes types, where deepcopy-gen is used.
func (in *Scheduler_WorkloadParameters) DeepCopyInto(out *Scheduler_WorkloadParameters) {
	p := proto.Clone(in).(*Scheduler_WorkloadParameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler_WorkloadParameters. Required by controller-gen.
func (in *Scheduler_WorkloadParameters) DeepCopy() *Scheduler_WorkloadParameters {
	if in == nil {
		return nil
	}
	out := new(Scheduler_WorkloadParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler_WorkloadParameters. Required by controller-gen.
func (in *Scheduler_WorkloadParameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Scheduler_Workload within kubernetes types, where deepcopy-gen is used.
func (in *Scheduler_Workload) DeepCopyInto(out *Scheduler_Workload) {
	p := proto.Clone(in).(*Scheduler_Workload)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler_Workload. Required by controller-gen.
func (in *Scheduler_Workload) DeepCopy() *Scheduler_Workload {
	if in == nil {
		return nil
	}
	out := new(Scheduler_Workload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler_Workload. Required by controller-gen.
func (in *Scheduler_Workload) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Scheduler_Outs within kubernetes types, where deepcopy-gen is used.
func (in *Scheduler_Outs) DeepCopyInto(out *Scheduler_Outs) {
	p := proto.Clone(in).(*Scheduler_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler_Outs. Required by controller-gen.
func (in *Scheduler_Outs) DeepCopy() *Scheduler_Outs {
	if in == nil {
		return nil
	}
	out := new(Scheduler_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler_Outs. Required by controller-gen.
func (in *Scheduler_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadActuator within kubernetes types, where deepcopy-gen is used.
func (in *LoadActuator) DeepCopyInto(out *LoadActuator) {
	p := proto.Clone(in).(*LoadActuator)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadActuator. Required by controller-gen.
func (in *LoadActuator) DeepCopy() *LoadActuator {
	if in == nil {
		return nil
	}
	out := new(LoadActuator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadActuator. Required by controller-gen.
func (in *LoadActuator) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadActuator_Ins within kubernetes types, where deepcopy-gen is used.
func (in *LoadActuator_Ins) DeepCopyInto(out *LoadActuator_Ins) {
	p := proto.Clone(in).(*LoadActuator_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadActuator_Ins. Required by controller-gen.
func (in *LoadActuator_Ins) DeepCopy() *LoadActuator_Ins {
	if in == nil {
		return nil
	}
	out := new(LoadActuator_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadActuator_Ins. Required by controller-gen.
func (in *LoadActuator_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PromQL within kubernetes types, where deepcopy-gen is used.
func (in *PromQL) DeepCopyInto(out *PromQL) {
	p := proto.Clone(in).(*PromQL)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromQL. Required by controller-gen.
func (in *PromQL) DeepCopy() *PromQL {
	if in == nil {
		return nil
	}
	out := new(PromQL)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PromQL. Required by controller-gen.
func (in *PromQL) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PromQL_Outs within kubernetes types, where deepcopy-gen is used.
func (in *PromQL_Outs) DeepCopyInto(out *PromQL_Outs) {
	p := proto.Clone(in).(*PromQL_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PromQL_Outs. Required by controller-gen.
func (in *PromQL_Outs) DeepCopy() *PromQL_Outs {
	if in == nil {
		return nil
	}
	out := new(PromQL_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PromQL_Outs. Required by controller-gen.
func (in *PromQL_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Constant within kubernetes types, where deepcopy-gen is used.
func (in *Constant) DeepCopyInto(out *Constant) {
	p := proto.Clone(in).(*Constant)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Constant. Required by controller-gen.
func (in *Constant) DeepCopy() *Constant {
	if in == nil {
		return nil
	}
	out := new(Constant)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Constant. Required by controller-gen.
func (in *Constant) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Constant_Outs within kubernetes types, where deepcopy-gen is used.
func (in *Constant_Outs) DeepCopyInto(out *Constant_Outs) {
	p := proto.Clone(in).(*Constant_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Constant_Outs. Required by controller-gen.
func (in *Constant_Outs) DeepCopy() *Constant_Outs {
	if in == nil {
		return nil
	}
	out := new(Constant_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Constant_Outs. Required by controller-gen.
func (in *Constant_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Sqrt within kubernetes types, where deepcopy-gen is used.
func (in *Sqrt) DeepCopyInto(out *Sqrt) {
	p := proto.Clone(in).(*Sqrt)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Sqrt. Required by controller-gen.
func (in *Sqrt) DeepCopy() *Sqrt {
	if in == nil {
		return nil
	}
	out := new(Sqrt)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Sqrt. Required by controller-gen.
func (in *Sqrt) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Sqrt_Ins within kubernetes types, where deepcopy-gen is used.
func (in *Sqrt_Ins) DeepCopyInto(out *Sqrt_Ins) {
	p := proto.Clone(in).(*Sqrt_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Sqrt_Ins. Required by controller-gen.
func (in *Sqrt_Ins) DeepCopy() *Sqrt_Ins {
	if in == nil {
		return nil
	}
	out := new(Sqrt_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Sqrt_Ins. Required by controller-gen.
func (in *Sqrt_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Sqrt_Outs within kubernetes types, where deepcopy-gen is used.
func (in *Sqrt_Outs) DeepCopyInto(out *Sqrt_Outs) {
	p := proto.Clone(in).(*Sqrt_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Sqrt_Outs. Required by controller-gen.
func (in *Sqrt_Outs) DeepCopy() *Sqrt_Outs {
	if in == nil {
		return nil
	}
	out := new(Sqrt_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Sqrt_Outs. Required by controller-gen.
func (in *Sqrt_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Extrapolator within kubernetes types, where deepcopy-gen is used.
func (in *Extrapolator) DeepCopyInto(out *Extrapolator) {
	p := proto.Clone(in).(*Extrapolator)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Extrapolator. Required by controller-gen.
func (in *Extrapolator) DeepCopy() *Extrapolator {
	if in == nil {
		return nil
	}
	out := new(Extrapolator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Extrapolator. Required by controller-gen.
func (in *Extrapolator) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Extrapolator_Ins within kubernetes types, where deepcopy-gen is used.
func (in *Extrapolator_Ins) DeepCopyInto(out *Extrapolator_Ins) {
	p := proto.Clone(in).(*Extrapolator_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Extrapolator_Ins. Required by controller-gen.
func (in *Extrapolator_Ins) DeepCopy() *Extrapolator_Ins {
	if in == nil {
		return nil
	}
	out := new(Extrapolator_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Extrapolator_Ins. Required by controller-gen.
func (in *Extrapolator_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Extrapolator_Outs within kubernetes types, where deepcopy-gen is used.
func (in *Extrapolator_Outs) DeepCopyInto(out *Extrapolator_Outs) {
	p := proto.Clone(in).(*Extrapolator_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Extrapolator_Outs. Required by controller-gen.
func (in *Extrapolator_Outs) DeepCopy() *Extrapolator_Outs {
	if in == nil {
		return nil
	}
	out := new(Extrapolator_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Extrapolator_Outs. Required by controller-gen.
func (in *Extrapolator_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Max within kubernetes types, where deepcopy-gen is used.
func (in *Max) DeepCopyInto(out *Max) {
	p := proto.Clone(in).(*Max)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Max. Required by controller-gen.
func (in *Max) DeepCopy() *Max {
	if in == nil {
		return nil
	}
	out := new(Max)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Max. Required by controller-gen.
func (in *Max) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Max_Ins within kubernetes types, where deepcopy-gen is used.
func (in *Max_Ins) DeepCopyInto(out *Max_Ins) {
	p := proto.Clone(in).(*Max_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Max_Ins. Required by controller-gen.
func (in *Max_Ins) DeepCopy() *Max_Ins {
	if in == nil {
		return nil
	}
	out := new(Max_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Max_Ins. Required by controller-gen.
func (in *Max_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Max_Outs within kubernetes types, where deepcopy-gen is used.
func (in *Max_Outs) DeepCopyInto(out *Max_Outs) {
	p := proto.Clone(in).(*Max_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Max_Outs. Required by controller-gen.
func (in *Max_Outs) DeepCopy() *Max_Outs {
	if in == nil {
		return nil
	}
	out := new(Max_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Max_Outs. Required by controller-gen.
func (in *Max_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Min within kubernetes types, where deepcopy-gen is used.
func (in *Min) DeepCopyInto(out *Min) {
	p := proto.Clone(in).(*Min)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Min. Required by controller-gen.
func (in *Min) DeepCopy() *Min {
	if in == nil {
		return nil
	}
	out := new(Min)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Min. Required by controller-gen.
func (in *Min) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Min_Ins within kubernetes types, where deepcopy-gen is used.
func (in *Min_Ins) DeepCopyInto(out *Min_Ins) {
	p := proto.Clone(in).(*Min_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Min_Ins. Required by controller-gen.
func (in *Min_Ins) DeepCopy() *Min_Ins {
	if in == nil {
		return nil
	}
	out := new(Min_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Min_Ins. Required by controller-gen.
func (in *Min_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Min_Outs within kubernetes types, where deepcopy-gen is used.
func (in *Min_Outs) DeepCopyInto(out *Min_Outs) {
	p := proto.Clone(in).(*Min_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Min_Outs. Required by controller-gen.
func (in *Min_Outs) DeepCopy() *Min_Outs {
	if in == nil {
		return nil
	}
	out := new(Min_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Min_Outs. Required by controller-gen.
func (in *Min_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FirstValid within kubernetes types, where deepcopy-gen is used.
func (in *FirstValid) DeepCopyInto(out *FirstValid) {
	p := proto.Clone(in).(*FirstValid)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirstValid. Required by controller-gen.
func (in *FirstValid) DeepCopy() *FirstValid {
	if in == nil {
		return nil
	}
	out := new(FirstValid)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FirstValid. Required by controller-gen.
func (in *FirstValid) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FirstValid_Ins within kubernetes types, where deepcopy-gen is used.
func (in *FirstValid_Ins) DeepCopyInto(out *FirstValid_Ins) {
	p := proto.Clone(in).(*FirstValid_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirstValid_Ins. Required by controller-gen.
func (in *FirstValid_Ins) DeepCopy() *FirstValid_Ins {
	if in == nil {
		return nil
	}
	out := new(FirstValid_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FirstValid_Ins. Required by controller-gen.
func (in *FirstValid_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FirstValid_Outs within kubernetes types, where deepcopy-gen is used.
func (in *FirstValid_Outs) DeepCopyInto(out *FirstValid_Outs) {
	p := proto.Clone(in).(*FirstValid_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirstValid_Outs. Required by controller-gen.
func (in *FirstValid_Outs) DeepCopy() *FirstValid_Outs {
	if in == nil {
		return nil
	}
	out := new(FirstValid_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FirstValid_Outs. Required by controller-gen.
func (in *FirstValid_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
