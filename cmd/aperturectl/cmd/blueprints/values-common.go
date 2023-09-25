package blueprints

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

func createValuesFile(blueprintName string, valuesFile string, dynamicConfig bool) error {
	valFileName := valuesFileName

	if dynamicConfig {
		valFileName = dynamicConfigValuesFileName
	}

	blueprintDir := filepath.Join(blueprintsDir, blueprintName)
	// Show a warning if the blueprint is deprecated
	ok, message := utils.IsBlueprintDeprecated(blueprintDir)
	if ok {
		if utils.AllowDeprecated {
			log.Warn().Msgf("Blueprint %s is deprecated: %s", blueprintName, message)
		} else {
			return fmt.Errorf("blueprint %s is deprecated: %s", blueprintName, message)
		}
	}

	blueprintGenDir := filepath.Join(blueprintDir, "gen")

	srcValuesFile := filepath.Join(blueprintGenDir, valFileName)
	if _, err := os.Stat(srcValuesFile); err != nil {
		return fmt.Errorf("values file not found for the blueprint at: %s", srcValuesFile)
	}

	in, err := os.Open(srcValuesFile)
	if err != nil {
		return err
	}
	defer in.Close()
	// Warn if the file already exists and ask to overwrite
	if !overwrite {
		if _, err = os.Stat(valuesFile); err == nil {
			fmt.Printf("File %s already exists. Overwrite? [y/N] ", valuesFile)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return nil
			}
		}
	}
	out, err := os.Create(valuesFile)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	// prepend YAML modeline to the file
	if !noYAMLModeline {
		var schemaURL string
		prefix := fmt.Sprintf("file:%s", blueprintGenDir)

		if strings.Contains(blueprintsURI, "github.com") {
			// Splitting at '@' to get the tag
			parts := strings.Split(blueprintsURI, "@")
			if len(parts) == 2 {

				tag := parts[1]

				// Removing github.com and tag
				trimmedURI := parts[0][len("github.com/"):]

				// Get org and repo
				orgRepoParts := strings.SplitN(trimmedURI, "/", 2)
				org := orgRepoParts[0]
				repoAndPath := orgRepoParts[1]

				// Split repo and path
				repoPathParts := strings.SplitN(repoAndPath, "/", 2)
				repo := repoPathParts[0]
				path := ""
				if len(repoPathParts) > 1 {
					path = repoPathParts[1]
				}

				// Construct the raw GitHub URL
				prefix = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s/%s/gen", org, repo, tag, path, blueprintName)
			}
		}

		if !dynamicConfig {
			schemaURL = fmt.Sprintf("%s/definitions.json", prefix)
		} else {
			schemaURL = fmt.Sprintf("%s/dynamic-config-definitions.json", prefix)
		}
		_, err = out.WriteString("# yaml-language-server: $schema=" + schemaURL + "\n")
		if err != nil {
			return err
		}
	}
	// check whether the policy is deprecated

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}

	command, args := getEnvEditorWithFallback()
	args = append(args, valuesFile)
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error opening file with editor: %s", err)
	}
	return nil
}