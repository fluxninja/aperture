package config

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/FluxNinja/aperture/pkg/log"
)

// Duration is encoded as a string message which represents a signed span of time.
// It holds *durationpb.Duration which is generated types for google/protobuf/duration.proto.
// swagger:strfmt string
type Duration struct {
	// swagger:ignore
	Duration *durationpb.Duration
}

// UnmarshalJSON unmarshals and reads given bytes into a new Duration proto message.
func (d *Duration) UnmarshalJSON(b []byte) error {
	d.Duration = durationpb.New(0)
	if err := protojson.Unmarshal(b, d.Duration); err != nil {
		log.Error().Err(err).Bytes("b", b).Msg("Unable to unmarshal duration")
		return err
	}
	return nil
}

// MarshalJSON writes a Duration in JSON format.
func (d Duration) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(d.Duration)
}

// String returns the string representation of a Duration.
func (d Duration) String() string {
	return d.Duration.String()
}
