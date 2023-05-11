package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
)

var _ = Describe("Readiness", func() {
	When("it is queried for readiness", func() {
		It("returns system readiness status", func() {
			Eventually(func() bool {
				resp, err := http.Get(fmt.Sprintf("http://%v/v1/status/system/readiness", addr))
				if err != nil {
					return false
				}

				// get response body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return false
				}

				// unmarshal response body
				groupStatus := &statusv1.GroupStatus{}
				err = groupStatus.UnmarshalJSON(body)
				if err != nil {
					return false
				}

				// check if overall group status has not error
				if groupStatus.Status.Error != nil {
					return false
				}

				// check if all component status has no error
				for _, group := range groupStatus.Groups {
					if group.Status.Error != nil {
						return false
					}
				}

				// finally check if response status code is 200
				if resp.StatusCode != http.StatusOK {
					return false
				}

				return true
			}).MustPassRepeatedly(2).WithTimeout(15 * time.Second).Should(BeTrue())
		})
	})
})
