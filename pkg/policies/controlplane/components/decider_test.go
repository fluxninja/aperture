package components_test

/*import (
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	cn "github.com/fluxninja/aperture/pkg/policies/controlloop/controller"
	"github.com/fluxninja/aperture/pkg/policies/mocks"
)

var _ = Describe("Timed", func() {
	var (
		timed   cn.Controller
		options fx.Option
		err     error

		t                      GinkgoTestReporter
		mockCtrl               *gomock.Controller
		mockControlLoopReadAPI *mocks.MockControlLoopReadAPI
		prev                   runtime.Reading
		curr                   runtime.Reading
	)

	BeforeEach(func() {
		t = GinkgoTestReporter{}
		mockCtrl = gomock.NewController(t)
		mockControlLoopReadAPI = mocks.NewMockControlLoopReadAPI(mockCtrl)
		prev = runtime.InvalidReading()
		curr = runtime.InvalidReading()

		duration := 5 * time.Second
		controller := &policylangv1.Controller{
			Controller: &policylangv1.Controller_Timed{
				Timed: &policylangv1.Timed{
					For:                         durationpb.New(duration),
					ControlValueOnPositiveError: 1.0,
					ControlValueOnNegativeError: 2.0,
				},
			},
		}
		timed, options, err = cn.NewControllerAndOptions(controller, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(options).To(BeNil())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("Maintains output", func() {
		timed.MaintainOutput(runtime.InvalidReading(), runtime.InvalidReading(), nil)
	})

	It("Winds output", func() {
		currentOutput := runtime.NewReading(0.5)
		targetOutput := runtime.NewReading(1.0)
		output := timed.WindOutput(currentOutput, targetOutput, nil)
		Expect(output).To(Equal(targetOutput))
	})

	It("Cannot wind output if target output invalid", func() {
		currentOutput := runtime.NewReading(1.0)
		targetOutput := runtime.InvalidReading()
		output := timed.WindOutput(currentOutput, targetOutput, nil)
		Expect(output.Valid).To(BeFalse())
	})

	// Timed is only using output from controlLoopReadAPI
	// Therefore prev & curr are invalid readings

	It("Cannot compute output if setpoint runtime invalid", func() {
		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(runtime.InvalidReading()).Times(1)

		output := timed.ComputeOutput(prev, curr, mockControlLoopReadAPI)
		Expect(output.Valid).To(BeFalse())
	})

	It("Cannot compute output if signal reading invalid", func() {
		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(runtime.NewReading(0.1)).Times(1)
		mockControlLoopReadAPI.EXPECT().GetSignalReading().Return(runtime.InvalidReading()).Times(1)

		output := timed.ComputeOutput(prev, curr, mockControlLoopReadAPI)
		Expect(output.Valid).To(BeFalse())
	})

	It("Compute output if setpoint > signal", func() {
		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(runtime.NewReading(5.0)).AnyTimes()
		mockControlLoopReadAPI.EXPECT().GetSignalReading().Return(runtime.NewReading(2.0)).AnyTimes()

		output := timed.ComputeOutput(prev, curr, mockControlLoopReadAPI)
		Expect(output.Value).To(Equal(1.0))
	})

	It("Compute output if setpoint < signal", func() {
		mockControlLoopReadAPI.EXPECT().GetSetpointReading().Return(runtime.NewReading(2.0)).AnyTimes()
		mockControlLoopReadAPI.EXPECT().GetSignalReading().Return(runtime.NewReading(5.0)).AnyTimes()

		output := timed.ComputeOutput(prev, curr, mockControlLoopReadAPI)
		Expect(output.Value).To(Equal(2.0))
	})
})*/
