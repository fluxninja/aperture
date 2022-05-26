package dataplane

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/utils"
)

type GinkgoTestReporter struct{}

func (g GinkgoTestReporter) Errorf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g GinkgoTestReporter) Fatalf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func TestPolicyConfig(t *testing.T) {
	log.SetGlobalLevel(log.WarnLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Dataplane Engine Suite")
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})
