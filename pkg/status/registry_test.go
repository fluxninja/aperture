package status

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
)

var _ = Describe("Status Registry", func() {
	var reg *Registry
	delim := ""

	BeforeEach(func() {
		reg = NewRegistry(delim)
		Expect(reg).NotTo(BeNil())

		key1_status := NewStatus(nil, errors.New("key1"))
		Expect(key1_status.GetError().Message).To(Equal("key1"))

		err := reg.Push("key1", key1_status)
		Expect(err).NotTo(HaveOccurred())
		err = reg.Push("key2", NewStatus(nil, errors.New("key2")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.Push("job1.key1", NewStatus(nil, errors.New("job1.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.Push("job1.key2", NewStatus(nil, errors.New("job1.key2")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.Push("job1.subjob1.key1", NewStatus(nil, errors.New("job1.subjob1.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.Push("job1.subjob1.key2", NewStatus(nil, errors.New("job1.subjob1.key2")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.Push("job1.subjob2.key1", NewStatus(nil, errors.New("job1.subjob2.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.Push("job1.subjob2.key2", NewStatus(nil, errors.New("job1.subjob2.key2")))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return correct delimiter", func() {
		Expect(reg.Delim()).To(Equal(defaultDelim))
	})

	It("should return all the keys pushed to the result map", func() {
		expected := []string{"job1", "job1.key1", "job1.key2", "job1.subjob1", "job1.subjob1.key1", "job1.subjob1.key2", "job1.subjob2", "job1.subjob2.key1", "job1.subjob2.key2", "key1", "key2"}
		Expect(reg.Keys()).To(Equal(expected))
		Expect(reg.Keys()).To(HaveLen(11))
	})

	It("pushed status exist in the result map", func() {
		Expect(reg.Exists("key1")).To(BeTrue())
		Expect(reg.Exists("job1.key2")).To(BeTrue())
		Expect(reg.Exists("job1.subjob2.key2")).To(BeTrue())

		Expect(reg.Exists("key3")).To(BeFalse())
		Expect(reg.Exists("job1.key3")).To(BeFalse())
		Expect(reg.Exists("job1.subjob2.key3")).To(BeFalse())
	})

	It("push without path should return error", func() {
		err := reg.Push("", NewStatus(nil, errors.New("key1")))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("path doesn't exist"))
	})

	It("push should overwrite the status", func() {
		err := reg.Push("job1.subjob1", NewStatus(nil, errors.New("job1.subjob1")))
		Expect(err).ToNot(HaveOccurred())

		groupStatus := reg.Get("job1.subjob1")
		Expect(groupStatus).ToNot(BeNil())

		status := groupStatus.GetStatus()
		Expect(status).ToNot(BeNil())

		statusDetail := status.GetDetails()
		Expect(statusDetail).To(Equal(testGroupStatus("job1.subjob1").Status.Details))
	})

	It("should delete the status from the result map", func() {
		Expect(reg.Exists("key2")).To(BeTrue())

		reg.Delete("job1.subjob1.key1")
		Expect(reg.Exists("job1.subjob1.key1")).To(BeFalse())
		reg.Delete("job1.subjob1.key2")
		Expect(reg.Exists("job1.subjob1.key2")).To(BeFalse())
		reg.Delete("job1.subjob1")
		Expect(reg.Exists("job1.subjob1")).To(BeFalse())
	})

	It("delete with no path provided should empty the whole registry", func() {
		reg.Delete("")
		Expect(reg.Keys()).To(HaveLen(0))
	})

	It("should get the status details for the given key", func() {
		groupStatus_key1 := reg.Get("key1")
		Expect(groupStatus_key1).ToNot(BeNil())
		status_key1 := groupStatus_key1.GetStatus()
		Expect(status_key1).ToNot(BeNil())
		statusDetail_key1 := status_key1.GetDetails()
		Expect(statusDetail_key1).To(Equal(testGroupStatus("key1").Status.Details))

		groupStatus_key2 := reg.Get("job1.subjob2.key1")
		Expect(groupStatus_key2).ToNot(BeNil())
		status_key2 := groupStatus_key2.GetStatus()
		Expect(status_key2).ToNot(BeNil())
		statusDetail_key2 := status_key2.GetDetails()
		Expect(statusDetail_key2).To(Equal(testGroupStatus("job1.subjob2.key1").Status.Details))
	})

	It("should return correct flattened group status map", func() {
		flatMap, err := reg.GetAllFlat()
		Expect(err).ToNot(HaveOccurred())

		Expect(flatMap["key1"].GetStatus().GetDetails()).To(Equal(testGroupStatus("key1").Status.Details))
	})
})

func testGroupStatus(msg string) *statusv1.GroupStatus {
	return &statusv1.GroupStatus{
		Status: &statusv1.Status{
			Details: &statusv1.Status_Error{
				Error: NewErrorDetails(errors.New(msg)),
			},
		},
	}
}
