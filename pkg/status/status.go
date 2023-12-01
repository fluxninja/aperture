package status

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	monitoringv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/monitoring/v1"
	statusv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/status/v1"
)

// NewStatus creates a new instance of Status to be pushed into status registry. Use this function for creating status instead of by hand.
func NewStatus(d proto.Message, e error) *statusv1.Status {
	s := &statusv1.Status{
		Timestamp: timestamppb.Now(),
	}

	if d != nil {
		messageAny, err := anypb.New(d)
		if err != nil {
			return nil
		}
		s.Message = messageAny
	}

	if e != nil {
		s.Error = &statusv1.Status_Error{
			Message: e.Error(),
		}
	}

	return s
}

// ParseGroupStatus parses GroupStatus response and returns a map of status path to value.
func ParseGroupStatus(statusMap map[string]string, parent string, resp *statusv1.GroupStatus) (map[string]string, error) {
	if resp.Groups == nil && resp.GetStatus() != nil {
		if respErr := resp.Status.GetError(); respErr != nil {
			statusMap[fmt.Sprintf("Error for %s", parent)] = respErr.Message
		}
		if respMsg := resp.Status.GetMessage(); respMsg != nil {
			value := respMsg.String()

			if respMsg.MessageIs(new(wrapperspb.StringValue)) {
				stringVal := &wrapperspb.StringValue{}
				err := respMsg.UnmarshalTo(stringVal)
				if err != nil {
					return nil, fmt.Errorf("error unmarshalling string value for key %s: %s", parent, err)
				}
				value = stringVal.Value
			}
			if respMsg.MessageIs(new(wrapperspb.DoubleValue)) {
				doubleVal := &wrapperspb.DoubleValue{}
				err := respMsg.UnmarshalTo(doubleVal)
				if err != nil {
					return nil, fmt.Errorf("error unmarshalling double value for key %s: %s", parent, err)
				}
				value = fmt.Sprintf("%v", doubleVal.Value)
			}
			if respMsg.MessageIs(new(monitoringv1.SignalMetricsInfo)) {
				signalMetricsVal := &monitoringv1.SignalMetricsInfo{}
				err := respMsg.UnmarshalTo(signalMetricsVal)
				if err != nil {
					return nil, fmt.Errorf("error unmarshalling signal metrics value for key %s: %s", parent, err)
				}
				value = fmt.Sprintf("%v", signalMetricsVal)
			}
			if respMsg.MessageIs(new(emptypb.Empty)) {
				value = ""
			}

			statusMap[fmt.Sprintf("Status for %s", parent)] = value
		}
	}

	for k, v := range resp.Groups {
		newParent := fmt.Sprintf("%s/%s", parent, k)
		_, err := ParseGroupStatus(statusMap, newParent, v)
		if err != nil {
			return nil, err
		}
	}
	return statusMap, nil
}
