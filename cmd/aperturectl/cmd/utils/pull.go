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

func PullSource(dir, uri string) error {
	err := WriterLock(dir)
	if err != nil {
		return err
	}
	defer Unlock(dir)

	spec := specv1.New()
	contents, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	spec.LegacyImports = false
	contents = append(contents, []byte("\n")...)

	filename := filepath.Join(dir, jsonnetfile.File)
	err = os.WriteFile(filename, contents, os.ModePerm)
	if err != nil {
		return err
	}

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

	if !depEqual(spec.Dependencies[d.Name()], *d) {
		spec.Dependencies[d.Name()] = *d
		delete(lockFile.Dependencies, d.Name())
	}

	locked, err := pkg.Ensure(spec, dir, lockFile.Dependencies)
	if err != nil {
		return err
	}

	err = writeChangedJsonnetFile(contents, &spec, filename)
	if err != nil {
		return err
	}
	err = writeChangedJsonnetFile(jbLockFileBytes, &specv1.JsonnetFile{Dependencies: locked}, filepath.Join(dir, jsonnetfile.LockFile))
	if err != nil {
		return err
	}

	return nil
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
