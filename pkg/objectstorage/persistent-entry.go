package objectstorage

import (
	"encoding/binary"

	olricstorage "github.com/buraksezer/olric/pkg/storage"
)

// PersistentEntry is an implementation of Olric entry used for persistent storage.
type PersistentEntry struct {
	key        string
	value      []byte
	ttl        int64
	timestamp  int64
	lastAccess int64
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
	p.value = bytes
}

// Value is a value getter.
func (p *PersistentEntry) Value() []byte {
	return p.value
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
	p.lastAccess = i
}

// LastAccess is a last access getter.
func (p *PersistentEntry) LastAccess() int64 {
	return p.lastAccess
}

// Encode encodes entry.
func (p *PersistentEntry) Encode() []byte {
	var offset int

	klen := uint8(len(p.Key()))
	vlen := len(p.Value())
	length := 29 + len(p.Key()) + vlen

	buf := make([]byte, length)

	// Set key length. It's 1 bytp.
	copy(buf[offset:], []byte{klen})
	offset++

	// Set the key.
	copy(buf[offset:], p.Key())
	offset += len(p.Key())

	// Set the TTL. It's 8 bytes.
	binary.BigEndian.PutUint64(buf[offset:], uint64(p.TTL()))
	offset += 8

	// Set the Timestamp. It's 8 bytes.
	binary.BigEndian.PutUint64(buf[offset:], uint64(p.Timestamp()))
	offset += 8

	// Set the LastAccess. It's 8 bytes.
	binary.BigEndian.PutUint64(buf[offset:], uint64(p.LastAccess()))
	offset += 8

	// Set the value length. It's 4 bytes.
	binary.BigEndian.PutUint32(buf[offset:], uint32(len(p.Value())))
	offset += 4

	// Set the value.
	copy(buf[offset:], p.Value())
	return buf
}

// Decode decodes entry.
func (p *PersistentEntry) Decode(buf []byte) {
	var offset int

	keyLength := int(buf[offset])
	offset++

	p.key = string(buf[offset : offset+keyLength])
	offset += keyLength

	p.ttl = int64(binary.BigEndian.Uint64(buf[offset : offset+8]))
	offset += 8

	p.timestamp = int64(binary.BigEndian.Uint64(buf[offset : offset+8]))
	offset += 8

	p.lastAccess = int64(binary.BigEndian.Uint64(buf[offset : offset+8]))
	offset += 8

	vlen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	p.value = buf[offset : offset+int(vlen)]
}

var _ olricstorage.Entry = (*PersistentEntry)(nil)
