// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package privatev1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using PodScaleActuator within kubernetes types, where deepcopy-gen is used.
func (in *PodScaleActuator) DeepCopyInto(out *PodScaleActuator) {
	p := proto.Clone(in).(*PodScaleActuator)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodScaleActuator. Required by controller-gen.
func (in *PodScaleActuator) DeepCopy() *PodScaleActuator {
	if in == nil {
		return nil
	}
	out := new(PodScaleActuator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PodScaleActuator. Required by controller-gen.
func (in *PodScaleActuator) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PodScaleActuator_Ins within kubernetes types, where deepcopy-gen is used.
func (in *PodScaleActuator_Ins) DeepCopyInto(out *PodScaleActuator_Ins) {
	p := proto.Clone(in).(*PodScaleActuator_Ins)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodScaleActuator_Ins. Required by controller-gen.
func (in *PodScaleActuator_Ins) DeepCopy() *PodScaleActuator_Ins {
	if in == nil {
		return nil
	}
	out := new(PodScaleActuator_Ins)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PodScaleActuator_Ins. Required by controller-gen.
func (in *PodScaleActuator_Ins) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PodScaleReporter within kubernetes types, where deepcopy-gen is used.
func (in *PodScaleReporter) DeepCopyInto(out *PodScaleReporter) {
	p := proto.Clone(in).(*PodScaleReporter)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodScaleReporter. Required by controller-gen.
func (in *PodScaleReporter) DeepCopy() *PodScaleReporter {
	if in == nil {
		return nil
	}
	out := new(PodScaleReporter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PodScaleReporter. Required by controller-gen.
func (in *PodScaleReporter) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PodScaleReporter_Outs within kubernetes types, where deepcopy-gen is used.
func (in *PodScaleReporter_Outs) DeepCopyInto(out *PodScaleReporter_Outs) {
	p := proto.Clone(in).(*PodScaleReporter_Outs)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodScaleReporter_Outs. Required by controller-gen.
func (in *PodScaleReporter_Outs) DeepCopy() *PodScaleReporter_Outs {
	if in == nil {
		return nil
	}
	out := new(PodScaleReporter_Outs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PodScaleReporter_Outs. Required by controller-gen.
func (in *PodScaleReporter_Outs) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
