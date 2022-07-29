// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/common/watchdog/v1/watchdog.proto

package watchdogv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HeapResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit        uint64 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	HeapMarked   uint64 `protobuf:"varint,2,opt,name=heap_marked,json=heapMarked,proto3" json:"heap_marked,omitempty"`
	Threshold    uint64 `protobuf:"varint,3,opt,name=threshold,proto3" json:"threshold,omitempty"`
	CurrGogc     int32  `protobuf:"varint,4,opt,name=curr_gogc,json=currGogc,proto3" json:"curr_gogc,omitempty"`
	OriginalGogc int32  `protobuf:"varint,5,opt,name=original_gogc,json=originalGogc,proto3" json:"original_gogc,omitempty"`
	TotalAlloc   uint64 `protobuf:"varint,6,opt,name=total_alloc,json=totalAlloc,proto3" json:"total_alloc,omitempty"`
	Sys          uint64 `protobuf:"varint,7,opt,name=sys,proto3" json:"sys,omitempty"`
	Mallocs      uint64 `protobuf:"varint,8,opt,name=mallocs,proto3" json:"mallocs,omitempty"`
	Frees        uint64 `protobuf:"varint,9,opt,name=frees,proto3" json:"frees,omitempty"`
	HeapAlloc    uint64 `protobuf:"varint,10,opt,name=heap_alloc,json=heapAlloc,proto3" json:"heap_alloc,omitempty"`
	HeapSys      uint64 `protobuf:"varint,11,opt,name=heap_sys,json=heapSys,proto3" json:"heap_sys,omitempty"`
	HeapIdle     uint64 `protobuf:"varint,12,opt,name=heap_idle,json=heapIdle,proto3" json:"heap_idle,omitempty"`
	HeapInuse    uint64 `protobuf:"varint,13,opt,name=heap_inuse,json=heapInuse,proto3" json:"heap_inuse,omitempty"`
	HeapReleased uint64 `protobuf:"varint,14,opt,name=heap_released,json=heapReleased,proto3" json:"heap_released,omitempty"`
	HeapObjects  uint64 `protobuf:"varint,15,opt,name=heap_objects,json=heapObjects,proto3" json:"heap_objects,omitempty"`
	NextGc       uint64 `protobuf:"varint,16,opt,name=next_gc,json=nextGc,proto3" json:"next_gc,omitempty"`
	LastGc       uint64 `protobuf:"varint,17,opt,name=last_gc,json=lastGc,proto3" json:"last_gc,omitempty"`
	PauseTotalNs uint64 `protobuf:"varint,18,opt,name=pause_total_ns,json=pauseTotalNs,proto3" json:"pause_total_ns,omitempty"`
	NumGc        uint32 `protobuf:"varint,19,opt,name=num_gc,json=numGc,proto3" json:"num_gc,omitempty"`
	NumForcedGc  uint32 `protobuf:"varint,20,opt,name=num_forced_gc,json=numForcedGc,proto3" json:"num_forced_gc,omitempty"`
}

func (x *HeapResult) Reset() {
	*x = HeapResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_watchdog_v1_watchdog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeapResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeapResult) ProtoMessage() {}

func (x *HeapResult) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_watchdog_v1_watchdog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeapResult.ProtoReflect.Descriptor instead.
func (*HeapResult) Descriptor() ([]byte, []int) {
	return file_aperture_common_watchdog_v1_watchdog_proto_rawDescGZIP(), []int{0}
}

func (x *HeapResult) GetLimit() uint64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *HeapResult) GetHeapMarked() uint64 {
	if x != nil {
		return x.HeapMarked
	}
	return 0
}

func (x *HeapResult) GetThreshold() uint64 {
	if x != nil {
		return x.Threshold
	}
	return 0
}

func (x *HeapResult) GetCurrGogc() int32 {
	if x != nil {
		return x.CurrGogc
	}
	return 0
}

func (x *HeapResult) GetOriginalGogc() int32 {
	if x != nil {
		return x.OriginalGogc
	}
	return 0
}

func (x *HeapResult) GetTotalAlloc() uint64 {
	if x != nil {
		return x.TotalAlloc
	}
	return 0
}

func (x *HeapResult) GetSys() uint64 {
	if x != nil {
		return x.Sys
	}
	return 0
}

func (x *HeapResult) GetMallocs() uint64 {
	if x != nil {
		return x.Mallocs
	}
	return 0
}

func (x *HeapResult) GetFrees() uint64 {
	if x != nil {
		return x.Frees
	}
	return 0
}

func (x *HeapResult) GetHeapAlloc() uint64 {
	if x != nil {
		return x.HeapAlloc
	}
	return 0
}

