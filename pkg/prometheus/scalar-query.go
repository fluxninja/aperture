package prometheus

import (
	"context"
	"fmt"
	"time"

	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusmodel "github.com/prometheus/common/model"
	"google.golang.org/protobuf/proto"

	"aperture.tech/aperture/pkg/jobs"
	"aperture.tech/aperture/pkg/log"
)

// ScalarResultCallback is a callback that gets invoked with the result of the scalar query.
type ScalarResultCallback func(context.Context, float64, ...interface{}) (proto.Message, error)

// ScalarQuery is a wrapper that holds prometheus query and the ScalarResultCallback that returns the result of the scalar query.
type ScalarQuery struct {
	scalarResultCallback ScalarResultCallback
	query                string
}

// NewScalarQueryJob creates a new job that executes a prometheus query job with given scalar query.
func NewScalarQueryJob(
	query string,
	endTimestamp time.Time,
	promAPI prometheusv1.API,
	timeout time.Duration,
	resultCallback ScalarResultCallback,
	errorCallback PromErrorCallback,
	cbArgs ...interface{},
) jobs.JobCallback {
	sq := &ScalarQuery{scalarResultCallback: resultCallback, query: query}

	return NewPromQueryJob(query, endTimestamp, promAPI, timeout, sq.execute, errorCallback, cbArgs...)
}

func (sq *ScalarQuery) execute(ctx context.Context, value prometheusmodel.Value, cbArgs ...interface{}) (proto.Message, error) {
	log.Trace().Msg("ScalarQuery execute")
	if scalar, ok := value.(*prometheusmodel.Scalar); ok {
		return sq.scalarResultCallback(ctx, float64(scalar.Value), cbArgs...)
	} else if vector, ok := value.(prometheusmodel.Vector); ok {
		if len(vector) == 0 {
			return nil, fmt.Errorf("no data returned for query: %s", sq.query)
		} else if vector.Len() == 1 {
			return sq.scalarResultCallback(ctx, float64(vector[0].Value), cbArgs...)
		} else {
			return nil, fmt.Errorf("query returned a vector with %d elements, expecting only 1 element. query: %s", vector.Len(), sq.query)
		}
	}
	return nil, fmt.Errorf("query returned non-scalar value: %v. query string: %s", value, sq.query)
}
