package ratetracker

import (
	"bytes"
	"encoding/gob"
)

const (
	addFunction = "Add"
)

type counterState struct {
	Value float64 `json:"value"`
}

func add(_ string, currentState, arg []byte) (newState []byte, result []byte, err error) {
	// unmarshal currentState
	var cs counterState
	if currentState != nil {
		buf := bytes.NewBuffer(currentState)
		err = gob.NewDecoder(buf).Decode(&cs)
		if err != nil {
			return nil, nil, err
		}
	}

	// unmarshal arg into float64
	var i float64
	buf := bytes.NewBuffer(arg)
	err = gob.NewDecoder(buf).Decode(&i)
	if err != nil {
		return nil, nil, err
	}

	// add the integer to the counter
	cs.Value += i

	// marshal the new state
	buf = new(bytes.Buffer)
	err = gob.NewEncoder(buf).Encode(cs)
	if err != nil {
		return nil, nil, err
	}
	newState = buf.Bytes()

	// marshal the result
	buf = new(bytes.Buffer)
	err = gob.NewEncoder(buf).Encode(cs.Value)
	if err != nil {
		return nil, nil, err
	}
	result = buf.Bytes()

	return newState, result, nil
}
