// +kubebuilder:validation:Optional
package config

import "github.com/buraksezer/olric/config"

// swagger:operation POST /dist_cache common-configuration DistCache
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/DistCacheConfig"

// DistCacheConfig configures distributed cache that holds per-label counters in distributed rate limiters.
// swagger:model
// +kubebuilder:object:generate=true
type DistCacheConfig struct {
	// BindAddr denotes the address that DistCache will bind to for communication with other peer nodes.
	BindAddr string `json:"bind_addr" default:":3320" validate:"hostname_port"`
	// Address to bind [`memberlist`](https://github.com/hashicorp/memberlist) server to.
	MemberlistBindAddr string `json:"memberlist_bind_addr" default:":3322" validate:"hostname_port"`
	// Address of [`memberlist`](https://github.com/hashicorp/memberlist) to advertise to other cluster members. Used for NAT traversal if provided.
	MemberlistAdvertiseAddr string `json:"memberlist_advertise_addr" validate:"omitempty,hostname_port"`
	// ReplicaCount is 1 by default.
	ReplicaCount int `json:"replica_count" default:"1"`
	// SyncReplication enables synchronous replication. By default the replication is asynchronous.
	SyncReplication bool `json:"sync_replication" default:"false"`
	// Maximum in-memory usage for the given cache node
	MaxInUse int `json:"max_in_use" default:"0"`
	// Maximum Key count for the given cache node
	MaxKeys int `json:"maximum_keys" default:"0"`
	// Eviction Policy determines whether LRU policy should be used or not.
	EvictionPolicy config.EvictionPolicy `json:"eviction_policy" validate:"oneof=NONE LRU" default:"NONE"`
	// Number of evicitions workers to run in parallel.
	NumEvictionWorkers int `json:"num_eviction_workers" validate:"gte=0" default:"1"`
	// Number of keys to randomly select for the approximated lru implementation
	LRUSamples int `json:"lru_samples" validate:"gte=1" default:"5"`
}
