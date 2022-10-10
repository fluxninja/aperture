package classifier

import (
	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	classificationv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/prometheus/client_golang/prometheus"
)

// Classifier implements iface.Classifier interface.
type Classifier struct {
	iface.Classifier
	classifierID    iface.ClassifierID
	classifierProto *classificationv1.Classifier
	counter         prometheus.Counter
}

// GetSelector returns the selector.
func (c *Classifier) GetSelector() *selectorv1.Selector {
	if c.classifierProto != nil {
		return c.classifierProto.GetSelector()
	}
	return nil
}

// GetClassifierID returns ClassifierID object that should uniquely identify classifier.
func (c *Classifier) GetClassifierID() iface.ClassifierID {
	return c.classifierID
}

// GetCounter returns the counter for the classifier.
func (c *Classifier) GetCounter() prometheus.Counter {
	return c.counter
}
