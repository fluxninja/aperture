package compiler

import (
	"context"
	"fmt"
	"strings"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	classificationv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/resources/classifier/extractors"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/open-policy-agent/opa/rego"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const defaultPackageName = "fluxninja.classification.extractors"

// CompiledRuleset is compiled form of Classifier proto.
type CompiledRuleset struct {
	ControlPointID selectors.ControlPointID
	Labelers       []LabelerWithSelector
	ReportedRules  []ReportedRule
}

// LabelerWithSelector is a labeler with its selector.
type LabelerWithSelector struct {
	Labeler       *Labeler
	LabelSelector multimatcher.Expr
}

// Labeler is used to create flow labels
//
// label can create either:
// * a single label – LabelName is non-empty or
// * multiple labels – LabelName is empty.
type Labeler struct {
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

// LabelFlags are flags for a flow label.
type LabelFlags struct {
	// Should the created label be applied to the whole flow (propagated in baggage)
	Propagate bool
	// Should the created flow label be hidden from telemetry
	Hidden bool
}

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

func labelFlagsFromRule(rule *classificationv1.Rule) LabelFlags {
	return LabelFlags{
		Propagate: boolValueOrTrue(rule.GetPropagate()),
		Hidden:    rule.GetHidden(),
	}
}

func boolValueOrTrue(bv *wrapperspb.BoolValue) bool { return bv == nil || bv.Value }

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

// CompileRuleset parses ruleset's selector and compiles its rules.
func CompileRuleset(ctx context.Context, name string, classifier *classificationv1.Classifier) (CompiledRuleset, error) {
	if classifier.Selector == nil {
		return CompiledRuleset{}, fmt.Errorf("%w: missing selector", BadSelector)
	}

	selector, err := selectors.FromProto(classifier.Selector)
	if err != nil {
		return CompiledRuleset{}, fmt.Errorf("%w: %v", BadSelector, err)
	}

	labelers, err := compileRules(ctx, selector.LabelMatcher, classifier.Rules)
	if err != nil {
		return CompiledRuleset{}, fmt.Errorf("failed to compile %q rules for %v: %w", name, selector, err)
	}

	cr := CompiledRuleset{
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
func compileRules(ctx context.Context, labelSelector multimatcher.Expr, labelRules map[string]*classificationv1.Rule) ([]LabelerWithSelector, error) {
	log.Trace().Msg("Classifier.compileRules starting")

	// Group all the extractor-based rules so that we can compile them to a
	// single rego query
	labelExtractors := map[string]*classificationv1.Extractor{}
	labelFlags := map[string]LabelFlags{} // flags for labels created by extractors

	rawRegoCount := 0
	var labelers []LabelerWithSelector

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
			labelers = append(labelers, LabelerWithSelector{
				LabelSelector: labelSelector,
				Labeler: &Labeler{
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

		labelers = append(labelers, LabelerWithSelector{
			LabelSelector: labelSelector,
			Labeler: &Labeler{
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
