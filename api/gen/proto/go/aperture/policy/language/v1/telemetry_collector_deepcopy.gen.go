// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package languagev1

import (
	proto "google.golang.org/protobuf/proto"
)

// DeepCopyInto supports using TelemetryCollector within kubernetes types, where deepcopy-gen is used.
func (in *TelemetryCollector) DeepCopyInto(out *TelemetryCollector) {
	p := proto.Clone(in).(*TelemetryCollector)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TelemetryCollector. Required by controller-gen.
func (in *TelemetryCollector) DeepCopy() *TelemetryCollector {
	if in == nil {
		return nil
	}
	out := new(TelemetryCollector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new TelemetryCollector. Required by controller-gen.
func (in *TelemetryCollector) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using InfraMeter within kubernetes types, where deepcopy-gen is used.
func (in *InfraMeter) DeepCopyInto(out *InfraMeter) {
	p := proto.Clone(in).(*InfraMeter)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InfraMeter. Required by controller-gen.
func (in *InfraMeter) DeepCopy() *InfraMeter {
	if in == nil {
		return nil
	}
	out := new(InfraMeter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new InfraMeter. Required by controller-gen.
func (in *InfraMeter) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using InfraMeter_MetricsPipeline within kubernetes types, where deepcopy-gen is used.
func (in *InfraMeter_MetricsPipeline) DeepCopyInto(out *InfraMeter_MetricsPipeline) {
	p := proto.Clone(in).(*InfraMeter_MetricsPipeline)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InfraMeter_MetricsPipeline. Required by controller-gen.
func (in *InfraMeter_MetricsPipeline) DeepCopy() *InfraMeter_MetricsPipeline {
	if in == nil {
		return nil
	}
	out := new(InfraMeter_MetricsPipeline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new InfraMeter_MetricsPipeline. Required by controller-gen.
func (in *InfraMeter_MetricsPipeline) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
