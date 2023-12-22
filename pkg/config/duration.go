package config

import (
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Duration is encoded as a string message which represents a signed span of time.
// It holds `*durationpb.Duration` which is generated type for google/protobuf/duration.proto.
// swagger:strfmt string
// +kubebuilder:validation:Type=string
type Duration struct {
	// swagger:ignore
	duration *durationpb.Duration
}

// MakeDuration returns a new Duration.
func MakeDuration(duration time.Duration) Duration {
	return Duration{
		duration: durationpb.New(duration),
	}
}

// AsDuration returns the Duration as a time.Duration.
func (d *Duration) AsDuration() time.Duration {
	if d.duration == nil {
		return 0
	}
	return d.duration.AsDuration()
}

// UnmarshalJSON unmarshals and reads given bytes into a new Duration proto message.
func (d *Duration) UnmarshalJSON(b []byte) error {
	d.duration = durationpb.New(0)
	// skip if bytes == "null"
	if string(b) == "null" {
		return nil
	}
	if err := protojson.Unmarshal(b, d.duration); err != nil {
		log.Error().Err(err).Bytes("bytes", b).Msg("Unable to unmarshal duration")
		return err
	}
	return nil
}

// MarshalJSON writes a Duration in JSON format.
func (d Duration) MarshalJSON() ([]byte, error) {
	if d.duration == nil {
		return []byte("null"), nil
	}
	return protojson.Marshal(d.duration)
}

// String returns the string representation of a Duration.
func (d Duration) String() string {
	if d.duration == nil {
		return ""
	}
	return d.duration.String()
}

// DeepCopyInto deepcopy function, copying the receiver, writing into out.
func (d *Duration) DeepCopyInto(out *Duration) {
	*out = *d
	if d.duration != nil {
		out.duration = proto.Clone(d.duration).(*durationpb.Duration)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Duration.
func (d *Duration) DeepCopy() *Duration {
	if d == nil {
		return nil
	}
	out := new(Duration)
	d.DeepCopyInto(out)
	return out
}
