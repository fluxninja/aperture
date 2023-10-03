package testutils

import (
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

// CompareConfigMap compares two ConfigMaps by recursively sorting the keys in each map and comparing the sorted maps for given key in Data.
func CompareConfigMap(result, expected *corev1.ConfigMap) {
	for key := range result.Data {
		Expect(result.Data[key]).To(MatchYAML(expected.Data[key]))
	}

	// Compare the rest of the fields in the ConfigMap object that are not in the Data field
	result.Data = nil
	expected.Data = nil
	Expect(result).To(Equal(expected))
}
