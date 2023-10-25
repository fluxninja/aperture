// Leader-only-receiver wraps any metrics receiver and starts it only when agent is a leader.
package leaderonlyreceiver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	etcd "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
)

const (
	type_            component.Type = "aperture_leader_only"
	stability                       = component.StabilityLevelDevelopment
	lateStartTimeout                = 15 * time.Second
)

// Module provides receiver factory.
func Module() fx.Option {
	return fx.Provide(
		fx.Annotate(
			NewFactory,
			fx.ResultTags(otelconsts.ReceiverFactoriesFxTag),
		),
	)
}

// NewFactory creates a new aperture_leader_only receiver factory using given leader election.
func NewFactory(etcdClient *etcd.Client) receiver.Factory {
	return receiver.NewFactory(
		type_,
		func() component.Config {
			return &Config{
				etcdClient: etcdClient,
			}
		},
		receiver.WithMetrics(createMetricsReceiver, stability))
}

// Config is a config for leader-only-receiver.
type Config struct {
	// Config for the wrapped receiver
	etcdClient *etcd.Client
	Inner      map[string]any `mapstructure:"config"`
	// Type of the wrapped receiver
	InnerType component.Type `mapstructure:"type"`
}

// Validate implements component.ConfigValidator.
func (c *Config) Validate() error {
	if c.InnerType == "" {
		return errors.New("type is required")
	}
	return nil
}

func createMetricsReceiver(
	_ context.Context,
	createSettings receiver.CreateSettings,
	rConf component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	// At this point we do not have access to Factories, so we cannot do anything with the config
	return &leaderOnlyReceiver{
		config:             *rConf.(*Config),
		consumer:           consumer,
		origCreateSettings: createSettings,
	}, nil
}

type leaderOnlyReceiver struct {
	consumer           consumer.Metrics
	factory            receiver.Factory
	host               component.Host
	inner              receiver.Metrics // nil if inner receiver not started
	origCreateSettings receiver.CreateSettings
	config             Config
}

// leaderOnlyReceiver implements LeaderWatcher interface.
var _ etcd.ElectionWatcher = (*leaderOnlyReceiver)(nil)

// Start implements component.Component.
func (r *leaderOnlyReceiver) Start(startCtx context.Context, host component.Host) error {
	factory := host.GetFactory(component.KindReceiver, r.config.InnerType)
	if factory == nil {
		return fmt.Errorf("factory for %s receiver not found", r.config.InnerType)
	}

	r.factory = factory.(receiver.Factory)
	r.host = host

	if r.config.etcdClient.IsLeader() {
		// If we already know we're the leader, we can skip creating background
		// goroutine and start inner receiver immediately.
		if err := r.startInnerReceiver(startCtx); err != nil {
			return fmt.Errorf("failed to start %s receiver: %w", r.config.InnerType, err)
		}
		return nil
	}

	r.config.etcdClient.AddElectionWatcher(r)

	return nil
}

// Shutdown implements component.Component.
func (r *leaderOnlyReceiver) Shutdown(ctx context.Context) error {
	r.config.etcdClient.RemoveElectionWatcher(r)
	if r.inner != nil {
		return r.inner.Shutdown(ctx)
	}
	return nil
}

// OnLeaderStart starts the inner receiver.
func (r *leaderOnlyReceiver) OnLeaderStart() {
	startCtx, cancel := context.WithTimeout(context.Background(), lateStartTimeout)
	defer cancel()
	if err := r.startInnerReceiver(startCtx); err != nil {
		r.host.ReportFatalError(fmt.Errorf(
			"failed to start %s receiver after becoming a leader: %w",
			r.config.InnerType, err,
		))
	}
}

// OnLeaderStop stops the inner receiver.
func (r *leaderOnlyReceiver) OnLeaderStop() {
	if r.inner != nil {
		err := r.inner.Shutdown(context.Background())
		if err != nil {
			r.host.ReportFatalError(fmt.Errorf(
				"failed to stop %s receiver after becoming a follower: %w",
				r.config.InnerType, err,
			))
		}
	}
}

func (r *leaderOnlyReceiver) startInnerReceiver(ctx context.Context) error {
	cfg := r.factory.CreateDefaultConfig()
	if err := component.UnmarshalConfig(confmap.NewFromStringMap(r.config.Inner), cfg); err != nil {
		return fmt.Errorf("error reading configuration: %w", err)
	}

	if err := component.ValidateConfig(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Setting ID for inner receiver to: "<innerType>/aperture_leader_only/<name>"
	set := r.origCreateSettings
	set.ID = component.NewIDWithName(r.config.InnerType, r.origCreateSettings.ID.String())
	inner, err := r.factory.CreateMetricsReceiver(ctx, set, cfg, r.consumer)
	if err != nil {
		return fmt.Errorf("error creating receiver: %w", err)
	}

	if err := inner.Start(ctx, r.host); err != nil {
		return err
	}

	r.inner = inner
	return nil
}
