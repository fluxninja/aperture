package enrichmentprocessor

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/utils"
)

func TestEnrichmentprocessor(t *testing.T) {
	log.SetGlobalLevel(log.FatalLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enrichmentprocessor Suite")
}

func assertAttributesEqual(act, exp pcommon.Map) {
	Expect(act.Len()).To(Equal(exp.Len()))
	Expect(act.Sort()).To(Equal(exp.Sort()))
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})
