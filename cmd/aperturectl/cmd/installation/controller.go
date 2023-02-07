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

// controllerInstallCmd is the command to install Aperture Controller on Kubernetes.
var controllerInstallCmd = &cobra.Command{
	Use:   "controller",
	Short: "Install Aperture Controller",
	Long: fmt.Sprintf(`
Use this command to install Aperture Controller and its dependencies on your Kubernetes cluster.
Refer https://github.com/fluxninja/aperture/blob/v%s/manifests/charts/aperture-controller/README.md for list of configurable parameters for preparing values file.`,
		utils.Version),
	SilenceErrors: true,
	Example: `aperturectl install controller --values-file=values.yaml

aperturectl install controller --values-file=values.yaml --namespace=aperture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			namespace = apertureControllerNS
		}

		if err := manageNamespace(); err != nil {
			return err
		}

		crds, _, manifests, err := getTemplets(controller, releaseutil.InstallOrder)
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

// controllerInstallCmd is the command to install Aperture Controller on Kubernetes.
var controllerUnInstallCmd = &cobra.Command{
	Use:   "controller",
	Short: "Uninstall Aperture Controller",
	Long: `
Use this command to uninstall Aperture Controller and its dependencies from your Kubernetes cluster`,
	SilenceErrors: true,
	Example: `aperturectl uninstall controller

aperturectl uninstall controller --namespace=aperture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" {
			namespace = apertureControllerNS
		}

		crds, hooks, manifests, err := getTemplets(controller, releaseutil.UninstallOrder)

		for _, hook := range hooks {
			log.Info().Msgf("Executing hook - %s", hook.Name)
			if err = applyManifest(hook.Manifest); err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
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
