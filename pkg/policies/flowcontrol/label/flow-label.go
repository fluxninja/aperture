package label

import (
	"errors"
	"sort"
	"strconv"

	"github.com/fluxninja/aperture/v2/pkg/labels"
	"golang.org/x/exp/maps"
)

// ErrLabelNotFound is returned when a flow label is not found.
var ErrLabelNotFound = errors.New("label not found")

// FlowLabels is a map from flow labels to their values.
// Each value also stores some metadata.
//
// FlowLabels implement labels.Labels to avoid conversion to
// plain map.
type FlowLabels map[string]FlowLabelValue

// FlowLabelValue is a value of a flow label with additional metadata.
type FlowLabelValue struct {
	Value     string
	Telemetry bool
}

// Get implements labels.Labels.
func (fl FlowLabels) Get(key string) (string, bool) {
	label, exists := fl[key]
	if !exists {
		return "", false
	}
	return label.Value, true
}

// SortedKeys implements labels.Labels.
func (fl FlowLabels) SortedKeys() []string {
	keys := maps.Keys(fl)
	sort.Strings(keys)
	return keys
}

// Copy implements labels.Labels.
func (fl FlowLabels) Copy() labels.PlainMap {
	plainMap := make(map[string]string, len(fl))
	for key, val := range fl {
		plainMap[key] = val.Value
	}
	return plainMap
}

// TelemetryLabels returns only labels that have Telemetry flag set.
func (fl FlowLabels) TelemetryLabels() map[string]string {
	telemetryLabels := make(map[string]string, len(fl))
	for key, val := range fl {
		if val.Telemetry {
			telemetryLabels[key] = val.Value
		}
	}
	return telemetryLabels
}

// Merge combines two flow labels maps into one. Overwrites overlapping keys with values from src.
func Merge(dst, src FlowLabels) {
	for key, val := range src {
		dst[key] = val
	}
}

// GetUint64 returns uint64 value of a flow label.
func (fl FlowLabels) GetUint64(key string) (uint64, error) {
	val, ok := fl[key]
	if !ok {
		return 0, ErrLabelNotFound
	}
	return strconv.ParseUint(val.Value, 10, 64)
}

// GetInt64 returns the int64 value of a flow label with the given key.
// If the key is not found in the FlowLabels, it returns an error.
func (fl FlowLabels) GetInt64(key string) (int64, error) {
	val, ok := fl[key]
	if !ok {
		return 0, ErrLabelNotFound
	}
	return strconv.ParseInt(val.Value, 10, 64)
}

// GetString returns the string value of a flow label with the given key.
// If the key is not found in the FlowLabels, it returns an error.
func (fl FlowLabels) GetString(key string) (string, error) {
	val, ok := fl[key]
	if !ok {
		return "", ErrLabelNotFound
	}
	return val.Value, nil
}

// GetFloat64 returns the float64 value of a flow label with the given key.
// If the key is not found in the FlowLabels, it returns an error.
func (fl FlowLabels) GetFloat64(key string) (float64, error) {
	val, ok := fl[key]
	if !ok {
		return 0, ErrLabelNotFound
	}
	return strconv.ParseFloat(val.Value, 64)
}

// GetBool returns the bool value of a flow label with the given key.
// If the key is not found in the FlowLabels, it returns an error.
func (fl FlowLabels) GetBool(key string) (bool, error) {
	val, ok := fl[key]
	if !ok {
		return false, ErrLabelNotFound
	}
	return strconv.ParseBool(val.Value)
}
