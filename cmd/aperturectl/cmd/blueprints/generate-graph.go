package blueprints

// import (
// 	"os"
// 	"path/filepath"

// 	"github.com/goccy/go-graphviz"
// 	"github.com/spf13/cobra"

// 	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
// 	"github.com/fluxninja/aperture/pkg/log"
// )

// func init() {
// 	generateGraphCmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory path where the generated manifests will be stored. If not provided, will be printed on console")
// 	generateGraphCmd.Flags().StringVar(&valuesFile, "values-file", "", "Path to the values file for blueprints input")
// }

// var generateGraphCmd = &cobra.Command{
// 	Use:           "graph",
// 	Short:         "Generate graph",
// 	SilenceErrors: true,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		if outputDir == "" {
// 			cwd, err := os.Getwd()
// 			if err != nil {
// 				return err
// 			}

// 			outputDir = cwd
// 		}

// 		updatedOutputDir, err := getOutputDir(outputDir)
// 		if err != nil {
// 			return err
// 		}

// 		err = generatePolicyCmd.RunE(cmd, args)
// 		if err != nil {
// 			return err
// 		}

// 		policyFile, err := utils.FetchPolicyFromCR(filepath.Join(updatedOutputDir, "policy.yaml"))
// 		if err != nil {
// 			return err
// 		}
// 		defer os.Remove(policyFile)

// 		circuit, err := utils.CompilePolicy(policyFile)
// 		if err != nil {
// 			log.Error().Err(err).Msg("error reading policy spec")
// 			return err
// 		}

// 		if err = utils.GenerateDotFile(circuit, filepath.Join(updatedOutputDir, "graph.dot")); err != nil {
// 			return err
// 		}

// 		graphBytes, err := os.ReadFile(filepath.Join(updatedOutputDir, "graph.dot"))
// 		if err != nil {
// 			return err
// 		}

// 		g := graphviz.New()
// 		graph, err := graphviz.ParseBytes(graphBytes)
// 		if err != nil {
// 			return err
// 		}

// 		if err := g.RenderFilename(graph, graphviz.SVG, filepath.Join(updatedOutputDir, "graph.svg")); err != nil {
// 			return err
// 		}
// 		log.Info().Msgf("Stored SVG file at '%s'", filepath.Join(updatedOutputDir, "graph.svg"))

// 		return nil
// 	},
// }
