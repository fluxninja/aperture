package aperture

import (
	"errors"

	checkv1 "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/check/v1"
)

// LookupStatus is the status of a cache lookup, either HIT or MISS.
type LookupStatus string

// OperationStatus is the status of a cache operation, either SUCCESS or ERROR.
type OperationStatus string

const (
	// LookupStatusHit indicates that the cache lookup was a hit.
	LookupStatusHit LookupStatus = "HIT"
	// LookupStatusMiss indicates that the cache lookup was a miss.
	LookupStatusMiss LookupStatus = "MISS"
)

const (
	// OperationStatusSuccess indicates that the cache operation was successful.
	OperationStatusSuccess OperationStatus = "SUCCESS"
	// OperationStatusError indicates that the cache operation was unsuccessful.
	OperationStatusError OperationStatus = "ERROR"
)

// convertCacheLookupStatus converts CacheLookupStatus to LookupStatus.
func convertCacheLookupStatus(status checkv1.CacheLookupStatus) LookupStatus {
	switch status {
	case checkv1.CacheLookupStatus_HIT:
		return LookupStatusHit
	case checkv1.CacheLookupStatus_MISS:
		return LookupStatusMiss
	default:
		return LookupStatusMiss
	}
}

// convertCacheError converts a string error message to Go's error type.
// Returns nil if the input string is empty.
func convertCacheError(errorMessage string) error {
	if errorMessage == "" {
		return nil
	}
	return errors.New(errorMessage)
}

// KeyLookupResponse is the interface to read the response from a get cached value operation.
type KeyLookupResponse interface {
	Value() []byte
	LookupStatus() LookupStatus
	Error() error
}

// KeyUpsertResponse is the interface to read the response from a set cached value operation.
type KeyUpsertResponse interface {
	Error() error
}

// KeyDeleteResponse is the interface to read the response from a delete cached value operation.
type KeyDeleteResponse interface {
	Error() error
}

type keyLookupResponse struct {
	error           error
	lookupStatus    LookupStatus
	operationStatus OperationStatus
	value           []byte
}

type keyUpsertResponse struct {
	error           error
	operationStatus OperationStatus
}

type keyDeleteResponse struct {
	error           error
	operationStatus OperationStatus
}

// Value returns the cached value.
func (g *keyLookupResponse) Value() []byte {
	return g.value
}

// LookupStatus returns the lookup status.
func (g *keyLookupResponse) LookupStatus() LookupStatus {
	return g.lookupStatus
}

// OperationStatus returns the operation status.
func (g *keyLookupResponse) OperationStatus() OperationStatus {
	return g.operationStatus
}

// Error returns the error.
func (g *keyLookupResponse) Error() error {
	return g.error
}

// Error returns the error.
func (s *keyUpsertResponse) Error() error {
	return s.error
}

// OperationStatus returns the operation status.
func (s *keyUpsertResponse) OperationStatus() OperationStatus {
	return s.operationStatus
}

// Error returns the error.
func (d *keyDeleteResponse) Error() error {
	return d.error
}

// OperationStatus returns the operation status.
func (d *keyDeleteResponse) OperationStatus() OperationStatus {
	return d.operationStatus
}

func newKeyLookupResponse(value []byte, lookupStatus LookupStatus, err error) KeyLookupResponse {
	return &keyLookupResponse{
		value:        value,
		lookupStatus: lookupStatus,
		error:        err,
	}
}

func newKeyUpsertResponse(err error) KeyUpsertResponse {
	return &keyUpsertResponse{
		error: err,
	}
}

func newKeyDeleteResponse(err error) KeyDeleteResponse {
	return &keyDeleteResponse{
		error: err,
	}
}
