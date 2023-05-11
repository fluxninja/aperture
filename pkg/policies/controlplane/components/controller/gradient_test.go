package controller_test

/*"github.com/golang/mock/gomock"
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"
"go.uber.org/fx"

policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
cn "github.com/fluxninja/aperture/v2/pkg/policies/controlloop/controller"
"github.com/fluxninja/aperture/v2/pkg/policies/mocks"*/

// TODO: Adapt this test to the new Circuit runtime
/*var _ = Describe("Gradient", func() {
	var (
		gradient cn.Controller
		options  fx.Option
		err      error

		t                      GinkgoTestReporter
		mockCtrl               *gomock.Controller
		mockControlLoopReadAPI *mocks.MockControlLoopReadAPI
		previousSignal         runtime.Reading
	)

	BeforeEach(func() {
		t = GinkgoTestReporter{}
		mockCtrl = gomock.NewController(t)
		mockControlLoopReadAPI = mocks.NewMockControlLoopReadAPI(mockCtrl)
		previousSignal = runtime.InvalidReading()

		controller := &policylangv1.Controller{
			Controller: &policylangv1.Controller_Gradient{
				Gradient: &policylangv1.Gradient{
					Inverted:    false,
					MinGradient: 0.2,
					MaxGradient: 2.0,
					Tolerance:   0.5,
				},
			},
		}
		gradient, options, err = cn.NewControllerAndOptions(controller, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(options).To(BeNil())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	// Just for coverage
	It("Maintains output", func() {
		gradient.MaintainOutput(runtime.InvalidReading(), runtime.InvalidReading(), nil)
	})

	It("Winds output", func() {
		currentOutput := runtime.NewReading(0.5)
		targetOutput := runtime.NewReading(1.0)
		output := gradient.WindOutput(currentOutput, targetOutput, nil)
		Expect(output).To(Equal(targetOutput))
	})

	It("Compute output returns invalid if current reading is invalid", func() {
		currentSignal := runtime.InvalidReading()

		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Times(1)
		mockControlLoopReadAPI.EXPECT().GetLastOutput().Times(1)
		mockControlLoopReadAPI.EXPECT().GetControlVariableReading().Times(1)

		output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
		Expect(output.Valid).To(BeFalse())
	})

	It("Compute output returns invalid if setpoint reading is invalid", func() {
		currentSignal := runtime.NewReading(2.0)

		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(runtime.InvalidReading()).Times(1)
		mockControlLoopReadAPI.EXPECT().GetLastOutput().Times(1)
		mockControlLoopReadAPI.EXPECT().GetControlVariableReading().Times(1)

		output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
		Expect(output.Valid).To(BeFalse())
	})

	It("Compute output returns invalid if last output reading is invalid", func() {
		currentSignal := runtime.NewReading(2.0)

		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(runtime.NewReading(1.0)).Times(1)
		mockControlLoopReadAPI.EXPECT().GetLastOutput().Return(runtime.InvalidReading()).Times(1)
		mockControlLoopReadAPI.EXPECT().GetControlVariableReading().Times(1)

		output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
		Expect(output.Valid).To(BeFalse())
	})

	Context("With non-inverted gradient", func() {
		DescribeTable("Computes output",
			func(currentSignal runtime.Reading, desiredOutput float64) {
				mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(runtime.NewReading(5.0)).Times(1)
				mockControlLoopReadAPI.EXPECT().GetLastOutput().Return(runtime.NewReading(3.0)).Times(1)

				output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
				Expect(output.Valid).To(BeTrue())
				Expect(output.Value).To(BeNumerically("~", desiredOutput))
			},
			Entry("normally", runtime.NewReading(4.0), 1.875),
			Entry("clamps to max value", runtime.NewReading(25.0), 0.6),
			Entry("clamps to min value", runtime.NewReading(0.1), 6.0),
		)
	})

	Context("With inverted gradient", func() {
		BeforeEach(func() {
			controller := &policylangv1.Controller{
				Controller: &policylangv1.Controller_Gradient{
					Gradient: &policylangv1.Gradient{
						Inverted:    true,
						MinGradient: 0.2,
						MaxGradient: 2.0,
						Tolerance:   0.5,
					},
				},
			}
			gradient, options, err = cn.NewControllerAndOptions(controller, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(options).To(BeNil())
		})

		DescribeTable("Computes output",
			func(currentSignal runtime.Reading, desiredOutput float64) {
				mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(runtime.NewReading(5.0)).Times(1)
				mockControlLoopReadAPI.EXPECT().GetLastOutput().Return(runtime.NewReading(3.0)).Times(1)

				output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
				Expect(output.Valid).To(BeTrue())
				Expect(output.Value).To(BeNumerically("~", desiredOutput))
			},
			Entry("normally", runtime.NewReading(4.0), 1.2),
			Entry("clamps to max value", runtime.NewReading(1.0), 0.6),
			Entry("clamps to min value", runtime.NewReading(25.0), 6.0),
		)
	})
})*/
