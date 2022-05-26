package controller_test

import (
/*"github.com/golang/mock/gomock"
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"
"go.uber.org/fx"

policylangv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/language/v1"
cn "aperture.tech/aperture/pkg/policies/controlloop/controller"
"aperture.tech/aperture/pkg/policies/mocks"
"aperture.tech/aperture/pkg/policies/controlplane/reading"*/
)

// TODO: Adapt this test to the new Circuit runtime
/*var _ = Describe("Gradient", func() {
	var (
		gradient cn.Controller
		options  fx.Option
		err      error

		t                      GinkgoTestReporter
		mockCtrl               *gomock.Controller
		mockControlLoopReadAPI *mocks.MockControlLoopReadAPI
		previousSignal         reading.Reading
	)

	BeforeEach(func() {
		t = GinkgoTestReporter{}
		mockCtrl = gomock.NewController(t)
		mockControlLoopReadAPI = mocks.NewMockControlLoopReadAPI(mockCtrl)
		previousSignal = reading.NewInvalid()

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
		gradient.MaintainOutput(reading.NewInvalid(), reading.NewInvalid(), nil)
	})

	It("Winds output", func() {
		currentOutput := reading.New(0.5)
		targetOutput := reading.New(1.0)
		output := gradient.WindOutput(currentOutput, targetOutput, nil)
		Expect(output).To(Equal(targetOutput))
	})

	It("Compute output returns invalid if current reading is invalid", func() {
		currentSignal := reading.NewInvalid()

		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Times(1)
		mockControlLoopReadAPI.EXPECT().GetLastOutput().Times(1)
		mockControlLoopReadAPI.EXPECT().GetControlVariableReading().Times(1)

		output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
		Expect(output.Valid).To(BeFalse())
	})

	It("Compute output returns invalid if setpoint reading is invalid", func() {
		currentSignal := reading.New(2.0)

		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(reading.NewInvalid()).Times(1)
		mockControlLoopReadAPI.EXPECT().GetLastOutput().Times(1)
		mockControlLoopReadAPI.EXPECT().GetControlVariableReading().Times(1)

		output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
		Expect(output.Valid).To(BeFalse())
	})

	It("Compute output returns invalid if last output reading is invalid", func() {
		currentSignal := reading.New(2.0)

		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(reading.New(1.0)).Times(1)
		mockControlLoopReadAPI.EXPECT().GetLastOutput().Return(reading.NewInvalid()).Times(1)
		mockControlLoopReadAPI.EXPECT().GetControlVariableReading().Times(1)

		output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
		Expect(output.Valid).To(BeFalse())
	})

	Context("With non-inverted gradient", func() {
		DescribeTable("Computes output",
			func(currentSignal reading.Reading, desiredOutput float64) {
				mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(reading.New(5.0)).Times(1)
				mockControlLoopReadAPI.EXPECT().GetLastOutput().Return(reading.New(3.0)).Times(1)

				output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
				Expect(output.Valid).To(BeTrue())
				Expect(output.Value).To(BeNumerically("~", desiredOutput))
			},
			Entry("normally", reading.New(4.0), 1.875),
			Entry("clamps to max value", reading.New(25.0), 0.6),
			Entry("clamps to min value", reading.New(0.1), 6.0),
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
			func(currentSignal reading.Reading, desiredOutput float64) {
				mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(reading.New(5.0)).Times(1)
				mockControlLoopReadAPI.EXPECT().GetLastOutput().Return(reading.New(3.0)).Times(1)

				output := gradient.ComputeOutput(previousSignal, currentSignal, mockControlLoopReadAPI)
				Expect(output.Valid).To(BeTrue())
				Expect(output.Value).To(BeNumerically("~", desiredOutput))
			},
			Entry("normally", reading.New(4.0), 1.2),
			Entry("clamps to max value", reading.New(1.0), 0.6),
			Entry("clamps to min value", reading.New(25.0), 6.0),
		)
	})
})*/
