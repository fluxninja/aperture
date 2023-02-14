package circuitfactory_test

import (
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/circuitfactory"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/controller"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

var _ = Describe("Component factory", func() {
	Context("With unimplemented component type", func() {
		compProto := &policylangv1.Component{}
		It("Returns error if component type is not one of specified", func() {
			_, _, _, err := circuitfactory.NewComponentAndOptions(compProto, runtime.NewComponentID("root.0"), nil)
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
			component, options, err := components.NewDeciderAndOptions(deciderProto, "root.0", nil)
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
			component, options, err := components.NewAlerterAndOptions(alerterProto, "root.0", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(alerterComponent)))
			Expect(options).NotTo(BeNil())
		})
	})

	Context("Switcher", func() {
		switcherProto := &policylangv1.Switcher{}
		It("Creates Switcher component", func() {
			switcherComponent := &components.Switcher{}
			component, options, err := components.NewSwitcherAndOptions(switcherProto, "root.0", nil)
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
				component, options, err := controller.NewGradientControllerAndOptions(gradientControllerProto, "root.0", nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(controllerComponent)))
				Expect(options).To(BeNil())
			})
		})
	})

  Context("ArithmeticCombinator", func() {
    arithmeticCombinatorProto := &policylangv1.ArithmeticCombinator{
      Operator: "add",
    }
    arithmeticCombinatorComponent := &components.ArithmeticCombinator{}
    It("Creates ArithmeticCombinator component", func() {
      component, options, err := components.NewArithmeticCombinatorAndOptions(arithmeticCombinatorProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(arithmeticCombinatorComponent)))
      Expect(options).To(BeNil())
    })
  })

  Context("Differentiator", func() {
    differentiatorProto := &policylangv1.Differentiator{}
    It("Creates Differentiator component", func() {
      differentiatorComponent := &components.Differentiator{}
      differentiator,options, err := components.NewDifferentiatorAndOptions(differentiatorProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(differentiator)).To(Equal(reflect.TypeOf(differentiatorComponent)))
      Expect(options).To(BeNil())
    })
  })

  Context("First-Valid", func() {
    firstValidProto := &policylangv1.FirstValid{}
    It("Creates First-Valid component", func() {
      firstValidComponent := &components.FirstValid{}
      firstValid, options, err := components.NewFirstValidAndOptions(firstValidProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(firstValid)).To(Equal(reflect.TypeOf(firstValidComponent)))
      Expect(options).To(BeNil())
    })
  })

  Context("Integrator", func() {
    integratorProto := &policylangv1.Integrator{}
    It("Creates Integrator component", func() {
      integratorComponent := &components.Integrator{}
      integrator, options, err := components.NewIntegratorAndOptions(integratorProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(integrator)).To(Equal(reflect.TypeOf(integratorComponent)))
      Expect(options).To(BeNil())
    })
  })

  Context("Max", func() {
    maxProto := &policylangv1.Max{}
    It("Creates Max component", func() {
      maxComponent := &components.Max{}
      max, options, err := components.NewMaxAndOptions(maxProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(max)).To(Equal(reflect.TypeOf(maxComponent)))
      Expect(options).To(BeNil())
    })
  })

  Context("Min", func() {
    minProto := &policylangv1.Min{}
    It("Creates Min component", func() {
      minComponent := &components.Min{}
      min, options, err := components.NewMinAndOptions(minProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(min)).To(Equal(reflect.TypeOf(minComponent)))
      Expect(options).To(BeNil())
    })
  })

  Context("Nested-Signals", func() {
    nestedIngressSignalProto := &policylangv1.NestedSignalIngress{}
    It("Creates Nested-Signal-Ingress component", func() {
      nestedSignalComponent := &components.NestedSignalIngress{}
      nestedSignal, options, err := components.NewNestedSignalIngressAndOptions(nestedIngressSignalProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(nestedSignal)).To(Equal(reflect.TypeOf(nestedSignalComponent)))
      Expect(options).To(BeNil())
    })
    nestedEgressSignalProto := &policylangv1.NestedSignalEgress{}
    It("Creates Nested-Signal-Egress component", func() {
      nestedSignalComponent := &components.NestedSignalEgress{}
      nestedSignal, options, err := components.NewNestedSignalEgressAndOptions(nestedEgressSignalProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(nestedSignal)).To(Equal(reflect.TypeOf(nestedSignalComponent)))
      Expect(options).To(BeNil())
    })
  })

  Context("Sqrt", func() {
    sqrtProto := &policylangv1.Sqrt{}
    It("Creates Sqrt component", func() {
      sqrtComponent := &components.Sqrt{}
      sqrt, options, err := components.NewSqrtAndOptions(sqrtProto, "root.0", nil)
      Expect(err).NotTo(HaveOccurred())
      Expect(reflect.TypeOf(sqrt)).To(Equal(reflect.TypeOf(sqrtComponent)))
      Expect(options).To(BeNil())
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
	// 		component, options, err := component.NewEMAAndOptions(emaProto, "root.0", nil)
	// 		Expect(err).NotTo(HaveOccurred())
	// 		Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(emaComponent)))
	// 		Expect(options).To(BeNil())
	// 	})
	// })
})
