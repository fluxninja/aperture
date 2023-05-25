package otelconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

func TestOtelcollector(t *testing.T) {
	log.SetGlobalLevel(log.ErrorLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "OTelcollector Config Suite")
}
