package config

import (
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fluxninja/aperture/pkg/log"
)

// Time is encoded as a string message which represents a point in time independent of any time zone.
// It holds *timestamppb.Time which is generated types for google/protobuf/timestamp.proto.
// swagger:strfmt string
type Time struct {
	// swagger:ignore
	timestamp *timestamppb.Timestamp
}

// MakeTime returns a new Timestamp.
func MakeTime(timestamp time.Time) Time {
	return Time{
		timestamp: timestamppb.New(timestamp),
	}
}

// AsTime returns the Timestamp as a time.Time.
func (t *Time) AsTime() time.Time {
	if t.timestamp == nil {
		return time.Time{}
	}
	return t.timestamp.AsTime()
}

// UnmarshalJSON unmarshals and reads given bytes into a new Timestamp proto message.
func (t *Time) UnmarshalJSON(b []byte) error {
	t.timestamp = timestamppb.Now()
	err := protojson.Unmarshal(b, t.timestamp)
	if err != nil {
		log.Error().Err(err).Bytes("bytes", b).Msg("Unable to unmarshal timestamp")
		return err
	}
	return nil
}

// MarshalJSON writes a Timestamp in JSON format.
func (t Time) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t.timestamp)
}

// String returns the string representation of a Timestamp.
func (t Time) String() string {
	return t.timestamp.String()
}

// DeepCopyInto deepcopy function for Timestamp.
func (t *Time) DeepCopyInto(out *Time) {
	*out = *t
	if t.timestamp != nil {
		out.timestamp = proto.Clone(t.timestamp).(*timestamppb.Timestamp)
	}
}

// DeepCopy deepcopy function for Timestamp.
func (t *Time) DeepCopy() *Time {
	if t == nil {
		return nil
	}
	out := new(Time)
	t.DeepCopyInto(out)
	return out
}
