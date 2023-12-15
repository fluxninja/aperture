package objectstorage

import olricstorage "github.com/buraksezer/olric/pkg/storage"

type PersistedEntry struct {
	key   string
	value *[]byte
	ttl   int64
}

func (p PersistedEntry) SetKey(s string) {
	p.key = s
}

func (p PersistedEntry) Key() string {
	return p.key
}

func (p PersistedEntry) SetValue(bytes []byte) {
	p.value = &bytes
}

func (p PersistedEntry) Value() []byte {
	return *p.value
}

func (p PersistedEntry) SetTTL(i int64) {
	p.ttl = i
}

func (p PersistedEntry) TTL() int64 {
	return p.ttl
}

func (p PersistedEntry) SetTimestamp(i int64) {

	panic("implement me")
}

func (p PersistedEntry) Timestamp() int64 {

	panic("implement me")
}

func (p PersistedEntry) SetLastAccess(i int64) {

	panic("implement me")
}

func (p PersistedEntry) LastAccess() int64 {

	panic("implement me")
}

func (p PersistedEntry) Encode() []byte {

	panic("implement me")
}

func (p PersistedEntry) Decode(bytes []byte) {

	panic("implement me")
}

var _ olricstorage.Entry = (*PersistedEntry)(nil)
