package election

import (
	"context"
	"sync/atomic"

	"github.com/fluxninja/aperture/v2/pkg/agentinfo"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcd "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/panichandler"
	"github.com/fluxninja/aperture/v2/pkg/utils"
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

	election := &Election{
		doneChan: make(chan struct{}),
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Create an election for this client
			election.Election = concurrencyv3.NewElection(in.Client.Session, "/election/"+in.AgentInfo.GetAgentGroup())
			// A goroutine to do leader election
			panichandler.Go(func() {
				defer close(election.doneChan)
				// Campaign for leadership
				err := election.Election.Campaign(ctx, info.GetHostInfo().Uuid)
				if err != nil {
					log.Error().Err(err).Msg("Unable to elect a leader")
					utils.Shutdown(in.Shutdowner)
				}
				// Check if canceled
				if ctx.Err() != nil {
					return
				}
				// This is the leader
				election.isLeader.Store(true)
				log.Info().Msg("Node is now a leader")
				in.Trackers.WriteEvent(ElectionResultKey, []byte("true"))
			})

			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			var err error
			cancel()
			// Wait for the election goroutine to finish
			<-election.Done()
			// resign from the election if we are the leader
			if election.IsLeader() {
				election.isLeader.Store(false)
				err = election.Election.Resign(stopCtx)
				if err != nil {
					log.Error().Err(err).Msg("Unable to resign from the election")
				}
			}
			return err
		},
	})

	return election, nil
}

// Election is a wrapper around etcd election.
type Election struct {
	Election *concurrencyv3.Election
	isLeader atomic.Bool

	// When closed, leader election has stopped (either due to becoming the
	// leader or due to cancellation).
	// Note: chan is used here instead of WaitGroup, so that calls to
	// WaitUntilLeader done before election is started do not immediately return.
	doneChan chan struct{}
}

// IsLeader returns true if the current node is the leader.
func (e *Election) IsLeader() bool {
	return e.isLeader.Load()
}

// Done returns a channel that could be used to wait for election results.
//
// When the channel is closed then either:
// * Node became a leader (IsLeader() == true),
// * Leader election was canceled (IsLeader() == false).
func (e *Election) Done() <-chan struct{} { return e.doneChan }
