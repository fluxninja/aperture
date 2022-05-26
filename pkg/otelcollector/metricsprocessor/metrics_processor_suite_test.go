package metricsprocessor

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMetricsProcessor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrics Processor Suite")
}
