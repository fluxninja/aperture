// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package cmdv1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using ListServicesRequest within kubernetes types, where deepcopy-gen is used.
func (in *ListServicesRequest) DeepCopyInto(out *ListServicesRequest) {
	p := proto.Clone(in).(*ListServicesRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListServicesRequest. Required by controller-gen.
func (in *ListServicesRequest) DeepCopy() *ListServicesRequest {
	if in == nil {
		return nil
	}
	out := new(ListServicesRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListServicesRequest. Required by controller-gen.
func (in *ListServicesRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListServicesAgentResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListServicesAgentResponse) DeepCopyInto(out *ListServicesAgentResponse) {
	p := proto.Clone(in).(*ListServicesAgentResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListServicesAgentResponse. Required by controller-gen.
func (in *ListServicesAgentResponse) DeepCopy() *ListServicesAgentResponse {
	if in == nil {
		return nil
	}
	out := new(ListServicesAgentResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListServicesAgentResponse. Required by controller-gen.
func (in *ListServicesAgentResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListServicesControllerResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListServicesControllerResponse) DeepCopyInto(out *ListServicesControllerResponse) {
	p := proto.Clone(in).(*ListServicesControllerResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListServicesControllerResponse. Required by controller-gen.
func (in *ListServicesControllerResponse) DeepCopy() *ListServicesControllerResponse {
	if in == nil {
		return nil
	}
	out := new(ListServicesControllerResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListServicesControllerResponse. Required by controller-gen.
func (in *ListServicesControllerResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListFlowControlPointsRequest within kubernetes types, where deepcopy-gen is used.
func (in *ListFlowControlPointsRequest) DeepCopyInto(out *ListFlowControlPointsRequest) {
	p := proto.Clone(in).(*ListFlowControlPointsRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListFlowControlPointsRequest. Required by controller-gen.
func (in *ListFlowControlPointsRequest) DeepCopy() *ListFlowControlPointsRequest {
	if in == nil {
		return nil
	}
	out := new(ListFlowControlPointsRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListFlowControlPointsRequest. Required by controller-gen.
func (in *ListFlowControlPointsRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListFlowControlPointsAgentResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListFlowControlPointsAgentResponse) DeepCopyInto(out *ListFlowControlPointsAgentResponse) {
	p := proto.Clone(in).(*ListFlowControlPointsAgentResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListFlowControlPointsAgentResponse. Required by controller-gen.
func (in *ListFlowControlPointsAgentResponse) DeepCopy() *ListFlowControlPointsAgentResponse {
	if in == nil {
		return nil
	}
	out := new(ListFlowControlPointsAgentResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListFlowControlPointsAgentResponse. Required by controller-gen.
func (in *ListFlowControlPointsAgentResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListFlowControlPointsControllerResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListFlowControlPointsControllerResponse) DeepCopyInto(out *ListFlowControlPointsControllerResponse) {
	p := proto.Clone(in).(*ListFlowControlPointsControllerResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListFlowControlPointsControllerResponse. Required by controller-gen.
func (in *ListFlowControlPointsControllerResponse) DeepCopy() *ListFlowControlPointsControllerResponse {
	if in == nil {
		return nil
	}
	out := new(ListFlowControlPointsControllerResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListFlowControlPointsControllerResponse. Required by controller-gen.
func (in *ListFlowControlPointsControllerResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListAutoScaleControlPointsRequest within kubernetes types, where deepcopy-gen is used.
func (in *ListAutoScaleControlPointsRequest) DeepCopyInto(out *ListAutoScaleControlPointsRequest) {
	p := proto.Clone(in).(*ListAutoScaleControlPointsRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListAutoScaleControlPointsRequest. Required by controller-gen.
func (in *ListAutoScaleControlPointsRequest) DeepCopy() *ListAutoScaleControlPointsRequest {
	if in == nil {
		return nil
	}
	out := new(ListAutoScaleControlPointsRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListAutoScaleControlPointsRequest. Required by controller-gen.
func (in *ListAutoScaleControlPointsRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListAutoScaleControlPointsAgentResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListAutoScaleControlPointsAgentResponse) DeepCopyInto(out *ListAutoScaleControlPointsAgentResponse) {
	p := proto.Clone(in).(*ListAutoScaleControlPointsAgentResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListAutoScaleControlPointsAgentResponse. Required by controller-gen.
func (in *ListAutoScaleControlPointsAgentResponse) DeepCopy() *ListAutoScaleControlPointsAgentResponse {
	if in == nil {
		return nil
	}
	out := new(ListAutoScaleControlPointsAgentResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListAutoScaleControlPointsAgentResponse. Required by controller-gen.
func (in *ListAutoScaleControlPointsAgentResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListAutoScaleControlPointsControllerResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListAutoScaleControlPointsControllerResponse) DeepCopyInto(out *ListAutoScaleControlPointsControllerResponse) {
	p := proto.Clone(in).(*ListAutoScaleControlPointsControllerResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListAutoScaleControlPointsControllerResponse. Required by controller-gen.
func (in *ListAutoScaleControlPointsControllerResponse) DeepCopy() *ListAutoScaleControlPointsControllerResponse {
	if in == nil {
		return nil
	}
	out := new(ListAutoScaleControlPointsControllerResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListAutoScaleControlPointsControllerResponse. Required by controller-gen.
func (in *ListAutoScaleControlPointsControllerResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListAgentsResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListAgentsResponse) DeepCopyInto(out *ListAgentsResponse) {
	p := proto.Clone(in).(*ListAgentsResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListAgentsResponse. Required by controller-gen.
func (in *ListAgentsResponse) DeepCopy() *ListAgentsResponse {
	if in == nil {
		return nil
	}
	out := new(ListAgentsResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListAgentsResponse. Required by controller-gen.
func (in *ListAgentsResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using GlobalFlowControlPoint within kubernetes types, where deepcopy-gen is used.
func (in *GlobalFlowControlPoint) DeepCopyInto(out *GlobalFlowControlPoint) {
	p := proto.Clone(in).(*GlobalFlowControlPoint)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalFlowControlPoint. Required by controller-gen.
func (in *GlobalFlowControlPoint) DeepCopy() *GlobalFlowControlPoint {
	if in == nil {
		return nil
	}
	out := new(GlobalFlowControlPoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new GlobalFlowControlPoint. Required by controller-gen.
func (in *GlobalFlowControlPoint) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using GlobalAutoScaleControlPoint within kubernetes types, where deepcopy-gen is used.
func (in *GlobalAutoScaleControlPoint) DeepCopyInto(out *GlobalAutoScaleControlPoint) {
	p := proto.Clone(in).(*GlobalAutoScaleControlPoint)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalAutoScaleControlPoint. Required by controller-gen.
func (in *GlobalAutoScaleControlPoint) DeepCopy() *GlobalAutoScaleControlPoint {
	if in == nil {
		return nil
	}
	out := new(GlobalAutoScaleControlPoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new GlobalAutoScaleControlPoint. Required by controller-gen.
func (in *GlobalAutoScaleControlPoint) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using GlobalService within kubernetes types, where deepcopy-gen is used.
func (in *GlobalService) DeepCopyInto(out *GlobalService) {
	p := proto.Clone(in).(*GlobalService)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalService. Required by controller-gen.
func (in *GlobalService) DeepCopy() *GlobalService {
	if in == nil {
		return nil
	}
	out := new(GlobalService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new GlobalService. Required by controller-gen.
func (in *GlobalService) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PreviewFlowLabelsRequest within kubernetes types, where deepcopy-gen is used.
func (in *PreviewFlowLabelsRequest) DeepCopyInto(out *PreviewFlowLabelsRequest) {
	p := proto.Clone(in).(*PreviewFlowLabelsRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewFlowLabelsRequest. Required by controller-gen.
func (in *PreviewFlowLabelsRequest) DeepCopy() *PreviewFlowLabelsRequest {
	if in == nil {
		return nil
	}
	out := new(PreviewFlowLabelsRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PreviewFlowLabelsRequest. Required by controller-gen.
func (in *PreviewFlowLabelsRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PreviewHTTPRequestsRequest within kubernetes types, where deepcopy-gen is used.
func (in *PreviewHTTPRequestsRequest) DeepCopyInto(out *PreviewHTTPRequestsRequest) {
	p := proto.Clone(in).(*PreviewHTTPRequestsRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewHTTPRequestsRequest. Required by controller-gen.
func (in *PreviewHTTPRequestsRequest) DeepCopy() *PreviewHTTPRequestsRequest {
	if in == nil {
		return nil
	}
	out := new(PreviewHTTPRequestsRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PreviewHTTPRequestsRequest. Required by controller-gen.
func (in *PreviewHTTPRequestsRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PreviewFlowLabelsControllerResponse within kubernetes types, where deepcopy-gen is used.
func (in *PreviewFlowLabelsControllerResponse) DeepCopyInto(out *PreviewFlowLabelsControllerResponse) {
	p := proto.Clone(in).(*PreviewFlowLabelsControllerResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewFlowLabelsControllerResponse. Required by controller-gen.
func (in *PreviewFlowLabelsControllerResponse) DeepCopy() *PreviewFlowLabelsControllerResponse {
	if in == nil {
		return nil
	}
	out := new(PreviewFlowLabelsControllerResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PreviewFlowLabelsControllerResponse. Required by controller-gen.
func (in *PreviewFlowLabelsControllerResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using PreviewHTTPRequestsControllerResponse within kubernetes types, where deepcopy-gen is used.
func (in *PreviewHTTPRequestsControllerResponse) DeepCopyInto(out *PreviewHTTPRequestsControllerResponse) {
	p := proto.Clone(in).(*PreviewHTTPRequestsControllerResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviewHTTPRequestsControllerResponse. Required by controller-gen.
func (in *PreviewHTTPRequestsControllerResponse) DeepCopy() *PreviewHTTPRequestsControllerResponse {
	if in == nil {
		return nil
	}
	out := new(PreviewHTTPRequestsControllerResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new PreviewHTTPRequestsControllerResponse. Required by controller-gen.
func (in *PreviewHTTPRequestsControllerResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
