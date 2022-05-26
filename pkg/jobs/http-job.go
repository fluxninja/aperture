package jobs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// HTTPJob wraps a basic job along with HTTPJobConfig to execute an HTTP job.
type HTTPJob struct {
	BasicJob
	config HTTPJobConfig
}

// Make sure HTTPJob complies with Job interface.
var _ Job = (*HTTPJob)(nil)

// HTTPJobConfig is the configuration for an HTTP job.
type HTTPJobConfig struct {
	Client         *http.Client
	Body           BodyProvider
	URL            string
	Method         string
	ExpectedBody   string
	Name           string
	ExpectedStatus int
}

// BodyProvider allows the users to provide a body to the HTTP jobs. For example for posting a payload as a job.
type BodyProvider func() io.Reader

// NewHTTPJob creates a new HTTPJob.
func NewHTTPJob(config HTTPJobConfig) *HTTPJob {
	job := &HTTPJob{}
	job.config = config
	job.JobName = config.Name
	job.JobFunc = func(ctx context.Context) (proto.Message, error) {
		resp, err := fetchURL(ctx, config.Method, config.URL, config.Client, config.Body())
		if err != nil {
			return nil, err
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != config.ExpectedStatus {
			return nil, errors.Errorf("unexpected status code: %v, expected: %v", resp.StatusCode,
				config.ExpectedStatus)
		}

		if config.ExpectedBody != "" {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, errors.Errorf("failed to read response body: %v", err)
			}

			if !strings.Contains(string(body), config.ExpectedBody) {
				return nil, errors.Errorf("body does not contain expected content '%v'", config.ExpectedBody)
			}
		}

		message := wrapperspb.String(fmt.Sprintf("%s is accessible", config.URL))
		return message, nil
	}
	return job
}

// Name returns the name of the job.
func (job *HTTPJob) Name() string {
	return job.BasicJob.Name()
}

// JobWatchers returns the job watchers for the job.
func (job *HTTPJob) JobWatchers() JobWatchers {
	return job.BasicJob.JobWatchers()
}

// Execute executes the job.
func (job *HTTPJob) Execute(ctx context.Context) (proto.Message, error) {
	return job.BasicJob.Execute(ctx)
}

func fetchURL(ctx context.Context, method string, url string, client *http.Client, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, errors.Errorf("unable to create job http request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Errorf("fail to execute %v request: %v", method, err)
	}
	return resp, nil
}
