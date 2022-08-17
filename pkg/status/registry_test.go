package status

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"

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

		err := reg.At("key1").Push(key1_status)
		Expect(err).NotTo(HaveOccurred())
		err = reg.At("key2").Push(NewStatus(nil, errors.New("key2")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.At("job1", "key1").Push(NewStatus(nil, errors.New("job1.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.At("job1", "key2").Push(NewStatus(nil, errors.New("job1.key2")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.At("job1", "subjob1", "key1").Push(NewStatus(nil, errors.New("job1.subjob1.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.At("job1", "subjob1", "key2").Push(NewStatus(nil, errors.New("job1.subjob1.key2")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.At("job1", "subjob2", "key1").Push(NewStatus(nil, errors.New("job1.subjob2.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = reg.At("job1", "subjob2", "key2").Push(NewStatus(nil, errors.New("job1.subjob2.key2")))
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

	It("should have paths that were pushed", func() {
		Expect(reg.At("key1").Exists()).To(BeTrue())
		Expect(reg.At("job1", "key2").Exists()).To(BeTrue())
		Expect(reg.At("job1", "subjob2", "key2").Exists()).To(BeTrue())

		Expect(reg.At("key3").Exists()).To(BeFalse())
		Expect(reg.At("job1", "key3").Exists()).To(BeFalse())
		Expect(reg.At("job1", "subjob2", "key3").Exists()).To(BeFalse())
	})

	It("should error when status is pushed without path", func() {
		err := reg.At("").Push(NewStatus(nil, errors.New("key1")))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("path doesn't exist"))
	})

	It("should overwrite the status when pushed to an existing path", func() {
		err := reg.At("job1", "subjob1").Push(NewStatus(nil, errors.New("job1.subjob1")))
		Expect(err).ToNot(HaveOccurred())

		statusDetail := reg.At("job1", "subjob1").Get().GetStatus().GetDetails()
		Expect(statusDetail).To(Equal(testGroupStatus("job1.subjob1").Status.Details))
	})

	It("should delete the status", func() {
		Expect(reg.At("job1", "subjob1").Exists()).To(BeTrue())
		reg.At("job1", "subjob1").Delete()
		Expect(reg.At("job1", "subjob1").Exists()).To(BeFalse())
	})

	It("should have no keys when root path is deleted", func() {
		reg.At("").Delete()
		Expect(reg.Keys()).To(HaveLen(0))
	})

	It("should get the status details for the given key", func() {
		statusDetail := reg.At("key1").Get().GetStatus().GetDetails()
		Expect(statusDetail).To(Equal(testGroupStatus("key1").Status.Details))

		statusDetail = reg.At("job1", "subjob2", "key1").Get().GetStatus().GetDetails()
		Expect(statusDetail).To(Equal(testGroupStatus("job1.subjob2.key1").Status.Details))
	})

	It("should return correct flattened group status map", func() {
		flatMap, err := reg.GetAllFlat()
		Expect(err).ToNot(HaveOccurred())

		Expect(flatMap["key1"].GetStatus().GetDetails()).To(Equal(testGroupStatus("key1").Status.Details))
	})

	It("should push, get, and delete statuses efficiently", Serial, Label("benchmark measurement"), func() {
		experiment := gmeasure.NewExperiment("status registry operations performance")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("one.two.measure%d", i)
			pathStatus := NewStatus(nil, errors.New(path))
			Expect(pathStatus.GetError().Message).To(Equal(path))
			experiment.MeasureDuration("Push", func() {
				err := reg.At(path).Push(pathStatus)
				Expect(err).ToNot(HaveOccurred())
			})
		}, gmeasure.SamplingConfig{N: 1000})

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("one.two.measure%d", i)
			experiment.MeasureDuration("Get", func() {
				statusDetail := reg.At(path).Get().GetStatus().GetDetails()
				Expect(statusDetail).To(Equal(testGroupStatus(path).Status.Details))
			})
		}, gmeasure.SamplingConfig{N: 1000})

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("one.two.measure%d", i)
			pathStatus := NewStatus(nil, errors.New(path))
			Expect(pathStatus.GetError().Message).To(Equal(path))
			experiment.MeasureDuration("Delete", func() {
				reg.At(path).Delete()
				ok := reg.At(path).Exists()
				Expect(ok).ToNot(BeTrue())
			})
		}, gmeasure.SamplingConfig{N: 1000})
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
