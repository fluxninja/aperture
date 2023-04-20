package classifier

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/consts"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier/compiler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/status"
)

type multiMatcherResult struct {
	labelers []*compiler.LabelerWithSelector
	previews []iface.HTTPRequestPreview
}

type (
	multiMatcherByControlPoint map[selectors.ControlPointID]*multimatcher.MultiMatcher[int, multiMatcherResult]
)

// rules is a helper struct to keep both compiled and uncompiled sets of rules in sync.
type rules struct {
	// rules compiled to map from ControlPointID to MultiMatcher
	MultiMatcherByControlPointID multiMatcherByControlPoint
	// non-compiled version of rules, used for reporting
	ReportedRules []compiler.ReportedRule
}

// ClassificationEngine receives classification policies and provides Classify method.
type ClassificationEngine struct {
	rulesMutex         sync.Mutex
	activeRules        atomic.Value
	classifierMapMutex sync.RWMutex
	registry           status.Registry
	activePreviews     map[iface.PreviewID]iface.HTTPRequestPreview
	activeRulesets     map[rulesetID]compiler.CompiledRuleset
	classifierMap      map[iface.ClassifierID]iface.Classifier
	counterVec         *prometheus.CounterVec
	nextRulesetID      rulesetID
}

type rulesetID = uint64

// NewClassificationEngine creates a new Classifier.
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
		activePreviews: make(map[iface.PreviewID]iface.HTTPRequestPreview),
		counterVec:     counterVector,
	}
}

var (
	evalFailedSampler         = log.NewRatelimitingSampler()
	emptyResultsetSampler     = log.NewRatelimitingSampler()
	ambiguousResultsetSampler = log.NewRatelimitingSampler()
	not1ExprSampler           = log.NewRatelimitingSampler()
	tokenNotStringSampler     = log.NewRatelimitingSampler()
)

func (c *ClassificationEngine) populateFlowLabels(ctx context.Context,
	flowLabels flowlabel.FlowLabels,
	mm *multimatcher.MultiMatcher[int, multiMatcherResult],
	labelsForMatching map[string]string,
	input ast.Value,
) (classifierMsgs []*flowcontrolv1.ClassifierInfo, tokens uint64) {
	tokens = 0
	logger := c.registry.GetLogger()
	appendNewClassifier := func(labelerWithSelector *compiler.LabelerWithSelector, error flowcontrolv1.ClassifierInfo_Error) {
		classifierMsgs = append(classifierMsgs, &flowcontrolv1.ClassifierInfo{
			PolicyName:      labelerWithSelector.ClassifierAttributes.PolicyName,
			PolicyHash:      labelerWithSelector.ClassifierAttributes.PolicyHash,
			ClassifierIndex: labelerWithSelector.ClassifierAttributes.ClassifierIndex,
			Error:           error,
		})
	}

	mmResult := mm.Match(labelsForMatching)

	requestParsedOK := false
	ifaceMap := make(map[string]interface{})
	previews := mmResult.previews
	if len(previews) > 0 {
		// Extract interface{} from ast.Value
		ifaceRequest, err := ast.ValueToInterface(input, valueResolver{})
		if err != nil {
			log.Bug().Msgf("failed to convert value to interface: %v", err)
		} else {
			ifaceMap, requestParsedOK = ifaceRequest.(map[string]interface{})
		}
	}
	for _, preview := range previews {
		if requestParsedOK {
			preview.AddHTTPRequestPreview(ifaceMap)
		}
	}

	labelers := mmResult.labelers

	for _, labelerWithSelector := range labelers {
		labeler := labelerWithSelector.Labeler
		resultSet, err := labeler.Query.Eval(ctx, rego.EvalParsedInput(input))
		if err != nil {
			logger.Sample(evalFailedSampler).Warn().Msg("Rego: Evaluation failed")
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_EVAL_FAILED)
			continue
		}

		if len(resultSet) == 0 {
			logger.Sample(emptyResultsetSampler).Warn().Msg("Rego: Empty resultSet")
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_EMPTY_RESULTSET)
			continue
		} else if len(resultSet) > 1 {
			logger.Sample(ambiguousResultsetSampler).Warn().Msg("Rego: Ambiguous resultSet")
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_AMBIGUOUS_RESULTSET)
			continue
		}

		if nExpressions := len(resultSet[0].Expressions); nExpressions != 1 {
			logger.Sample(not1ExprSampler).Warn().Int("n", nExpressions).Msg("Rego: Expected exactly one expression")
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_MULTI_EXPRESSION)
			continue
		}

		if labeler.LabelName != "" {
			// single-label-query
			flowLabels[labeler.LabelName] = flowlabel.FlowLabelValue{
				Value: resultSet[0].Expressions[0].String(),
				// nolint
				Telemetry: labeler.Telemetry,
			}
			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_NONE)
		} else {
			// multi-label-query
			variables, isMap := resultSet[0].Expressions[0].Value.(map[string]interface{})
			if !isMap {
				logger.Bug().Msg("bug: Rego: Expression is not a map")
				appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_EXPRESSION_NOT_MAP)
				continue
			}

			appendNewClassifier(labelerWithSelector, flowcontrolv1.ClassifierInfo_ERROR_NONE)
			for key, value := range variables {
				// if key is tokens, set tokens instead of flowLabels
				if key == consts.TokensLabel {
					// parse tokens to uint64
					tokens, err = strconv.ParseUint(fmt.Sprint(value), 10, 64)
					if err != nil {
						logger.Sample(tokenNotStringSampler).Warn().Msg("Rego: tokens is not a string")
						continue
					}
				} else {
					// copy this variable to labels
					if l, ok := labeler.Labels[key]; ok {
						flowLabels[key] = flowlabel.FlowLabelValue{
							Value:     fmt.Sprint(value),
							Telemetry: l.Telemetry,
						}
					}
				}
			}
		}
	}
	return
}

