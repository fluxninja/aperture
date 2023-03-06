package notifiers

// NotifyFunc is a signature for basic notifier function.
type NotifyFunc func(Event)

// BasicKeyNotifier holds fields for basic key notifier.
type BasicKeyNotifier struct {
	NotifyFunc NotifyFunc
	KeyBase
}

// Make sure BasicKeyNotifier implements KeyNotifier.
var _ KeyNotifier = (*BasicKeyNotifier)(nil)

// NewBasicKeyNotifier returns a new basic key notifier.
func NewBasicKeyNotifier(key Key, notifyFunc NotifyFunc) *BasicKeyNotifier {
	notifier := &BasicKeyNotifier{
		KeyBase:    NewKeyBase(key),
		NotifyFunc: notifyFunc,
	}
	return notifier
}

// Notify calls the registered notifier function with the given event.
func (bfn *BasicKeyNotifier) Notify(event Event) {
	if bfn.NotifyFunc != nil {
		bfn.NotifyFunc(event)
	}
}

// basicPrefixNotifier holds fields for basic prefix notifier.
type basicPrefixNotifier struct {
	NotifyFunc NotifyFunc
	PrefixBase
}

// Make sure BasicPrefixNotifier implements PrefixNotifier.
var _ PrefixNotifier = (*basicPrefixNotifier)(nil)

// GetKeyNotifier returns a basic key notifier for the given key.
func (bdn *basicPrefixNotifier) GetKeyNotifier(key Key) (KeyNotifier, error) {
	return NewBasicKeyNotifier(key, bdn.NotifyFunc), nil
}

// NewBasicPrefixNotifier returns a new basic prefix notifier.
func NewBasicPrefixNotifier(prefix string, notifyFunc NotifyFunc) PrefixNotifier {
	notifier := &basicPrefixNotifier{
		PrefixBase: NewPrefixBase(prefix),
		NotifyFunc: notifyFunc,
	}
	return notifier
}
