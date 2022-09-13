package otelcollector

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// Rollup represents single rollup operation. It describes Type of operation to be
// done on all `FromField`s from logs/traces. Result of operation is stored in
// `ToField`.
type Rollup struct {
	FromField   string     `mapstructure:"from"`
	ToField     string     `mapstructure:"to"`
	Type        RollupType `mapstructure:"type"`
	TreatAsZero []string   `mapstructure:"treat_as_zero"`
}

// GetFromFieldValue returns value of `FromField` from attributes as float64.
func (rollup *Rollup) GetFromFieldValue(attributes pcommon.Map) (float64, bool) {
	return GetFloat64(attributes, rollup.FromField, rollup.TreatAsZero)
}

// GetToFieldValue returns value of `ToField` from attributes as float64.
func (rollup *Rollup) GetToFieldValue(attributes pcommon.Map) (float64, bool) {
	return GetFloat64(attributes, rollup.ToField, rollup.TreatAsZero)
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

// RollupTypes contains all available rollup types.
var RollupTypes = []RollupType{RollupSum, RollupDatasketch, RollupMax, RollupMin, RollupSumOfSquares}
