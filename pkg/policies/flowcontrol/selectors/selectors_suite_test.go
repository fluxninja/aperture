package selectors_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSelectors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Selectors Suite")
}
