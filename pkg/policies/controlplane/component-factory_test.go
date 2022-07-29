package controlplane_test

import (
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/component"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/component/controller"
)

var _ = Describe("Component factory", func() {
	Context("With unimplemented component type", func() {
		compProto := &policylangv1.Component{}
		It("Returns error if component type is not one of specified", func() {
			_, _, _, _, _, err := controlplane.NewComponentAndOptions(compProto, 0, nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Decider", func() {
		duration := 5 * time.Second
		deciderProto := &policylangv1.Decider{
			Operator: "gte",
			TrueFor:  durationpb.New(duration),
			FalseFor: durationpb.New(duration),
		}
		It("Creates Decider component", func() {
			deciderComponent := &component.Decider{}
			component, options, err := component.NewDeciderAndOptions(deciderProto, 0, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(deciderComponent)))
			Expect(options).To(BeNil())
		})
	})

	Context("Gradient", func() {
		Context("With correct gradient", func() {
			gradientControllerProto := &policylangv1.GradientController{
				Tolerance: 0.5,
			}
			It("Creates Gradient controller", func() {
				controllerComponent := &controller.ControllerComponent{}
				component, options, err := controller.NewGradientControllerAndOptions(gradientControllerProto, 0, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(controllerComponent)))
				Expect(options).To(BeNil())
			})
		})
	})

	// Context("EMA", func() {
	// 	emaProto := &policylangv1.EMA{
	// 		EmaWindow:                      durationpb.New(2000 * time.Millisecond),
	// 		WarmUpWindow:                   durationpb.New(500 * time.Millisecond),
	// 		CorrectionFactorOnMinViolation: 2.0,
	// 		CorrectionFactorOnMaxViolation: 0.9,
	// 	}
	// 	It("Should not return error", func() {
	// 		emaComponent := &component.EMA{}
	// 		component, options, err := component.NewEMAAndOptions(emaProto, 0, nil)
	// 		Expect(err).NotTo(HaveOccurred())
	// 		Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(emaComponent)))
	// 		Expect(options).To(BeNil())
	// 	})
	// })
})
