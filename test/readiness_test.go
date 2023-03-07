package test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Readiness", func() {
	When("it is queried for readiness", func() {
		It("returns system readiness status", func() {
			resp, err := http.Get(fmt.Sprintf("http://%v/v1/status", addr))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
})
