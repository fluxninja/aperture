package classifier

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier/compiler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/status"
)

type multiMatcherByControlPoint map[selectors.ControlPointID]*multimatcher.MultiMatcher[int, []*compiler.LabelerWithSelector]

// rules is a helper struct to keep both compiled and uncompiled sets of rules in sync.
type rules struct {
	// rules compiled to map from ControlPointID to MultiMatcher
	MultiMatcherByControlPointID multiMatcherByControlPoint
	// non-compiled version of rules, used for reporting
	ReportedRules []compiler.ReportedRule
}

// ClassificationEngine receives classification policies and provides Classify method.
type ClassificationEngine struct {
	mu                 sync.Mutex
	activeRules        atomic.Value
	registry           status.Registry
	activeRulesets     map[rulesetID]compiler.CompiledRuleset
	nextRulesetID      rulesetID
	classifierMapMutex sync.RWMutex
	classifierMap      map[iface.ClassifierID]iface.Classifier
	counterVec         *prometheus.CounterVec
}

type rulesetID = uint64

// NewClassificationEngine creates a new Flow Classifier.
func NewClassificationEngine(registry status.Registry) *ClassificationEngine {
	counterVector := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.ClassifierCounterMetricName,
		Help: "A counter measuring the number of times classifier was triggered",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ClassifierIndexLabel,
	})

	return &ClassificationEngine{
		activeRulesets: make(map[rulesetID]compiler.CompiledRuleset),
		registry:       registry,
		classifierMap:  make(map[iface.ClassifierID]iface.Classifier),
		counterVec:     counterVector,
	}
}

func (c *ClassificationEngine) populateFlowLabels(ctx context.Context,
	flowLabels flowlabel.FlowLabels,
	mm *multimatcher.MultiMatcher[int, []*compiler.LabelerWithSelector],
	labelsForMatching map[string]string,
	input ast.Value,
) (classifierMsgs []*flowcontrolv1.ClassifierInfo) {
	logger := c.registry.GetLogger()
	logSampled := logger.Sample(zerolog.Sometimes)
	appendNewClassifier := func(labelerWithSelector *compiler.LabelerWithSelector, error flowcontrolv1.ClassifierInfo_Error) {
		classifierMsgs = append(classifierMsgs, &flowcontrolv1.ClassifierInfo{
			PolicyName:      labelerWithSelector.CommonAttributes.PolicyName,
			PolicyHash:      labelerWithSelector.CommonAttributes.PolicyHash,
			ClassifierIndex: labelerWithSelector.CommonAttributes.ComponentIndex,
			LabelKey:        labelerWithSelector.Labeler.LabelName,
			Error:           error,
		})
	}

	for _, labelerWithSelector := range mm.Match(labelsForMatching) {
		labeler := labelerWithSelector.Labeler
		resultSet, err := labeler.Query.Eval(ctx, rego.EvalParsedInput(input))
		if err != nil {
			logSampled.Warn().Msg("Rego: Evaluation failed")
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_EVAL_FAILED)
			continue
		}

		if len(resultSet) == 0 {
			logSampled.Warn().Msg("Rego: Empty resultSet")
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_EMPTY_RESULTSET)
			continue
		} else if len(resultSet) > 1 {
			logSampled.Warn().Msg("Rego: Ambiguous resultSet")
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_AMBIGUOUS_RESULTSET)
			continue
		}

		if len(resultSet[0].Expressions) != 1 {
			logger.Warn().Msg("Rego: Expected exactly one expression")
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_MULTI_EXPRESSION)
			continue
		}

		if labeler.LabelName != "" {
			// single-label-query
			flowLabels[labeler.LabelName] = flowlabel.FlowLabelValue{
				Value:     resultSet[0].Expressions[0].String(),
				Telemetry: labeler.Telemetry,
			}
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_NONE)
		} else {
			// multi-label-query
			variables, isMap := resultSet[0].Expressions[0].Value.(map[string]interface{})
			if !isMap {
				logger.Error().Msg("Rego: Expression's not a map (bug)")
				appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_EXPRESSION_NOT_MAP)
				continue
			}

			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_NONE)
			for key, value := range variables {
				flowLabels[key] = flowlabel.FlowLabelValue{
					Value:     fmt.Sprint(value),
					Telemetry: labeler.LabelsTelemetry[key],
				}
			}
		}
	}
	return
}

