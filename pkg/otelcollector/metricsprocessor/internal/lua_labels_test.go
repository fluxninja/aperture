package internal

import (
	"testing"

	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
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
				m.PutDouble(otelconsts.BytesSentLabel, 10)
				m.PutDouble(otelconsts.BytesReceivedLabel, 20)
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble(otelconsts.BytesSentLabel, 10)
				m.PutDouble(otelconsts.BytesReceivedLabel, 20)
				m.PutDouble(otelconsts.HTTPRequestContentLength, 10)
				m.PutDouble(otelconsts.HTTPResponseContentLength, 20)
				m.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedFalse)
				return m
			}(),
		},
		{
			name: "Adds ResponseReceivedLabel with value true",
			attributes: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble(otelconsts.ApertureFlowEndTimestampLabel, 123456799)
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble(otelconsts.ApertureFlowEndTimestampLabel, 123456799)
				m.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedTrue)
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
				m.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedFalse)
				return m
			}(),
		},
		{
			name: "Adds FlowDurationLabel",
			attributes: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble(otelconsts.ApertureFlowEndTimestampLabel, 123456799)
				m.PutDouble(otelconsts.ApertureFlowStartTimestampLabel, 123456789)
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble(otelconsts.ApertureFlowEndTimestampLabel, 123456799)
				m.PutDouble(otelconsts.ApertureFlowStartTimestampLabel, 123456789)
				m.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedTrue)
				m.PutDouble(otelconsts.FlowDurationLabel, 10)
				return m
			}(),
		},
		{
			name: "Adds WorkloadDurationLabel",
			attributes: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble(otelconsts.ApertureFlowEndTimestampLabel, 123456799)
				m.PutDouble(otelconsts.ApertureFlowStartTimestampLabel, 123456789)
				m.PutDouble(otelconsts.ApertureWorkloadStartTimestampLabel, 123456790)
				return m
			}(),
			expected: func() pcommon.Map {
				m := pcommon.NewMap()
				m.PutDouble(otelconsts.ApertureFlowEndTimestampLabel, 123456799)
				m.PutDouble(otelconsts.ApertureFlowStartTimestampLabel, 123456789)
				m.PutDouble(otelconsts.ApertureWorkloadStartTimestampLabel, 123456790)
				m.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedTrue)
				m.PutDouble(otelconsts.FlowDurationLabel, 10)
				m.PutDouble(otelconsts.WorkloadDurationLabel, 9)
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
