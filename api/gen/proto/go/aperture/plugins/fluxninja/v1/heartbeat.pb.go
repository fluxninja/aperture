// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/plugins/fluxninja/v1/heartbeat.proto

package fluxninjav1

import (
	v15 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	v14 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
	v1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/info/v1"
	v11 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/peers/v1"
	v13 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	v12 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/status/v1"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReportRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VersionInfo                      *v1.VersionInfo                       `protobuf:"bytes,1,opt,name=version_info,json=versionInfo,proto3" json:"version_info,omitempty"`
	ProcessInfo                      *v1.ProcessInfo                       `protobuf:"bytes,2,opt,name=process_info,json=processInfo,proto3" json:"process_info,omitempty"`
	HostInfo                         *v1.HostInfo                          `protobuf:"bytes,3,opt,name=host_info,json=hostInfo,proto3" json:"host_info,omitempty"`
	AgentGroup                       string                                `protobuf:"bytes,4,opt,name=agent_group,json=agentGroup,proto3" json:"agent_group,omitempty"`
	ControllerInfo                   *ControllerInfo                       `protobuf:"bytes,5,opt,name=controller_info,json=controllerInfo,proto3" json:"controller_info,omitempty"`
	Peers                            *v11.Peers                            `protobuf:"bytes,6,opt,name=peers,proto3" json:"peers,omitempty"`
	ServicesList                     *ServicesList                         `protobuf:"bytes,8,opt,name=services_list,json=servicesList,proto3" json:"services_list,omitempty"`
	AllStatuses                      *v12.GroupStatus                      `protobuf:"bytes,9,opt,name=all_statuses,json=allStatuses,proto3" json:"all_statuses,omitempty"`
	Policies                         *v13.PolicyWrappers                   `protobuf:"bytes,10,opt,name=policies,proto3" json:"policies,omitempty"`
	FlowControlPoints                *v14.FlowControlPoints                `protobuf:"bytes,11,opt,name=flow_control_points,json=flowControlPoints,proto3" json:"flow_control_points,omitempty"`
	AutoscaleKubernetesControlPoints *v15.AutoscaleKubernetesControlPoints `protobuf:"bytes,12,opt,name=autoscale_kubernetes_control_points,json=autoscaleKubernetesControlPoints,proto3" json:"autoscale_kubernetes_control_points,omitempty"`
	InstallationMode                 string                                `protobuf:"bytes,13,opt,name=installation_mode,json=installationMode,proto3" json:"installation_mode,omitempty"`
}

func (x *ReportRequest) Reset() {
	*x = ReportRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportRequest) ProtoMessage() {}

func (x *ReportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportRequest.ProtoReflect.Descriptor instead.
func (*ReportRequest) Descriptor() ([]byte, []int) {
	return file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescGZIP(), []int{0}
}

func (x *ReportRequest) GetVersionInfo() *v1.VersionInfo {
	if x != nil {
		return x.VersionInfo
	}
	return nil
}

func (x *ReportRequest) GetProcessInfo() *v1.ProcessInfo {
	if x != nil {
		return x.ProcessInfo
	}
	return nil
}

func (x *ReportRequest) GetHostInfo() *v1.HostInfo {
	if x != nil {
		return x.HostInfo
	}
	return nil
}

func (x *ReportRequest) GetAgentGroup() string {
	if x != nil {
		return x.AgentGroup
	}
	return ""
}

func (x *ReportRequest) GetControllerInfo() *ControllerInfo {
	if x != nil {
		return x.ControllerInfo
	}
	return nil
}

func (x *ReportRequest) GetPeers() *v11.Peers {
	if x != nil {
		return x.Peers
	}
	return nil
}

func (x *ReportRequest) GetServicesList() *ServicesList {
	if x != nil {
		return x.ServicesList
	}
	return nil
}

func (x *ReportRequest) GetAllStatuses() *v12.GroupStatus {
	if x != nil {
		return x.AllStatuses
	}
	return nil
}

func (x *ReportRequest) GetPolicies() *v13.PolicyWrappers {
	if x != nil {
		return x.Policies
	}
	return nil
}

func (x *ReportRequest) GetFlowControlPoints() *v14.FlowControlPoints {
	if x != nil {
		return x.FlowControlPoints
	}
	return nil
}

func (x *ReportRequest) GetAutoscaleKubernetesControlPoints() *v15.AutoscaleKubernetesControlPoints {
	if x != nil {
		return x.AutoscaleKubernetesControlPoints
	}
	return nil
}

func (x *ReportRequest) GetInstallationMode() string {
	if x != nil {
		return x.InstallationMode
	}
	return ""
}

// ReportResponse is empty for now.
type ReportResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReportResponse) Reset() {
	*x = ReportResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportResponse) ProtoMessage() {}

