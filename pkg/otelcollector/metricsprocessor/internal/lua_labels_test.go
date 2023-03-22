package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestAddLuaSpecificLabels(t *testing.T) {
	tests := []struct {
		name       string
		attributes pcommon.Map
		expected   pcommon.Map
	}{
		{
			name: "Adds HTTPRequestContentLength and HTTPResponseContentLength labels",
			attributes: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble("BYTES_SENT", 10)
				m.PutDouble("BYTES_RECEIVED", 20)
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble("BYTES_SENT", 10)
				m.PutDouble("BYTES_RECEIVED", 20)
				m.PutDouble("http.request_content_length", 10)
				m.PutDouble("http.response_content_length", 20)
				m.PutStr("response_received", "false")
				return m
			}(),
		},
		{
			name: "Adds ResponseReceivedLabel with value true",
			attributes: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble("RESPONSE_DURATION", 1)
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble("RESPONSE_DURATION", 1)
				m.PutStr("response_received", "true")
				m.PutDouble("flow_duration_ms", 1)
				return m
			}(),
		},
		{
			name: "Adds ResponseReceivedLabel with value false",
			attributes: func() pcommon.Map {
				m := pcommon.NewMap()
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutStr("response_received", "false")
				return m
			}(),
		},
		{
			name: "Adds FlowDurationLabel",
			attributes: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble("RESPONSE_DURATION", 2)
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble("RESPONSE_DURATION", 2)
				m.PutStr("response_received", "true")
				m.PutDouble("flow_duration_ms", 2)
				return m
			}(),
		},
		{
			name: "Adds WorkloadDurationLabel",
			attributes: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble("RESPONSE_DURATION", 3)
				m.PutDouble("checkhttp_duration", 1)
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble("RESPONSE_DURATION", 3)
				m.PutDouble("checkhttp_duration", 1)
				m.PutStr("response_received", "true")
				m.PutDouble("flow_duration_ms", 3)
				m.PutDouble("workload_duration_ms", 2)
				return m
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddLuaSpecificLabels(tt.attributes)
			assert.Equal(t, tt.expected, tt.attributes)
		})
	}
}
