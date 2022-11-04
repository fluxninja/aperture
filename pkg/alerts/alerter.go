package alerts

// Alerter is responsible for receiving alerts and propagating them to the channel
// returned by AlertsChan().
type Alerter interface {
	AddAlert(*Alert)
	AlertsChan() <-chan *Alert
}

// SimpleAlerter implements Alerter interface. It just simple propagates alerts
// to the channel.
type SimpleAlerter struct {
	alertsCh chan *Alert
}

// NewSimpleAlerter returns new instance of SimpleAlerter with channel of given size.
func NewSimpleAlerter(channelSize int) *SimpleAlerter {
	return &SimpleAlerter{
		alertsCh: make(chan *Alert, channelSize),
	}
}

// AddAlert adds alert to the channel.
func (a *SimpleAlerter) AddAlert(alert *Alert) {
	a.alertsCh <- alert
}

// AlertsChan returns the alerts channel.
func (a *SimpleAlerter) AlertsChan() <-chan *Alert {
	return a.alertsCh
}
