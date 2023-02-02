package blueprints

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	specv1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:           "pull",
	Short:         "Pull a blueprint",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		apertureBlueprintsDir := filepath.Join(blueprintsDir, blueprintsVersion)
		err := os.MkdirAll(apertureBlueprintsDir, os.ModePerm)
		if err != nil {
			return err
		}

		spec := specv1.New()
		contents, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return err
		}
		spec.LegacyImports = false
		contents = append(contents, []byte("\n")...)

		filename := filepath.Join(apertureBlueprintsDir, jsonnetfile.File)
		err = os.WriteFile(filename, contents, os.ModePerm)
		if err != nil {
			return err
		}

		jbLockFileBytes, err := os.ReadFile(filepath.Join(apertureBlueprintsDir, jsonnetfile.LockFile))
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		lockFile, err := jsonnetfile.Unmarshal(jbLockFileBytes)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Join(apertureBlueprintsDir, ".tmp"), os.ModePerm)
		if err != nil {
			return err
		}

		uri := fmt.Sprintf("%s@%s", apertureBlueprintsURI, blueprintsVersion)
		d := deps.Parse(apertureBlueprintsDir, uri)
		if !depEqual(spec.Dependencies[d.Name()], *d) {
			spec.Dependencies[d.Name()] = *d
			delete(lockFile.Dependencies, d.Name())
		}

		locked, err := pkg.Ensure(spec, apertureBlueprintsDir, lockFile.Dependencies)
		if err != nil {
			_ = os.RemoveAll(apertureBlueprintsDir)
			return err
		}

		err = writeChangedJsonnetFile(contents, &spec, filename)
		if err != nil {
			return err
		}
		err = writeChangedJsonnetFile(jbLockFileBytes, &specv1.JsonnetFile{Dependencies: locked}, filepath.Join(apertureBlueprintsDir, jsonnetfile.LockFile))
		if err != nil {
			return err
		}

		return nil
	},
}

func depEqual(d1, d2 deps.Dependency) bool {
	name := d1.Name() == d2.Name()
	version := d1.Version == d2.Version
	source := reflect.DeepEqual(d1.Source, d2.Source)

	return name && version && source
}

func writeChangedJsonnetFile(originalBytes []byte, modified *specv1.JsonnetFile, path string) error {
	origJsonnetFile, err := jsonnetfile.Unmarshal(originalBytes)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(origJsonnetFile, *modified) {
		return nil
	}

	return writeJSONFile(path, *modified)
}

func writeJSONFile(name string, d interface{}) error {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	b = append(b, []byte("\n")...)

	// nolint: gosec
	return os.WriteFile(name, b, 0o644)
}
