// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package languagev1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using Selector within kubernetes types, where deepcopy-gen is used.
func (in *Selector) DeepCopyInto(out *Selector) {
	p := proto.Clone(in).(*Selector)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Selector. Required by controller-gen.
func (in *Selector) DeepCopy() *Selector {
	if in == nil {
		return nil
	}
	out := new(Selector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Selector. Required by controller-gen.
func (in *Selector) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FlowControlResources within kubernetes types, where deepcopy-gen is used.
func (in *FlowControlResources) DeepCopyInto(out *FlowControlResources) {
	p := proto.Clone(in).(*FlowControlResources)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlowControlResources. Required by controller-gen.
func (in *FlowControlResources) DeepCopy() *FlowControlResources {
	if in == nil {
		return nil
	}
	out := new(FlowControlResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FlowControlResources. Required by controller-gen.
func (in *FlowControlResources) DeepCopyInterface() interface{} {
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

// DeepCopyInto supports using FluxMeter_StaticBuckets within kubernetes types, where deepcopy-gen is used.
func (in *FluxMeter_StaticBuckets) DeepCopyInto(out *FluxMeter_StaticBuckets) {
	p := proto.Clone(in).(*FluxMeter_StaticBuckets)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter_StaticBuckets. Required by controller-gen.
func (in *FluxMeter_StaticBuckets) DeepCopy() *FluxMeter_StaticBuckets {
	if in == nil {
		return nil
	}
	out := new(FluxMeter_StaticBuckets)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter_StaticBuckets. Required by controller-gen.
func (in *FluxMeter_StaticBuckets) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FluxMeter_LinearBuckets within kubernetes types, where deepcopy-gen is used.
func (in *FluxMeter_LinearBuckets) DeepCopyInto(out *FluxMeter_LinearBuckets) {
	p := proto.Clone(in).(*FluxMeter_LinearBuckets)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter_LinearBuckets. Required by controller-gen.
func (in *FluxMeter_LinearBuckets) DeepCopy() *FluxMeter_LinearBuckets {
	if in == nil {
		return nil
	}
	out := new(FluxMeter_LinearBuckets)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter_LinearBuckets. Required by controller-gen.
func (in *FluxMeter_LinearBuckets) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FluxMeter_ExponentialBuckets within kubernetes types, where deepcopy-gen is used.
func (in *FluxMeter_ExponentialBuckets) DeepCopyInto(out *FluxMeter_ExponentialBuckets) {
	p := proto.Clone(in).(*FluxMeter_ExponentialBuckets)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter_ExponentialBuckets. Required by controller-gen.
func (in *FluxMeter_ExponentialBuckets) DeepCopy() *FluxMeter_ExponentialBuckets {
	if in == nil {
		return nil
	}
	out := new(FluxMeter_ExponentialBuckets)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter_ExponentialBuckets. Required by controller-gen.
func (in *FluxMeter_ExponentialBuckets) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FluxMeter_ExponentialBucketsRange within kubernetes types, where deepcopy-gen is used.
func (in *FluxMeter_ExponentialBucketsRange) DeepCopyInto(out *FluxMeter_ExponentialBucketsRange) {
	p := proto.Clone(in).(*FluxMeter_ExponentialBucketsRange)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter_ExponentialBucketsRange. Required by controller-gen.
func (in *FluxMeter_ExponentialBucketsRange) DeepCopy() *FluxMeter_ExponentialBucketsRange {
	if in == nil {
		return nil
	}
	out := new(FluxMeter_ExponentialBucketsRange)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FluxMeter_ExponentialBucketsRange. Required by controller-gen.
func (in *FluxMeter_ExponentialBucketsRange) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Classifier within kubernetes types, where deepcopy-gen is used.
func (in *Classifier) DeepCopyInto(out *Classifier) {
	p := proto.Clone(in).(*Classifier)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Classifier. Required by controller-gen.
func (in *Classifier) DeepCopy() *Classifier {
	if in == nil {
		return nil
	}
	out := new(Classifier)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Classifier. Required by controller-gen.
func (in *Classifier) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Rule within kubernetes types, where deepcopy-gen is used.
func (in *Rule) DeepCopyInto(out *Rule) {
	p := proto.Clone(in).(*Rule)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule. Required by controller-gen.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Rule. Required by controller-gen.
func (in *Rule) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Rego within kubernetes types, where deepcopy-gen is used.
func (in *Rego) DeepCopyInto(out *Rego) {
	p := proto.Clone(in).(*Rego)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rego. Required by controller-gen.
func (in *Rego) DeepCopy() *Rego {
	if in == nil {
		return nil
	}
	out := new(Rego)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Rego. Required by controller-gen.
func (in *Rego) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Rego_LabelProperties within kubernetes types, where deepcopy-gen is used.
func (in *Rego_LabelProperties) DeepCopyInto(out *Rego_LabelProperties) {
	p := proto.Clone(in).(*Rego_LabelProperties)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rego_LabelProperties. Required by controller-gen.
func (in *Rego_LabelProperties) DeepCopy() *Rego_LabelProperties {
	if in == nil {
		return nil
	}
	out := new(Rego_LabelProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Rego_LabelProperties. Required by controller-gen.
func (in *Rego_LabelProperties) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Extractor within kubernetes types, where deepcopy-gen is used.
func (in *Extractor) DeepCopyInto(out *Extractor) {
	p := proto.Clone(in).(*Extractor)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Extractor. Required by controller-gen.
func (in *Extractor) DeepCopy() *Extractor {
	if in == nil {
		return nil
	}
	out := new(Extractor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Extractor. Required by controller-gen.
func (in *Extractor) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using JSONExtractor within kubernetes types, where deepcopy-gen is used.
func (in *JSONExtractor) DeepCopyInto(out *JSONExtractor) {
	p := proto.Clone(in).(*JSONExtractor)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONExtractor. Required by controller-gen.
func (in *JSONExtractor) DeepCopy() *JSONExtractor {
	if in == nil {
		return nil
	}
	out := new(JSONExtractor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new JSONExtractor. Required by controller-gen.
func (in *JSONExtractor) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AddressExtractor within kubernetes types, where deepcopy-gen is used.
func (in *AddressExtractor) DeepCopyInto(out *AddressExtractor) {
	p := proto.Clone(in).(*AddressExtractor)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AddressExtractor. Required by controller-gen.
func (in *AddressExtractor) DeepCopy() *AddressExtractor {
	if in == nil {
		return nil
	}
	out := new(AddressExtractor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AddressExtractor. Required by controller-gen.
func (in *AddressExtractor) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using JWTExtractor within kubernetes types, where deepcopy-gen is used.
func (in *JWTExtractor) DeepCopyInto(out *JWTExtractor) {
	p := proto.Clone(in).(*JWTExtractor)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JWTExtractor. Required by controller-gen.
func (in *JWTExtractor) DeepCopy() *JWTExtractor {
	if in == nil {
		return nil
	}
	out := new(JWTExtractor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new JWTExtractor. Required by controller-gen.
func (in *JWTExtractor) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PathTemplateMatcher within kubernetes types, where deepcopy-gen is used.
func (in *PathTemplateMatcher) DeepCopyInto(out *PathTemplateMatcher) {
	p := proto.Clone(in).(*PathTemplateMatcher)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PathTemplateMatcher. Required by controller-gen.
func (in *PathTemplateMatcher) DeepCopy() *PathTemplateMatcher {
	if in == nil {
		return nil
	}
	out := new(PathTemplateMatcher)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PathTemplateMatcher. Required by controller-gen.
func (in *PathTemplateMatcher) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FlowControl within kubernetes types, where deepcopy-gen is used.
func (in *FlowControl) DeepCopyInto(out *FlowControl) {
	p := proto.Clone(in).(*FlowControl)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlowControl. Required by controller-gen.
func (in *FlowControl) DeepCopy() *FlowControl {
	if in == nil {
		return nil
	}
	out := new(FlowControl)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FlowControl. Required by controller-gen.
func (in *FlowControl) DeepCopyInterface() interface{} {
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

// DeepCopyInto supports using RateLimiter_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiter_Parameters) DeepCopyInto(out *RateLimiter_Parameters) {
	p := proto.Clone(in).(*RateLimiter_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_Parameters. Required by controller-gen.
func (in *RateLimiter_Parameters) DeepCopy() *RateLimiter_Parameters {
	if in == nil {
		return nil
	}
	out := new(RateLimiter_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_Parameters. Required by controller-gen.
func (in *RateLimiter_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using RateLimiter_Parameters_LazySync within kubernetes types, where deepcopy-gen is used.
func (in *RateLimiter_Parameters_LazySync) DeepCopyInto(out *RateLimiter_Parameters_LazySync) {
	p := proto.Clone(in).(*RateLimiter_Parameters_LazySync)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_Parameters_LazySync. Required by controller-gen.
func (in *RateLimiter_Parameters_LazySync) DeepCopy() *RateLimiter_Parameters_LazySync {
	if in == nil {
		return nil
	}
	out := new(RateLimiter_Parameters_LazySync)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter_Parameters_LazySync. Required by controller-gen.
func (in *RateLimiter_Parameters_LazySync) DeepCopyInterface() interface{} {
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

// DeepCopyInto supports using LeakyBucketRateLimiter within kubernetes types, where deepcopy-gen is used.
func (in *LeakyBucketRateLimiter) DeepCopyInto(out *LeakyBucketRateLimiter) {
	p := proto.Clone(in).(*LeakyBucketRateLimiter)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeakyBucketRateLimiter. Required by controller-gen.
func (in *LeakyBucketRateLimiter) DeepCopy() *LeakyBucketRateLimiter {
	if in == nil {
		return nil
	}
	out := new(LeakyBucketRateLimiter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LeakyBucketRateLimiter. Required by controller-gen.
func (in *LeakyBucketRateLimiter) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LeakyBucketRateLimiter_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *LeakyBucketRateLimiter_Parameters) DeepCopyInto(out *LeakyBucketRateLimiter_Parameters) {
	p := proto.Clone(in).(*LeakyBucketRateLimiter_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeakyBucketRateLimiter_Parameters. Required by controller-gen.
func (in *LeakyBucketRateLimiter_Parameters) DeepCopy() *LeakyBucketRateLimiter_Parameters {
	if in == nil {
		return nil
	}
	out := new(LeakyBucketRateLimiter_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LeakyBucketRateLimiter_Parameters. Required by controller-gen.
func (in *LeakyBucketRateLimiter_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LeakyBucketRateLimiter_Ins within kubernetes types, where deepcopy-gen is used.
func (in *LeakyBucketRateLimiter_Ins) DeepCopyInto(out *LeakyBucketRateLimiter_Ins) {
	p := proto.Clone(in).(*LeakyBucketRateLimiter_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeakyBucketRateLimiter_Ins. Required by controller-gen.
func (in *LeakyBucketRateLimiter_Ins) DeepCopy() *LeakyBucketRateLimiter_Ins {
	if in == nil {
		return nil
	}
	out := new(LeakyBucketRateLimiter_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LeakyBucketRateLimiter_Ins. Required by controller-gen.
func (in *LeakyBucketRateLimiter_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadScheduler within kubernetes types, where deepcopy-gen is used.
func (in *LoadScheduler) DeepCopyInto(out *LoadScheduler) {
	p := proto.Clone(in).(*LoadScheduler)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler. Required by controller-gen.
func (in *LoadScheduler) DeepCopy() *LoadScheduler {
	if in == nil {
		return nil
	}
	out := new(LoadScheduler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler. Required by controller-gen.
func (in *LoadScheduler) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadScheduler_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *LoadScheduler_Parameters) DeepCopyInto(out *LoadScheduler_Parameters) {
	p := proto.Clone(in).(*LoadScheduler_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler_Parameters. Required by controller-gen.
func (in *LoadScheduler_Parameters) DeepCopy() *LoadScheduler_Parameters {
	if in == nil {
		return nil
	}
	out := new(LoadScheduler_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler_Parameters. Required by controller-gen.
func (in *LoadScheduler_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadScheduler_DynamicConfig within kubernetes types, where deepcopy-gen is used.
func (in *LoadScheduler_DynamicConfig) DeepCopyInto(out *LoadScheduler_DynamicConfig) {
	p := proto.Clone(in).(*LoadScheduler_DynamicConfig)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler_DynamicConfig. Required by controller-gen.
func (in *LoadScheduler_DynamicConfig) DeepCopy() *LoadScheduler_DynamicConfig {
	if in == nil {
		return nil
	}
	out := new(LoadScheduler_DynamicConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler_DynamicConfig. Required by controller-gen.
func (in *LoadScheduler_DynamicConfig) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadScheduler_Ins within kubernetes types, where deepcopy-gen is used.
func (in *LoadScheduler_Ins) DeepCopyInto(out *LoadScheduler_Ins) {
	p := proto.Clone(in).(*LoadScheduler_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler_Ins. Required by controller-gen.
func (in *LoadScheduler_Ins) DeepCopy() *LoadScheduler_Ins {
	if in == nil {
		return nil
	}
	out := new(LoadScheduler_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler_Ins. Required by controller-gen.
func (in *LoadScheduler_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadScheduler_Outs within kubernetes types, where deepcopy-gen is used.
func (in *LoadScheduler_Outs) DeepCopyInto(out *LoadScheduler_Outs) {
	p := proto.Clone(in).(*LoadScheduler_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler_Outs. Required by controller-gen.
func (in *LoadScheduler_Outs) DeepCopy() *LoadScheduler_Outs {
	if in == nil {
		return nil
	}
	out := new(LoadScheduler_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadScheduler_Outs. Required by controller-gen.
func (in *LoadScheduler_Outs) DeepCopyInterface() interface{} {
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

// DeepCopyInto supports using Scheduler_Workload_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *Scheduler_Workload_Parameters) DeepCopyInto(out *Scheduler_Workload_Parameters) {
	p := proto.Clone(in).(*Scheduler_Workload_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler_Workload_Parameters. Required by controller-gen.
func (in *Scheduler_Workload_Parameters) DeepCopy() *Scheduler_Workload_Parameters {
	if in == nil {
		return nil
	}
	out := new(Scheduler_Workload_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Scheduler_Workload_Parameters. Required by controller-gen.
func (in *Scheduler_Workload_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AdaptiveLoadScheduler within kubernetes types, where deepcopy-gen is used.
func (in *AdaptiveLoadScheduler) DeepCopyInto(out *AdaptiveLoadScheduler) {
	p := proto.Clone(in).(*AdaptiveLoadScheduler)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdaptiveLoadScheduler. Required by controller-gen.
func (in *AdaptiveLoadScheduler) DeepCopy() *AdaptiveLoadScheduler {
	if in == nil {
		return nil
	}
	out := new(AdaptiveLoadScheduler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AdaptiveLoadScheduler. Required by controller-gen.
func (in *AdaptiveLoadScheduler) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AdaptiveLoadScheduler_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *AdaptiveLoadScheduler_Parameters) DeepCopyInto(out *AdaptiveLoadScheduler_Parameters) {
	p := proto.Clone(in).(*AdaptiveLoadScheduler_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdaptiveLoadScheduler_Parameters. Required by controller-gen.
func (in *AdaptiveLoadScheduler_Parameters) DeepCopy() *AdaptiveLoadScheduler_Parameters {
	if in == nil {
		return nil
	}
	out := new(AdaptiveLoadScheduler_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AdaptiveLoadScheduler_Parameters. Required by controller-gen.
func (in *AdaptiveLoadScheduler_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AdaptiveLoadScheduler_Ins within kubernetes types, where deepcopy-gen is used.
func (in *AdaptiveLoadScheduler_Ins) DeepCopyInto(out *AdaptiveLoadScheduler_Ins) {
	p := proto.Clone(in).(*AdaptiveLoadScheduler_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdaptiveLoadScheduler_Ins. Required by controller-gen.
func (in *AdaptiveLoadScheduler_Ins) DeepCopy() *AdaptiveLoadScheduler_Ins {
	if in == nil {
		return nil
	}
	out := new(AdaptiveLoadScheduler_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AdaptiveLoadScheduler_Ins. Required by controller-gen.
func (in *AdaptiveLoadScheduler_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AdaptiveLoadScheduler_Outs within kubernetes types, where deepcopy-gen is used.
func (in *AdaptiveLoadScheduler_Outs) DeepCopyInto(out *AdaptiveLoadScheduler_Outs) {
	p := proto.Clone(in).(*AdaptiveLoadScheduler_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdaptiveLoadScheduler_Outs. Required by controller-gen.
func (in *AdaptiveLoadScheduler_Outs) DeepCopy() *AdaptiveLoadScheduler_Outs {
	if in == nil {
		return nil
	}
	out := new(AdaptiveLoadScheduler_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AdaptiveLoadScheduler_Outs. Required by controller-gen.
func (in *AdaptiveLoadScheduler_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Regulator within kubernetes types, where deepcopy-gen is used.
func (in *Regulator) DeepCopyInto(out *Regulator) {
	p := proto.Clone(in).(*Regulator)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Regulator. Required by controller-gen.
func (in *Regulator) DeepCopy() *Regulator {
	if in == nil {
		return nil
	}
	out := new(Regulator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Regulator. Required by controller-gen.
func (in *Regulator) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Regulator_DynamicConfig within kubernetes types, where deepcopy-gen is used.
func (in *Regulator_DynamicConfig) DeepCopyInto(out *Regulator_DynamicConfig) {
	p := proto.Clone(in).(*Regulator_DynamicConfig)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Regulator_DynamicConfig. Required by controller-gen.
func (in *Regulator_DynamicConfig) DeepCopy() *Regulator_DynamicConfig {
	if in == nil {
		return nil
	}
	out := new(Regulator_DynamicConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Regulator_DynamicConfig. Required by controller-gen.
func (in *Regulator_DynamicConfig) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Regulator_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *Regulator_Parameters) DeepCopyInto(out *Regulator_Parameters) {
	p := proto.Clone(in).(*Regulator_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Regulator_Parameters. Required by controller-gen.
func (in *Regulator_Parameters) DeepCopy() *Regulator_Parameters {
	if in == nil {
		return nil
	}
	out := new(Regulator_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Regulator_Parameters. Required by controller-gen.
func (in *Regulator_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using Regulator_Ins within kubernetes types, where deepcopy-gen is used.
func (in *Regulator_Ins) DeepCopyInto(out *Regulator_Ins) {
	p := proto.Clone(in).(*Regulator_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Regulator_Ins. Required by controller-gen.
func (in *Regulator_Ins) DeepCopy() *Regulator_Ins {
	if in == nil {
		return nil
	}
	out := new(Regulator_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Regulator_Ins. Required by controller-gen.
func (in *Regulator_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRamp within kubernetes types, where deepcopy-gen is used.
func (in *LoadRamp) DeepCopyInto(out *LoadRamp) {
	p := proto.Clone(in).(*LoadRamp)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp. Required by controller-gen.
func (in *LoadRamp) DeepCopy() *LoadRamp {
	if in == nil {
		return nil
	}
	out := new(LoadRamp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp. Required by controller-gen.
func (in *LoadRamp) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRamp_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *LoadRamp_Parameters) DeepCopyInto(out *LoadRamp_Parameters) {
	p := proto.Clone(in).(*LoadRamp_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp_Parameters. Required by controller-gen.
func (in *LoadRamp_Parameters) DeepCopy() *LoadRamp_Parameters {
	if in == nil {
		return nil
	}
	out := new(LoadRamp_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp_Parameters. Required by controller-gen.
func (in *LoadRamp_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRamp_Parameters_Step within kubernetes types, where deepcopy-gen is used.
func (in *LoadRamp_Parameters_Step) DeepCopyInto(out *LoadRamp_Parameters_Step) {
	p := proto.Clone(in).(*LoadRamp_Parameters_Step)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp_Parameters_Step. Required by controller-gen.
func (in *LoadRamp_Parameters_Step) DeepCopy() *LoadRamp_Parameters_Step {
	if in == nil {
		return nil
	}
	out := new(LoadRamp_Parameters_Step)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp_Parameters_Step. Required by controller-gen.
func (in *LoadRamp_Parameters_Step) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRamp_Ins within kubernetes types, where deepcopy-gen is used.
func (in *LoadRamp_Ins) DeepCopyInto(out *LoadRamp_Ins) {
	p := proto.Clone(in).(*LoadRamp_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp_Ins. Required by controller-gen.
func (in *LoadRamp_Ins) DeepCopy() *LoadRamp_Ins {
	if in == nil {
		return nil
	}
	out := new(LoadRamp_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp_Ins. Required by controller-gen.
func (in *LoadRamp_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRamp_Outs within kubernetes types, where deepcopy-gen is used.
func (in *LoadRamp_Outs) DeepCopyInto(out *LoadRamp_Outs) {
	p := proto.Clone(in).(*LoadRamp_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp_Outs. Required by controller-gen.
func (in *LoadRamp_Outs) DeepCopy() *LoadRamp_Outs {
	if in == nil {
		return nil
	}
	out := new(LoadRamp_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRamp_Outs. Required by controller-gen.
func (in *LoadRamp_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRampSeries within kubernetes types, where deepcopy-gen is used.
func (in *LoadRampSeries) DeepCopyInto(out *LoadRampSeries) {
	p := proto.Clone(in).(*LoadRampSeries)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRampSeries. Required by controller-gen.
func (in *LoadRampSeries) DeepCopy() *LoadRampSeries {
	if in == nil {
		return nil
	}
	out := new(LoadRampSeries)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRampSeries. Required by controller-gen.
func (in *LoadRampSeries) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRampSeries_LoadRampInstance within kubernetes types, where deepcopy-gen is used.
func (in *LoadRampSeries_LoadRampInstance) DeepCopyInto(out *LoadRampSeries_LoadRampInstance) {
	p := proto.Clone(in).(*LoadRampSeries_LoadRampInstance)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRampSeries_LoadRampInstance. Required by controller-gen.
func (in *LoadRampSeries_LoadRampInstance) DeepCopy() *LoadRampSeries_LoadRampInstance {
	if in == nil {
		return nil
	}
	out := new(LoadRampSeries_LoadRampInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRampSeries_LoadRampInstance. Required by controller-gen.
func (in *LoadRampSeries_LoadRampInstance) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRampSeries_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *LoadRampSeries_Parameters) DeepCopyInto(out *LoadRampSeries_Parameters) {
	p := proto.Clone(in).(*LoadRampSeries_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRampSeries_Parameters. Required by controller-gen.
func (in *LoadRampSeries_Parameters) DeepCopy() *LoadRampSeries_Parameters {
	if in == nil {
		return nil
	}
	out := new(LoadRampSeries_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRampSeries_Parameters. Required by controller-gen.
func (in *LoadRampSeries_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using LoadRampSeries_Ins within kubernetes types, where deepcopy-gen is used.
func (in *LoadRampSeries_Ins) DeepCopyInto(out *LoadRampSeries_Ins) {
	p := proto.Clone(in).(*LoadRampSeries_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadRampSeries_Ins. Required by controller-gen.
func (in *LoadRampSeries_Ins) DeepCopy() *LoadRampSeries_Ins {
	if in == nil {
		return nil
	}
	out := new(LoadRampSeries_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new LoadRampSeries_Ins. Required by controller-gen.
func (in *LoadRampSeries_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
