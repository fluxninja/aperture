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

package watchdog

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HeapConfig) DeepCopyInto(out *HeapConfig) {
	*out = *in
	in.WatchdogPolicyType.DeepCopyInto(&out.WatchdogPolicyType)
	out.HeapLimit = in.HeapLimit
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HeapConfig.
func (in *HeapConfig) DeepCopy() *HeapConfig {
	if in == nil {
		return nil
	}
	out := new(HeapConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HeapLimit) DeepCopyInto(out *HeapLimit) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HeapLimit.
func (in *HeapLimit) DeepCopy() *HeapLimit {
	if in == nil {
		return nil
	}
	out := new(HeapLimit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyCommon) DeepCopyInto(out *PolicyCommon) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyCommon.
func (in *PolicyCommon) DeepCopy() *PolicyCommon {
	if in == nil {
		return nil
	}
	out := new(PolicyCommon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WatchdogConfig) DeepCopyInto(out *WatchdogConfig) {
	*out = *in
	in.Job.DeepCopyInto(&out.Job)
	in.CGroup.DeepCopyInto(&out.CGroup)
	in.System.DeepCopyInto(&out.System)
	in.Heap.DeepCopyInto(&out.Heap)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WatchdogConfig.
func (in *WatchdogConfig) DeepCopy() *WatchdogConfig {
	if in == nil {
		return nil
	}
	out := new(WatchdogConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WatchdogPolicyType) DeepCopyInto(out *WatchdogPolicyType) {
	*out = *in
	in.WatermarksPolicy.DeepCopyInto(&out.WatermarksPolicy)
	out.AdaptivePolicy = in.AdaptivePolicy
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WatchdogPolicyType.
func (in *WatchdogPolicyType) DeepCopy() *WatchdogPolicyType {
	if in == nil {
		return nil
	}
	out := new(WatchdogPolicyType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WatermarksPolicy) DeepCopyInto(out *WatermarksPolicy) {
	*out = *in
	if in.Watermarks != nil {
		in, out := &in.Watermarks, &out.Watermarks
		*out = make([]float64, len(*in))
		copy(*out, *in)
	}
	if in.thresholds != nil {
		in, out := &in.thresholds, &out.thresholds
		*out = make([]uint64, len(*in))
		copy(*out, *in)
	}
	out.PolicyCommon = in.PolicyCommon
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WatermarksPolicy.
func (in *WatermarksPolicy) DeepCopy() *WatermarksPolicy {
	if in == nil {
		return nil
	}
	out := new(WatermarksPolicy)
	in.DeepCopyInto(out)
	return out
}
