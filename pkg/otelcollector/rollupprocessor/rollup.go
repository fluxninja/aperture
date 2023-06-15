package rollupprocessor

import (
	"fmt"

	"github.com/fluxninja/aperture/v2/pkg/otelcollector"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// Rollup represents single rollup operation. It describes Type of operation to be
// done on all `FromField`s from logs/traces. Result of operation is stored in
// `ToField`.
type Rollup struct {
	FromField      string
	ToField        string
	Type           RollupType
	TreatAsMissing []string
}

// RollupGroup represents a group of rollup operations of different types that
// all use the same FromField.
//
// By default all basic rollup types will be used (sum, max, min, sum of
// squares). WithDatasketch enables also the Datasketch rollup type.
//
// The name of ToField of Rollup will be inferred from the type of the rollup.
type RollupGroup struct {
	FromField      string
	WithDatasketch bool
	TreatAsMissing []string
}

var allRollupTypes = []RollupType{
	RollupSum,
	RollupDatasketch,
	RollupMax,
	RollupMin,
	RollupSumOfSquares,
}

// NewRollups creates individual rollups based on rollup groups.
func NewRollups(groups []RollupGroup) []*Rollup {
	var rollups []*Rollup
	fromFields := make(map[string]struct{})

	for _, group := range groups {
		for _, rollupType := range allRollupTypes {
			if rollupType == RollupDatasketch && !group.WithDatasketch {
				continue
			}

			fromFields[group.FromField] = struct{}{}
			rollups = append(rollups, &Rollup{
				FromField:      group.FromField,
				ToField:        AggregateField(group.FromField, rollupType),
				Type:           rollupType,
				TreatAsMissing: group.TreatAsMissing,
			})
		}
	}

	return rollups
}

// GetFromFieldValue returns value of `FromField` from attributes as float64.
func (rollup *Rollup) GetFromFieldValue(attributes pcommon.Map) (float64, bool) {
	return otelcollector.GetFloat64(attributes, rollup.FromField, rollup.TreatAsMissing)
}

// GetToFieldValue returns value of `ToField` from attributes as float64.
func (rollup *Rollup) GetToFieldValue(attributes pcommon.Map) (float64, bool) {
	return otelcollector.GetFloat64(attributes, rollup.ToField, rollup.TreatAsMissing)
}

// RollupType represents rollup type available in the processor.
type RollupType string

const (
	// RollupSum rolls up fields by adding them.
	RollupSum RollupType = "sum"
	// RollupMax rolls up fields by getting max value of them.
	RollupMax RollupType = "max"
	// RollupMin rolls up fields by getting min value of them.
	RollupMin RollupType = "min"
	// RollupSumOfSquares rolls up fields by summing squares of them.
	RollupSumOfSquares RollupType = "sumOfSquares"
	// RollupDatasketch rolls up fields by creating datasketch from them.
	RollupDatasketch RollupType = "datasketch"
)

// AggregateField returns the aggregate field name for the given field and rollup type.
func AggregateField(field string, rollupType RollupType) string {
	return fmt.Sprintf("%s_%s", field, rollupType)
}

const (
	// RollupCountKey is the key used to store the count of the rollup.
	RollupCountKey = "rollup_count"
)
