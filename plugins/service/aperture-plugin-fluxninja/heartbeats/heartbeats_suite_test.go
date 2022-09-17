package heartbeats

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/utils"
)

func TestHeartbeats(t *testing.T) {
	// Disable logs for cleaner tests output
	log.SetGlobalLevel(log.FatalLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Heartbeats Suite")
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})
