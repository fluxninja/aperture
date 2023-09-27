package utils

import (
	"context"
	"errors"
	"strings"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/spf13/cobra"
)

var decisionTypes = []string{
	"load_scheduler",
	"rate_limiter",
	"quota_scheduler",
	"pod_scaler",
	"sampler",
}

// DecisionsPreRun validates the decisions command.
func DecisionsPreRun(all bool, decisionType string) (string, error) {
	if !all {
		if decisionType == "" {
			return decisionType, errors.New("decision type is required or use --all to get all decisions")
		} else {
			var found bool
			for _, v := range decisionTypes {
				if v == decisionType {
					found = true
					break
				}
			}
			if !found {
				return decisionType, errors.New("invalid decision type, use one of the valid types (" + strings.Join(decisionTypes, ", ") + ") or use --all to get all decisions")
			}
		}
	} else {
		decisionType = ""
	}

	return decisionType, nil
}

// ParseDecisions parses the decisions.
func ParseDecisions(cmd *cobra.Command, client PolicyClient, all bool, decisionType string) error {
	getDecisionsReq := &languagev1.GetDecisionsRequest{
		DecisionType: decisionType,
	}
	if all {
		getDecisionsReq = nil
	}
	decisionsResp, err := client.GetDecisions(
		context.Background(),
		getDecisionsReq,
	)
	if err != nil {
		return err
	}

	for k, v := range decisionsResp.Decisions {
		cmd.Printf("%s:\n%s\n\n", k, v)
	}

	return nil
}
