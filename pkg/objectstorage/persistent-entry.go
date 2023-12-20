package objectstorage

import olricstorage "github.com/buraksezer/olric/pkg/storage"

// PersistentEntry is an implementation of Olric entry used for persistent storage.
type PersistentEntry struct {
	key       string
	value     *[]byte
	ttl       int64
	timestamp int64
}

// SetKey is a key setter.
func (p *PersistentEntry) SetKey(s string) {
	p.key = s
}

// Key is a key getter.
func (p *PersistentEntry) Key() string {
	return p.key
}

// SetValue is a value setter.
func (p *PersistentEntry) SetValue(bytes []byte) {
	p.value = &bytes
}

// Value is a value getter.
func (p *PersistentEntry) Value() []byte {
	return *p.value
}

// SetTTL is a ttl setter.
func (p *PersistentEntry) SetTTL(i int64) {
	p.ttl = i
}

// TTL is a ttl getter.
func (p *PersistentEntry) TTL() int64 {
	return p.ttl
}

// SetTimestamp is a timestamp setter.
func (p *PersistentEntry) SetTimestamp(i int64) {
	p.timestamp = i
}

// Timestamp is a timestamp getter.
func (p *PersistentEntry) Timestamp() int64 {
	return p.timestamp
}

// SetLastAccess is a last access setter.
func (p *PersistentEntry) SetLastAccess(i int64) {
	panic("implement me")
}

// LastAccess is a last access getter.
func (p *PersistentEntry) LastAccess() int64 {
	panic("implement me")
}

// Encode encodes entry.
func (p *PersistentEntry) Encode() []byte {
	panic("implement me")
}

// Decode decodes entry.
func (p *PersistentEntry) Decode(bytes []byte) {
	panic("implement me")
}

var _ olricstorage.Entry = (*PersistentEntry)(nil)
