package circuitfactory_test

import (
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/circuitfactory"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/controller"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

var componentId = runtime.NewComponentID("root.0")

var _ = Describe("Component factory", func() {
	Context("With unimplemented component type", func() {
		compProto := &policylangv1.Component{}
		It("Returns error if component type is not one of specified", func() {
			_, _, _, err := circuitfactory.NewComponentAndOptions(compProto, componentId, nil)
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
			deciderComponent := &components.Decider{}
			component, options, err := components.NewDeciderAndOptions(deciderProto, componentId, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(deciderComponent)))
			Expect(options).To(BeNil())
		})
	})

	Context("Alerter", func() {
		alerterProto := &policylangv1.Alerter{
			Parameters: &policylangv1.Alerter_Parameters{
				AlertName: "testName",
				Severity:  "crit",
			},
		}
		It("Creates Alerter component", func() {
			alerterComponent := &components.Alerter{}
			component, options, err := components.NewAlerterAndOptions(alerterProto, componentId, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(alerterComponent)))
			Expect(options).NotTo(BeNil())
		})
	})

	Context("Switcher", func() {
		switcherProto := &policylangv1.Switcher{}
		It("Creates Switcher component", func() {
			switcherComponent := &components.Switcher{}
			component, options, err := components.NewSwitcherAndOptions(switcherProto, componentId, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(switcherComponent)))
			Expect(options).To(BeNil())
		})
	})

	Context("Gradient", func() {
		Context("With correct gradient", func() {
			gradientControllerProto := &policylangv1.GradientController{
				Parameters: &policylangv1.GradientController_Parameters{
					Slope: -0.5,
				},
			}
			It("Creates Gradient controller", func() {
				controllerComponent := &controller.ControllerComponent{}
				component, options, err := controller.NewGradientControllerAndOptions(gradientControllerProto, componentId, nil)
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
	// 		component, options, err := component.NewEMAAndOptions(emaProto, componentId, nil)
	// 		Expect(err).NotTo(HaveOccurred())
	// 		Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(emaComponent)))
	// 		Expect(options).To(BeNil())
	// 	})
	// })
})
