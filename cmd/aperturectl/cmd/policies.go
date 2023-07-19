package cmd

import (
	"context"
	"fmt"
	"strings"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	controller.InitFlags(policiesCmd.PersistentFlags())
}

var policiesCmd = &cobra.Command{
	Use:               "policies",
	Short:             "List applied policies",
	Long:              `List applied policies`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
	RunE: func(*cobra.Command, []string) error {
		client, err := controller.Client()
		if err != nil {
			return err
		}

		policies, err := client.ListPolicies(context.Background(), &emptypb.Empty{})
		if err != nil {
			return err
		}

		for name, body := range policies.GetPolicies().Policies {
			fmt.Printf("%v:\n", name)
			if body.GetStatus() == languagev1.GetPolicyResponse_INVALID {
				fmt.Println("\tStatus: INVALID")
				reason := strings.ReplaceAll(body.GetReason(), "\n", "\n\n\t\t")
				reason = strings.ReplaceAll(reason, " Error", "\n\t\tError")
				fmt.Printf("\tReason: %s\n", reason)
				fmt.Println("\t\t---")
				continue
			}

			if len(body.Policy.Resources.InfraMeters) > 0 {
				fmt.Println("\tInfra Meters:")
				for im := range body.Policy.Resources.InfraMeters {
					fmt.Printf("\t\t%v\n", im)
				}
			}
			if body.Policy.Resources.FlowControl != nil {
				if len(body.Policy.Resources.FlowControl.FluxMeters) > 0 {
					fmt.Println("\tFlux Meters:")
					for fm := range body.Policy.Resources.FlowControl.FluxMeters {
						fmt.Printf("\t\t%v\n", fm)
					}
				}
				if len(body.Policy.Resources.FlowControl.Classifiers) > 0 {
					fmt.Println("\tClassifiers:")
					for _, c := range body.Policy.Resources.FlowControl.Classifiers {
						if len(c.Selectors) > 0 {
							fmt.Println("\t\tSelectors:")
							for _, s := range c.Selectors {
								fmt.Printf("\t\t\t%v\n", s)
							}
						}
						if len(c.Rules) > 0 {
							fmt.Println("\t\tRules:")
							for r := range c.Rules {
								fmt.Printf("\t\t\t%v\n", r)
							}
						}
					}
					fmt.Println("\t\t---")
				}
			}
		}

		return nil
	},
}
