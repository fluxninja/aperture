package panichandler

import (
	"runtime"
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

// RegisterPanicHandler calls global registry's internal register panic handler function to panic handler registry.
func RegisterPanicHandler(ph PanicHandler) {
	getPanicHandlerRegistry().RegisterPanicHandler(ph)
}

// Recover invokes each of the registered panic handlers, and then rethrows the panic.
func Recover() {
	getPanicHandlerRegistry().RecoverAndCrash()
}

// crashOnce prevents multiple panics to interfere each other, process single panic.
var crashOnce = sync.Once{}

// Callstack is a full stacktrace.
type Callstack []uintptr

const stackLimit = 50

// Capture returns a full stacktrace.
func Capture() Callstack {
	callers := make([]uintptr, stackLimit)
	count := runtime.Callers(2, callers)
	stack := callers[:count]
	return Callstack(stack)
}

// RegisterPanicHandler appends panic handler to list of global registry's panic handler.
func (r *PanicHandlerRegistry) RegisterPanicHandler(ph PanicHandler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.Handlers == nil {
		r.Handlers = make([]PanicHandler, 0)
	}
	r.Handlers = append(r.Handlers, ph)
}

// Crash invokes each of the registered panic handler and then rethrows panic to shut down the app.
func (r *PanicHandlerRegistry) Crash(e interface{}) {
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
				h(e, stackTrace)
			}()
		}
		wg.Wait()
		close(waitCh)
	})

	select {
	case <-waitCh:
		panic(e)
	case <-time.After(5 * time.Second):
		panic(e)
	}
}

// RecoverAndCrash invokes each of the registered panic handler and then shuts down the app.
func (r *PanicHandlerRegistry) RecoverAndCrash() {
	if v := recover(); v != nil {
    // Note: DiodeFlush -- Check
		time.Sleep(time.Second * 1)
    r.Crash(v)
  }
}

// Go calls f on a new go-routine, reporting panics to the registered handlers.
func (r *PanicHandlerRegistry) Go(f func()) {
	go func() {
		// SetPanicOnFault allows the runtime trigger only a panic, not a crash
		debug.SetPanicOnFault(true)
		defer r.RecoverAndCrash()
		f()
	}()
}

// PanicHandlerOut enables registering panic handlers via Fx.
type PanicHandlerOut struct {
	fx.Out
	PanicHandler PanicHandler `group:"panic-handlers"`
}
