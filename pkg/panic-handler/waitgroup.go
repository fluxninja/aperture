package panichandler

import (
	"context"
	"sync"

	"go.uber.org/fx"
)

// WaitGroup allows spawning and waiting for goroutines.
//
// (This is similar to https://pkg.go.dev/github.com/sourcegraph/conc#WaitGroup,
// but using the panichandler.Go to spawn goroutines).
//
// The zero value of WaitGroup is usable, just like sync.WaitGroup.
type WaitGroup struct{ rawWG sync.WaitGroup }

// Go spawns a new goroutine in the WaitGroup using panichandler.Go.
func (wg *WaitGroup) Go(f func()) {
	wg.rawWG.Add(1)
	Go(func() {
		defer wg.rawWG.Done()
		f()
	})
}

// Wait waits for all previously spawned goroutines.
func (wg *WaitGroup) Wait() { wg.rawWG.Wait() }

// CancellableWaitGroup combines a WaitGroup with a cancellable context.
//
// The zero value of CancellableWaitGroup is not usable.
type CancellableWaitGroup struct {
	wg     WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCancellableWaitGroup creates a new CancellableWaitGroup.
func NewCancellableWaitGroup() *CancellableWaitGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &CancellableWaitGroup{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Ctx returns the context that will be canceled on Cancel.
func (cwg *CancellableWaitGroup) Ctx() context.Context { return cwg.ctx }

// Go spawns a new goroutine in the WaitGroup using panichandler.Go.
func (cwg *CancellableWaitGroup) Go(f func(context.Context)) {
	cwg.wg.Go(func() { f(cwg.ctx) })
}

// CancelAndWait cancels all the running goroutines and waits for them to finish.
func (cwg *CancellableWaitGroup) CancelAndWait() {
	cwg.cancel()
	cwg.wg.Wait()
}

// GoOnDone registers a callback to be called when the given channel is closed.
// Waiting can be canceled using CancelAndWait.
func (cwg *CancellableWaitGroup) GoOnDone(doneChan <-chan struct{}, onDone func(context.Context)) {
	cwg.Go(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case <-doneChan:
			onDone(ctx)
		}
	})
}

// FxOnDone is a helper that wraps NewCancellableWaitGroup + GoOnDone in an fx lifecycle.
// It registers a callback to be called when the given channel is closed.
//
// The callback will be called only within period between Start and Stop (if
// channel is closed before Start, callback will be delayed to until Start).
func FxOnDone(doneChan <-chan struct{}, onDone func(context.Context), lifecycle fx.Lifecycle) {
	var wg *CancellableWaitGroup
	lifecycle.Append(fx.StartStopHook(
		func() {
			wg = NewCancellableWaitGroup()
			wg.GoOnDone(doneChan, onDone)
		},
		func() { wg.CancelAndWait() },
	))
}
