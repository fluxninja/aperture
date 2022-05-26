package otelcollector_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOtelcollector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Otelcollector Suite")
}
