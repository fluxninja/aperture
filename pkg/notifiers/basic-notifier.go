package notifiers

// NotifyFunc is a signature for basic notifier function.
type NotifyFunc func(Event)

// BasicKeyNotifier holds fields for basic key notifier.
type BasicKeyNotifier struct {
	NotifyFunc NotifyFunc
	KeyNotifierBase
}

// Make sure BasicKeyNotifier implements KeyNotifier.
var _ KeyNotifier = (*BasicKeyNotifier)(nil)

// Notify calls the registered notifier function with the given event.
func (bfn *BasicKeyNotifier) Notify(event Event) {
	if bfn.NotifyFunc != nil {
		bfn.NotifyFunc(event)
	}
}

// BasicPrefixNotifier holds fields for basic prefix notifier.
type BasicPrefixNotifier struct {
	NotifyFunc NotifyFunc
	PrefixNotifierBase
}

// Make sure BasicPrefixNotifier implements PrefixNotifier.
var _ PrefixNotifier = (*BasicPrefixNotifier)(nil)

// GetKeyNotifier returns a basic key notifier for the given key.
func (bdn *BasicPrefixNotifier) GetKeyNotifier(key Key) KeyNotifier {
	notifier := BasicKeyNotifier{
		NotifyFunc: bdn.NotifyFunc,
	}
	return &notifier
}
