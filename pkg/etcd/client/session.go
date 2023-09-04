package etcd

import (
	"context"
	"errors"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	concurrencyv3 "go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/log"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// Session wraps concurrencyv3.Session.
//
// Session may be available at OnStart, but that's not guaranteed.
// Session will shutdown the app when it expires.
type Session struct {
	session      *concurrencyv3.Session // should be read only after creationDone is closed
	creationDone chan struct{}          // when closed, session is either non-nil, or nil (failed to create)
}

// SessionScopedKV implements clientv3.KV by attaching the session's lease to
// all Put requests, effectively scoping all created keys to the session.
type SessionScopedKV struct {
	KVWrapper
}

// WaitSession waits (up to context deadline) for etcd session to be established or errors.
func (s *Session) WaitSession(ctx context.Context) (*concurrencyv3.Session, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.creationDone:
		if s.session == nil {
			return nil, errNoSessionFailed
		}
		return s.session, nil
	}
}

var errNoSessionFailed = errors.New("etcd session unavailable: failed to establish")

// ProvideSession provides Session.
func ProvideSession(
	client *Client,
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
) (*Session, error) {
	sessionWrapper := Session{
		creationDone: make(chan struct{}),
	}

	lc.Append(fx.StartHook(func() {
		// We don't join this goroutine anywhere, as it should finish when
		// client's context is canceled anyway.
		panichandler.Go(func() {
			session, err := createSession(client.Client, int(client.leaseTTL.AsDuration().Seconds()))
			if err != nil {
				close(sessionWrapper.creationDone)
				log.Error().Err(err).Msg("Unable to establish etcd session")
				utils.Shutdown(shutdowner)
				return
			}
			sessionWrapper.session = session
			close(sessionWrapper.creationDone)
			log.Info().Msg("etcd session established")

			// wait for session to be closed
			<-session.Done()

			select {
			case <-session.Client().Ctx().Done():
				return
			default:
				// session close not caused by client shutdown
				log.Error().Msg("etcd session expired, request shutdown")
				utils.Shutdown(shutdowner)
			}
		})
	}))

	// Try to establish session already in OnStart, but only if it takes reasonable time.
	lc.Append(fx.StartHook(func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, maxBlockingSessionWaitTime)
		defer cancel()

		select {
		case <-sessionWrapper.creationDone:
			return
		case <-ctx.Done():
			log.Warn().Msg("Establishing etcd session takes long time, putting in background")
			return
		}
	}))

	return &sessionWrapper, nil
}

func createSession(client *clientv3.Client, ttl int) (*concurrencyv3.Session, error) {
	// Normally, session uses its long-lived context to obtain its
	// lease, which in case of errors may results in "hanging" for a
	// long time without printing anything.
	// To have more chance to show errors, we first try to manually
	// grab a lease with a short timeout.
	initialLeaseGrantCtx, cancelInitialCtx := context.WithTimeout(client.Ctx(), 4*time.Second)
	defer cancelInitialCtx()
	lease, err := client.Lease.Grant(initialLeaseGrantCtx, int64(ttl))
	if err != nil {
		if client.Ctx().Err() != nil {
			// Client shutdown, don't retry.
			return nil, err
		}

		// Actual error will be printed out by etcd as a warning.
		log.Error().Err(err).Msg("Initial attempt to establish etcd session failed, retrying")
		return concurrencyv3.NewSession(client, concurrencyv3.WithTTL(ttl))
	}

	return concurrencyv3.NewSession(
		client,
		concurrencyv3.WithTTL(ttl),
		concurrencyv3.WithLease(lease.ID),
	)
}

// Max time to wait on session creation, before it's deferred to background.
const maxBlockingSessionWaitTime = 600 * time.Millisecond

// ProvideSessionScopedKV provides SessionScopedKV.
//
// Note: This requires Session, so any usage of SessionScopedKV will cause app to shut down.
func ProvideSessionScopedKV(client *Client, session *Session, lc fx.Lifecycle) *SessionScopedKV {
	var scopedKV SessionScopedKV
	lc.Append(fx.StartHook(func() {
		scopedKV.KV = kvWithLease{
			rawKV:   client.KV,
			session: session,
		}
	}))
	return &scopedKV
}
