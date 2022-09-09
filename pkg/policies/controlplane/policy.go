package controlplane

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	goObjectHash "github.com/benlaurie/objecthash/go/objecthash"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	wrappersv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/wrappers/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/resources/fluxmeter"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// policyModule returns Fx options of Policy for the Main App.
func policyModule() fx.Option {
	// Circuit module options
	componentFactoryOptions := componentFactoryModule()

	return fx.Options(
		circuitFactoryModule(),
		componentFactoryOptions,
	)
}

// Policy invokes the Circuit runtime at tick frequency.
type Policy struct {
	iface.PolicyBase
	// Circuit
	circuit *runtime.Circuit
	// job group
	circuitJobGroup *jobs.JobGroup
	// Job name
	jobName string
	// Evaluation interval determines how often the Circuit gets executed
	evaluationInterval time.Duration
}

// Make sure Policy complies with PolicyAPI interface.
var _ iface.Policy = (*Policy)(nil)

// newPolicyOptions creates a new Policy object and returns its Fx options for the per Policy App.
func newPolicyOptions(
	wrapperMessage *wrappersv1.PolicyWrapper,
) (fx.Option, error) {
	// List of options for the policy.
	policyOptions := []fx.Option{}
	policy, compiledCircuit, partialPolicyOption, err := compilePolicyWrapper(wrapperMessage)
	if err != nil {
		return nil, err
	}
	policyOptions = append(policyOptions, partialPolicyOption)
	policyOptions = append(policyOptions, fx.Supply(
		fx.Annotate(policy, fx.As(new(iface.Policy))),
	))

	compWithPortsList := make([]runtime.CompiledComponentAndPorts, 0, len(compiledCircuit))
	for _, compiledComponent := range compiledCircuit {
		// Skip nil component
		if compiledComponent.CompiledComponent.Component != nil {
			compWithPortsList = append(compWithPortsList, compiledComponent.CompiledComponentAndPorts)
		}
	}

	// Create circuit
	circuit, circuitOption := runtime.NewCircuitAndOptions(compWithPortsList, policy)
	policyOptions = append(policyOptions, circuitOption)

	policyOptions = append(policyOptions, componentFactoryModuleForPolicyApp(circuit))

	policyOptions = append(policyOptions, fx.Supply(fx.Annotate(circuit, fx.As(new(runtime.CircuitAPI)))))
	policy.circuit = circuit

	return fx.Options(policyOptions...), nil
}

// CompilePolicy takes policyMessage and returns a compiled policy. This is a helper method for standalone consumption of policy compiler.
func CompilePolicy(policyMessage *policylangv1.Policy) (CompiledCircuit, error) {
	wrapperMessage, err := hashAndPolicyWrap(policyMessage, "DoesNotMatter")
	if err != nil {
		return nil, err
	}
	_, compWithPortsList, _, err := compilePolicyWrapper(wrapperMessage)
	if err != nil {
		return nil, err
	}
	return compWithPortsList, nil
}

