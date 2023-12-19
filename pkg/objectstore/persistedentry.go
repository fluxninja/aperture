package objectstorage

import olricstorage "github.com/buraksezer/olric/pkg/storage"

// PersistedEntry is an implementation of Olric entry used for persistent storage.
type PersistedEntry struct {
	key       string
	value     *[]byte
	ttl       int64
	timestamp int64
}

// SetKey is a key setter.
func (p *PersistedEntry) SetKey(s string) {
	p.key = s
}

// Key is a key getter.
func (p *PersistedEntry) Key() string {
	return p.key
}

// SetValue is a value setter.
func (p *PersistedEntry) SetValue(bytes []byte) {
	p.value = &bytes
}

// Value is a value getter.
func (p *PersistedEntry) Value() []byte {
	return *p.value
}

// SetTTL is a ttl setter.
func (p *PersistedEntry) SetTTL(i int64) {
	p.ttl = i
}

// TTL is a ttl getter.
func (p *PersistedEntry) TTL() int64 {
	return p.ttl
}

// SetTimestamp is a timestamp setter.
func (p *PersistedEntry) SetTimestamp(i int64) {
	p.timestamp = i
}

// Timestamp is a timestamp getter.
func (p *PersistedEntry) Timestamp() int64 {
	return p.timestamp
}

// SetLastAccess is a last access setter.
func (p *PersistedEntry) SetLastAccess(i int64) {
	panic("implement me")
}

// LastAccess is a last access getter.
func (p *PersistedEntry) LastAccess() int64 {
	panic("implement me")
}

// Encode encodes entry.
func (p *PersistedEntry) Encode() []byte {
	panic("implement me")
}

// Decode decodes entry.
func (p *PersistedEntry) Decode(bytes []byte) {
	panic("implement me")
}

var _ olricstorage.Entry = (*PersistedEntry)(nil)
