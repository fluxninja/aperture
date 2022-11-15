package label

// FlowLabels is a map from flow labels to their values.
type FlowLabels map[string]FlowLabelValue

// FlowLabelValue is a value of a flow label with additional metadata.
type FlowLabelValue struct {
	Value     string
	Telemetry bool
}

// NewFromPlainMap returns flow labels from normal map[string]string. Telemetry flag is set to true for all flow labels.
func NewFromPlainMap(input map[string]string) FlowLabels {
	flowLabels := make(FlowLabels, len(input))
	for key, val := range input {
		flowLabels[key] = FlowLabelValue{
			Value:     val,
			Telemetry: true,
		}
	}
	return flowLabels
}

// ToPlainMap returns flow labels as normal map[string]string.
func (fl FlowLabels) ToPlainMap() map[string]string {
	plainMap := make(map[string]string, len(fl))
	for key, val := range fl {
		plainMap[key] = val.Value
	}
	return plainMap
}

// Merge combines two flow labels maps into one. Overwrites overlapping keys with values from src.
func Merge(dst, src FlowLabels) {
	for key, val := range src {
		dst[key] = val
	}
}
