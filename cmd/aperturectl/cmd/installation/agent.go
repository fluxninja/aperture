package installation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/releaseutil"
)

// agentInstallCmd is the command to install Aperture Agent on Kubernetes.
var agentInstallCmd = &cobra.Command{
	Use:   "agent",
	Short: "Install Aperture Agent",
	Long: fmt.Sprintf(`
Use this command to install Aperture Agent on your Kubernetes cluster.
Refer https://github.com/fluxninja/aperture/blob/v%s/manifests/charts/aperture-agent/README.md for list of configurable parameters for preparing values file.`,
		utils.Version),
	SilenceErrors: true,
	Example: `aperturectl install agent --values-file=values.yaml

aperturectl install agent --values-file=values.yaml --namespace=aperture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if valuesFile == "" {
			return fmt.Errorf("--values-file must be provided")
		}

		if namespace == "" {
			namespace = apertureAgentNS
		}

		if err := manageNamespace(); err != nil {
			return err
		}

		crds, _, manifests, err := getTemplets(agent, releaseutil.InstallOrder)
		for _, crd := range crds {
			if err = applyManifest(string(crd.File.Data)); err != nil {
				return err
			}
		}

		for _, manifest := range manifests {
			if err = applyManifest(manifest.Content); err != nil {
				return err
			}
		}
		return err
	},
}

// agentUnInstallCmd is the command to uninstall Aperture Agent from Kubernetes.
var agentUnInstallCmd = &cobra.Command{
	Use:   "agent",
	Short: "Uninstall Aperture Agent",
	Long: `
Use this command to uninstall Aperture Agent from your Kubernetes cluster`,
	SilenceErrors: true,
	Example: `aperturectl uninstall agent

aperturectl uninstall agent --namespace=aperture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			namespace = apertureAgentNS
		}

		crds, hooks, manifests, err := getTemplets(agent, releaseutil.UninstallOrder)

		for _, hook := range hooks {
			log.Info().Msgf("Executing hook - %s", hook.Name)
			if err = applyManifest(hook.Manifest); err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
			defer cancel()
			if err = waitForHook(hook.Name, ctx); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return fmt.Errorf("timed out waiting for pre-delete hook completion")
				}
				return err
			}

			if err = deleteManifest(hook.Manifest); err != nil {
				return err
			}
		}

		for _, manifest := range manifests {
			if err = deleteManifest(manifest.Content); err != nil {
				return err
			}
		}

		for _, crd := range crds {
			if err = deleteManifest(string(crd.File.Data)); err != nil {
				return err
			}
		}
		return err
	},
}
