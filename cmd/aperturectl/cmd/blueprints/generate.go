package blueprints

import (
	"github.com/spf13/cobra"
)

func init() {
	generateCmd.Flags().StringVar(&policyType, "policy_type", "", "Type of policy to generate e.g. static-rate-limiting, latency-aimd-concurrency-limiting")
	generateCmd.Flags().StringVar(&outputDir, "output_dir", "", "Directory path where the generated manifests will be stored. If not provided, will be printed on console")
	generateCmd.Flags().BoolVar(&apply, "apply", false, "Apply policy on the Kubernetes cluster")
	generateCmd.Flags().StringVar(&kubeConfig, "kube_config", "", "Path to the Kubernets cluster config. Defaults to '~/.kube/config'")
	generateCmd.Flags().StringVar(&valuesFile, "values_file", "", "Path to the values file for blueprints input")

	generateCmd.AddCommand(generatePolicyCmd)
	generateCmd.AddCommand(generateDashboardCmd)
	generateCmd.AddCommand(generateGraphCmd)
}

var generateCmd = &cobra.Command{
	Use:           "generate",
	Short:         "Generate manifests for the given policy type",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateGraphCmd.RunE(cmd, args); err != nil {
			return err
		}
		if err := generateDashboardCmd.RunE(cmd, args); err != nil {
			return err
		}
		return nil
	},
}
