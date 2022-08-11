package sentry_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSentry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sentry Suite")
}
