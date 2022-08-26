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
	var reg, j1, j1s1, j1s2, j1s1k1, j1s1k2, j1s2k1, j1s2k2 Registry

	BeforeEach(func() {
		reg = NewRegistry(nil, "")
		Expect(reg).NotTo(BeNil())

		j1 = NewRegistry(reg, "job1")
		err := j1.Push(NewStatus(nil, errors.New("job1")))
		Expect(err).NotTo(HaveOccurred())

		j1s1 = NewRegistry(j1, "subjob1")
		err = j1s1.Push(NewStatus(nil, errors.New("job1.subjob1")))
		Expect(err).NotTo(HaveOccurred())
		j1s2 = NewRegistry(j1, "subjob2")
		err = j1s2.Push(NewStatus(nil, errors.New("job1.subjob2")))
		Expect(err).NotTo(HaveOccurred())

		j1s1k1 = NewRegistry(j1s1, "key1")
		err = j1s1k1.Push(NewStatus(nil, errors.New("job1.subjob1.key1")))
		Expect(err).NotTo(HaveOccurred())
		j1s1k2 = NewRegistry(j1s1, "key2")
		err = j1s1k2.Push(NewStatus(nil, errors.New("job1.subjob1.key2")))
		Expect(err).NotTo(HaveOccurred())

		j1s2k1 = NewRegistry(j1s2, "key1")
		err = j1s2k1.Push(NewStatus(nil, errors.New("job1.subjob2.key1")))
		Expect(err).NotTo(HaveOccurred())
		j1s2k2 = NewRegistry(j1s2, "key2")
		err = j1s2k2.Push(NewStatus(nil, errors.New("job1.subjob2.key2")))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return correct delimiters for root and sub registries", func() {
		Expect(reg.Delim()).To(Equal(defaultDelim))
		Expect(j1.Delim()).To(Equal(defaultDelim))
		Expect(j1s1.Delim()).To(Equal(defaultDelim))
		Expect(j1s2k1.Delim()).To(Equal(defaultDelim))
	})

	It("should return correct prefixes for root and sub registries", func() {
		Expect(reg.Path()).To(Equal(""))
		Expect(j1.Path()).To(Equal("job1"))
		Expect(j1s1.Path()).To(Equal("job1.subjob1"))
		Expect(j1s2k1.Path()).To(Equal("job1.subjob2.key1"))
	})

	It("should return all the keys pushed to the result map", func() {
		regKeys := reg.Keys()
		regExpected := []string{"job1", "job1.subjob1", "job1.subjob1.key1", "job1.subjob1.key2", "job1.subjob2", "job1.subjob2.key1", "job1.subjob2.key2"}
		Expect(regKeys).To(Equal(regExpected))
		Expect(regKeys).To(HaveLen(7))

		j1regKeys := j1.Keys()
		j1regExpected := []string{"job1", "job1.subjob1", "job1.subjob1.key1", "job1.subjob1.key2", "job1.subjob2", "job1.subjob2.key1", "job1.subjob2.key2"}
		Expect(j1regKeys).To(Equal(j1regExpected))
		Expect(j1regKeys).To(HaveLen(7))

		j1s1regKeys := j1s1.Keys()
		j1s1regKeysExpected := []string{"job1.subjob1", "job1.subjob1.key1", "job1.subjob1.key2"}
		Expect(j1s1regKeys).To(Equal(j1s1regKeysExpected))
		Expect(j1s1regKeys).To(HaveLen(3))

		j1s2k1regKeys := j1s2k1.Keys()
		j1s2k1regKeysExpected := []string{"job1.subjob2.key1"}
		Expect(j1s2k1regKeys).To(Equal(j1s2k1regKeysExpected))
		Expect(j1s2k1regKeys).To(HaveLen(1))
	})

	It("should have paths that were pushed", func() {
		Expect(j1.Exists()).To(BeTrue())
		Expect(j1s1k2.Exists()).To(BeTrue())

		r := NewRegistry(reg, "r1")
		Expect(r.Exists()).To(BeFalse())
	})

	It("should error when status is pushed without path", func() {
		err := reg.Push(NewStatus(nil, errors.New("key1")))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("path doesn't exist"))

		n1 := NewRegistry(j1, "")
		err = n1.Push(NewStatus(nil, errors.New("newjob1")))
		Expect(err).ToNot(HaveOccurred())
	})

	It("should overwrite the status when pushed to an existing path", func() {
		err := j1s1k1.Push(NewStatus(nil, errors.New("new")))
		Expect(err).ToNot(HaveOccurred())
		statusDetail := j1s1k1.Get().GetStatus().GetDetails()
		Expect(statusDetail).To(Equal(testGroupStatus("new").Status.Details))
	})

	It("should delete the status", func() {
		Expect(j1s1k1.Exists()).To(BeTrue())
		j1s1k1.Delete()
		Expect(j1s1k1.Exists()).To(BeFalse())
	})

	It("should have no keys when root path is deleted", func() {
		j1s1.Delete()
		Expect(reg.Keys()).To(HaveLen(4))
		Expect(j1s1.Keys()).To(HaveLen(0))

		reg.Delete()
		Expect(reg.Keys()).To(HaveLen(0))
	})

	It("should get the status details for the given key", func() {
		s := j1.Get().GetStatus().GetDetails()
		Expect(s).To(Equal(testGroupStatus("job1").Status.Details))

		s = j1s2k1.Get().GetStatus().GetDetails()
		Expect(s).To(Equal(testGroupStatus("job1.subjob2.key1").Status.Details))
	})

	It("should return correct flattened group status map", func() {
		flatMap, err := reg.GetAllFlat()
		Expect(err).ToNot(HaveOccurred())

		Expect(flatMap["job1"].GetStatus().GetDetails()).To(Equal(testGroupStatus("job1").Status.Details))
	})

	It("should push, get, and delete statuses efficiently", Serial, Label("benchmark measurement"), func() {
		experiment := gmeasure.NewExperiment("status registry operations performance")
		AddReportEntry(experiment.Name, experiment)

		r1 := NewRegistry(reg, "one")
		r2 := NewRegistry(r1, "two")

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("measure%d", i)
			m := NewRegistry(r2, path)
			pathStatus := NewStatus(nil, errors.New(path))
			Expect(pathStatus.GetError().Message).To(Equal(path))
			experiment.MeasureDuration("Push", func() {
				err := m.Push(pathStatus)
				Expect(err).ToNot(HaveOccurred())
			})
		}, gmeasure.SamplingConfig{N: 1000})

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("measure%d", i)
			m := NewRegistry(r2, path)
			experiment.MeasureDuration("Get", func() {
				statusDetail := m.Get().GetStatus().GetDetails()
				Expect(statusDetail).To(Equal(testGroupStatus(path).Status.Details))
			})
		}, gmeasure.SamplingConfig{N: 1000})

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("measure%d", i)
			m := NewRegistry(r2, path)
			pathStatus := NewStatus(nil, errors.New(path))
			Expect(pathStatus.GetError().Message).To(Equal(path))
			experiment.MeasureDuration("Delete", func() {
				err := m.Delete()
				Expect(err).ToNot(HaveOccurred())
				ok := m.Exists()
				Expect(ok).To(BeFalse())
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