type valueResolver struct{}

// Resolve implements ast.ValueResolver interface.
func (valueResolver) Resolve(ref ast.Ref) (interface{}, error) {
	return make(map[string]interface{}), nil
}

// Classify takes rego input, performs classification, and returns a map of flow labels and tokens.
// LabelsForMatching are additional labels to use for selector matching.
// Request is passed as ast.Value directly instead of map[string]interface{} to avoid unnecessary json conversion.
func (c *ClassificationEngine) Classify(
	ctx context.Context,
	svcs []string,
	ctrlPt string,
	labelsForMatching map[string]string,
	input ast.Value,
) ([]*flowcontrolv1.ClassifierInfo, flowlabel.FlowLabels, uint64) {
	tokens := uint64(0)
	flowLabels := make(flowlabel.FlowLabels)

	r, ok := c.activeRules.Load().(rules)
	if !ok {
		return nil, flowLabels, 0
	}

	var classifierMsgs []*flowcontrolv1.ClassifierInfo

	// Catch all Service
	cpID := selectors.NewControlPointID(consts.AnyService, ctrlPt)
	mm, ok := r.MultiMatcherByControlPointID[cpID]
	if ok {
		classifierInfos, t := c.populateFlowLabels(ctx, flowLabels, mm, labelsForMatching, input)
		classifierMsgs = append(classifierMsgs, classifierInfos...)
		tokens = t
	}

	// TODO (krdln): update prometheus metrics upon classification errors.

	// Specific Service
	for _, svc := range svcs {
		cpID := selectors.NewControlPointID(svc, ctrlPt)
		mm, ok := r.MultiMatcherByControlPointID[cpID]
		if !ok {
			c.registry.GetLogger().Trace().Interface("controlPointID", cpID).Msg("No labelers for controlPointID")
			continue
		}
		classifierInfos, t := c.populateFlowLabels(ctx, flowLabels, mm, labelsForMatching, input)
		classifierMsgs = append(classifierMsgs, classifierInfos...)
		// check if t is greater than tokens
		if t > tokens {
			tokens = t
		}
	}

	return classifierMsgs, flowLabels, tokens
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

	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()
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
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()
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

	// function to add rules and previews to multimatcher
	addToMatcher := func(controlPointID selectors.ControlPointID, labelSelector multimatcher.Expr, callback multimatcher.MatchCallback[multiMatcherResult]) error {
		mm, ok := combined.MultiMatcherByControlPointID[controlPointID]
		if !ok {
			mm = multimatcher.New[int, multiMatcherResult]()
			combined.MultiMatcherByControlPointID[controlPointID] = mm
		}
		matcherID := controlPointKeys[controlPointID]
		controlPointKeys[controlPointID]++
		err := mm.AddEntry(matcherID, labelSelector, callback)
		if err != nil {
			log.Error().Err(err).Msg("Failed to add entry to multimatcher")
			return err
		}
		return nil
	}

	for _, ruleset := range c.activeRulesets {
		combined.ReportedRules = append(combined.ReportedRules, ruleset.ReportedRules...)
		for i := range ruleset.Labelers {
			labelerWithSelector := &ruleset.Labelers[i]
			err := addToMatcher(ruleset.ControlPointID, labelerWithSelector.LabelSelector, func(mmr multiMatcherResult) multiMatcherResult {
				mmr.labelers = append(mmr.labelers, labelerWithSelector)
				return mmr
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to add entry to multimatcher")
				return rules{}
			}
		}
	}

	// add activePreviews
	for _, preview := range c.activePreviews {
		selector, err := selectors.FromProto(preview.GetFlowSelector())
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse selector")
			continue
		}
		controlPointID := selector.ControlPointID()
		err = addToMatcher(controlPointID, selector.LabelMatcher(), func(mmr multiMatcherResult) multiMatcherResult {
			mmr.previews = append(mmr.previews, preview)
			return mmr
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to add preview entry to multimatcher")
			continue
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

// AddPreview adds a preview to the active previews.
func (c *ClassificationEngine) AddPreview(preview iface.HTTPRequestPreview) {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()
	c.activePreviews[preview.GetPreviewID()] = preview
	c.activateRulesets()
}

// DropPreview removes a preview from the active previews.
func (c *ClassificationEngine) DropPreview(preview iface.HTTPRequestPreview) {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()
	delete(c.activePreviews, preview.GetPreviewID())
	c.activateRulesets()
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
