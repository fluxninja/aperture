package classification

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"google.golang.org/protobuf/types/known/wrapperspb"

	classificationv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/classification/v1"
	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	"github.com/fluxninja/aperture/pkg/classification/extractors"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
)

const defaultPackageName = "fluxninja.classification.extractors"

type multiMatcherByControlPoint map[selectors.ControlPointID]*multimatcher.MultiMatcher[int, []*labeler]

// rules is a helper struct to keep both compiled and uncompiled sets of rules in sync.
type rules struct {
	// rules compiled to map from ControlPointID to MultiMatcher
	MultiMatcherByControlPointID multiMatcherByControlPoint
	// non-compiled version of rules, used for reporting
	ReportedRules []ReportedRule
}

// compiledRuleset is compiled form of Classifier proto.
type compiledRuleset struct {
	ControlPointID selectors.ControlPointID
	Labelers       []labelerWithSelector
	ReportedRules  []ReportedRule
}

type labelerWithSelector struct {
	Labeler       *labeler
	LabelSelector multimatcher.Expr
}

// labeler is used to create flow labels
//
// label can create either:
// * a single label – LabelName is non-empty or
// * multiple labels – LabelName is empty.
type labeler struct {
	// rego query that's prepared to take envoy authz request as an input.
	// Result expression should be a single value (if LabelName is set) or a
	// map[string]interface{} otherwise.
	Query rego.PreparedEvalQuery
	// flags for created flow labels:
	LabelsFlags map[string]LabelFlags // multi-label variant
	// flow label that the result should be assigned to (single-label variant)
	LabelName  string
	LabelFlags LabelFlags // single-label variant
}

// Classifier receives classification policies and provides Classify method.
type Classifier struct {
	// storing activeRules underneath
	mu             sync.Mutex
	activeRules    atomic.Value
	activeRulesets map[rulesetID]compiledRuleset // protected by mu
	nextRulesetID  rulesetID                     // protected by mu
}

type rulesetID = uint64

// ReportedRule is a rule along with its selector and label name.
type ReportedRule struct {
	Selector    *selectorv1.Selector
	Rule        *classificationv1.Rule
	RulesetName string
	LabelName   string
}

func rulesetToReportedRules(rs *classificationv1.Classifier, rulesetName string) []ReportedRule {
	out := make([]ReportedRule, 0, len(rs.Rules))
	for label, rule := range rs.Rules {
		out = append(out, ReportedRule{
			RulesetName: rulesetName,
			LabelName:   label,
			Rule:        rule,
			Selector:    rs.Selector,
		})
	}
	return out
}

// FlowLabels is a map from flow labels to their values.
type FlowLabels map[string]FlowLabelValue

// FlowLabelValue is a value of a flow label with additional metadata.
type FlowLabelValue struct {
	Value string
	Flags LabelFlags
}

// LabelFlags are flags for a flow label.
type LabelFlags struct {
	// Should the created label be applied to the whole flow (propagated in baggage)
	Propagate bool
	// Should the created flow label be hidden from telemetry
	Hidden bool
}

// ToPlainMap returns flow labels as normal map[string]string.
func (fl FlowLabels) ToPlainMap() map[string]string {
	plainMap := make(map[string]string, len(fl))
	for key, val := range fl {
		plainMap[key] = val.Value
	}
	return plainMap
}

func labelFlagsFromRule(rule *classificationv1.Rule) LabelFlags {
	return LabelFlags{
		Propagate: boolValueOrTrue(rule.GetPropagate()),
		Hidden:    rule.GetHidden(),
	}
}

func boolValueOrTrue(bv *wrapperspb.BoolValue) bool { return bv == nil || bv.Value }

// New creates a new Flow Classifier.
func New() *Classifier {
	return &Classifier{
		activeRulesets: make(map[rulesetID]compiledRuleset),
	}
}

