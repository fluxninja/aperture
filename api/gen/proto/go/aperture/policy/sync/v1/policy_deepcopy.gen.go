// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package syncv1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using PolicyWrapper within kubernetes types, where deepcopy-gen is used.
func (in *PolicyWrapper) DeepCopyInto(out *PolicyWrapper) {
	p := proto.Clone(in).(*PolicyWrapper)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyWrapper. Required by controller-gen.
func (in *PolicyWrapper) DeepCopy() *PolicyWrapper {
	if in == nil {
		return nil
	}
	out := new(PolicyWrapper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PolicyWrapper. Required by controller-gen.
func (in *PolicyWrapper) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PolicyWrappers within kubernetes types, where deepcopy-gen is used.
func (in *PolicyWrappers) DeepCopyInto(out *PolicyWrappers) {
	p := proto.Clone(in).(*PolicyWrappers)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyWrappers. Required by controller-gen.
func (in *PolicyWrappers) DeepCopy() *PolicyWrappers {
	if in == nil {
		return nil
	}
	out := new(PolicyWrappers)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PolicyWrappers. Required by controller-gen.
func (in *PolicyWrappers) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
