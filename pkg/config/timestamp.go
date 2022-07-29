package config

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fluxninja/aperture/pkg/log"
)

// Timestamp is encoded as a string message which represents a point in time independent of any time zone.
// It holds *timestamppb.Timestamp which is generated types for google/protobuf/timestamp.proto.
// swagger:strfmt string
type Timestamp struct {
	// swagger:ignore
	Timestamp *timestamppb.Timestamp
}

// UnmarshalJSON unmarshals and reads given bytes into a new Timestamp proto message.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	t.Timestamp = timestamppb.Now()
	err := protojson.Unmarshal(b, t.Timestamp)
	if err != nil {
		log.Error().Err(err).Bytes("b", b).Msg("Unable to unmarshal timestamp")
		return err
	}
	return nil
}

// MarshalJSON writes a Timestamp in JSON format.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t.Timestamp)
}

// String returns the string representation of a Timestamp.
func (t Timestamp) String() string {
	return t.Timestamp.String()
}
