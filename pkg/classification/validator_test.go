package classification_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/FluxNinja/aperture/pkg/classification"
	"github.com/FluxNinja/aperture/pkg/webhooks/validation"
)

var _ = Describe("Validator", func() {
	cmFileValidator := &classification.CMFileValidator{}

	It("rejects non-classification configmap name", func() {
		ok := cmFileValidator.CheckCMName("foo")
		Expect(ok).To(BeFalse())
	})

	It("rejects empty ruleset file", func() {
		ok, _, err := cmFileValidator.ValidateFile(context.TODO(), "foo", []byte{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeFalse())
	})

	cmValidator := validation.NewCMValidator()
	cmValidator.RegisterCMFileValidator(cmFileValidator)

	validateExample := func(name string) {
		contents, err := os.ReadFile(filepath.Join(
			"..", // pkg
			"..", // aperture
			fmt.Sprintf("example-classification-cm-%s.yaml", name),
		))
		Expect(err).NotTo(HaveOccurred())

		var cm corev1.ConfigMap
		err = yaml.Unmarshal(contents, &cm)
		Expect(err).NotTo(HaveOccurred())

		ok, msg, err := cmValidator.ValidateConfigMap(context.TODO(), cm)
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
