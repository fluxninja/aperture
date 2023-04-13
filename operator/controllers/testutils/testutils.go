package testutils

import (
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

// CompareConfigMap compares two ConfigMaps by recursively sorting the keys in each map and comparing the sorted maps for given key in Data.
func CompareConfigMap(result, expected *corev1.ConfigMap) {
	for key := range result.Data {
		var obj1, obj2 map[string]interface{}
		err1 := yaml.Unmarshal([]byte(result.Data[key]), &obj1)
		err2 := yaml.Unmarshal([]byte(expected.Data[key]), &obj2)

		if err1 == nil && err2 == nil {
			Expect(obj1).To(Equal(obj2))
		} else {
			Expect(result.Data[key]).To(Equal(expected.Data[key]))
		}
	}

	// Compare the rest of the fields in the ConfigMap object that are not in the Data field
	result.Data = nil
	expected.Data = nil
	Expect(result).To(Equal(expected))
}