// Classify takes rego input, performs classification, and returns a map of flow labels.
// LabelsForMatching are additional labels to use for selector matching.
func (c *ClassificationEngine) Classify(
	ctx context.Context,
	svcs []string,
	ctrlPt selectors.ControlPoint,
	labelsForMatching map[string]string,
	input ast.Value,
) ([]*flowcontrolv1.ClassifierInfo, flowlabel.FlowLabels, error) {
	logSampled := c.registry.GetLogger().Sample(zerolog.Sometimes)
	flowLabels := make(flowlabel.FlowLabels)

	r, ok := c.activeRules.Load().(rules)
	if !ok {
		return nil, flowLabels, nil
	}

	var classifierMsgs []*flowcontrolv1.ClassifierInfo

	// Catch all Service
	cpID := selectors.NewControlPointID("", ctrlPt)
	mm, ok := r.MultiMatcherByControlPointID[cpID]
	if ok {
		classifierMsgs = append(classifierMsgs, c.populateFlowLabels(ctx, flowLabels, mm, labelsForMatching, input)...)
	}

	// TODO (krdln): update prometheus metrics upon classification errors.

	// Specific Service
	for _, svc := range svcs {
		cpID := selectors.NewControlPointID(svc, ctrlPt)
		mm, ok := r.MultiMatcherByControlPointID[cpID]
		if !ok {
			logSampled.Trace().Interface("controlPointID", cpID).Msg("No labelers for controlPointID")
			continue
		}
		classifierMsgs = append(classifierMsgs, c.populateFlowLabels(ctx, flowLabels, mm, labelsForMatching, input)...)
	}

	return classifierMsgs, flowLabels, nil
}

// ActiveRules returns a slice of uncompiled Rules which are currently active.
func (c *ClassificationEngine) ActiveRules() []compiler.ReportedRule {
	ac, _ := c.activeRules.Load().(rules)
	return ac.ReportedRules
}

// AddRules compiles a ruleset and adds it to the active rules
//
// # The name will be used for reporting
//
// To retract the rules, call Classifier.Drop.
func (c *ClassificationEngine) AddRules(
	ctx context.Context,
	name string,
	classifierWrapper *policysyncv1.ClassifierWrapper,
) (ActiveRuleset, error) {
	compiledRuleset, err := compiler.CompileRuleset(ctx, name, classifierWrapper)
	if err != nil {
		return ActiveRuleset{}, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	// Why index activeRulesets via ID instead of provided name?
	// * more robust if caller provides non-unique names
	// * when modifying file, one approach would be to first unload old ruleset
	//   and load a new one – in this case duplicated name is kinda expected.
	// So the name is used only for reporting.
	id := c.nextRulesetID
	c.nextRulesetID++

	c.activeRulesets[id] = compiledRuleset
	c.activateRulesets()
	return ActiveRuleset{id: id, classificationEngine: c}, nil
}

// ActiveRuleset represents one of currently active set of rules.
type ActiveRuleset struct {
	classificationEngine *ClassificationEngine
	id                   rulesetID
}

// Drop retracts all the rules belonging to a ruleset.
func (rs ActiveRuleset) Drop() {
	if rs.classificationEngine == nil {
		return
	}
	c := rs.classificationEngine
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.activeRulesets, rs.id)
	c.activateRulesets()
}

// needs to be called with activeRulesets mutex held.
func (c *ClassificationEngine) activateRulesets() {
	logger := c.registry.GetLogger()
	c.activeRules.Store(c.combineRulesets())
	logger.Info().Int("rulesets", len(c.activeRulesets)).Msg("Rules updated")
}

func (c *ClassificationEngine) combineRulesets() rules {
	combined := rules{
		MultiMatcherByControlPointID: make(multiMatcherByControlPoint),
		ReportedRules:                make([]compiler.ReportedRule, 0),
	}

	// to have unique keys to AddEntry
	controlPointKeys := make(map[selectors.ControlPointID]int)

	for _, ruleset := range c.activeRulesets {
		combined.ReportedRules = append(combined.ReportedRules, ruleset.ReportedRules...)
		for i := range ruleset.Labelers {
			labelerWithSelector := &ruleset.Labelers[i]
			mm, ok := combined.MultiMatcherByControlPointID[ruleset.ControlPointID]
			if !ok {
				mm = multimatcher.New[int, []*compiler.LabelerWithSelector]()
				combined.MultiMatcherByControlPointID[ruleset.ControlPointID] = mm
			}

			matcherID := controlPointKeys[ruleset.ControlPointID]
			controlPointKeys[ruleset.ControlPointID]++

			err := mm.AddEntry(matcherID, labelerWithSelector.LabelSelector, multimatcher.Appender(labelerWithSelector))
			if err != nil {
				log.Error().Err(err).Msg("Failed to add entry to catchall multimatcher")
				return rules{}
			}
		}
	}

	return combined
}

// RegisterClassifier adds classifier to map.
func (c *ClassificationEngine) RegisterClassifier(classifier iface.Classifier) error {
	c.classifierMapMutex.Lock()
	defer c.classifierMapMutex.Unlock()
	if _, ok := c.classifierMap[classifier.GetClassifierID()]; !ok {
		c.classifierMap[classifier.GetClassifierID()] = classifier
	} else {
		return fmt.Errorf("classifier id already registered")
	}

	return nil
}

// UnregisterClassifier removes classifier from map.
func (c *ClassificationEngine) UnregisterClassifier(classifier iface.Classifier) error {
	c.classifierMapMutex.Lock()
	defer c.classifierMapMutex.Unlock()
	delete(c.classifierMap, classifier.GetClassifierID())

	return nil
}

// GetClassifier Lookup function for getting classifier.
func (c *ClassificationEngine) GetClassifier(classifierID iface.ClassifierID) iface.Classifier {
	c.classifierMapMutex.RLock()
	defer c.classifierMapMutex.RUnlock()
	return c.classifierMap[classifierID]
}
