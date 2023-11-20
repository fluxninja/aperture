package aperture

import (
	"errors"

	checkv1 "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/check/v1"
)

// ErrCacheKeyNotSet is returned when empty cache key is provided by the caller.
var ErrCacheKeyNotSet = errors.New("cache key not set")

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

// convertCacheOperationStatus converts CacheOperationStatus to OperationStatus.
func convertCacheOperationStatus(status checkv1.CacheOperationStatus) OperationStatus {
	switch status {
	case checkv1.CacheOperationStatus_SUCCESS:
		return OperationStatusSuccess
	case checkv1.CacheOperationStatus_ERROR:
		return OperationStatusError
	default:
		return OperationStatusError
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

// GetCachedValueResponse is the interface to read the response from a get cached value operation.
type GetCachedValueResponse interface {
	GetValue() []byte
	GetLookupStatus() LookupStatus
	GetOperationStatus() OperationStatus
	GetError() error
}

// SetCachedValueResponse is the interface to read the response from a set cached value operation.
type SetCachedValueResponse interface {
	GetError() error
	GetOperationStatus() OperationStatus
}

// DeleteCachedValueResponse is the interface to read the response from a delete cached value operation.
type DeleteCachedValueResponse interface {
	GetError() error
	GetOperationStatus() OperationStatus
}

type getCachedValueResponse struct {
	value           []byte
	lookupStatus    LookupStatus
	operationStatus OperationStatus
	error           error
}

type setCachedValueResponse struct {
	error           error
	operationStatus OperationStatus
}

type deleteCachedValueResponse struct {
	error           error
	operationStatus OperationStatus
}

// GetValue returns the cached value.
func (g *getCachedValueResponse) GetValue() []byte {
	return g.value
}

// GetLookupStatus returns the lookup status.
func (g *getCachedValueResponse) GetLookupStatus() LookupStatus {
	return g.lookupStatus
}

// GetOperationStatus returns the operation status.
func (g *getCachedValueResponse) GetOperationStatus() OperationStatus {
	return g.operationStatus
}

// GetError returns the error.
func (g *getCachedValueResponse) GetError() error {
	return g.error
}

// GetError returns the error.
func (s *setCachedValueResponse) GetError() error {
	return s.error
}

// GetOperationStatus returns the operation status.
func (s *setCachedValueResponse) GetOperationStatus() OperationStatus {
	return s.operationStatus
}

// GetError returns the error.
func (d *deleteCachedValueResponse) GetError() error {
	return d.error
}

// GetOperationStatus returns the operation status.
func (d *deleteCachedValueResponse) GetOperationStatus() OperationStatus {
	return d.operationStatus
}

func newGetCachedValueResponse(value []byte, lookupStatus LookupStatus, operationStatus OperationStatus, err error) GetCachedValueResponse {
	return &getCachedValueResponse{
		value:           value,
		lookupStatus:    lookupStatus,
		operationStatus: operationStatus,
		error:           err,
	}
}

func newSetCachedValueResponse(operationStatus OperationStatus, err error) SetCachedValueResponse {
	return &setCachedValueResponse{
		error:           err,
		operationStatus: operationStatus,
	}
}

func newDeleteCachedValueResponse(operationStatus OperationStatus, err error) DeleteCachedValueResponse {
	return &deleteCachedValueResponse{
		error:           err,
		operationStatus: operationStatus,
	}
}
