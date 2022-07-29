package validation_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/FluxNinja/aperture/pkg/webhooks/validation"
)

var _ = Describe("Validator", func() {
	testedVatidator := validation.NewCMValidator()
	It("accepts empty configmap", func() {
		ok, _, err := testedVatidator.ValidateConfigMap(context.TODO(), corev1.ConfigMap{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	validateExample := func(name string) {
		contents, err := os.ReadFile(filepath.Join(
			"../../..", // aperture
			fmt.Sprintf("example-classification-cm-%s.yaml", name),
		))
		Expect(err).NotTo(HaveOccurred())

		var cm corev1.ConfigMap
		err = yaml.Unmarshal(contents, &cm)
		Expect(err).NotTo(HaveOccurred())

		ok, msg, err := testedVatidator.ValidateConfigMap(context.TODO(), cm)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg).To(BeEmpty())
		Expect(ok).To(BeTrue())
	}

	It("accepts example configmap for demoapp", func() {
		validateExample("demoapp")
	})

	It("accepts example configmap for bookinfo", func() {
		validateExample("bookinfo")
	})
})