func populateFlowLabels(ctx context.Context, flowLabels FlowLabels, mm *multimatcher.MultiMatcher[int, []*labeler], labelsForMatching selectors.Labels, input ast.Value) {
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
func (c *Classifier) ActiveRules() []ReportedRule {
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
	compiled, err := compileRuleset(ctx, name, classifier)
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

// BadExtractor is an error occurring when extractor is invalid.
var BadExtractor = extractors.BadExtractor

// BadRego is an error occurring when rego compilation fails.
var BadRego = badRego{}

type badRego struct{}

func (b badRego) Error() string { return "failed to compile rego" }

// BadSelector is an error occurring when selector is invalid.
var BadSelector = badSelector{}

type badSelector struct{}

func (b badSelector) Error() string { return "invalid ruleset selector" }

// BadLabelName is an error occurring when label name is invalid.
var BadLabelName = badLabelName{}

type badLabelName struct{}

func (b badLabelName) Error() string { return "invalid label name" }

// compileRuleset parses ruleset's selector and compiles its rules.
func compileRuleset(ctx context.Context, name string, classifier *classificationv1.Classifier) (compiledRuleset, error) {
	if classifier.Selector == nil {
		return compiledRuleset{}, fmt.Errorf("%w: missing selector", BadSelector)
	}

	selector, err := selectors.FromProto(classifier.Selector)
	if err != nil {
		return compiledRuleset{}, fmt.Errorf("%w: %v", BadSelector, err)
	}

	labelers, err := compileRules(ctx, selector.LabelMatcher, classifier.Rules)
	if err != nil {
		return compiledRuleset{}, fmt.Errorf("failed to compile %q rules for %v: %w", name, selector, err)
	}

	cr := compiledRuleset{
		ControlPointID: selector.ControlPointID,
		Labelers:       labelers,
		ReportedRules:  rulesetToReportedRules(classifier, name),
	}

	return cr, nil
}

// compileRules compiles a set of rules into set of rego queries
//
// Raw rego rules are compiled 1:1 to rego queries. High-level extractor-based
// rules are compiled into a single rego query.
func compileRules(ctx context.Context, labelSelector multimatcher.Expr, labelRules map[string]*classificationv1.Rule) ([]labelerWithSelector, error) {
	log.Trace().Msg("Classifier.compileRules starting")

	// Group all the extractor-based rules so that we can compile them to a
	// single rego query
	labelExtractors := map[string]*classificationv1.Extractor{}
	labelFlags := map[string]LabelFlags{} // flags for labels created by extractors

	rawRegoCount := 0
	var labelers []labelerWithSelector

	for labelName, rule := range labelRules {
		if strings.Contains(labelName, "/") {
			// Forbidding '/' in case we want to support multiple rules for the
			// same label:
			// labels:
			//   user/1: <snip>
			//   user/2: <snip>
			return nil, fmt.Errorf("%w: cannot contain '/'", BadLabelName)
		}

		switch source := rule.GetSource().(type) {
		case *classificationv1.Rule_Extractor:
			labelExtractors[labelName] = source.Extractor
			labelFlags[labelName] = labelFlagsFromRule(rule)
		case *classificationv1.Rule_Rego_:
			query, err := rego.New(
				rego.Query(source.Rego.Query),
				rego.Module("tmp.rego", source.Rego.Source),
			).PrepareForEval(ctx)
			if err != nil {
				log.Trace().Str("src", source.Rego.Source).Str("query", source.Rego.Query).
					Msg("Failed to prepare for eval")
				return nil, fmt.Errorf(
					"failed to compile raw rego module, label: %s, query: %s: %w: %v",
					labelName,
					source.Rego.Query,
					BadRego,
					err,
				)
			}
			labelers = append(labelers, labelerWithSelector{
				LabelSelector: labelSelector,
				Labeler: &labeler{
					Query:      query,
					LabelName:  labelName,
					LabelFlags: labelFlagsFromRule(rule),
				},
			})
			rawRegoCount++
		}
	}

	if len(labelExtractors) != 0 {
		regoSrc, err := extractors.CompileToRego(defaultPackageName, labelExtractors)
		if err != nil {
			return nil, fmt.Errorf("failed to compile extractors to rego: %w", err)
		}
		query, err := rego.New(
			rego.Query("data."+defaultPackageName),
			rego.Module("tmp.rego", regoSrc),
		).PrepareForEval(ctx)
		if err != nil {
			// Note: Not wrapping BadRego error here – the rego returned by
			// compileExtractors should always be valid, otherwise it's a
			// bug, and not user's fault.
			log.Trace().Str("src", regoSrc).Msg("Failed to prepare for eval")
			return nil, fmt.Errorf("(bug) failed to compile classification rules: %w", err)
		}

		labelers = append(labelers, labelerWithSelector{
			LabelSelector: labelSelector,
			Labeler: &labeler{
				Query:       query,
				LabelsFlags: labelFlags,
			},
		})
	}

	log.Info().
		Int("modules", len(labelers)).
		Int("raw rego modules", rawRegoCount).
		Int("extractors", len(labelExtractors)).
		Msg("Compilation of rules finished")

	return labelers, nil
}

// needs to be called with activeRulesets mutex held.
func (c *Classifier) activateRulesets() {
	c.activeRules.Store(c.combineRulesets())
	log.Info().Int("rulesets", len(c.activeRulesets)).Msg("Rules updated")
}

func (c *Classifier) combineRulesets() rules {
	combined := rules{
		MultiMatcherByControlPointID: make(multiMatcherByControlPoint),
		ReportedRules:                make([]ReportedRule, 0),
	}

	// to have unique keys to AddEntry
	controlPointKeys := make(map[selectors.ControlPointID]int)

	for _, ruleset := range c.activeRulesets {
		combined.ReportedRules = append(combined.ReportedRules, ruleset.ReportedRules...)
		for _, labelerWithSelector := range ruleset.Labelers {
			mm, ok := combined.MultiMatcherByControlPointID[ruleset.ControlPointID]
			if !ok {
				mm = multimatcher.New[int, []*labeler]()
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
