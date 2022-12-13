package alertsreceiver

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAlertsreceiver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alertsreceiver Suite")
}
