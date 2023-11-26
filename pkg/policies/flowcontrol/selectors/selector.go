// Companion package for github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1
// containing conversions of proto-generated struct into golang ones and other helpers.
package selectors

import (
	"fmt"

	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/utils"

	mm "github.com/fluxninja/aperture/v2/pkg/multi-matcher"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
)

// UniqueAgentGroups returns the unique agent groups of selectors.
func UniqueAgentGroups(selectorProto []*policylangv1.Selector) []string {
	agentGroups := []string{}
	for _, selector := range selectorProto {
		agentGroup := selector.GetAgentGroup()
		if !utils.SliceContains(agentGroups, agentGroup) {
			agentGroups = append(agentGroups, agentGroup)
		}
	}
	return agentGroups
}

type selector struct {
	labelMatcher mm.Expr
	ctrlPtID     ControlPointID
}

// FromSelectors creates a Selector from a "raw" proto-based Selector.
func FromSelectors(selectorsProto []*policylangv1.Selector, agentGroup string) ([]selector, error) {
	s := []selector{}
	for _, selectorProto := range selectorsProto {
		if selectorProto.GetAgentGroup() != agentGroup {
			continue
		}
		labelMatcher, err := MMExprFromLabelMatcher(selectorProto.GetLabelMatcher())
		if err != nil {
			return []selector{}, fmt.Errorf("invalid label matcher: %w", err)
		}
		ctrlPtID := ControlPointID{
			ControlPoint: selectorProto.GetControlPoint(),
			Service:      selectorProto.GetService(),
		}
		sel := selector{
			ctrlPtID:     ctrlPtID,
			labelMatcher: labelMatcher,
		}

		s = append(s, sel)
	}
	return s, nil
}

// LabelMatcher returns the label matcher of the selector.
func (s *selector) LabelMatcher() mm.Expr {
	return s.labelMatcher
}

// ControlPointID returns the control point ID of the selector.
func (s *selector) ControlPointID() ControlPointID {
	return s.ctrlPtID
}

// MMExprFromLabelMatcher translates proto definition of label matcher into
// a // single multimatcher expression
//
// LabelMatcher can be nil or a validated LabelMatcher.
func MMExprFromLabelMatcher(lm *policylangv1.LabelMatcher) (mm.Expr, error) {
	var reqExprs []mm.Expr

	for k, v := range lm.GetMatchLabels() {
		reqExprs = append(reqExprs, mm.LabelEquals(k, v))
	}

	if protoExpr := lm.GetExpression(); protoExpr != nil {
		expr, err := MMExprFromProto(protoExpr)
		if err != nil {
			return nil, err
		}
		reqExprs = append(reqExprs, expr)
	}

	return mm.All(reqExprs), nil
}

// MMExprFromProto converts proto definition of expression into multimatcher Expression
//
// The expr is assumed to be validated and nonnil.
func MMExprFromProto(expr *policylangv1.MatchExpression) (mm.Expr, error) {
	switch e := expr.Variant.(type) {
	case *policylangv1.MatchExpression_Not:
		expr, err := MMExprFromProto(e.Not)
		if err != nil {
			return nil, err
		}
		return mm.Not(expr), nil
	case *policylangv1.MatchExpression_All:
		exprs, err := mmExprsFromProtoList(e.All)
		if err != nil {
			return nil, err
		}
		return mm.All(exprs), nil
	case *policylangv1.MatchExpression_Any:
		exprs, err := mmExprsFromProtoList(e.Any)
		if err != nil {
			return nil, err
		}
		return mm.Any(exprs), nil
	case *policylangv1.MatchExpression_LabelExists:
		return mm.LabelExists(e.LabelExists), nil
	case *policylangv1.MatchExpression_LabelEquals:
		return mm.LabelEquals(e.LabelEquals.Label, e.LabelEquals.Value), nil
	case *policylangv1.MatchExpression_LabelMatches:
		return mm.LabelMatchesRegex(e.LabelMatches.Label, e.LabelMatches.Regex)
	default:
		log.Error().Msg("unknown/unset expression variant")
		return nil, nil
	}
}

func mmExprsFromProtoList(list *policylangv1.MatchExpression_List) ([]mm.Expr, error) {
	exprs := make([]mm.Expr, 0, len(list.Of))
	for _, protoExpr := range list.Of {
		expr, err := MMExprFromProto(protoExpr)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}
	return exprs, nil
}
