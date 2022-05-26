package test

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {
	When("it is queried for version details", func() {
		It("returns version details", func() {
			resp, err := http.Get(fmt.Sprintf("http://%v/v1/info/version", addr))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var version map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&version)
			Expect(err).NotTo(HaveOccurred())
			Expect(version).To(HaveKey("build_host"))
			Expect(version).To(HaveKey("build_os"))
			Expect(version).To(HaveKey("build_time"))
			Expect(version).To(HaveKey("git_branch"))
			Expect(version).To(HaveKey("git_commit_hash"))
			Expect(version).To(HaveKey("version"))
		})
	})
})
