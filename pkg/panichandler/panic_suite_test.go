package panichandler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPanic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Panic Suite")
}
