package election

import (
	"context"
	"time"

	"github.com/fluxninja/aperture/pkg/agentinfo"
	etcd "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panichandler"
	concurrencyv3 "go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/fx"
)

// Module is a fx module that provides etcd based leader election per agent group.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(ProvideElection),
	)
}

// ElectionIn holds parameters for ProvideElection.
type ElectionIn struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Client     *etcd.Client
	AgentInfo  *agentinfo.AgentInfo
}

// ProvideElection provides a wrapper around etcd based leader election.
func ProvideElection(in ElectionIn) (*Election, error) {
	ctx, cancel := context.WithCancel(context.Background())

	election := &Election{}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Create an election for this client
			election.Election = concurrencyv3.NewElection(in.Client.Session, "/election/"+in.AgentInfo.GetAgentGroup())
			// A goroutine to do leader election
			panichandler.Go(func() {
				// Campaign for leadership
				err := election.Election.Campaign(ctx, info.GetHostInfo().Uuid)
				if err != nil {
					log.Error().Err(err).Msg("Unable to elect a leader")
					shutdownErr := in.Shutdowner.Shutdown()
					if shutdownErr != nil {
						log.Error().Err(shutdownErr).Msg("Error on invoking shutdown")
					}
				}
				// Check if canceled
				if ctx.Err() != nil {
					return
				}
				// This is the leader
				election.isLeader = true
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			// resign from the election if we are the leader
			if election.IsLeader() {
				stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
				err := election.Election.Resign(stopCtx)
				stopCancel()
				if err != nil {
					log.Error().Err(err).Msg("Unable to resign from the election")
				}
				return err
			}
			return nil
		},
	})

	return election, nil
}

// Election is a wrapper around etcd election.
type Election struct {
	Election *concurrencyv3.Election
	isLeader bool
}

// IsLeader returns true if the current node is the leader.
func (e *Election) IsLeader() bool {
	return e.isLeader
}
