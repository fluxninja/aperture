package rollupprocessor

import (
	"fmt"

	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for rollup processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	Rollups []Rollup `mapstructure:"rollups"`
}

// Rollup represents singe rollup operation. It describes Type of operation to be
// done on all `FromField`s from logs/traces. Result of operation is stored in
// `ToField`.
type Rollup struct {
	FromField string     `mapstructure:"from"`
	ToField   string     `mapstructure:"to"`
	Type      RollupType `mapstructure:"type"`
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

var _ config.Processor = (*Config)(nil)

// Validate checks if the processor configuration is valid.
func (cfg *Config) Validate() error {
	message := ""
	for i, rollup := range cfg.Rollups {
		if len(rollup.FromField) == 0 {
			message += fmt.Sprintf("Rollup %v 'from' not set", i)
		}
		if len(rollup.ToField) == 0 {
			message += fmt.Sprintf("Rollup %v 'to' not set", i)
		}
		if len(rollup.Type) == 0 {
			message += fmt.Sprintf("Rollup %v 'type' not set", i)
			continue
		}
		switch rollup.Type {
		case RollupSum, RollupDatasketch, RollupMax, RollupMin, RollupSumOfSquares:
		default:
			message += fmt.Sprintf("Rollup %v 'type' not valid. Expected one of %v, got %v",
				i, RollupTypes, rollup.Type)
		}
	}
	return nil
}
