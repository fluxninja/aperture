package prometheus

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusmodel "github.com/prometheus/common/model"
	"google.golang.org/protobuf/proto"

	"github.com/FluxNinja/aperture/pkg/jobs"
	"github.com/FluxNinja/aperture/pkg/log"
)

// PromResultCallback is a callback that gets invoked with the result of the prom query.
type PromResultCallback func(context.Context, prometheusmodel.Value, ...interface{}) (proto.Message, error)

// PromErrorCallback is a callback that gets invoked when there's an error from running PromQL.
type PromErrorCallback func(error, ...interface{}) (proto.Message, error)

type promQuery struct {
	promAPI        prometheusv1.API
	resultCallback PromResultCallback
	errorCallback  PromErrorCallback
	cbArgs         []interface{}
	query          string
	timeout        time.Duration
	endTimestamp   time.Time
}

// NewPromQueryJob creates a new job that executes a prometheus query.
// It takes a PromResultCallback which gets invoked periodically with results of the query and
// an error callback which gets invoked when there's an error from running PromQL then it returns
// a callback compatible with scheduler BasicJob.
func NewPromQueryJob(
	query string,
	endTimestamp time.Time,
	promAPI prometheusv1.API,
	timeout time.Duration,
	resultCallback PromResultCallback,
	errorCallback PromErrorCallback,
	cbArgs ...interface{},
) jobs.JobCallback {
	pQuery := &promQuery{query: query, promAPI: promAPI, timeout: timeout, endTimestamp: endTimestamp, resultCallback: resultCallback, errorCallback: errorCallback, cbArgs: cbArgs}
	return pQuery.execute
}

func (pq *promQuery) execute(jobCtxt context.Context) (proto.Message, error) {
	var result prometheusmodel.Value
	var warnings prometheusv1.Warnings
	var msg proto.Message
	var err error

	operation := func() error {
		ctx, cancel := context.WithTimeout(jobCtxt, pq.timeout)
		defer cancel()

		result, warnings, err = pq.promAPI.Query(ctx, pq.query, pq.endTimestamp)
		if err != nil {
			log.Error().Err(err).Str("query", pq.query).Msg("Encountered error while executing promQL query")
			return err
		}
		for _, warning := range warnings {
			log.Warn().Str("query", pq.query).Str("warning", warning).Msg("Encountered warning while executing promQL query")
		}
		log.Trace().Str("query", pq.query).Time("end timestamp", pq.endTimestamp).Interface("result", result).Msg("Running prometheus query")
		return nil
	}

	err = backoff.Retry(operation, backoff.WithContext(backoff.NewExponentialBackOff(), jobCtxt))
	if err != nil {
		msg, err = pq.errorCallback(err)
		if err != nil {
			return msg, err
		}
	}

	return pq.resultCallback(jobCtxt, result, pq.cbArgs...)
}
