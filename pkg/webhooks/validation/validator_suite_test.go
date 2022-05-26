package validation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"aperture.tech/aperture/pkg/utils"
)

func TestValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Suite")
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})
