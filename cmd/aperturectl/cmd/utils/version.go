package utils

import (
	"bytes"
	"os"
	"path/filepath"

	"golang.org/x/mod/semver"
)

var versionFilePath string

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	versionFilePath = filepath.Join(userHomeDir, ".aperturectl", "version")
}

// CreateVersionFileIfNotExists creates a version file if it does not exist.
func CreateVersionFileIfNotExists(version string) error {
	if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
		versionFile, err := os.Create(versionFilePath)
		if err != nil {
			return err
		}
		defer versionFile.Close()

		_, err = versionFile.WriteString(version)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsCurrentVersionNewer checks if the version of aperturectl is newer than on disk.
func IsCurrentVersionNewer(version string) (bool, error) {
	versionFile, err := os.Open(versionFilePath)
	if err != nil {
		return false, err
	}
	defer versionFile.Close()

	// read version from file into bytes buffer
	versionOnFile := new(bytes.Buffer)
	_, err = versionOnFile.ReadFrom(versionFile)
	if err != nil {
		return false, err
	}

	return semver.Compare(versionOnFile.String(), version) <= 0, nil
}
