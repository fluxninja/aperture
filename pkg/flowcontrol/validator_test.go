package flowcontrol_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"aperture.tech/aperture/pkg/flowcontrol"
	"aperture.tech/aperture/pkg/webhooks/validation"
)

var _ = Describe("Validator", func() {
	cmFileValidator := &flowcontrol.CMFileValidator{}
	It("rejects non-flowcontrol configmap name", func() {
		ok := cmFileValidator.CheckCMName("foo")
		Expect(ok).To(BeFalse())
	})

	It("rejects empty policy file", func() {
		ok, _, err := cmFileValidator.ValidateFile(context.TODO(), "foo", []byte{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeFalse())
	})

	cmVatidator := validation.NewCMValidator()
	cmVatidator.RegisterCMFileValidator(cmFileValidator)

	validateExample := func(name string) {
		contents, err := os.ReadFile(filepath.Join(
			"..", // pkg
			"..", // aperture
			fmt.Sprintf("%s.yaml", name),
		))
		Expect(err).NotTo(HaveOccurred())

		var cm corev1.ConfigMap
		err = yaml.Unmarshal(contents, &cm)
		Expect(err).NotTo(HaveOccurred())

		ok, msg, err := cmVatidator.ValidateConfigMap(context.TODO(), cm)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg).To(BeEmpty())
		Expect(ok).To(BeTrue())
	}

	It("accepts example configmap for demoapp", func() {
		validateExample("latency-gradient-policy")
	})
})
