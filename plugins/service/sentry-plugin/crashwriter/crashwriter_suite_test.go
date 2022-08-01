package crashwriter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCrashWriter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CrashWriter Suite")
}
