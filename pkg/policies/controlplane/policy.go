package controlplane

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	goObjectHash "github.com/benlaurie/objecthash/go/objecthash"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/circuitfactory"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/resources/classifier"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/resources/fluxmeter"
	telemetrycollectors "github.com/fluxninja/aperture/v2/pkg/policies/controlplane/resources/telemetry-collectors"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// policyModule returns Fx options of Policy for the Main App.
func policyModule() fx.Option { return circuitfactory.Module() }

// Policy invokes the Circuit runtime at tick frequency.
type Policy struct {
	iface.PolicyBase
	// status registry
	registry status.Registry
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
func newPolicyOptions(wrapperMessage *policysyncv1.PolicyWrapper, registry status.Registry) (fx.Option, error) {
	// List of options for the policy.
	policyOptions := []fx.Option{}
	policy, compiledCircuit, partialPolicyOption, err := compilePolicyWrapper(wrapperMessage, registry)
	if err != nil {
		return nil, err
	}
	policyOptions = append(policyOptions, partialPolicyOption)
	policyOptions = append(policyOptions, fx.Supply(
		fx.Annotate(policy, fx.As(new(iface.Policy))),
	))

	// Create circuit
	circuit, circuitOption := runtime.NewCircuitAndOptions(compiledCircuit.Components(), policy)
	policyOptions = append(policyOptions, circuitOption)

	policyOptions = append(policyOptions, circuitfactory.FactoryModuleForPolicyApp(circuit))

	policyOptions = append(policyOptions, fx.Supply(fx.Annotate(circuit, fx.As(new(runtime.CircuitAPI)))))
	policy.circuit = circuit

	return fx.Options(policyOptions...), nil
}

// CompilePolicy takes policyMessage and returns a compiled policy. This is a helper method for standalone consumption of policy compiler.
func CompilePolicy(policyMessage *policylangv1.Policy, registry status.Registry) (*circuitfactory.Circuit, error) {
	wrapperMessage, err := hashAndPolicyWrap(policyMessage, "DoesNotMatter", nil)
	if err != nil {
		return nil, err
	}
	_, circuit, _, err := compilePolicyWrapper(wrapperMessage, registry)
	if err != nil {
		return nil, err
	}
	return circuit, nil
}

// compilePolicyWrapper takes policyProto and returns a compiled policy.
func compilePolicyWrapper(wrapperMessage *policysyncv1.PolicyWrapper, registry status.Registry) (*Policy, *circuitfactory.Circuit, fx.Option, error) {
	if wrapperMessage == nil {
		return nil, nil, nil, fmt.Errorf("nil policy wrapper message")
	}

	policy := &Policy{
		PolicyBase: wrapperMessage.GetCommonAttributes(),
		registry:   registry,
	}

	// Get Policy Proto
	policyProto := wrapperMessage.GetPolicy()
	if policyProto == nil {
		return nil, nil, nil, fmt.Errorf("nil policy proto")
	}

	var resourceOptions []fx.Option
	resources := policyProto.GetResources()
	if resources != nil {
		flowControl := resources.GetFlowControl()
		if flowControl != nil {
			// Initialize flux meters
			fluxMeters := flowControl.GetFluxMeters()
			for name, fluxMeterProto := range fluxMeters {
				fluxMeterOption, err := fluxmeter.NewFluxMeterOptions(name, fluxMeterProto, policy)
				if err != nil {
					return nil, nil, nil, err
				}
				resourceOptions = append(resourceOptions, fluxMeterOption)
			}
			// Initialize classifiers
			classifiers := flowControl.GetClassifiers()
			for index, classifierProto := range classifiers {
				classifierOption, err := classifier.NewClassifierOptions(index, classifierProto, policy)
				if err != nil {
					return nil, nil, nil, err
				}
				resourceOptions = append(resourceOptions, classifierOption)
			}
		}
		telemetryCollectors := resources.GetTelemetryCollectors()
		if telemetryCollectors != nil {
			tcOption, err := telemetrycollectors.NewTelemetryCollectorsOptions(telemetryCollectors, policy)
			if err != nil {
				return nil, nil, nil, err
			}
			resourceOptions = append(resourceOptions, tcOption)
		}
	}
	var compiledCircuit *circuitfactory.Circuit
	partialCircuitOption := fx.Options()
	var err error

	if policyProto.GetCircuit() != nil {
		// Read evaluation interval
		policy.evaluationInterval = policyProto.GetCircuit().GetEvaluationInterval().AsDuration()

		compiledCircuit, partialCircuitOption, err = circuitfactory.CompileFromProto(
			policyProto,
			policy,
		)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return policy, compiledCircuit, fx.Options(
		fx.Options(resourceOptions...),
		partialCircuitOption,
		fx.Invoke(
			policy.setupCircuitJob,
			policy.setupDynamicConfig,
		),
	), nil
}

// GetEvaluationInterval returns the ID of the policy.
func (policy *Policy) GetEvaluationInterval() time.Duration {
	return policy.evaluationInterval
}

func (policy *Policy) setupCircuitJob(lifecycle fx.Lifecycle, circuitJobGroup *jobs.JobGroup) error {
	logger := policy.GetStatusRegistry().GetLogger()
	if policy.evaluationInterval > 0 {
		// Job name
		policy.jobName = fmt.Sprintf("Policy-%s", policy.GetPolicyName())
		// Job group
		policy.circuitJobGroup = circuitJobGroup

		lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				// Create a job that runs every tick i.e. evaluation_interval. Set timeout duration to half of evaluation_interval
				job := jobs.NewBasicJob(policy.jobName, policy.executeTick)
				executionPeriod := config.MakeDuration(policy.evaluationInterval)
				executionTimeout := config.MakeDuration(time.Millisecond * 100)
				jobConfig := jobs.JobConfig{
					InitiallyHealthy: true,
					ExecutionPeriod:  executionPeriod,
					ExecutionTimeout: executionTimeout,
				}
				// Register job with registry
				err := policy.circuitJobGroup.RegisterJob(job, jobConfig)
				if err != nil {
					logger.Error().Err(err).Str("job", policy.jobName).Msg("Error registering job")
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

func (policy *Policy) setupDynamicConfig(dynamicConfigWatcher notifiers.Watcher, lifecycle fx.Lifecycle) error {
	unmarshaller, _ := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller([]byte{})
	unmarshalNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(policy.GetPolicyName()),
		unmarshaller,
		policy.dynamicConfigUpdate,
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return dynamicConfigWatcher.AddKeyNotifier(unmarshalNotifier)
		},
		OnStop: func(_ context.Context) error {
			return dynamicConfigWatcher.RemoveKeyNotifier(unmarshalNotifier)
		},
	})

	return nil
}

func (policy *Policy) dynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	policy.circuit.DynamicConfigUpdate(event, unmarshaller)
}

func (policy *Policy) executeTick(jobCtxt context.Context) (proto.Message, error) {
	// Get JobInfo
	jobInfo, err := policy.circuitJobGroup.JobInfo(policy.jobName)
	if err != nil {
		return nil, err
	}

	tickInfo := runtime.NewTickInfo(jobInfo.LastExecuteTime,
		jobInfo.ExecuteCount,
		policy.evaluationInterval)
	// Execute Circuit
	err = policy.circuit.Execute(tickInfo)
	// TODO: return tick info (publish to health framework) instead of returning nil proto.Message
	return nil, err
}

// GetStatusRegistry returns the status registry of the policy.
func (policy *Policy) GetStatusRegistry() status.Registry {
	return policy.registry
}

// hashAndPolicyWrap wraps a proto message with a config properties wrapper and hashes it.
func hashAndPolicyWrap(policyMessage *policylangv1.Policy, policyName string, metadata *policylangv1.PolicyMetadata) (*policysyncv1.PolicyWrapper, error) {
	jsonDat, marshalErr := json.Marshal(policyMessage)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msgf("Failed to marshal proto message %+v", policyMessage)
		return nil, marshalErr
	}
	// convert dat to yaml format
	dat, marshalErr := yaml.JSONToYAML(jsonDat)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msgf("Failed to convert json to yaml %+v", jsonDat)
		return nil, marshalErr
	}
	log.Trace().Msgf("Policy message: %s", string(dat))
	hashBytes, hashErr := goObjectHash.ObjectHash(dat)
	if hashErr != nil {
		log.Warn().Err(hashErr).Msgf("Failed to hash json serialized proto message %s", string(dat))
		return nil, hashErr
	}
	hash := base64.StdEncoding.EncodeToString(hashBytes[:])

	return &policysyncv1.PolicyWrapper{
		Policy: policyMessage,
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName: policyName,
			PolicyHash: hash,
		},
		PolicyMetadata: metadata,
	}, nil
}
