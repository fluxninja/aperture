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
	var reg, j1, sj1, sj2 Registry
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

		j1 = NewRegistryPrefix(reg, "job1")

		err = j1.Push("key1", NewStatus(nil, errors.New("job1.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = j1.Push("key2", NewStatus(nil, errors.New("job1.key2")))
		Expect(err).NotTo(HaveOccurred())

		sj1 = NewRegistryPrefix(j1, "subjob1")
		err = sj1.Push("key1", NewStatus(nil, errors.New("job1.subjob1.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = sj1.Push("key2", NewStatus(nil, errors.New("job1.subjob1.key2")))
		Expect(err).NotTo(HaveOccurred())

		sj2 = NewRegistryPrefix(j1, "subjob2")
		err = sj2.Push("key1", NewStatus(nil, errors.New("job1.subjob2.key1")))
		Expect(err).NotTo(HaveOccurred())
		err = sj2.Push("key2", NewStatus(nil, errors.New("job1.subjob2.key2")))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return correct delimiters for root and sub registries", func() {
		Expect(reg.Delim()).To(Equal(defaultDelim))
		Expect(j1.Delim()).To(Equal(defaultDelim))
		Expect(sj1.Delim()).To(Equal(defaultDelim))
		Expect(sj2.Delim()).To(Equal(defaultDelim))
	})

	It("should return correct prefixes for root and sub registries", func() {
		Expect(reg.Prefix()).To(Equal(""))
		Expect(j1.Prefix()).To(Equal("job1"))
		Expect(sj1.Prefix()).To(Equal("job1.subjob1"))
		Expect(sj2.Prefix()).To(Equal("job1.subjob2"))
	})

	It("should return all the keys pushed to the result map", func() {
		regKeys := reg.Keys()
		regExpected := []string{"job1", "job1.key1", "job1.key2", "job1.subjob1", "job1.subjob1.key1", "job1.subjob1.key2", "job1.subjob2", "job1.subjob2.key1", "job1.subjob2.key2", "key1", "key2"}
		Expect(regKeys).To(Equal(regExpected))
		Expect(regKeys).To(HaveLen(11))

		j1regKeys := j1.Keys()
		j1regExpected := []string{"job1", "job1.key1", "job1.key2", "job1.subjob1", "job1.subjob1.key1", "job1.subjob1.key2", "job1.subjob2", "job1.subjob2.key1", "job1.subjob2.key2"}
		Expect(j1regKeys).To(Equal(j1regExpected))
		Expect(j1regKeys).To(HaveLen(9))

		sj1regKeys := sj1.Keys()
		sj1regExpected := []string{"job1.subjob1", "job1.subjob1.key1", "job1.subjob1.key2"}
		Expect(sj1regKeys).To(Equal(sj1regExpected))
		Expect(sj1regKeys).To(HaveLen(3))

		sj2regKeys := sj2.Keys()
		sj2regExpected := []string{"job1.subjob2", "job1.subjob2.key1", "job1.subjob2.key2"}
		Expect(sj2regKeys).To(Equal(sj2regExpected))
		Expect(sj2regKeys).To(HaveLen(3))
	})

	It("should have paths that were pushed", func() {
		Expect(reg.Exists("key1")).To(BeTrue())
		Expect(reg.Exists("key2")).To(BeTrue())
		Expect(sj2.Exists("key2")).To(BeTrue())

		Expect(reg.Exists("key3")).To(BeFalse())
		Expect(sj2.Exists("key3")).To(BeFalse())
	})

	It("should error when status is pushed without path", func() {
		err := reg.Push("", NewStatus(nil, errors.New("key1")))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("path doesn't exist"))

		err = j1.Push("", NewStatus(nil, errors.New("job1")))
		Expect(err).ToNot(HaveOccurred())
	})

	It("should overwrite the status when pushed to an existing path", func() {
		err := sj1.Push("key1", NewStatus(nil, errors.New("new")))
		Expect(err).ToNot(HaveOccurred())
		statusDetail := sj1.Get("key1").GetStatus().GetDetails()
		Expect(statusDetail).To(Equal(testGroupStatus("new").Status.Details))
	})

	It("should write to a sub group when pushed from root registry", func() {
		err := reg.Push("job1.subjob1.key1", NewStatus(nil, errors.New("new")))
		Expect(err).ToNot(HaveOccurred())
		statusDetail := reg.Get("job1.subjob1.key1").GetStatus().GetDetails()
		Expect(statusDetail).To(Equal(testGroupStatus("new").Status.Details))
	})

	It("should delete the status", func() {
		Expect(sj1.Exists("key1")).To(BeTrue())
		sj1.Delete("key1")
		Expect(sj1.Exists("key1")).To(BeFalse())

		Expect(sj1.Exists("")).To(BeTrue())
		sj1.Delete("")
		Expect(sj1.Exists("")).To(BeFalse())
	})

	It("should have no keys when root path is deleted", func() {
		reg.Delete("")
		Expect(reg.Keys()).To(HaveLen(0))
	})

	It("should get the status details for the given key", func() {
		s := reg.Get("key1").GetStatus().GetDetails()
		Expect(s).To(Equal(testGroupStatus("key1").Status.Details))

		s = sj2.Get("key1").GetStatus().GetDetails()
		Expect(s).To(Equal(testGroupStatus("job1.subjob2.key1").Status.Details))
	})

	It("should return correct flattened group status map", func() {
		flatMap, err := reg.GetAllFlat()
		Expect(err).ToNot(HaveOccurred())

		Expect(flatMap["key1"].GetStatus().GetDetails()).To(Equal(testGroupStatus("key1").Status.Details))
	})

	It("should push, get, and delete statuses efficiently", Serial, Label("benchmark measurement"), func() {
		experiment := gmeasure.NewExperiment("status registry operations performance")
		AddReportEntry(experiment.Name, experiment)

		r1 := NewRegistryPrefix(reg, "one")
		r2 := NewRegistryPrefix(r1, "two")

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("measure%d", i)
			pathStatus := NewStatus(nil, errors.New(path))
			Expect(pathStatus.GetError().Message).To(Equal(path))
			experiment.MeasureDuration("Push", func() {
				err := r2.Push(path, pathStatus)
				Expect(err).ToNot(HaveOccurred())
			})
		}, gmeasure.SamplingConfig{N: 1000})

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("measure%d", i)
			experiment.MeasureDuration("Get", func() {
				statusDetail := r2.Get(path).GetStatus().GetDetails()
				Expect(statusDetail).To(Equal(testGroupStatus(path).Status.Details))
			})
		}, gmeasure.SamplingConfig{N: 1000})

		experiment.Sample(func(i int) {
			path := fmt.Sprintf("measure%d", i)
			pathStatus := NewStatus(nil, errors.New(path))
			Expect(pathStatus.GetError().Message).To(Equal(path))
			experiment.MeasureDuration("Delete", func() {
				err := r2.Delete(path)
				Expect(err).ToNot(HaveOccurred())
				ok := r2.Exists(path)
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
