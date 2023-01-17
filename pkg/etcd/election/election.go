package election

import (
	"context"
	"sync"
	"time"

	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcd "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/panichandler"
	concurrencyv3 "go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/fx"
)

var (
	// FxTagBase is the tag's base used to identify the election result Tracker.
	FxTagBase = "etcd_election"
	// FxTag is the tag used to identify the election result Tracker.
	FxTag = config.NameTag(FxTagBase)
	// ElectionResultKey is the key used to identify the election result in the election Tracker.
	ElectionResultKey = notifiers.Key("election_result")
)

// Module is a fx module that provides etcd based leader election per agent group.
func Module() fx.Option {
	return fx.Options(
		notifiers.TrackersConstructor{Name: FxTagBase}.Annotate(),
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
	Trackers   notifiers.Trackers `name:"etcd_election"`
}

// ProvideElection provides a wrapper around etcd based leader election.
func ProvideElection(in ElectionIn) (*Election, error) {
	ctx, cancel := context.WithCancel(context.Background())

	election := &Election{}

	var waitGroup sync.WaitGroup

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Create an election for this client
			election.Election = concurrencyv3.NewElection(in.Client.Session, "/election/"+in.AgentInfo.GetAgentGroup())
			waitGroup.Add(1)
			// A goroutine to do leader election
			panichandler.Go(func() {
				defer waitGroup.Done()
				// Campaign for leadership
				err := election.Election.Campaign(ctx, info.GetHostInfo().Uuid)
				if err != nil {
					log.Error().Err(err).Msg("Unable to elect a leader")
					shutdownErr := in.Shutdowner.Shutdown()
					if shutdownErr != nil {
						log.Error().Err(shutdownErr).Msg("Error on invoking shutdown")
						return
					}
				}
				// Check if canceled
				if ctx.Err() != nil {
					return
				}
				// This is the leader
				election.isLeader = true
				log.Info().Msg("Propagate Election result")
				in.Trackers.WriteEvent(ElectionResultKey, []byte("true"))
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			var err error
			cancel()
			// resign from the election if we are the leader
			if election.IsLeader() {
				stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
				err = election.Election.Resign(stopCtx)
				stopCancel()
				if err != nil {
					log.Error().Err(err).Msg("Unable to resign from the election")
				}
			}
			// Wait for the election goroutine to finish
			waitGroup.Wait()
			return err
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