func (x *ReportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportResponse.ProtoReflect.Descriptor instead.
func (*ReportResponse) Descriptor() ([]byte, []int) {
	return file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescGZIP(), []int{1}
}

type ControllerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ControllerInfo) Reset() {
	*x = ControllerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ControllerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ControllerInfo) ProtoMessage() {}

func (x *ControllerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ControllerInfo.ProtoReflect.Descriptor instead.
func (*ControllerInfo) Descriptor() ([]byte, []int) {
	return file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescGZIP(), []int{2}
}

func (x *ControllerInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ServicesList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Services            []*Service            `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	OverlappingServices []*OverlappingService `protobuf:"bytes,2,rep,name=overlapping_services,json=overlappingServices,proto3" json:"overlapping_services,omitempty"`
}

func (x *ServicesList) Reset() {
	*x = ServicesList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServicesList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServicesList) ProtoMessage() {}

func (x *ServicesList) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServicesList.ProtoReflect.Descriptor instead.
func (*ServicesList) Descriptor() ([]byte, []int) {
	return file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescGZIP(), []int{3}
}

func (x *ServicesList) GetServices() []*Service {
	if x != nil {
		return x.Services
	}
	return nil
}

func (x *ServicesList) GetOverlappingServices() []*OverlappingService {
	if x != nil {
		return x.OverlappingServices
	}
	return nil
}

// Service contains information about single service discovered in agent group by a
// particular agent.
type Service struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	EntitiesCount int32  `protobuf:"varint,2,opt,name=entities_count,json=entitiesCount,proto3" json:"entities_count,omitempty"`
}

func (x *Service) Reset() {
	*x = Service{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service.ProtoReflect.Descriptor instead.
func (*Service) Descriptor() ([]byte, []int) {
	return file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescGZIP(), []int{4}
}

func (x *Service) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Service) GetEntitiesCount() int32 {
	if x != nil {
		return x.EntitiesCount
	}
	return 0
}

// OverlappingService contains info about a service that overlaps with another one.
type OverlappingService struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service1      string `protobuf:"bytes,1,opt,name=service1,proto3" json:"service1,omitempty"`
	Service2      string `protobuf:"bytes,2,opt,name=service2,proto3" json:"service2,omitempty"`
	EntitiesCount int32  `protobuf:"varint,3,opt,name=entities_count,json=entitiesCount,proto3" json:"entities_count,omitempty"`
}

func (x *OverlappingService) Reset() {
	*x = OverlappingService{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OverlappingService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OverlappingService) ProtoMessage() {}

func (x *OverlappingService) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OverlappingService.ProtoReflect.Descriptor instead.
func (*OverlappingService) Descriptor() ([]byte, []int) {
	return file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescGZIP(), []int{5}
}

func (x *OverlappingService) GetService1() string {
	if x != nil {
		return x.Service1
	}
	return ""
}

func (x *OverlappingService) GetService2() string {
	if x != nil {
		return x.Service2
	}
	return ""
}

func (x *OverlappingService) GetEntitiesCount() int32 {
	if x != nil {
		return x.EntitiesCount
	}
	return 0
}

var File_aperture_plugins_fluxninja_v1_heartbeat_proto protoreflect.FileDescriptor

var file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x73, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x76, 0x31, 0x2f,
	0x68, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x1d, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x73, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x76, 0x31, 0x1a, 0x42,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61,
	0x6c, 0x65, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2f, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x39, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x66, 0x6c, 0x6f,
	0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x76, 0x31, 0x2f,
	0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x65,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89, 0x07, 0x0a, 0x0d,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x40, 0x0a,
	0x0c, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x69,
	0x6e, 0x66, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x0b, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x40, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x37, 0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x56, 0x0a, 0x0f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a,
	0x61, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x2e, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x65,
	0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x73, 0x52, 0x05, 0x70, 0x65,
	0x65, 0x72, 0x73, 0x12, 0x50, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x5f,
	0x6c, 0x69, 0x73, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2e, 0x66, 0x6c,
	0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x42, 0x0a, 0x0c, 0x61, 0x6c, 0x6c, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0b, 0x61, 0x6c,
	0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x12, 0x43, 0x0a, 0x08, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x69, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x57, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x73, 0x52, 0x08, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x12, 0x68,
	0x0a, 0x13, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50,
	0x6f, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x11, 0x66, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x9f, 0x01, 0x0a, 0x23, 0x61, 0x75, 0x74,
	0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x5f, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65,
	0x73, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x50, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x2e, 0x6b, 0x75, 0x62, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c,
	0x65, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x20, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63,
	0x61, 0x6c, 0x65, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x69, 0x6e,
	0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x4d, 0x6f, 0x64, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x20, 0x0a, 0x0e, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xb8, 0x01, 0x0a, 0x0c,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x42, 0x0a, 0x08,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26,
	0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x73, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x12, 0x64, 0x0a, 0x14, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31,
	0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x73, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4f,
	0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x52, 0x13, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x22, 0x44, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x73, 0x0a, 0x12,
	0x4f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x31, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x31, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x32, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0d, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x32, 0xa3, 0x01, 0x0a, 0x10, 0x46, 0x6c, 0x75, 0x78, 0x4e, 0x69, 0x6e, 0x6a, 0x61, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x8e, 0x01, 0x0a, 0x06, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x73, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x73, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x27,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21, 0x3a, 0x01, 0x2a, 0x22, 0x1c, 0x2f, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x73, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x76, 0x31,
	0x2f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x32, 0xa2, 0x01, 0x0a, 0x15, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x88, 0x01, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x6c, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x2d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x73, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x2c,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x26, 0x12, 0x24, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73,
	0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x42, 0xb8, 0x02, 0x0a,
	0x35, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69,
	0x6e, 0x6a, 0x61, 0x2e, 0x76, 0x31, 0x42, 0x0e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61,
	0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x58, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69,
	0x6e, 0x6a, 0x61, 0x2f, 0x76, 0x31, 0x3b, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61,
	0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x46, 0xaa, 0x02, 0x1d, 0x41, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2e, 0x46, 0x6c, 0x75, 0x78,
	0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1d, 0x41, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x5c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x5c, 0x46, 0x6c, 0x75, 0x78,
	0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x29, 0x41, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x5c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x5c, 0x46, 0x6c, 0x75, 0x78,
	0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x20, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a,
	0x3a, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x3a, 0x3a, 0x46, 0x6c, 0x75, 0x78, 0x6e, 0x69,
	0x6e, 0x6a, 0x61, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescOnce sync.Once
	file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescData = file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDesc
)

func file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescGZIP() []byte {
	file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescOnce.Do(func() {
		file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescData)
	})
	return file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDescData
}

var file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_aperture_plugins_fluxninja_v1_heartbeat_proto_goTypes = []interface{}{
	(*ReportRequest)(nil),                        // 0: aperture.plugins.fluxninja.v1.ReportRequest
	(*ReportResponse)(nil),                       // 1: aperture.plugins.fluxninja.v1.ReportResponse
	(*ControllerInfo)(nil),                       // 2: aperture.plugins.fluxninja.v1.ControllerInfo
	(*ServicesList)(nil),                         // 3: aperture.plugins.fluxninja.v1.ServicesList
	(*Service)(nil),                              // 4: aperture.plugins.fluxninja.v1.Service
	(*OverlappingService)(nil),                   // 5: aperture.plugins.fluxninja.v1.OverlappingService
	(*v1.VersionInfo)(nil),                       // 6: aperture.info.v1.VersionInfo
	(*v1.ProcessInfo)(nil),                       // 7: aperture.info.v1.ProcessInfo
	(*v1.HostInfo)(nil),                          // 8: aperture.info.v1.HostInfo
	(*v11.Peers)(nil),                            // 9: aperture.peers.v1.Peers
	(*v12.GroupStatus)(nil),                      // 10: aperture.status.v1.GroupStatus
	(*v13.PolicyWrappers)(nil),                   // 11: aperture.policy.sync.v1.PolicyWrappers
	(*v14.FlowControlPoints)(nil),                // 12: aperture.flowcontrol.controlpoints.v1.FlowControlPoints
	(*v15.AutoscaleKubernetesControlPoints)(nil), // 13: aperture.autoscale.kubernetes.controlpoints.v1.AutoscaleKubernetesControlPoints
	(*emptypb.Empty)(nil),                        // 14: google.protobuf.Empty
}
var file_aperture_plugins_fluxninja_v1_heartbeat_proto_depIdxs = []int32{
	6,  // 0: aperture.plugins.fluxninja.v1.ReportRequest.version_info:type_name -> aperture.info.v1.VersionInfo
	7,  // 1: aperture.plugins.fluxninja.v1.ReportRequest.process_info:type_name -> aperture.info.v1.ProcessInfo
	8,  // 2: aperture.plugins.fluxninja.v1.ReportRequest.host_info:type_name -> aperture.info.v1.HostInfo
	2,  // 3: aperture.plugins.fluxninja.v1.ReportRequest.controller_info:type_name -> aperture.plugins.fluxninja.v1.ControllerInfo
	9,  // 4: aperture.plugins.fluxninja.v1.ReportRequest.peers:type_name -> aperture.peers.v1.Peers
	3,  // 5: aperture.plugins.fluxninja.v1.ReportRequest.services_list:type_name -> aperture.plugins.fluxninja.v1.ServicesList
	10, // 6: aperture.plugins.fluxninja.v1.ReportRequest.all_statuses:type_name -> aperture.status.v1.GroupStatus
	11, // 7: aperture.plugins.fluxninja.v1.ReportRequest.policies:type_name -> aperture.policy.sync.v1.PolicyWrappers
	12, // 8: aperture.plugins.fluxninja.v1.ReportRequest.flow_control_points:type_name -> aperture.flowcontrol.controlpoints.v1.FlowControlPoints
	13, // 9: aperture.plugins.fluxninja.v1.ReportRequest.autoscale_kubernetes_control_points:type_name -> aperture.autoscale.kubernetes.controlpoints.v1.AutoscaleKubernetesControlPoints
	4,  // 10: aperture.plugins.fluxninja.v1.ServicesList.services:type_name -> aperture.plugins.fluxninja.v1.Service
	5,  // 11: aperture.plugins.fluxninja.v1.ServicesList.overlapping_services:type_name -> aperture.plugins.fluxninja.v1.OverlappingService
	0,  // 12: aperture.plugins.fluxninja.v1.FluxNinjaService.Report:input_type -> aperture.plugins.fluxninja.v1.ReportRequest
	14, // 13: aperture.plugins.fluxninja.v1.ControllerInfoService.GetControllerInfo:input_type -> google.protobuf.Empty
	1,  // 14: aperture.plugins.fluxninja.v1.FluxNinjaService.Report:output_type -> aperture.plugins.fluxninja.v1.ReportResponse
	2,  // 15: aperture.plugins.fluxninja.v1.ControllerInfoService.GetControllerInfo:output_type -> aperture.plugins.fluxninja.v1.ControllerInfo
	14, // [14:16] is the sub-list for method output_type
	12, // [12:14] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_aperture_plugins_fluxninja_v1_heartbeat_proto_init() }
func file_aperture_plugins_fluxninja_v1_heartbeat_proto_init() {
	if File_aperture_plugins_fluxninja_v1_heartbeat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ControllerInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServicesList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Service); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OverlappingService); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_aperture_plugins_fluxninja_v1_heartbeat_proto_goTypes,
		DependencyIndexes: file_aperture_plugins_fluxninja_v1_heartbeat_proto_depIdxs,
		MessageInfos:      file_aperture_plugins_fluxninja_v1_heartbeat_proto_msgTypes,
	}.Build()
	File_aperture_plugins_fluxninja_v1_heartbeat_proto = out.File
	file_aperture_plugins_fluxninja_v1_heartbeat_proto_rawDesc = nil
	file_aperture_plugins_fluxninja_v1_heartbeat_proto_goTypes = nil
	file_aperture_plugins_fluxninja_v1_heartbeat_proto_depIdxs = nil
}
