package alerts

// Alerter is responsible for receiving alerts and propagating them to the channel
// returned by AlertsChan().
type Alerter interface {
	AddAlert(*Alert)
	AlertsChan() <-chan *Alert
}
