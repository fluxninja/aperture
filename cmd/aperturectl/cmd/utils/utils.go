package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ghodss/yaml"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hashicorp/go-multierror"
	"github.com/xeipuuv/gojsonschema"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	languagev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/circuitfactory"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// GenerateDotFile generates a DOT file from the given circuit with the specified depth.
// The depth determines how many levels of components in the tree should be expanded in the graph.
// If maxDepth is set to -1, the function will expand components up to the maximum possible depth.
//
// Parameters:
//   - circuit: A pointer to the circuitfactory.Circuit object to be used for generating the DOT file.
//   - dotFilePath: The file path where the generated DOT file should be saved.
//   - maxDepth: The maximum depth the graph should be expanded to.
//     If set to -1, the function will expand components up to the maximum possible depth.
//
// Returns:
//   - An error if any issues occur during the file creation or writing process, otherwise nil.
//
// Example usage:
//
//	err := GenerateDotFile(circuit, "output.dot", 3)
//	// This will generate a DOT file with components expanded up to a depth of 3.
//
//	err := GenerateDotFile(circuit, "output.dot", -1)
//	// This will generate a DOT file with components expanded up to the maximum possible depth.
func GenerateDotFile(circuit *circuitfactory.Circuit, dotFilePath string, depth int) error {
	graph, err := circuit.Tree.GetSubGraph(runtime.NewComponentID(runtime.RootComponentID), depth)
	if err != nil {
		return err
	}

	d := circuitfactory.DOTGraph(graph)
	f, err := os.Create(dotFilePath)
	if err != nil {
		log.Error().Err(err).Msg("error creating file")
		return err
	}
	defer f.Close()

	_, err = f.WriteString(d)
	if err != nil {
		log.Error().Err(err).Msg("error writing to file")
		return err
	}
	return nil
}

// GenerateMermaidFile generates a Mermaid file from the given circuit with the specified depth.
// The depth determines how many levels of components in the tree should be expanded in the graph.
// If maxDepth is set to -1, the function will expand components up to the maximum possible depth.
//
// Parameters:
//   - circuit: A pointer to the circuitfactory.Circuit object to be used for generating the Mermaid file.
//   - mermaidFile: The file path where the generated Mermaid file should be saved.
//   - maxDepth: The maximum depth the graph should be expanded to.
//     If set to -1, the function will expand components up to the maximum possible depth.
//
// Returns:
//   - An error if any issues occur during the file creation or writing process, otherwise nil.
//
// Example usage:
//
//	err := GenerateMermaidFile(circuit, "output.mmd", 3)
//	// This will generate a Mermaid file with components expanded up to a depth of 3.
//
//	err := GenerateMermaidFile(circuit, "output.mmd", -1)
//	// This will generate a Mermaid file with components expanded up to the maximum possible depth.
func GenerateMermaidFile(circuit *circuitfactory.Circuit, mermaidFile string, depth int) error {
	graph, err := circuit.Tree.GetSubGraph(runtime.NewComponentID(runtime.RootComponentID), depth)
	if err != nil {
		return err
	}

	m := circuitfactory.MermaidGraph(graph)
	f, err := os.Create(mermaidFile)
	if err != nil {
		log.Error().Err(err).Msg("error creating file")
		return err
	}
	defer f.Close()

	_, err = f.WriteString(m)
	if err != nil {
		log.Error().Err(err).Msg("error writing to file")
		return err
	}
	return nil
}

// CompilePolicy compiles the policy and returns the circuit.
func CompilePolicy(path string) (*circuitfactory.Circuit, *languagev1.Policy, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()

	// FIXME This ValidateAndCompile function validates the policy as a whole –
	// circuit, but also the other resource classifiers, fluxmeters.  This
	// command is called "circuit-compiler" though, so it is bit... surprising.
	// If we compiled just a circuit, we could drop dependency on
	// `controlplane` package.
	circuit, policy, err := controlplane.ValidateAndCompile(ctx, filepath.Base(path), yamlFile)
	if err != nil {
		return nil, nil, err
	}
	return circuit, policy, nil
}

// FetchPolicyFromCR extracts the spec key from a CR and saves it to a temp file.
func FetchPolicyFromCR(crPath string) (string, error) {
	// extract spec key from CR and save it a temp file
	// call compilePolicy with the temp file
	// delete the temp file
	crFile, err := os.ReadFile(crPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to read CR file")
		return "", err
	}
	// unmarshal yaml to map struct and extract spec key
	var cr map[string]interface{}
	err = yaml.Unmarshal(crFile, &cr)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal CR file")
		return "", err
	}
	spec, ok := cr["spec"]
	if !ok {
		log.Error().Msg("failed to find spec key in CR file")
		return "", err
	}
	// marshal spec to yaml
	specYaml, err := yaml.Marshal(spec)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal spec key in CR file")
		return "", err
	}
	// get filename from path
	filename := filepath.Base(crPath)
	// create temp file
	tmpfile, err := os.CreateTemp("", filename)
	if err != nil {
		log.Error().Err(err).Msg("failed to create temp file")
		return "", err
	}
	// write spec to temp file
	_, err = tmpfile.Write(specYaml)
	if err != nil {
		log.Error().Err(err).Msg("failed to write to temp file")
		return "", err
	}
	// close temp file
	err = tmpfile.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close temp file")
		return "", err
	}

	return tmpfile.Name(), nil
}

