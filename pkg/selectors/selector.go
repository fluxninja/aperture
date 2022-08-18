// Companion package for github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1
// containing convertions of proto-generated struct into golang ones and other helpers.
package selectors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fluxninja/aperture/pkg/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mm "github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/services"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

// Selector is a parsed/preprocessed version of policylangv1.Selector
//
// refer to proto definition for docs.
type Selector struct {
	// Additionally, arbitrary label matcher can be used to match labels.
	// For policies this matcher can _also_ match flow labels.
	LabelMatcher mm.Expr
	// ServiceID and control point are required
	ControlPointID
}

// FromProto creates a Selector from a "raw" proto-based Selector
//
// The selector is assumed to be already validated and non-nil.
func FromProto(selector *policylangv1.Selector) (Selector, error) {
	labelMatcher, err := MMExprFromLabelMatcher(selector.GetLabelMatcher())
	if err != nil {
		return Selector{}, fmt.Errorf("invalid label matcher: %w", err)
	}
	return Selector{
		ControlPointID: ControlPointIDFromProto(selector),
		LabelMatcher:   labelMatcher,
	}, nil
}

// ControlPoint identifies control point within a service that the rule or
// policy should apply to
//
// ControlPoint is either a library feature name or one of ingress / egress
// traffic control point.
type ControlPoint struct {
	// FIXME(FLUX-2362) Stop hardcoding "envoy vs feature" (?)
	// either
	Feature string
	// or
	Traffic TrafficDirection
}

// TrafficDirection indicates enumerated traffic direction.
type TrafficDirection int

const (
	// TrafficDirectionUndefined is a placeholder for undefined traffic direction.
	TrafficDirectionUndefined TrafficDirection = iota
	// Ingress is a traffic direction for inbound traffic.
	Ingress
	// Egress is a traffic direction for outbound traffic.
	Egress
)

// String returns a string representation of the traffic direction.
func (d TrafficDirection) String() string {
	switch d {
	case Ingress:
		return "ingress"
	case Egress:
		return "egress"
	default:
		return ""
	}
}

// String returns a string representation of either a library feature name or a traffic direction.
func (p *ControlPoint) String() string {
	if p.Feature != "" {
		return p.Feature
	}
	return fmt.Sprintf("traffic:%s", p.Traffic)
}

// ControlPointFromProto creates a ControlPoint from "raw" proto-based ControlPoint
//
// The controlPoint is assumed to be already validated and nonnil.
func ControlPointFromProto(controlPoint *policylangv1.ControlPoint) ControlPoint {
	switch cp := controlPoint.Controlpoint.(type) {
	case *policylangv1.ControlPoint_Feature:
		return ControlPoint{Feature: cp.Feature}
	case *policylangv1.ControlPoint_Traffic:
		switch cp.Traffic {
		case "ingress":
			return ControlPoint{Traffic: Ingress}
		case "egress":
			return ControlPoint{Traffic: Egress}
		default:
			log.Error().Msg("invalid traffic direction")
			return ControlPoint{}
		}
	default:
		log.Error().Msg("unknown/missing control point")
		return ControlPoint{}
	}
}

// ControlPointID uniquely identifies the control point within a cluster – so
// it's a ServiceID and ControlPoint combined
//
// Control Point.
type ControlPointID struct {
	ServiceID    services.ServiceID
	ControlPoint ControlPoint
}

// String returns a string representation of control point and service.
func (p ControlPointID) String() string {
	return fmt.Sprintf("%v@%v", p.ControlPoint, p.ServiceID)
}

// ControlPointIDFromProto extracts a ControlPointID from proto-based selector
// (ignoring LabelMatcher)
//
// Selector is assumed to be validated and non-nil.
func ControlPointIDFromProto(selector *policylangv1.Selector) ControlPointID {
	return ControlPointID{
		ServiceID: services.ServiceID{
			Service: selector.Service,
		},
		ControlPoint: ControlPointFromProto(selector.ControlPoint),
	}
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

// valuesRegex returns regex expression that'll match any of given values.
func valuesRegex(values []string) string {
	escaped := make([]string, 0, len(values))
	for _, v := range values {
		escaped = append(escaped, regexp.QuoteMeta(v))
	}
	return "^(" + strings.Join(escaped, "|") + ")$"
}
