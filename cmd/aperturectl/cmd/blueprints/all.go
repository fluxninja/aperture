package blueprints

import (
	"github.com/spf13/cobra"
)

func init() {
	allCmd.Flags().StringVar(&policyType, "policy_type", "", "Type of policy to generate e.g. static-rate-limiting, latency-aimd-concurrency-limiting")
	allCmd.Flags().StringVar(&outputDir, "output_dir", "", "Directory path where the generated manifests will be stored. If not provided, will be printed on console")
	allCmd.Flags().BoolVar(&apply, "apply", false, "Apply policy on the Kubernetes cluster")
	allCmd.Flags().StringVar(&kubeConfig, "kube_config", "", "Path to the Kubernets cluster config. Defaults to '~/.kube/config'")
	allCmd.Flags().StringVar(&valuesFile, "values_file", "", "Path to the values file for blueprints input")
}

var allCmd = &cobra.Command{
	Use:           "generate-all",
	Short:         "Generate all",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := graphCmd.RunE(cmd, args); err != nil {
			return err
		}

		if err := dashboardCmd.RunE(cmd, args); err != nil {
			return err
		}

		return nil
	},
}
