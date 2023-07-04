package prometheus

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusmodel "github.com/prometheus/common/model"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// PromResultCallback is a callback that gets invoked with the result of the prom query.
type PromResultCallback func(context.Context, prometheusmodel.Value, ...interface{}) (proto.Message, error)

// PromErrorCallback is a callback that gets invoked when there's an error from running PromQL.
type PromErrorCallback func(error, ...interface{}) (proto.Message, error)

type promQuery struct {
	endTimestamp   time.Time
	promAPI        prometheusv1.API
	enforcer       *PrometheusEnforcer
	resultCallback PromResultCallback
	errorCallback  PromErrorCallback
	query          string
	cbArgs         []interface{}
	timeout        time.Duration
}

// NewPromQueryJob creates a new job that executes a prometheus query.
// It takes a PromResultCallback which gets invoked periodically with results of the query and
// an error callback which gets invoked when there's an error from running PromQL then it returns
// a callback compatible with scheduler BasicJob.
func NewPromQueryJob(
	query string,
	endTimestamp time.Time,
	promAPI prometheusv1.API,
	enforcer *PrometheusEnforcer,
	timeout time.Duration,
	resultCallback PromResultCallback,
	errorCallback PromErrorCallback,
	cbArgs ...interface{},
) jobs.JobCallback {
	pQuery := &promQuery{query: query, promAPI: promAPI, enforcer: enforcer, timeout: timeout, endTimestamp: endTimestamp, resultCallback: resultCallback, errorCallback: errorCallback, cbArgs: cbArgs}
	return pQuery.execute
}

func (pq *promQuery) execute(jobCtxt context.Context) (proto.Message, error) {
	var result prometheusmodel.Value
	var warnings prometheusv1.Warnings
	var err error

	operation := func() error {
		ctx, cancel := context.WithTimeout(jobCtxt, pq.timeout)
		defer cancel()

		query, innerErr := pq.enforcer.EnforceLabels(pq.query)
		if innerErr != nil {
			return innerErr
		}

		result, warnings, err = pq.promAPI.Query(ctx, query, pq.endTimestamp)
		// if jobCtxt is closed, return PermanentError
		if jobCtxt.Err() != nil {
			return backoff.Permanent(jobCtxt.Err())
		}
		if err != nil {
			log.Error().Err(err).Str("query", query).Msg("Encountered error while executing promQL query")
			return err
		}
		for _, warning := range warnings {
			log.Warn().Str("query", query).Str("warning", warning).Msg("Encountered warning while executing promQL query")
		}
		log.Trace().Str("query", query).Time("end timestamp", pq.endTimestamp).Interface("result", result).Msg("Running prometheus query")
		return nil
	}

	merr := backoff.Retry(operation, backoff.WithContext(backoff.NewExponentialBackOff(), jobCtxt))
	if merr != nil {
		msg, cbErr := pq.errorCallback(err)
		if cbErr != nil {
			merr = multierr.Combine(merr, cbErr)
		}
		log.Error().Err(merr).Msg("Context canceled while executing promQL query")
		return msg, merr
	}

	return pq.resultCallback(jobCtxt, result, pq.cbArgs...)
}
