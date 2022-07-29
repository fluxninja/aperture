package uuid_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUUID(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UUID Suite")
}