// GetKubeConfig prepares Kubernetes config to connect with the cluster using provided or default kube config file location.
func GetKubeConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig == "" {
		if kubeConfigEnv, exists := os.LookupEnv("KUBECONFIG"); exists {
			kubeConfig = kubeConfigEnv
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			kubeConfig = filepath.Join(homeDir, ".kube", "config")
		}
		log.Info().Msgf("Using Kubernetes config '%s'", kubeConfig)
	}
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kubernetes. Error: %s", err.Error())
	}
	kubeRestConfig := restConfig
	return kubeRestConfig, nil
}

// ResolveLatestVersion returns the latest release version of Aperture.
func ResolveLatestVersion() (string, error) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{apertureRepo},
	})

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	var latestRelease *semver.Version

	tagsRefPrefix := "refs/tags/v"

	for _, ref := range refs {
		reference := ref.Name().String()
		if ref.Name().IsTag() && strings.HasPrefix(reference, tagsRefPrefix) {
			version := strings.TrimPrefix(reference, tagsRefPrefix)

			release, err := semver.NewVersion(version)
			if err != nil {
				return "", err
			}

			if release.Prerelease() != "" {
				continue
			}

			if latestRelease == nil || release.GreaterThan(latestRelease) {
				latestRelease = release
			}
		}
	}

	if latestRelease == nil {
		return "", errors.New("unable to resolve release tags to find latest release")
	}
	return fmt.Sprintf("v%s", latestRelease.String()), nil
}

// ValidateWithJSONSchema validates the given document (YAML) against the given JSON schema.
func ValidateWithJSONSchema(rootSchema string, schemas []string, documentFile string) error {
	// load schema
	schemaLoader := gojsonschema.NewSchemaLoader()
	// check whether schemas exist
	for _, schema := range schemas {
		if _, err := os.Stat(schema); os.IsNotExist(err) {
			log.Warn().Msgf("Schema %s does not exist. Skipping validation", schema)
			return nil
		}
		err := schemaLoader.AddSchemas(gojsonschema.NewReferenceLoader("file://" + schema))
		if err != nil {
			log.Error().Err(err).Msgf("Failed to add schema %s", schema)
			return err
		}
	}

	if _, err := os.Stat(rootSchema); os.IsNotExist(err) {
		log.Warn().Msgf("Schema %s does not exist. Skipping validation", rootSchema)
		return nil
	}

	schema, err := schemaLoader.Compile(gojsonschema.NewReferenceLoader("file://" + rootSchema))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to compile schema %s", rootSchema)
		return err
	}
	// marshal documentFile to json and load it
	documentYamlBytes, err := os.ReadFile(documentFile)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read document %s", documentFile)
		return err
	}
	documentJSON, err := yaml.YAMLToJSON(documentYamlBytes)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to convert document %s to JSON", documentFile)
		return err
	}

	documentLoader := gojsonschema.NewBytesLoader(documentJSON)

	// validate document
	result, err := schema.Validate(documentLoader)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to validate document %s", documentFile)
		return err
	}
	if !result.Valid() {
		merr := fmt.Errorf("the document %s is not valid", documentFile)
		for _, desc := range result.Errors() {
			errorMessage := fmt.Sprintf("- %s", desc)
			merr = multierror.Append(merr, errors.New(errorMessage))
			log.Error().Err(err).Msg(errorMessage)
		}
		return merr
	}
	return nil
}

// IsBlueprintDeprecated whether the policyDir is deprecated
// it reads metadata.yaml and checks for deprecated key
// the value of that key is the deprecation message.
func IsBlueprintDeprecated(policyDir string) (bool, string) {
	metadataPath := filepath.Join(policyDir, "metadata.yaml")
	metadataFile, err := os.ReadFile(metadataPath)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to read metadata.yaml file in %s", policyDir)
		return false, ""
	}
	var metadata map[string]interface{}
	err = yaml.Unmarshal(metadataFile, &metadata)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to unmarshal metadata.yaml file in %s", policyDir)
		return false, ""
	}
	deprecated, ok := metadata["deprecated"]
	if !ok {
		return false, ""
	}
	return true, deprecated.(string)
}

// GetControllerDeployment returns the deployment of the Aperture Controller.
func GetControllerDeployment(kubeRestConfig *rest.Config, namespace string) (*appsv1.Deployment, error) {
	clientSet, err := kubernetes.NewForConfig(kubeRestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new ClientSet: %w", err)
	}

	deployment, err := clientSet.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.Set{"app.kubernetes.io/component": "aperture-controller"}.String(),
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf(
				"no deployment with name 'aperture-controller' found on the Kubernetes cluster. The policy can be only applied in the namespace where the Aperture Controller is running")
		}
		return nil, fmt.Errorf("failed to fetch namespace of Aperture Controller in Kubernetes: %w", err)
	}

	if len(deployment.Items) != 1 {
		return nil, errors.New("no deployment with name 'aperture-controller' found on the Kubernetes cluster. The policy can be only applied in the namespace where the Aperture Controller is running")
	}

	return &deployment.Items[0], nil
}
