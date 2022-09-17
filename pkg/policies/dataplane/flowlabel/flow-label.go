package flowlabel

// FlowLabels is a map from flow labels to their values.
type FlowLabels map[string]FlowLabelValue

// FlowLabelValue is a value of a flow label with additional metadata.
type FlowLabelValue struct {
	Value     string
	Telemetry bool
}

// ToPlainMap returns flow labels as normal map[string]string.
func (fl FlowLabels) ToPlainMap() map[string]string {
	plainMap := make(map[string]string, len(fl))
	for key, val := range fl {
		plainMap[key] = val.Value
	}
	return plainMap
}

// CombineWith combines two flow labels maps into one. Overwrites overlapping key with values from other.
func (fl FlowLabels) CombineWith(other FlowLabels) FlowLabels {
	for k, v := range other {
		fl[k] = v
	}
	return fl
}
