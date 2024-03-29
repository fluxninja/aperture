// +kubebuilder:validation:Optional
package config

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
}
