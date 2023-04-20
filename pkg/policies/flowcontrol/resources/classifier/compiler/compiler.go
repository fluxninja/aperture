package compiler

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/open-policy-agent/opa/rego"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier/extractors"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
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
	Labeler              *Labeler
	LabelSelector        multimatcher.Expr
	ClassifierAttributes *policysyncv1.ClassifierAttributes
}

// Labeler is used to create flow labels
//
// label can create either:
// * a single label – LabelName is non-empty or
// * multiple labels – LabelName is empty.
type Labeler struct {
	// rego query that is prepared to take envoy authz request as an input.
	// Result expression should be a single value (if LabelName is set) or a
	// map[string]interface{} otherwise.
	Query rego.PreparedEvalQuery
	// flags for created flow labels:
	Labels map[string]LabelProperties
	// flow label that the result should be assigned to (single-label variant)
	// Deprecated: 1.5.0
	LabelName string
	// Deprecated: 1.5.0
	Telemetry bool // single-label variant
}

// LabelProperties is a set of properties for a label.
type LabelProperties struct {
	Telemetry bool
}

// ReportedRule is a rule along with its selector and label name.
type ReportedRule struct {
	FlowSelector *policylangv1.FlowSelector
	Rule         *policylangv1.Rule
	RulesetName  string
	LabelName    string
}

func rulesetToReportedRules(rs *policylangv1.Classifier, rulesetName string) []ReportedRule {
	out := make([]ReportedRule, 0, len(rs.Rules))
	for label, rule := range rs.Rules {
		out = append(out, ReportedRule{
			RulesetName:  rulesetName,
			LabelName:    label,
			Rule:         rule,
			FlowSelector: rs.FlowSelector,
		})
	}
	return out
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
var BadLabelName = extractors.BadLabelName

// CompileRuleset parses ruleset's selector and compiles its rules.
func CompileRuleset(ctx context.Context, name string, classifierWrapper *policysyncv1.ClassifierWrapper) (CompiledRuleset, error) {
	classifierMsg := classifierWrapper.GetClassifier()
	if classifierMsg.FlowSelector == nil {
		return CompiledRuleset{}, fmt.Errorf("%w: missing selector", BadSelector)
	}

	selector, err := selectors.FromProto(classifierMsg.FlowSelector)
	if err != nil {
		return CompiledRuleset{}, fmt.Errorf("%w: %v", BadSelector, err)
	}

	labelers, err := compileRules(ctx, selector.LabelMatcher(), classifierWrapper)
	if err != nil {
		return CompiledRuleset{}, fmt.Errorf("failed to compile %q rules for %v: %w", name, selector, err)
	}

	cr := CompiledRuleset{
		ControlPointID: selector.ControlPointID(),
		Labelers:       labelers,
		ReportedRules:  rulesetToReportedRules(classifierMsg, name),
	}

	return cr, nil
}

// compileRules compiles a set of rules into set of rego queries
//
// Raw rego rules are compiled 1:1 to rego queries. High-level extractor-based
// rules are compiled into a single rego query.
func compileRules(ctx context.Context, labelSelector multimatcher.Expr, classifierWrapper *policysyncv1.ClassifierWrapper) ([]LabelerWithSelector, error) {
	log.Trace().Msg("Classifier.compileRules starting")

	classifierAttributes := classifierWrapper.GetClassifierAttributes()
	if classifierAttributes == nil {
		return nil, fmt.Errorf("commonAttributes is nil")
	}

	var labelers []LabelerWithSelector

	labelRules := classifierWrapper.GetClassifier().GetRules()

	if len(labelRules) > 0 {
		// Group all the extractor-based rules so that we can compile them to a
		// single rego query
		labelExtractors := map[string]*policylangv1.Extractor{}
		labelsProperties := map[string]LabelProperties{} // Telemetry flag for labels created by extractors

		rawRegoCount := 0

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
			case *policylangv1.Rule_Extractor:
				labelExtractors[labelName] = source.Extractor
				labelsProperties[labelName] = LabelProperties{
					Telemetry: rule.GetTelemetry(),
				}
			case *policylangv1.Rule_Rego_:
				query, err := rego.New(
					//nolint
					rego.Query(source.Rego.Query),
					//nolint
					rego.Module("tmp.rego", source.Rego.Source),
				).PrepareForEval(ctx)
				if err != nil {
					//nolint
					log.Trace().Str("src", source.Rego.Source).Str("query", source.Rego.Query).
						Msg("Failed to prepare for eval")
					return nil, fmt.Errorf(
						"failed to compile raw rego module, label: %s, query: %s: %w: %v",
						labelName,
						// nolint
						source.Rego.Query,
						BadRego,
						err,
					)
				}
				labelers = append(labelers, LabelerWithSelector{
					LabelSelector: labelSelector,
					Labeler: &Labeler{
						Query:     query,
						LabelName: labelName,
						Telemetry: rule.GetTelemetry(),
					},
					ClassifierAttributes: classifierAttributes,
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
				// compileExtractors should always be valid, otherwise it is a
				// bug, and not user's fault.
				log.Trace().Str("src", regoSrc).Msg("Failed to prepare for eval")
				return nil, fmt.Errorf("(bug) failed to compile classification rules: %w", err)
			}

			labelers = append(labelers, LabelerWithSelector{
				LabelSelector: labelSelector,
				Labeler: &Labeler{
					Query:  query,
					Labels: labelsProperties,
				},
				ClassifierAttributes: classifierAttributes,
			})
		}
		log.Debug().
			Int("raw rego modules", rawRegoCount).
			Int("extractors", len(labelExtractors)).
			Msg("Compilation of extractor rules finished")
	}

	// compile rego rules
	r := classifierWrapper.GetClassifier().GetRego()
	if r != nil {
		module := r.GetModule()
		labels := r.GetLabels()

		labelsProperties := map[string]LabelProperties{}
		for labelKey, lp := range labels {
			if !extractors.IsRegoIdent(labelKey) {
				return nil, fmt.Errorf("%w: %q is not a valid label name", BadLabelName, labelKey)
			}
			labelsProperties[labelKey] = LabelProperties{
				Telemetry: lp.Telemetry,
			}
		}

		// get package name in module
		// package name is specified in a line like "package foo"
		p := regexp.MustCompile(`package\s+(\w+)`)
		m := p.FindStringSubmatch(module)
		if len(m) != 2 {
			return nil, fmt.Errorf("failed to get package name from rego module")
		}
		packageName := m[1]

		// compile rego module and queries
		query, err := rego.New(rego.Query("data."+packageName),
			rego.Module("tmp.rego", module)).PrepareForEval(ctx)
		if err != nil {
			log.Trace().Str("src", r.GetModule()).Msg("Failed to prepare for eval")
			return nil, fmt.Errorf(
				"failed to compile raw rego module, query: %s: %w: %v",
				r.GetModule(),
				BadRego,
				err,
			)
		}
		// add to labelers
		labelers = append(labelers, LabelerWithSelector{
			LabelSelector: labelSelector,
			Labeler: &Labeler{
				Query:  query,
				Labels: labelsProperties,
			},
			ClassifierAttributes: classifierAttributes,
		})
		log.Debug().
			Int("extractors", len(labels)).
			Msg("Compilation of rego finished")
	}
	log.Debug().
		Int("labelers", len(labelers)).
		Msg("Compilation finished")

	return labelers, nil
}
