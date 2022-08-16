package panichandler

import (
	"runtime/debug"
	"sync"
	"time"

	"go.uber.org/fx"
)

var globalPanicHandlerRegistry *PanicHandlerRegistry

func init() {
	globalPanicHandlerRegistry = getPanicHandlerRegistry()
}

// PanicHandlerRegistry defines a list of panic handlers.
type PanicHandlerRegistry struct {
	mutex    sync.RWMutex
	Handlers []PanicHandler
}

// PanicHandlerRegistryIn holds parameters, list of panic handlers, for RegisterPanicHandlers.
type PanicHandlerRegistryIn struct {
	fx.In
	Handlers []PanicHandler `group:"panic-handlers"`
}

// RegisterPanicHandlers register panic handlers to panic handler registry.
func RegisterPanicHandlers(in PanicHandlerRegistryIn) {
	// loop the handlers
	for _, handler := range in.Handlers {
		RegisterPanicHandler(handler)
	}
}

// PanicHandler is a panic handling function that is called when a panic occurs with full stacktrace.
type PanicHandler func(interface{}, Callstack)

// getPanicHandlerRegistry returns the global panic handler registry.
func getPanicHandlerRegistry() *PanicHandlerRegistry {
	if globalPanicHandlerRegistry == nil {
		globalPanicHandlerRegistry = &PanicHandlerRegistry{
			Handlers: make([]PanicHandler, 0),
		}
	}
	return globalPanicHandlerRegistry
}

// Go calls registry's internal Go function to start f on a new go routine.
func Go(f func()) {
	getPanicHandlerRegistry().Go(f)
}

// Crash calls registry's internal Crash function to invoke registered panic handler.
func Crash(v interface{}) {
	getPanicHandlerRegistry().Crash(v)
}

// RegisterPanicHandler calls global registry's internal register panic handler function to panic handler registry.
func RegisterPanicHandler(ph PanicHandler) {
	getPanicHandlerRegistry().RegisterPanicHandler(ph)
}

// crashOnce prevents multiple panics to interfere each other, process single panic.
var crashOnce = sync.Once{}

// RegisterPanicHandler appends panic handler to list of global registry's panic handler.
func (r *PanicHandlerRegistry) RegisterPanicHandler(ph PanicHandler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.Handlers == nil {
		r.Handlers = make([]PanicHandler, 0)
	}
	r.Handlers = append(r.Handlers, ph)
}

// Crash invokes each of the registered panic handler and then rethrows panic - shutting down the app.
func (r *PanicHandlerRegistry) Crash(v interface{}) {
	stackTrace := Capture()
	waitCh := make(chan struct{})

	crashOnce.Do(func() {
		r.mutex.RLock()
		defer r.mutex.RUnlock()
		wg := sync.WaitGroup{}
		wg.Add(len(r.Handlers))

		for _, handler := range r.Handlers {
			h := handler
			go func() {
				defer wg.Done()
				h(v, stackTrace)
			}()
		}
		wg.Wait()
		close(waitCh)
	})

	select {
	case <-waitCh:
	case <-time.After(5 * time.Second):
	}
	panic(v)
}

// Go calls f on a new go-routine, reporting panics to the registered handlers.
func (r *PanicHandlerRegistry) Go(f func()) {
	go func() {
		// SetPanicOnFault allows the runtime trigger only a panic.
		debug.SetPanicOnFault(true)
		defer func() {
			if v := recover(); v != nil {
				r.Crash(v)
			}
		}()
		f()
	}()
}

// PanicHandlerOut enables registering panic handlers via Fx.
type PanicHandlerOut struct {
	fx.Out
	PanicHandler PanicHandler `group:"panic-handlers"`
}
