package main

import (
	"context"
	"errors"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/etcd/election"
)

// enforceSingleControllerModule enforces that only a single controller runs
// concurrently within a single project.
//
// This is needed because various components of controller are writing to etcd
// assuming they're the only writer for a given prefix.
//
// If the previous controller has exited cleanly, we should obtain leadership
// immediately.  Otherwise (if in a pathological case the previous controller
// is still running or has crashed), we'll wait until its lease expires
// (potentially, entering crash loop).
//
// Enforcing single controller is done using etcd election, which is backed by
// the same session as SessionScopedKV. This guarantees that when we lose
// leadership because of session expiration, we won't overwrite entries written
// by the new leader (as our session would expire).
//
// Note: Would be nice to guarantee via types that controller doesn't actually
// use SessionScopedKV before leader is elected, but right now we just rely on
// this hook to run before controlplane.setup, which should be good enough to
// prevent controller writing anything before becoming a leader.
var enforceSingleControllerModule = fx.Options(
	fx.Provide(fx.Annotate(
		provideControllerElection,
		fx.ResultTags(config.NameTag("controller")),
	)),
	fx.Invoke(fx.Annotate(
		enforceSingleController,
		fx.ParamTags(config.NameTag("controller")),
	)),
)

func provideControllerElection(in election.ElectionIn) *election.Election {
	return election.ProvideElection("/controller-election", in)
}

func enforceSingleController(election *election.Election, lc fx.Lifecycle) {
	lc.Append(fx.StartHook(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-election.Done():
			if election.IsLeader() {
				return nil
			} else {
				return errors.New("failed to become a leader")
			}
		}
	}))
}
