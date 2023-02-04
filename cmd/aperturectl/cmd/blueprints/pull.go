package blueprints

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	specv1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
	"github.com/spf13/cobra"
)

func resolveLatestVersion() (string, error) {
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

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull Aperture Blueprints",
	Long: `
Use this command to pull the Aperture Blueprints in local system to use for generating Aperture Policies and Grafana Dashboards.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Example: `aperturectl blueprints pull

aperturectl blueprints pull --version v0.22.0`,
	RunE: func(_ *cobra.Command, _ []string) error {
		spec := specv1.New()
		contents, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return err
		}
		spec.LegacyImports = false
		contents = append(contents, []byte("\n")...)

		filename := filepath.Join(blueprintsDir, jsonnetfile.File)
		err = os.WriteFile(filename, contents, os.ModePerm)
		if err != nil {
			return err
		}

		jbLockFileBytes, err := os.ReadFile(filepath.Join(blueprintsDir, jsonnetfile.LockFile))
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		lockFile, err := jsonnetfile.Unmarshal(jbLockFileBytes)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Join(blueprintsDir, ".tmp"), os.ModePerm)
		if err != nil {
			return err
		}

		d := deps.Parse("", blueprintsURI)
		if d == nil {
			return errors.New("unable to parse blueprints URI: " + blueprintsURI)
		}

		// read d and based on source write uri to uriFilename
		source := ""
		relPath := ""
		version := d.Version

		if d.Source.GitSource != nil {
			source = d.Source.GitSource.Name()
			relPath = source
		} else if d.Source.LocalSource != nil {
			source = d.Source.LocalSource.Directory
			relPath = filepath.Base(source)
		} else {
			return errors.New("unable to parse blueprints URI: " + blueprintsURI)
		}

		err = os.WriteFile(filepath.Join(blueprintsDir, sourceFilename), []byte(source), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(blueprintsDir, relPathFilename), []byte(relPath), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(blueprintsDir, versionFilename), []byte(version), os.ModePerm)
		if err != nil {
			return err
		}

		if !depEqual(spec.Dependencies[d.Name()], *d) {
			spec.Dependencies[d.Name()] = *d
			delete(lockFile.Dependencies, d.Name())
		}

		locked, err := pkg.Ensure(spec, blueprintsDir, lockFile.Dependencies)
		if err != nil {
			_ = os.RemoveAll(blueprintsDir)
			return err
		}

		err = writeChangedJsonnetFile(contents, &spec, filename)
		if err != nil {
			return err
		}
		err = writeChangedJsonnetFile(jbLockFileBytes, &specv1.JsonnetFile{Dependencies: locked}, filepath.Join(blueprintsDir, jsonnetfile.LockFile))
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

func getSource(blueprintsDir string) string {
	source := ""
	// if it doesn't exist, continue
	if _, err := os.Stat(filepath.Join(blueprintsDir, sourceFilename)); err == nil {
		sourceBytes, err := os.ReadFile(filepath.Join(blueprintsDir, sourceFilename))
		if err == nil {
			source = string(sourceBytes)
		}
	}
	return source
}

func getRelPath(blueprintsDir string) string {
	relPath := ""
	// if it doesn't exist, continue
	if _, err := os.Stat(filepath.Join(blueprintsDir, relPathFilename)); err == nil {
		relPathBytes, err := os.ReadFile(filepath.Join(blueprintsDir, relPathFilename))
		if err == nil {
			relPath = string(relPathBytes)
		}
	}
	return relPath
}
