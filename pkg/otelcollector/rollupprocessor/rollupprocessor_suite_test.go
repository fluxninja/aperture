package rollupprocessor_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/FluxNinja/aperture/pkg/utils"
)

func TestRollupprocessor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rollupprocessor Suite")
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})
