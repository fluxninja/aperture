// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package annotations

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using StatusAnnotation within kubernetes types, where deepcopy-gen is used.
func (in *StatusAnnotation) DeepCopyInto(out *StatusAnnotation) {
	p := proto.Clone(in).(*StatusAnnotation)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatusAnnotation. Required by controller-gen.
func (in *StatusAnnotation) DeepCopy() *StatusAnnotation {
	if in == nil {
		return nil
	}
	out := new(StatusAnnotation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new StatusAnnotation. Required by controller-gen.
func (in *StatusAnnotation) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