func (x *HeapResult) GetHeapSys() uint64 {
	if x != nil {
		return x.HeapSys
	}
	return 0
}

func (x *HeapResult) GetHeapIdle() uint64 {
	if x != nil {
		return x.HeapIdle
	}
	return 0
}

func (x *HeapResult) GetHeapInuse() uint64 {
	if x != nil {
		return x.HeapInuse
	}
	return 0
}

func (x *HeapResult) GetHeapReleased() uint64 {
	if x != nil {
		return x.HeapReleased
	}
	return 0
}

func (x *HeapResult) GetHeapObjects() uint64 {
	if x != nil {
		return x.HeapObjects
	}
	return 0
}

func (x *HeapResult) GetNextGc() uint64 {
	if x != nil {
		return x.NextGc
	}
	return 0
}

func (x *HeapResult) GetLastGc() uint64 {
	if x != nil {
		return x.LastGc
	}
	return 0
}

func (x *HeapResult) GetPauseTotalNs() uint64 {
	if x != nil {
		return x.PauseTotalNs
	}
	return 0
}

func (x *HeapResult) GetNumGc() uint32 {
	if x != nil {
		return x.NumGc
	}
	return 0
}

func (x *HeapResult) GetNumForcedGc() uint32 {
	if x != nil {
		return x.NumForcedGc
	}
	return 0
}

type WatchdogResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total       uint64               `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Used        uint64               `protobuf:"varint,2,opt,name=used,proto3" json:"used,omitempty"`
	Threshold   uint64               `protobuf:"varint,3,opt,name=threshold,proto3" json:"threshold,omitempty"`
	ForceGcTook *durationpb.Duration `protobuf:"bytes,4,opt,name=force_gc_took,json=forceGcTook,proto3" json:"force_gc_took,omitempty"`
}

func (x *WatchdogResult) Reset() {
	*x = WatchdogResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_watchdog_v1_watchdog_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchdogResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchdogResult) ProtoMessage() {}

func (x *WatchdogResult) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_watchdog_v1_watchdog_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchdogResult.ProtoReflect.Descriptor instead.
func (*WatchdogResult) Descriptor() ([]byte, []int) {
	return file_aperture_common_watchdog_v1_watchdog_proto_rawDescGZIP(), []int{1}
}

func (x *WatchdogResult) GetTotal() uint64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *WatchdogResult) GetUsed() uint64 {
	if x != nil {
		return x.Used
	}
	return 0
}

func (x *WatchdogResult) GetThreshold() uint64 {
	if x != nil {
		return x.Threshold
	}
	return 0
}

func (x *WatchdogResult) GetForceGcTook() *durationpb.Duration {
	if x != nil {
		return x.ForceGcTook
	}
	return nil
}

var File_aperture_common_watchdog_v1_watchdog_proto protoreflect.FileDescriptor

var file_aperture_common_watchdog_v1_watchdog_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x77, 0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd7, 0x04, 0x0a, 0x0a, 0x48, 0x65,
	0x61, 0x70, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x68, 0x65, 0x61, 0x70, 0x5f, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x68, 0x65, 0x61, 0x70, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x1b, 0x0a,
	0x09, 0x63, 0x75, 0x72, 0x72, 0x5f, 0x67, 0x6f, 0x67, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x47, 0x6f, 0x67, 0x63, 0x12, 0x23, 0x0a, 0x0d, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x67, 0x6f, 0x67, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x47, 0x6f, 0x67, 0x63, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x63, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x41, 0x6c, 0x6c, 0x6f, 0x63,
	0x12, 0x10, 0x0a, 0x03, 0x73, 0x79, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x73,
	0x79, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x61, 0x6c, 0x6c, 0x6f, 0x63, 0x73, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x07, 0x6d, 0x61, 0x6c, 0x6c, 0x6f, 0x63, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x66, 0x72, 0x65, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x66, 0x72, 0x65,
	0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x68, 0x65, 0x61, 0x70, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x63,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x68, 0x65, 0x61, 0x70, 0x41, 0x6c, 0x6c, 0x6f,
	0x63, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x65, 0x61, 0x70, 0x5f, 0x73, 0x79, 0x73, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x07, 0x68, 0x65, 0x61, 0x70, 0x53, 0x79, 0x73, 0x12, 0x1b, 0x0a, 0x09,
	0x68, 0x65, 0x61, 0x70, 0x5f, 0x69, 0x64, 0x6c, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x08, 0x68, 0x65, 0x61, 0x70, 0x49, 0x64, 0x6c, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x68, 0x65, 0x61,
	0x70, 0x5f, 0x69, 0x6e, 0x75, 0x73, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x68,
	0x65, 0x61, 0x70, 0x49, 0x6e, 0x75, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x68, 0x65, 0x61, 0x70,
	0x5f, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0c, 0x68, 0x65, 0x61, 0x70, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x64, 0x12, 0x21, 0x0a,
	0x0c, 0x68, 0x65, 0x61, 0x70, 0x5f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x18, 0x0f, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x68, 0x65, 0x61, 0x70, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73,
	0x12, 0x17, 0x0a, 0x07, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x67, 0x63, 0x18, 0x10, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x6e, 0x65, 0x78, 0x74, 0x47, 0x63, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x67, 0x63, 0x18, 0x11, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6c, 0x61, 0x73, 0x74,
	0x47, 0x63, 0x12, 0x24, 0x0a, 0x0e, 0x70, 0x61, 0x75, 0x73, 0x65, 0x5f, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x5f, 0x6e, 0x73, 0x18, 0x12, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x70, 0x61, 0x75, 0x73,
	0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x4e, 0x73, 0x12, 0x15, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x5f,
	0x67, 0x63, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6e, 0x75, 0x6d, 0x47, 0x63, 0x12,
	0x22, 0x0a, 0x0d, 0x6e, 0x75, 0x6d, 0x5f, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x64, 0x5f, 0x67, 0x63,
	0x18, 0x14, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x46, 0x6f, 0x72, 0x63, 0x65,
	0x64, 0x47, 0x63, 0x22, 0x97, 0x01, 0x0a, 0x0e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04,
	0x75, 0x73, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x75, 0x73, 0x65, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x3d,
	0x0a, 0x0d, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x5f, 0x67, 0x63, 0x5f, 0x74, 0x6f, 0x6f, 0x6b, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x47, 0x63, 0x54, 0x6f, 0x6f, 0x6b, 0x42, 0x96, 0x02,
	0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x77, 0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x2e, 0x76,
	0x31, 0x42, 0x0d, 0x57, 0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x55, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x46,
	0x6c, 0x75, 0x78, 0x4e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x77, 0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x77,
	0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x43, 0x57, 0xaa,
	0x02, 0x1b, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1b,
	0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c,
	0x57, 0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x27, 0x41, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x57, 0x61,
	0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1e, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x3a, 0x3a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x57, 0x61, 0x74, 0x63, 0x68, 0x64,
	0x6f, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_common_watchdog_v1_watchdog_proto_rawDescOnce sync.Once
	file_aperture_common_watchdog_v1_watchdog_proto_rawDescData = file_aperture_common_watchdog_v1_watchdog_proto_rawDesc
)

func file_aperture_common_watchdog_v1_watchdog_proto_rawDescGZIP() []byte {
	file_aperture_common_watchdog_v1_watchdog_proto_rawDescOnce.Do(func() {
		file_aperture_common_watchdog_v1_watchdog_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_common_watchdog_v1_watchdog_proto_rawDescData)
	})
	return file_aperture_common_watchdog_v1_watchdog_proto_rawDescData
}

var file_aperture_common_watchdog_v1_watchdog_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_aperture_common_watchdog_v1_watchdog_proto_goTypes = []interface{}{
	(*HeapResult)(nil),          // 0: aperture.common.watchdog.v1.HeapResult
	(*WatchdogResult)(nil),      // 1: aperture.common.watchdog.v1.WatchdogResult
	(*durationpb.Duration)(nil), // 2: google.protobuf.Duration
}
var file_aperture_common_watchdog_v1_watchdog_proto_depIdxs = []int32{
	2, // 0: aperture.common.watchdog.v1.WatchdogResult.force_gc_took:type_name -> google.protobuf.Duration
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_aperture_common_watchdog_v1_watchdog_proto_init() }
func file_aperture_common_watchdog_v1_watchdog_proto_init() {
	if File_aperture_common_watchdog_v1_watchdog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_common_watchdog_v1_watchdog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeapResult); i {
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
		file_aperture_common_watchdog_v1_watchdog_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchdogResult); i {
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
			RawDescriptor: file_aperture_common_watchdog_v1_watchdog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_common_watchdog_v1_watchdog_proto_goTypes,
		DependencyIndexes: file_aperture_common_watchdog_v1_watchdog_proto_depIdxs,
		MessageInfos:      file_aperture_common_watchdog_v1_watchdog_proto_msgTypes,
	}.Build()
	File_aperture_common_watchdog_v1_watchdog_proto = out.File
	file_aperture_common_watchdog_v1_watchdog_proto_rawDesc = nil
	file_aperture_common_watchdog_v1_watchdog_proto_goTypes = nil
	file_aperture_common_watchdog_v1_watchdog_proto_depIdxs = nil
}
