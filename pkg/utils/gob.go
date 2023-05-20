package utils

import (
	"bytes"
	"encoding/gob"
)

// MarshalGob encodes an interface into a byte array.
func MarshalGob(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalGob decodes a byte array into an interface.
func UnmarshalGob(data []byte, value interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(value)
	if err != nil {
		return err
	}
	return nil
}
