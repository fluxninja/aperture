package concurrency

import (
	"context"
	"errors"
	"fmt"
	"path"
	"sync"

	"go.uber.org/fx"

	configv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/common/config/v1"
	policydecisionsv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	"github.com/FluxNinja/aperture/pkg/config"
	etcdclient "github.com/FluxNinja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/FluxNinja/aperture/pkg/etcd/watcher"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/notifiers"
	"github.com/FluxNinja/aperture/pkg/paths"
)

type autoTokensFactory struct {
	tokensDecisionWatcher notifiers.Watcher
}

func newAutoTokensFactory(lifecycle fx.Lifecycle, client *etcdclient.Client) (*autoTokensFactory, error) {
	var err error
	etcdPath := path.Join(paths.AutoTokenResultsPath)
	tokensDecisionWatcher, err := etcdwatcher.NewWatcher(client, etcdPath)
	if err != nil {
		return nil, err
	}
	autoTokensFactory := &autoTokensFactory{
		tokensDecisionWatcher: tokensDecisionWatcher,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := tokensDecisionWatcher.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			err := tokensDecisionWatcher.Stop()
			if err != nil {
				return err
			}
			return nil
		},
	})

	return autoTokensFactory, nil
}

// newAutoTokens provides creates new autoTokens.
func (atFactory *autoTokensFactory) newAutoTokens(
	agentGroupName,
	policyName,
	policyHash string,
	lc fx.Lifecycle,
	componentIdx int64,
) (*autoTokens, error) {
	at := &autoTokens{
		avgTokensForPolicy: 1,
		tokensDecision: &policydecisionsv1.TokensDecision{
			DefaultWorkloadTokens: 1,
			TokensByWorkload:      make(map[string]uint64),
		},
		agentGroup:            agentGroupName,
		policyName:            policyName,
		policyHash:            policyHash,
		tokensDecisionWatcher: atFactory.tokensDecisionWatcher,
		componentIdx:          componentIdx,
	}

	unmarshaller, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create protobuf unmarshaller")
		return nil, err
	}

	tokensNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(paths.IdentifierForComponent(at.agentGroup, at.policyName, at.componentIdx)),
		unmarshaller,
		at.tokenUpdateCallback,
	)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := at.tokensDecisionWatcher.AddKeyNotifier(tokensNotifier)
			return err
		},
		OnStop: func(context.Context) error {
			err := at.tokensDecisionWatcher.RemoveKeyNotifier(tokensNotifier)
			return err
		},
	})

	return at, nil
}

// autoTokens struct tokens per workload.
type autoTokens struct {
	mutex                 sync.RWMutex
	tokensDecision        *policydecisionsv1.TokensDecision
	tokensDecisionWatcher notifiers.Watcher
	agentGroup            string
	policyName            string
	policyHash            string
	avgTokensForPolicy    uint64
	componentIdx          int64
}

func (at *autoTokens) tokenUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	at.mutex.Lock()
	defer at.mutex.Unlock()
	if event.Type == notifiers.Remove {
		log.Trace().Msg("Tokens were removed")
		return
	}

	var wrapperMessage configv1.ConfigPropertiesWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.Config == nil {
		log.Error().Err(err).Msg("Failed to unmarshal config wrapper")
		return
	}

	if wrapperMessage.AgentGroupName != at.agentGroup {
		log.Trace().Msg("Tokens not updated - agent group mismatch")
		return
	}

	// check if this token value is for the same policy id as what we have
	if wrapperMessage.PolicyName != at.policyName || wrapperMessage.PolicyHash != at.policyHash {
		err = errors.New("policy id mismatch")
		statusMsg := fmt.Sprintf("Expected policy: %s, %s, Got: %s, %s",
			at.policyName, at.policyHash, wrapperMessage.PolicyName, wrapperMessage.PolicyHash)
		log.Warn().Err(err).Msg(statusMsg)
		return
	}

	tokensDecision := &policydecisionsv1.TokensDecision{}
	err = wrapperMessage.Config.UnmarshalTo(tokensDecision)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal avg tokens value")
		return
	}
	at.tokensDecision = tokensDecision
}

// GetTokensForDefaultWorkload returns average tokens for default workload.
func (at *autoTokens) GetTokensForDefaultWorkload() uint64 {
	at.mutex.RLock()
	defer at.mutex.RUnlock()

	return at.avgTokensForPolicy
}

// GetTokensForWorkload returns tokens per workflow.
func (at *autoTokens) GetTokensForWorkload(workload *policydecisionsv1.WorkloadDesc) uint64 {
	at.mutex.RLock()
	defer at.mutex.RUnlock()

	if val, ok := at.tokensDecision.TokensByWorkload[workload.String()]; ok {
		return val
	}
	return at.tokensDecision.DefaultWorkloadTokens
}
