package iface

import (
	"context"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
)

// CacheType is the type of cache.
type CacheType int

//go:generate enumer -type=CacheType -transform=lower -output=cache-type-string.go
const (
	// Result is the type of cache for saving results.
	Result CacheType = iota
	// CacheTypeState is the type of cache for saving state.
	Global
)

// Cache is an interface for the cache.
type Cache interface {
	// Lookup looks up specified keys in cache. It takes flowcontrolv1.LookupRequest and returns flowcontrolv1.LookupResponse.
	Lookup(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) *flowcontrolv1.CacheLookupResponse
	// Upsert inserts or updates specified cache entries. It takes flowcontrolv1.UpsertRequest and returns flowcontrolv1.UpsertResponse.
	Upsert(ctx context.Context, req *flowcontrolv1.CacheUpsertRequest) *flowcontrolv1.CacheUpsertResponse
	// Delete deletes specified keys from cache. It takes flowcontrolv1.DeleteRequest and returns flowcontrolv1.DeleteResponse.
	Delete(ctx context.Context, req *flowcontrolv1.CacheDeleteRequest) *flowcontrolv1.CacheDeleteResponse
}
