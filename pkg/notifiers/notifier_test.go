package notifiers

import (
	"testing"
	"time"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	pns     []PrefixNotifier
	kns     []KeyNotifier
	evs     []Event
	results [][]byte
}

func createKeyNotifier(t *testing.T, key string, value []byte) (PrefixNotifier, KeyNotifier) {
	bpn := BasicPrefixNotifier{
		NotifyFunc: notifierFunc,
	}
	bkn := bpn.GetKeyNotifier(Key(key))
	bkn.SetKey(Key(key))
	bkn.setID(key + "-" + string(value))
	bkn.SetTransformFunc(transformFunc)

	t.Log("KeyNotifier with ID: ", bkn.getID(), ", Func Address: ", bkn.GetTransformFunc(), ", Key: ", bkn.GetKey(), " has been created.")

	return &bpn, bkn
}

func createEvent(t *testing.T, key Key, value []byte, action EventType) Event {
	event := Event{
		Key:   Key(key),
		Value: value,
		Type:  action,
	}

	t.Log("Event: ", event.String(), " has been created.")

	return event
}

// notifierFunc just prints the action that is being performed on a given notifier.
// Recalled event.String() for getting more coverage but printing the function leads to numerous log undesired log messages.
func notifierFunc(event Event) {
	event.String()
}

// transformFunc is a mock func that reverses the content of any bytes given to it if the action is valid.
func transformFunc(key Key, bytes []byte, eventType EventType) (Key, []byte, error) {
	if eventType.IsAEventType() {
		for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
			bytes[i], bytes[j] = bytes[j], bytes[i]
		}
	}
	return key, bytes, nil
}

func runTests(t *testing.T, config testConfig) {
	dt := NewDefaultTrackers()

	err := dt.Start()
	require.NoError(t, err)

	for i, bpn := range config.pns {
		err = dt.AddPrefixNotifier(bpn)
		require.NoError(t, err)

		err = dt.AddKeyNotifier(config.kns[i])
		require.NoError(t, err)

		dt.WriteEvent(config.kns[i].GetKey(), config.evs[i].Value)
	}

	// Allow some sleep so that events can be executed
	time.Sleep(time.Millisecond * 100)

	// Check results by matching to the expected results
	for i, kns := range config.kns {
		bytes := dt.GetCurrentValue(kns.GetKey())
		require.Equal(t, config.results[i], bytes)
	}

	// if multiple key notifiers were implemented just removed by using Purge func
	if len(config.kns) > 1 {
		dt.Purge("")
	}

	for _, pns := range config.pns {
		err = dt.RemovePrefixNotifier(pns)
		require.NoError(t, err)
	}

	for _, kns := range config.kns {
		err = dt.RemoveKeyNotifier(kns)
		require.NoError(t, err)

		dt.RemoveEvent(kns.GetKey())
	}

	dt.Stop()
}

// TestBasicKeyNotifier creates key notifiers, prefix notifiers and events that are passed to default trackers and should reverse the content of the bytes given to the tracker.
func TestBasicKeyNotifier(t *testing.T) {
	bpn, bkn := createKeyNotifier(t, "Notifier-test-1", []byte("testBytes"))
	bpn2, bkn2 := createKeyNotifier(t, "", []byte("newTestBytes"))
	bkn2.setID("") // set ID empty to generate random id throw guuid
	bpn3, bkn3 := createKeyNotifier(t, "Notifier-test-3", []byte("newTestBytes"))
	bkn3.setID(bkn.getID()) // setting id same as basicKeyNotifier1 to check special case scenario when two notifiers have the same id

	event := createEvent(t, bkn.GetKey(), []byte("testBytes1"), 1)
	event2 := createEvent(t, bkn.GetKey(), []byte("testBytes2"), 1)
	event3 := createEvent(t, bkn3.GetKey(), []byte("testBytes3"), 2)

	config := testConfig{
		pns: []PrefixNotifier{bpn, bpn2, bpn3},
		kns: []KeyNotifier{bkn, bkn2, bkn3},
		evs: []Event{event, event2, event3},
		results: [][]byte{
			[]byte("1setyBtset"),
			[]byte("2setyBtset"),
			[]byte("3setyBtset"),
		},
	}

	runTests(t, config)
}

// wrapper func for unmarshaller notifier func, implementation of UnmarshalNotifyFunc
func unmarshalNotifyFunc(event Event, unmarshaller config.Unmarshaller) {
	notifierFunc(event)
}

// implementation of GetUnmarshallerFunc for unmarshaller notifier
func createUnmarshaller(bytes []byte) (config.Unmarshaller, error) {
	unmarshaller, err := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
	return unmarshaller, err
}

// TestUnmarshallerKeyNotifier creates unmarshaller prefix notifiers,  key notifiers and events that are passed to default trackers and should unmarshal a configuration given to the tracker.
func TestUnmarshallerNotifier(t *testing.T) {
	bytes := []byte(`
          configs:
            - Name: UnmarshallerNotifierTest
            - Value: 1
  `)
	upn := UnmarshalPrefixNotifier{
		UnmarshalNotifyFunc: unmarshalNotifyFunc,
		GetUnmarshallerFunc: createUnmarshaller,
	}
	kn := upn.GetKeyNotifier("Unmarshal-Notifier-Key")

	event := createEvent(t, kn.GetKey(), bytes, 1)

	config := testConfig{
		pns: []PrefixNotifier{&upn},
		kns: []KeyNotifier{kn},
		evs: []Event{event},
		results: [][]byte{
			bytes,
		},
	}
	runTests(t, config)
}

// TestUnmarshallerNotifierConstructor has the same functionality of TestUnmarshallerNotifier but uses the NewUnmarshalKeyNotifier constructor to create the unmarshaller.
func TestUnmarshallerNotifierConstructor(t *testing.T) {
	bytes := []byte(`
          configs:
            - Name: TestUnmarshallerNotifierConstructorTest
            - Value: 1`)
	upn := UnmarshalPrefixNotifier{
		UnmarshalNotifyFunc: unmarshalNotifyFunc,
		GetUnmarshallerFunc: createUnmarshaller,
	}

	unmarshaller, err := createUnmarshaller(bytes)
	require.NoError(t, err)

	kn := NewUnmarshalKeyNotifier("Unmarshal-Notifier-Key-2", unmarshaller, unmarshalNotifyFunc)

	event := createEvent(t, kn.GetKey(), bytes, 1)

	config := testConfig{
		pns: []PrefixNotifier{&upn},
		kns: []KeyNotifier{kn},
		evs: []Event{event},
		results: [][]byte{
			bytes,
		},
	}

	runTests(t, config)

	val := unmarshaller.IsSet("configs") // check if the unmarshaller has loaded the configuration correctly
	require.True(t, val)
}
