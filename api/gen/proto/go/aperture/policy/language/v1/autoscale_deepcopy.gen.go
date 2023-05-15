// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package languagev1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using KubernetesObjectSelector within kubernetes types, where deepcopy-gen is used.
func (in *KubernetesObjectSelector) DeepCopyInto(out *KubernetesObjectSelector) {
	p := proto.Clone(in).(*KubernetesObjectSelector)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesObjectSelector. Required by controller-gen.
func (in *KubernetesObjectSelector) DeepCopy() *KubernetesObjectSelector {
	if in == nil {
		return nil
	}
	out := new(KubernetesObjectSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesObjectSelector. Required by controller-gen.
func (in *KubernetesObjectSelector) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AutoScale within kubernetes types, where deepcopy-gen is used.
func (in *AutoScale) DeepCopyInto(out *AutoScale) {
	p := proto.Clone(in).(*AutoScale)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoScale. Required by controller-gen.
func (in *AutoScale) DeepCopy() *AutoScale {
	if in == nil {
		return nil
	}
	out := new(AutoScale)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AutoScale. Required by controller-gen.
func (in *AutoScale) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PodScaler within kubernetes types, where deepcopy-gen is used.
func (in *PodScaler) DeepCopyInto(out *PodScaler) {
	p := proto.Clone(in).(*PodScaler)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodScaler. Required by controller-gen.
func (in *PodScaler) DeepCopy() *PodScaler {
	if in == nil {
		return nil
	}
	out := new(PodScaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PodScaler. Required by controller-gen.
func (in *PodScaler) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PodScaler_Ins within kubernetes types, where deepcopy-gen is used.
func (in *PodScaler_Ins) DeepCopyInto(out *PodScaler_Ins) {
	p := proto.Clone(in).(*PodScaler_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodScaler_Ins. Required by controller-gen.
func (in *PodScaler_Ins) DeepCopy() *PodScaler_Ins {
	if in == nil {
		return nil
	}
	out := new(PodScaler_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PodScaler_Ins. Required by controller-gen.
func (in *PodScaler_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PodScaler_Outs within kubernetes types, where deepcopy-gen is used.
func (in *PodScaler_Outs) DeepCopyInto(out *PodScaler_Outs) {
	p := proto.Clone(in).(*PodScaler_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodScaler_Outs. Required by controller-gen.
func (in *PodScaler_Outs) DeepCopy() *PodScaler_Outs {
	if in == nil {
		return nil
	}
	out := new(PodScaler_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PodScaler_Outs. Required by controller-gen.
func (in *PodScaler_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using IncreasingGradient within kubernetes types, where deepcopy-gen is used.
func (in *IncreasingGradient) DeepCopyInto(out *IncreasingGradient) {
	p := proto.Clone(in).(*IncreasingGradient)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IncreasingGradient. Required by controller-gen.
func (in *IncreasingGradient) DeepCopy() *IncreasingGradient {
	if in == nil {
		return nil
	}
	out := new(IncreasingGradient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new IncreasingGradient. Required by controller-gen.
func (in *IncreasingGradient) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using IncreasingGradient_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *IncreasingGradient_Parameters) DeepCopyInto(out *IncreasingGradient_Parameters) {
	p := proto.Clone(in).(*IncreasingGradient_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IncreasingGradient_Parameters. Required by controller-gen.
func (in *IncreasingGradient_Parameters) DeepCopy() *IncreasingGradient_Parameters {
	if in == nil {
		return nil
	}
	out := new(IncreasingGradient_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new IncreasingGradient_Parameters. Required by controller-gen.
func (in *IncreasingGradient_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using IncreasingGradient_Ins within kubernetes types, where deepcopy-gen is used.
func (in *IncreasingGradient_Ins) DeepCopyInto(out *IncreasingGradient_Ins) {
	p := proto.Clone(in).(*IncreasingGradient_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IncreasingGradient_Ins. Required by controller-gen.
func (in *IncreasingGradient_Ins) DeepCopy() *IncreasingGradient_Ins {
	if in == nil {
		return nil
	}
	out := new(IncreasingGradient_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new IncreasingGradient_Ins. Required by controller-gen.
func (in *IncreasingGradient_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DecreasingGradient within kubernetes types, where deepcopy-gen is used.
func (in *DecreasingGradient) DeepCopyInto(out *DecreasingGradient) {
	p := proto.Clone(in).(*DecreasingGradient)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DecreasingGradient. Required by controller-gen.
func (in *DecreasingGradient) DeepCopy() *DecreasingGradient {
	if in == nil {
		return nil
	}
	out := new(DecreasingGradient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DecreasingGradient. Required by controller-gen.
func (in *DecreasingGradient) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DecreasingGradient_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *DecreasingGradient_Parameters) DeepCopyInto(out *DecreasingGradient_Parameters) {
	p := proto.Clone(in).(*DecreasingGradient_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DecreasingGradient_Parameters. Required by controller-gen.
func (in *DecreasingGradient_Parameters) DeepCopy() *DecreasingGradient_Parameters {
	if in == nil {
		return nil
	}
	out := new(DecreasingGradient_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DecreasingGradient_Parameters. Required by controller-gen.
func (in *DecreasingGradient_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DecreasingGradient_Ins within kubernetes types, where deepcopy-gen is used.
func (in *DecreasingGradient_Ins) DeepCopyInto(out *DecreasingGradient_Ins) {
	p := proto.Clone(in).(*DecreasingGradient_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DecreasingGradient_Ins. Required by controller-gen.
func (in *DecreasingGradient_Ins) DeepCopy() *DecreasingGradient_Ins {
	if in == nil {
		return nil
	}
	out := new(DecreasingGradient_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DecreasingGradient_Ins. Required by controller-gen.
func (in *DecreasingGradient_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PeriodicDecrease within kubernetes types, where deepcopy-gen is used.
func (in *PeriodicDecrease) DeepCopyInto(out *PeriodicDecrease) {
	p := proto.Clone(in).(*PeriodicDecrease)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeriodicDecrease. Required by controller-gen.
func (in *PeriodicDecrease) DeepCopy() *PeriodicDecrease {
	if in == nil {
		return nil
	}
	out := new(PeriodicDecrease)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PeriodicDecrease. Required by controller-gen.
func (in *PeriodicDecrease) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PeriodicDecrease_Parameters within kubernetes types, where deepcopy-gen is used.
func (in *PeriodicDecrease_Parameters) DeepCopyInto(out *PeriodicDecrease_Parameters) {
	p := proto.Clone(in).(*PeriodicDecrease_Parameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeriodicDecrease_Parameters. Required by controller-gen.
func (in *PeriodicDecrease_Parameters) DeepCopy() *PeriodicDecrease_Parameters {
	if in == nil {
		return nil
	}
	out := new(PeriodicDecrease_Parameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PeriodicDecrease_Parameters. Required by controller-gen.
func (in *PeriodicDecrease_Parameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ScaleOutController within kubernetes types, where deepcopy-gen is used.
func (in *ScaleOutController) DeepCopyInto(out *ScaleOutController) {
	p := proto.Clone(in).(*ScaleOutController)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScaleOutController. Required by controller-gen.
func (in *ScaleOutController) DeepCopy() *ScaleOutController {
	if in == nil {
		return nil
	}
	out := new(ScaleOutController)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ScaleOutController. Required by controller-gen.
func (in *ScaleOutController) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ScaleOutController_Controller within kubernetes types, where deepcopy-gen is used.
func (in *ScaleOutController_Controller) DeepCopyInto(out *ScaleOutController_Controller) {
	p := proto.Clone(in).(*ScaleOutController_Controller)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScaleOutController_Controller. Required by controller-gen.
func (in *ScaleOutController_Controller) DeepCopy() *ScaleOutController_Controller {
	if in == nil {
		return nil
	}
	out := new(ScaleOutController_Controller)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ScaleOutController_Controller. Required by controller-gen.
func (in *ScaleOutController_Controller) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ScaleInController within kubernetes types, where deepcopy-gen is used.
func (in *ScaleInController) DeepCopyInto(out *ScaleInController) {
	p := proto.Clone(in).(*ScaleInController)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScaleInController. Required by controller-gen.
func (in *ScaleInController) DeepCopy() *ScaleInController {
	if in == nil {
		return nil
	}
	out := new(ScaleInController)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ScaleInController. Required by controller-gen.
func (in *ScaleInController) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ScaleInController_Controller within kubernetes types, where deepcopy-gen is used.
func (in *ScaleInController_Controller) DeepCopyInto(out *ScaleInController_Controller) {
	p := proto.Clone(in).(*ScaleInController_Controller)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScaleInController_Controller. Required by controller-gen.
func (in *ScaleInController_Controller) DeepCopy() *ScaleInController_Controller {
	if in == nil {
		return nil
	}
	out := new(ScaleInController_Controller)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ScaleInController_Controller. Required by controller-gen.
func (in *ScaleInController_Controller) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AutoScaler within kubernetes types, where deepcopy-gen is used.
func (in *AutoScaler) DeepCopyInto(out *AutoScaler) {
	p := proto.Clone(in).(*AutoScaler)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler. Required by controller-gen.
func (in *AutoScaler) DeepCopy() *AutoScaler {
	if in == nil {
		return nil
	}
	out := new(AutoScaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler. Required by controller-gen.
func (in *AutoScaler) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AutoScaler_ScalingParameters within kubernetes types, where deepcopy-gen is used.
func (in *AutoScaler_ScalingParameters) DeepCopyInto(out *AutoScaler_ScalingParameters) {
	p := proto.Clone(in).(*AutoScaler_ScalingParameters)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler_ScalingParameters. Required by controller-gen.
func (in *AutoScaler_ScalingParameters) DeepCopy() *AutoScaler_ScalingParameters {
	if in == nil {
		return nil
	}
	out := new(AutoScaler_ScalingParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler_ScalingParameters. Required by controller-gen.
func (in *AutoScaler_ScalingParameters) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AutoScaler_ScalingBackend within kubernetes types, where deepcopy-gen is used.
func (in *AutoScaler_ScalingBackend) DeepCopyInto(out *AutoScaler_ScalingBackend) {
	p := proto.Clone(in).(*AutoScaler_ScalingBackend)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler_ScalingBackend. Required by controller-gen.
func (in *AutoScaler_ScalingBackend) DeepCopy() *AutoScaler_ScalingBackend {
	if in == nil {
		return nil
	}
	out := new(AutoScaler_ScalingBackend)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler_ScalingBackend. Required by controller-gen.
func (in *AutoScaler_ScalingBackend) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AutoScaler_ScalingBackend_KubernetesReplicas within kubernetes types, where deepcopy-gen is used.
func (in *AutoScaler_ScalingBackend_KubernetesReplicas) DeepCopyInto(out *AutoScaler_ScalingBackend_KubernetesReplicas) {
	p := proto.Clone(in).(*AutoScaler_ScalingBackend_KubernetesReplicas)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler_ScalingBackend_KubernetesReplicas. Required by controller-gen.
func (in *AutoScaler_ScalingBackend_KubernetesReplicas) DeepCopy() *AutoScaler_ScalingBackend_KubernetesReplicas {
	if in == nil {
		return nil
	}
	out := new(AutoScaler_ScalingBackend_KubernetesReplicas)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler_ScalingBackend_KubernetesReplicas. Required by controller-gen.
func (in *AutoScaler_ScalingBackend_KubernetesReplicas) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using AutoScaler_ScalingBackend_KubernetesReplicas_Outs within kubernetes types, where deepcopy-gen is used.
func (in *AutoScaler_ScalingBackend_KubernetesReplicas_Outs) DeepCopyInto(out *AutoScaler_ScalingBackend_KubernetesReplicas_Outs) {
	p := proto.Clone(in).(*AutoScaler_ScalingBackend_KubernetesReplicas_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler_ScalingBackend_KubernetesReplicas_Outs. Required by controller-gen.
func (in *AutoScaler_ScalingBackend_KubernetesReplicas_Outs) DeepCopy() *AutoScaler_ScalingBackend_KubernetesReplicas_Outs {
	if in == nil {
		return nil
	}
	out := new(AutoScaler_ScalingBackend_KubernetesReplicas_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new AutoScaler_ScalingBackend_KubernetesReplicas_Outs. Required by controller-gen.
func (in *AutoScaler_ScalingBackend_KubernetesReplicas_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