// compilePolicyWrapper takes policyProto and returns a compiled policy.
func compilePolicyWrapper(wrapperMessage *wrappersv1.PolicyWrapper) (*Policy, CompiledCircuit, fx.Option, error) {
	if wrapperMessage == nil {
		return nil, nil, nil, fmt.Errorf("nil policy wrapper message")
	}

	policy := &Policy{
		PolicyBase: wrapperMessage,
	}

	// Get Policy Proto
	policyProto := wrapperMessage.GetPolicy()
	if policyProto == nil {
		return nil, nil, nil, fmt.Errorf("nil policy proto")
	}

	var resourceOptions []fx.Option
	if policyProto.GetResources() != nil {
		// Initialize flux meters
		for name, fluxMeterProto := range policyProto.GetResources().FluxMeters {
			fluxMeterOption, err := fluxmeter.NewFluxMeterOptions(name, fluxMeterProto, policy)
			if err != nil {
				return nil, nil, nil, err
			}
			resourceOptions = append(resourceOptions, fluxMeterOption)
		}
		// Initialize classifiers
		for index, classifierProto := range policyProto.GetResources().Classifiers {
			classifierOption, err := classifier.NewClassifierOptions(int64(index), classifierProto, policy)
			if err != nil {
				return nil, nil, nil, err
			}
			resourceOptions = append(resourceOptions, classifierOption)
		}
	}
	var compiledCircuit CompiledCircuit
	partialCircuitOption := fx.Options()
	var err error

	if policyProto.GetCircuit() != nil {
		// Read evaluation interval
		policy.evaluationInterval = policyProto.GetCircuit().GetEvaluationInterval().AsDuration()

		compiledCircuit, partialCircuitOption, err = compileCircuit(policyProto.GetCircuit().Components, policy)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return policy, compiledCircuit, fx.Options(
		fx.Options(resourceOptions...),
		partialCircuitOption,
		fx.Invoke(policy.setupCircuitJob),
	), nil
}

// GetEvaluationInterval returns the ID of the policy.
func (policy *Policy) GetEvaluationInterval() time.Duration {
	return policy.evaluationInterval
}

func (policy *Policy) setupCircuitJob(
	lifecycle fx.Lifecycle,
	circuitJobGroup *jobs.JobGroup,
) error {
	if policy.evaluationInterval > 0 {
		// Job name
		policy.jobName = fmt.Sprintf("Policy-%s", policy.GetPolicyName())
		// Job group
		policy.circuitJobGroup = circuitJobGroup

		lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				// Create a job that runs every tick i.e. evaluation_interval. Set timeout duration to half of evaluation_interval
				job := jobs.BasicJob{
					JobFunc: policy.executeTick,
				}
				job.JobName = policy.jobName
				initialDelay := config.MakeDuration(0)
				executionPeriod := config.MakeDuration(policy.evaluationInterval)
				executionTimeout := config.MakeDuration(time.Millisecond * 100)
				jobConfig := jobs.JobConfig{
					InitiallyHealthy: true,
					InitialDelay:     initialDelay,
					ExecutionPeriod:  executionPeriod,
					ExecutionTimeout: executionTimeout,
				}
				// Register job with registry
				err := policy.circuitJobGroup.RegisterJob(&job, jobConfig)
				if err != nil {
					log.Error().Err(err).Str("job", policy.jobName).Msg("Error registering job")
					return err
				}
				return nil
			},
			OnStop: func(_ context.Context) error {
				// Deregister job from registry
				_ = policy.circuitJobGroup.DeregisterJob(policy.jobName)
				return nil
			},
		})
	}

	return nil
}

func (policy *Policy) executeTick(jobCtxt context.Context) (proto.Message, error) {
	// Get JobInfo
	jobInfo := policy.circuitJobGroup.JobInfo(policy.jobName)
	if jobInfo == nil {
		return nil, fmt.Errorf("job info not found for job %s", policy.jobName)
	}
	tickInfo := runtime.TickInfo{
		Timestamp:     jobInfo.LastRunTime,
		NextTimestamp: jobInfo.NextRunTime,
		Tick:          jobInfo.RunCount,
		Interval:      policy.evaluationInterval,
	}
	// Execute Circuit
	err := policy.circuit.Execute(tickInfo)
	// TODO: return tick info (publish to health framework) instead of returning nil proto.Message
	return nil, err
}

// hashAndPolicyWrap wraps a proto message with a config properties wrapper and hashes it.
func hashAndPolicyWrap(policyMessage *policylangv1.Policy, policyName string) (*wrappersv1.PolicyWrapper, error) {
	dat, marshalErr := yaml.Marshal(policyMessage)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msgf("Failed to marshal proto message %+v", policyMessage)
		return nil, marshalErr
	}
	hashBytes, hashErr := goObjectHash.ObjectHash(dat)
	if hashErr != nil {
		log.Warn().Err(hashErr).Msgf("Failed to hash json serialized proto message %s", string(dat))
		return nil, hashErr
	}
	hash := base64.StdEncoding.EncodeToString(hashBytes[:])

	return &wrappersv1.PolicyWrapper{
		Policy:     policyMessage,
		PolicyName: policyName,
		PolicyHash: hash,
	}, nil
}
