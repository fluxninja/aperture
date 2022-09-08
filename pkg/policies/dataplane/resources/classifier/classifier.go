package classifier

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	classificationv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/resources/classifier/compiler"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
)

type multiMatcherByControlPoint map[selectors.ControlPointID]*multimatcher.MultiMatcher[int, []*compiler.Labeler]

// rules is a helper struct to keep both compiled and uncompiled sets of rules in sync.
type rules struct {
	// rules compiled to map from ControlPointID to MultiMatcher
	MultiMatcherByControlPointID multiMatcherByControlPoint
	// non-compiled version of rules, used for reporting
	ReportedRules []compiler.ReportedRule
}

// Classifier receives classification policies and provides Classify method.
type Classifier struct {
	mu              sync.Mutex
	activeRules     atomic.Value
	activeRulesets  map[rulesetID]compiler.CompiledRuleset
	classifierProto *classificationv1.Classifier
	policyName      string
	policyHash      string
	nextRulesetID   rulesetID
	classifierIndex int64
}

type rulesetID = uint64

// FlowLabels is a map from flow labels to their values.
type FlowLabels map[string]FlowLabelValue

// FlowLabelValue is a value of a flow label with additional metadata.
type FlowLabelValue struct {
	Value string
	Flags compiler.LabelFlags
}

// ToPlainMap returns flow labels as normal map[string]string.
func (fl FlowLabels) ToPlainMap() map[string]string {
	plainMap := make(map[string]string, len(fl))
	for key, val := range fl {
		plainMap[key] = val.Value
	}
	return plainMap
}

// New creates a new Flow Classifier.
func New() *Classifier {
	return &Classifier{
		activeRulesets: make(map[rulesetID]compiler.CompiledRuleset),
	}
}

func populateFlowLabels(ctx context.Context, flowLabels FlowLabels, mm *multimatcher.MultiMatcher[int, []*compiler.Labeler], labelsForMatching selectors.Labels, input ast.Value) {
	for _, query := range mm.Match(labelsForMatching.ToPlainMap()) {
		resultSet, err := query.Query.Eval(ctx, rego.EvalParsedInput(input))
		if err != nil {
			log.Warn().Msg("Rego: Evaluation failed")
			continue
		}

		if len(resultSet) == 0 {
			log.Warn().Msg("Rego: Empty resultSet")
			continue
		}

		if len(resultSet) > 1 {
			log.Warn().Msg("Rego: Ambiguous resultSet")
		}

		if len(resultSet[0].Expressions) != 1 {
			log.Warn().Msg("Rego: Expected exactly one expression")
			continue
		}

		if query.LabelName != "" {
			// single-label-query
			flowLabels[query.LabelName] = FlowLabelValue{
				Value: resultSet[0].Expressions[0].String(),
				Flags: query.LabelFlags,
			}
		} else {
			// multi-label-query
			variables, isMap := resultSet[0].Expressions[0].Value.(map[string]interface{})
			if !isMap {
				log.Error().Msg("Rego: Expression's not a map (bug)")
				continue
			}

			for key, value := range variables {
				flowLabels[key] = FlowLabelValue{
					Value: fmt.Sprint(value),
					Flags: query.LabelsFlags[key],
				}
			}
		}
	}
}

// Classify takes rego input, performs classification, and returns a map of flow labels.
// LabelsForMatching are additional labels to use for selector matching.
func (c *Classifier) Classify(
	ctx context.Context,
	svcs []services.ServiceID,
	labelsForMatching selectors.Labels,
	direction selectors.TrafficDirection,
	input ast.Value,
) (FlowLabels, error) {
	flowLabels := make(FlowLabels)

	r, ok := c.activeRules.Load().(rules)
	if !ok {
		return flowLabels, nil
	}

	cp := selectors.ControlPoint{
		Traffic: direction,
	}

	cpID := selectors.ControlPointID{
		ServiceID: services.ServiceID{
			Service: "",
		},
		ControlPoint: cp,
	}
	camm, ok := r.MultiMatcherByControlPointID[cpID]
	if ok {
		populateFlowLabels(ctx, flowLabels, camm, labelsForMatching, input)
	}

	// TODO (krdln): update prometheus metrics upon classification errors.

	for _, svc := range svcs {
		cpID := selectors.ControlPointID{
			ServiceID:    svc,
			ControlPoint: cp,
		}
		mm, ok := r.MultiMatcherByControlPointID[cpID]
		if !ok {
			log.Trace().Str("controlPoint", cpID.String()).Msg("No labelers for controlPoint")
			continue
		}
		populateFlowLabels(ctx, flowLabels, mm, labelsForMatching, input)
	}

	return flowLabels, nil
}

// ActiveRules returns a slice of uncompiled Rules which are currently active.
func (c *Classifier) ActiveRules() []compiler.ReportedRule {
	ac, _ := c.activeRules.Load().(rules)
	return ac.ReportedRules
}

// AddRules compiles a ruleset and adds it to the active rules
//
// # The name will be used for reporting
//
// To retract the rules, call Classifier.Drop.
func (c *Classifier) AddRules(
	ctx context.Context,
	name string,
	classifier *classificationv1.Classifier,
) (ActiveRuleset, error) {
	compiled, err := compiler.CompileRuleset(ctx, name, classifier)
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

	c.activeRulesets[id] = compiled
	c.activateRulesets()
	return ActiveRuleset{id: id, classifier: c}, nil
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
	return iface.ClassifierID{
		PolicyName:      c.policyName,
		PolicyHash:      c.policyHash,
		ClassifierIndex: c.classifierIndex,
	}
}

// ActiveRuleset represents one of currently active set of rules.
type ActiveRuleset struct {
	classifier *Classifier
	id         rulesetID
}

// Drop retracts all the rules belonging to a ruleset.
func (rs ActiveRuleset) Drop() {
	if rs.classifier == nil {
		return
	}
	c := rs.classifier
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.activeRulesets, rs.id)
	c.activateRulesets()
}

// needs to be called with activeRulesets mutex held.
func (c *Classifier) activateRulesets() {
	c.activeRules.Store(c.combineRulesets())
	log.Info().Int("rulesets", len(c.activeRulesets)).Msg("Rules updated")
}

func (c *Classifier) combineRulesets() rules {
	combined := rules{
		MultiMatcherByControlPointID: make(multiMatcherByControlPoint),
		ReportedRules:                make([]compiler.ReportedRule, 0),
	}

	// to have unique keys to AddEntry
	controlPointKeys := make(map[selectors.ControlPointID]int)

	for _, ruleset := range c.activeRulesets {
		combined.ReportedRules = append(combined.ReportedRules, ruleset.ReportedRules...)
		for _, labelerWithSelector := range ruleset.Labelers {
			mm, ok := combined.MultiMatcherByControlPointID[ruleset.ControlPointID]
			if !ok {
				mm = multimatcher.New[int, []*compiler.Labeler]()
				combined.MultiMatcherByControlPointID[ruleset.ControlPointID] = mm
			}

			matcherID := controlPointKeys[ruleset.ControlPointID]
			controlPointKeys[ruleset.ControlPointID]++

			err := mm.AddEntry(matcherID, labelerWithSelector.LabelSelector, multimatcher.Appender(labelerWithSelector.Labeler))
			if err != nil {
				log.Error().Err(err).Msg("Failed to add entry to catchall multimatcher")
				return rules{}
			}
		}
	}

	return combined
}
