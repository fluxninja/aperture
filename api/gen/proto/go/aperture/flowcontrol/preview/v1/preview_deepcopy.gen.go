// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package previewv1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using PreviewRequest within kubernetes types, where deepcopy-gen is used.
func (in *PreviewRequest) DeepCopyInto(out *PreviewRequest) {
	p := proto.Clone(in).(*PreviewRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewRequest. Required by controller-gen.
func (in *PreviewRequest) DeepCopy() *PreviewRequest {
	if in == nil {
		return nil
	}
	out := new(PreviewRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PreviewRequest. Required by controller-gen.
func (in *PreviewRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PreviewFlowLabelsResponse within kubernetes types, where deepcopy-gen is used.
func (in *PreviewFlowLabelsResponse) DeepCopyInto(out *PreviewFlowLabelsResponse) {
	p := proto.Clone(in).(*PreviewFlowLabelsResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewFlowLabelsResponse. Required by controller-gen.
func (in *PreviewFlowLabelsResponse) DeepCopy() *PreviewFlowLabelsResponse {
	if in == nil {
		return nil
	}
	out := new(PreviewFlowLabelsResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PreviewFlowLabelsResponse. Required by controller-gen.
func (in *PreviewFlowLabelsResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PreviewFlowLabelsResponse_FlowLabels within kubernetes types, where deepcopy-gen is used.
func (in *PreviewFlowLabelsResponse_FlowLabels) DeepCopyInto(out *PreviewFlowLabelsResponse_FlowLabels) {
	p := proto.Clone(in).(*PreviewFlowLabelsResponse_FlowLabels)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewFlowLabelsResponse_FlowLabels. Required by controller-gen.
func (in *PreviewFlowLabelsResponse_FlowLabels) DeepCopy() *PreviewFlowLabelsResponse_FlowLabels {
	if in == nil {
		return nil
	}
	out := new(PreviewFlowLabelsResponse_FlowLabels)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PreviewFlowLabelsResponse_FlowLabels. Required by controller-gen.
func (in *PreviewFlowLabelsResponse_FlowLabels) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PreviewHTTPRequestsResponse within kubernetes types, where deepcopy-gen is used.
func (in *PreviewHTTPRequestsResponse) DeepCopyInto(out *PreviewHTTPRequestsResponse) {
	p := proto.Clone(in).(*PreviewHTTPRequestsResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewHTTPRequestsResponse. Required by controller-gen.
func (in *PreviewHTTPRequestsResponse) DeepCopy() *PreviewHTTPRequestsResponse {
	if in == nil {
		return nil
	}
	out := new(PreviewHTTPRequestsResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PreviewHTTPRequestsResponse. Required by controller-gen.
func (in *PreviewHTTPRequestsResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
