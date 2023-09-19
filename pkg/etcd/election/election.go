package election

import (
	"context"
	"sync/atomic"

	etcd "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/aperture/v2/pkg/utils"
	concurrencyv3 "go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/fx"
)

// ElectionIn holds parameters for ProvideElection.
type ElectionIn struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Session    *etcd.Session
}

// ProvideElection provides a wrapper around etcd based leader election for arbitrary key.
//
// Note: This is not exposed directly by any module – controller and agent have
// their own leader election fx Options.
func ProvideElection(key string, in ElectionIn) *Election {
	ctx, cancel := context.WithCancel(context.Background())

	election := &Election{
		doneChan: make(chan struct{}),
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// A goroutine to do leader election
			panichandler.Go(func() {
				defer close(election.doneChan)

				session, err := in.Session.WaitSession(ctx)
				if err != nil {
					log.Error().Err(err).Msg("Failed to get etcd session for leader election")
					utils.Shutdown(in.Shutdowner)
					return
				}

				// Create an election for this client
				election.election = concurrencyv3.NewElection(session, key)
				// Campaign for leadership
				err = election.election.Campaign(ctx, info.GetHostInfo().Uuid)
				if err != nil {
					log.Error().Err(err).Msg("Unable to elect a leader")
					utils.Shutdown(in.Shutdowner)
					return
				}
				// Check if canceled
				if ctx.Err() != nil {
					return
				}
				// This is the leader
				election.isLeader.Store(true)
				log.Info().Msg("Node is now a leader")
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
				if election.election == nil {
					return nil
				}
				err = election.election.Resign(stopCtx)
				if err != nil {
					log.Error().Err(err).Msg("Unable to resign from the election")
				}
			}
			return err
		},
	})

	return election
}

// Election is a wrapper around etcd election.
type Election struct {
	election *concurrencyv3.Election
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
