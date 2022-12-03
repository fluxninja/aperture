package controlplane_test

import (
	"reflect"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/controller"
	"github.com/fluxninja/aperture/pkg/policies/mocks"
)

var _ = Describe("Component factory", func() {
	Context("With unimplemented component type", func() {
		compProto := &policylangv1.Component{}
		It("Returns error if component type is not one of specified", func() {
			_, _, _, err := controlplane.NewComponentAndOptions(compProto, 0, nil)
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
			component, options, err := components.NewDeciderAndOptions(deciderProto, 0, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(deciderComponent)))
			Expect(options).To(BeNil())
		})
	})

	Context("Alerter", func() {
		var policyAPI *mocks.MockPolicy
		var alerterProto *policylangv1.Alerter
		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			policyAPI = mocks.NewMockPolicy(ctrl)
			alerterProto = &policylangv1.Alerter{
				AlertName: "testName",
				Severity:  "crit",
			}
			policyAPI.EXPECT().GetPolicyName().Return("test1").AnyTimes()
		})

		It("Creates Alerter component", func() {
			alerterComponent := &components.Alerter{}
			component, options, err := components.NewAlerterAndOptions(alerterProto, 0, policyAPI)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(alerterComponent)))
			Expect(options).NotTo(BeNil())
		})
	})

	Context("Switcher", func() {
		switcherProto := &policylangv1.Switcher{}
		It("Creates Switcher component", func() {
			switcherComponent := &components.Switcher{}
			component, options, err := components.NewSwitcherAndOptions(switcherProto, 0, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(reflect.TypeOf(component)).To(Equal(reflect.TypeOf(switcherComponent)))
			Expect(options).To(BeNil())
		})
	})

	Context("Gradient", func() {
		Context("With correct gradient", func() {
			gradientControllerProto := &policylangv1.GradientController{
				Slope: -0.5,
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
})
