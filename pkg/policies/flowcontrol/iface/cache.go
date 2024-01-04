package iface

import (
	"context"
	"sync"

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
	// LookupWait starts lookup for specified keys in cache. It does not wait for response. It takes flowcontrolv1.LookupRequest and returns flowcontrolv1.LookupResponse and result and global wait groups.
	LookupNoWait(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) (*flowcontrolv1.CacheLookupResponse, *sync.WaitGroup, *sync.WaitGroup)
	// LookupGlobal looks up specified keys in global cache. It takes flowcontrolv1.LookupRequest and returns flowcontrolv1.KeyLookupResponse.
	LookupGlobal(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) map[string]*flowcontrolv1.KeyLookupResponse
	// LookupGlobalNoWait starts lookup for specified keys in global cache. It does not wait for response. It takes flowcontrolv1.LookupRequest and returns flowcontrolv1.KeyLookupResponse and global wait group.
	LookupGlobalNoWait(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) (map[string]*flowcontrolv1.KeyLookupResponse, *sync.WaitGroup)
	// LookupResult looks up specified keys in result cache. It takes flowcontrolv1.LookupRequest and returns flowcontrolv1.KeyLookupResponse.
	LookupResult(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) *flowcontrolv1.KeyLookupResponse
	// LookupResultNoWait starts lookup for specified keys in result cache. It does not wait for response. It takes flowcontrolv1.LookupRequest and returns flowcontrolv1.KeyLookupResponse and result wait group.
	LookupResultNoWait(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) (*flowcontrolv1.KeyLookupResponse, *sync.WaitGroup)
	// Upsert inserts or updates specified cache entries. It takes flowcontrolv1.UpsertRequest and returns flowcontrolv1.UpsertResponse.
	Upsert(ctx context.Context, req *flowcontrolv1.CacheUpsertRequest) *flowcontrolv1.CacheUpsertResponse
	// Delete deletes specified keys from cache. It takes flowcontrolv1.DeleteRequest and returns flowcontrolv1.DeleteResponse.
	Delete(ctx context.Context, req *flowcontrolv1.CacheDeleteRequest) *flowcontrolv1.CacheDeleteResponse
}
