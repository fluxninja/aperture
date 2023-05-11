package utils

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/gofrs/flock"
)

const lockFilename = ".flock"

var lock *flock.Flock

func newLock(dir string) {
	lock = flock.New(filepath.Join(dir, lockFilename))
}

// WriterLock acquires a writer lock on the directory.
func WriterLock(dir string) error {
	if lock == nil {
		newLock(dir)
	}
	// Get writer lock
	locked, err := lock.TryLockContext(context.Background(), 10)
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("unable to acquire lock on directory")
	}
	return nil
}

// ReaderLock acquires a reader lock on the directory.
func ReaderLock(dir string) error {
	if lock == nil {
		newLock(dir)
	}
	// Get reader lock
	locked, err := lock.TryRLockContext(context.Background(), 10)
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("unable to acquire lock on directory")
	}
	return nil
}

// Unlock releases the lock on the directory.
func Unlock(dir string) {
	if lock == nil {
		newLock(dir)
	}
	err := lock.Unlock()
	if err != nil {
		log.Error().Err(err).Msg("unable to release lock on directory")
	}
}
