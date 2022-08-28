package concurrency

import (
	"context"
	"errors"
	"fmt"
	"path"
	"sync"

	"go.uber.org/fx"

	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
	policydecisionsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/paths"
)

type autoTokensFactory struct {
	tokensDecisionWatcher notifiers.Watcher
	agentGroup            string
}

func newAutoTokensFactory(
	lifecycle fx.Lifecycle,
	client *etcdclient.Client,
	agentGroup string,
) (*autoTokensFactory, error) {
	var err error
	etcdPath := path.Join(paths.AutoTokenResultsPath, paths.AgentGroupPrefix(agentGroup))
	tokensDecisionWatcher, err := etcdwatcher.NewWatcher(client, etcdPath)
	if err != nil {
		return nil, err
	}
	autoTokensFactory := &autoTokensFactory{
		tokensDecisionWatcher: tokensDecisionWatcher,
		agentGroup:            agentGroup,
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

// newAutoTokens creates new autoTokens.
func (atFactory *autoTokensFactory) newAutoTokens(
	policyName,
	policyHash string,
	lc fx.Lifecycle,
	componentIdx int64,
) (*autoTokens, error) {
	at := &autoTokens{
		tokensDecision: &policydecisionsv1.TokensDecision{
			TokensByWorkloadIndex: make(map[string]uint64),
		},
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
		notifiers.Key(paths.DataplaneComponentKey(atFactory.agentGroup, at.policyName, at.componentIdx)),
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
	policyName            string
	policyHash            string
	componentIdx          int64
}

func (at *autoTokens) tokenUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	at.mutex.Lock()
	defer at.mutex.Unlock()
	if event.Type == notifiers.Remove {
		log.Trace().Msg("Tokens were removed")
		return
	}

	var wrapperMessage configv1.TokensDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.TokensDecision == nil {
		log.Error().Err(err).Msg("Failed to unmarshal config wrapper")
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

	at.tokensDecision = wrapperMessage.TokensDecision
}

// GetTokensForWorkload returns tokens per workflow.
func (at *autoTokens) GetTokensForWorkload(workloadIndex string) (uint64, bool) {
	at.mutex.RLock()
	defer at.mutex.RUnlock()
	val, ok := at.tokensDecision.TokensByWorkloadIndex[workloadIndex]
	return val, ok
}
