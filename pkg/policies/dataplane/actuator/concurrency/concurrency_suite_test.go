package concurrency_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/utils"
)

func TestActuationStrategy(t *testing.T) {
	log.SetGlobalLevel(log.WarnLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "ActuationStrategy Suite")
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})
