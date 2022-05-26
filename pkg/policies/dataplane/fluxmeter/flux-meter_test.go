package fluxmeter_test

/*import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"

	policylangv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"aperture.tech/aperture/pkg/policies/dataplane/fluxmeter"
)

var _ = Describe("Fluxmeter", func() {
	Context("With proto", func() {
		selector := &policylangv1.Selector{
			Namespace: "testnamespace",
			Service:   "testservice",
			ControlPoint: &policylangv1.ControlPoint{
				Controlpoint: &policylangv1.ControlPoint_Feature{
					Feature: "awesomeFeature",
				},
			},
		}
		fluxMeterProto := &policylangv1.FluxMeter{
			Name:     "fluxmetername",
			Selector: selector,
		}

		var (
			newFluxmeter *fluxmeter.FluxMeter
			options      fx.Option
			err          error
		)

		BeforeEach(func() {
			newFluxmeter, options, err = fluxmeter.NewFluxMeterAndOptions(fluxMeterProto)

			Expect(err).NotTo(HaveOccurred())
			Expect(options).NotTo(BeNil())
			Expect(newFluxmeter).NotTo(BeNil())
		})

		It("Sets basic fields", func() {
			Expect(newFluxmeter.GetMetricName()).To(Equal("fluxmetername"))
			Expect(newFluxmeter.GetFluxMeterProto()).To(Equal(fluxMeterProto))
			Expect(newFluxmeter.GetSelector()).To(Equal(selector))
			Expect(newFluxmeter.GetBuckets()).To(BeNil())
			Expect(newFluxmeter.GetHistogram()).To(BeNil())
			Expect(newFluxmeter.GetMetricID()).To(BeEmpty())
		})
	})
})*/
