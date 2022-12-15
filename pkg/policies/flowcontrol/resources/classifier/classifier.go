package classifier

import (
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/prometheus/client_golang/prometheus"
)

// Classifier implements iface.Classifier interface.
type Classifier struct {
	iface.Classifier
	counter         prometheus.Counter
	classifierProto *policylangv1.Classifier
	classifierID    iface.ClassifierID
}

// GetFlowSelector returns the flow selector.
func (c *Classifier) GetFlowSelector() *policylangv1.FlowSelector {
	if c.classifierProto != nil {
		return c.classifierProto.GetFlowSelector()
	}
	return nil
}

// GetClassifierID returns ClassifierID object that should uniquely identify classifier.
func (c *Classifier) GetClassifierID() iface.ClassifierID {
	return c.classifierID
}

// GetRequestCounter returns the counter for the classifier.
func (c *Classifier) GetRequestCounter() prometheus.Counter {
	return c.counter
}
