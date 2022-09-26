// Companion package for github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1
// containing conversions of proto-generated struct into golang ones and other helpers.
package selectors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fluxninja/aperture/pkg/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mm "github.com/fluxninja/aperture/pkg/multimatcher"

	labelmatcherv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/labelmatcher/v1"
	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
)

type selector struct {
	labelMatcher mm.Expr
	ctrlPtID     ControlPointID
}

// FromProto creates a Selector from a "raw" proto-based Selector
//
// The selector is assumed to be already validated and non-nil.
func FromProto(selectorMsg *selectorv1.Selector) (selector, error) {
	labelMatcher, err := MMExprFromLabelMatcher(selectorMsg.FlowSelector.GetLabelMatcher())
	if err != nil {
		return selector{}, fmt.Errorf("invalid label matcher: %w", err)
	}
	ctrlPtID, err := controlPointIDFromSelectorProto(selectorMsg)
	if err != nil {
		return selector{}, fmt.Errorf("invalid control point: %w", err)
	}
	return selector{
		ctrlPtID:     ctrlPtID,
		labelMatcher: labelMatcher,
	}, nil
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
func MMExprFromLabelMatcher(lm *labelmatcherv1.LabelMatcher) (mm.Expr, error) {
	var reqExprs []mm.Expr

	for k, v := range lm.GetMatchLabels() {
		reqExprs = append(reqExprs, mm.LabelEquals(k, v))
	}

	for _, req := range lm.GetMatchExpressions() {
		switch metav1.LabelSelectorOperator(req.Operator) {
		case metav1.LabelSelectorOpIn:
			matchExpr, err := mm.LabelMatchesRegex(req.Key, valuesRegex(req.Values))
			if err != nil {
				// shouldn't happen as we're in control of the regex, but who knows
				return nil, err
			}
			reqExprs = append(reqExprs, matchExpr)
		case metav1.LabelSelectorOpNotIn:
			matchExpr, err := mm.LabelMatchesRegex(req.Key, valuesRegex(req.Values))
			if err != nil {
				// shouldn't happen as we're in control of the regex, but who knows
				return nil, err
			}
			reqExprs = append(reqExprs, mm.Not(matchExpr))
		case metav1.LabelSelectorOpExists:
			reqExprs = append(reqExprs, mm.LabelExists(req.Key))
		case metav1.LabelSelectorOpDoesNotExist:
			reqExprs = append(reqExprs, mm.Not(mm.LabelExists(req.Key)))
		default:
			log.Panic().Msg("unknown match expression operator")
		}
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
func MMExprFromProto(expr *labelmatcherv1.MatchExpression) (mm.Expr, error) {
	switch e := expr.Variant.(type) {
	case *labelmatcherv1.MatchExpression_Not:
		expr, err := MMExprFromProto(e.Not)
		if err != nil {
			return nil, err
		}
		return mm.Not(expr), nil
	case *labelmatcherv1.MatchExpression_All:
		exprs, err := mmExprsFromProtoList(e.All)
		if err != nil {
			return nil, err
		}
		return mm.All(exprs), nil
	case *labelmatcherv1.MatchExpression_Any:
		exprs, err := mmExprsFromProtoList(e.Any)
		if err != nil {
			return nil, err
		}
		return mm.Any(exprs), nil
	case *labelmatcherv1.MatchExpression_LabelExists:
		return mm.LabelExists(e.LabelExists), nil
	case *labelmatcherv1.MatchExpression_LabelEquals:
		return mm.LabelEquals(e.LabelEquals.Label, e.LabelEquals.Value), nil
	case *labelmatcherv1.MatchExpression_LabelMatches:
		return mm.LabelMatchesRegex(e.LabelMatches.Label, e.LabelMatches.Regex)
	default:
		log.Error().Msg("unknown/unset expression variant")
		return nil, nil
	}
}

func mmExprsFromProtoList(list *labelmatcherv1.MatchExpression_List) ([]mm.Expr, error) {
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

// valuesRegex returns regex expression that'll match any of given values.
func valuesRegex(values []string) string {
	escaped := make([]string, 0, len(values))
	for _, v := range values {
		escaped = append(escaped, regexp.QuoteMeta(v))
	}
	return "^(" + strings.Join(escaped, "|") + ")$"
}
