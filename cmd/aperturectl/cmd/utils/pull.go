package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	specv1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
)

const (
	sourceFilename  = ".source"
	versionFilename = ".version"
	relPathFilename = ".relpath"
)

// PullSource pulls the source of the dependency and updates the lock file.
func PullSource(dir, uri string) error {
	d := deps.Parse("", uri)
	if d == nil {
		return errors.New("unable to parse URI: " + uri)
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
		version = "local"
	} else {
		return errors.New("unable to parse URI: " + uri)
	}

	spec := specv1.New()
	spec.LegacyImports = false
	spec.Dependencies[d.Name()] = *d

	jbLockFileBytes, err := os.ReadFile(filepath.Join(dir, jsonnetfile.LockFile))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lockFile, err := jsonnetfile.Unmarshal(jbLockFileBytes)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(dir, ".tmp"), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(dir, sourceFilename), []byte(source), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(dir, relPathFilename), []byte(relPath), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(dir, versionFilename), []byte(version), os.ModePerm)
	if err != nil {
		return err
	}

	locked, err := pkg.Ensure(spec, dir, lockFile.Dependencies)
	if err != nil {
		return err
	}

	err = writeChangedJsonnetFile(jbLockFileBytes, &specv1.JsonnetFile{Dependencies: locked}, filepath.Join(dir, jsonnetfile.LockFile))
	if err != nil {
		return err
	}

	return nil
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

// GetSource returns the source of the dependency.
func GetSource(dir string) string {
	source := ""
	// if it does not exist, continue
	if _, err := os.Stat(filepath.Join(dir, sourceFilename)); err == nil {
		sourceBytes, err := os.ReadFile(filepath.Join(dir, sourceFilename))
		if err == nil {
			source = string(sourceBytes)
		}
	}
	return source
}

// GetVersion returns the version of the dependency.
func GetVersion(dir string) string {
	version := ""
	// if it does not exist, continue
	if _, err := os.Stat(filepath.Join(dir, versionFilename)); err == nil {
		versionBytes, err := os.ReadFile(filepath.Join(dir, versionFilename))
		if err == nil {
			version = string(versionBytes)
		}
	}
	return version
}

// GetRelPath returns the relative path to the dependency.
func GetRelPath(dir string) string {
	relPath := ""
	// if it does not exist, continue
	if _, err := os.Stat(filepath.Join(dir, relPathFilename)); err == nil {
		relPathBytes, err := os.ReadFile(filepath.Join(dir, relPathFilename))
		if err == nil {
			relPath = string(relPathBytes)
		}
	}
	return relPath
}
