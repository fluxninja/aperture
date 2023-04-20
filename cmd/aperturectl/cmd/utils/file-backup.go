package utils

import (
	"os"
)

// BackupFile backs up a file. If backup exists, it will be restored first.
func BackupFile(file string) error {
	var stat os.FileInfo
	var err error
	// check if backup exists and if it does, then restore it
	// as we might not have exited cleanly last time
	RestoreFile(file)

	backupFile := file + ".bak"
	// make sure the file exists
	if stat, err = os.Stat(file); err != nil {
		return err
	}
	// get file permissions
	mode := stat.Mode()

	// read the file
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// write the file to backup
	err = os.WriteFile(backupFile, fileBytes, mode)
	if err != nil {
		return err
	}

	return nil
}

// RestoreFile restores a file from its backup.
func RestoreFile(file string) {
	backupFile := file + ".bak"
	if _, err := os.Stat(backupFile); err == nil {
		// backup exists, restore it
		_ = os.Rename(backupFile, file)
	}
}
