package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/log"
)

type recordedRequest struct {
	host        string
	requestBody string
}

type recorderHTTPClient struct {
	recordedRequests []recordedRequest
}

func newRecorderHTTPClient() *recorderHTTPClient {
	return &recorderHTTPClient{
		recordedRequests: nil,
	}
}

type requestHandlerFactory struct {
	hostname  string
	envoyPort int
}

func (f requestHandlerFactory) withRecorderHTTPClient(client *recorderHTTPClient) *RequestHandler {
	return &RequestHandler{
		hostname:   f.hostname,
		httpClient: client,
	}
}

func (c *recorderHTTPClient) Do(req *http.Request) (*http.Response, error) {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Panic().Err(err).Msgf("Failed to read all: %v", err)
	}

	rr := recordedRequest{
		host:        req.URL.Host,
		requestBody: string(requestBody),
	}
	c.recordedRequests = append(c.recordedRequests, rr)

	response := &http.Response{
		StatusCode: http.StatusOK,
	}

	return response, nil
}

func postRequest(requestHandler *RequestHandler, requestBody string) *http.Response {
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "localhost:1000/request", strings.NewReader(requestBody))
	requestHandler.ServeHTTP(responseRecorder, req)

	return responseRecorder.Result()
}

func TestSimpleService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simple Service Suite")
}

var _ = Describe("Simple Service request handler", func() {
	var (
		serviceHostname           = "A"
		someOtherHostname1        = "B"
		someOtherHostname2        = "C"
		someEnvoyPort             = 2000
		someRequestHandlerFactory = requestHandlerFactory{
			serviceHostname,
			someEnvoyPort,
		}
	)

	Context("when it receives an empty request", func() {
		recorderClient := newRecorderHTTPClient()
		requestHandler := someRequestHandlerFactory.withRecorderHTTPClient(recorderClient)
		requestBody := `{"request":[]}`
		resp := postRequest(requestHandler, requestBody)

		It("should return a success response", func() {
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})
		It("should forward no http requests", func() {
			Expect(recorderClient.recordedRequests).Should(HaveLen(0))
		})
	})

	Context("when it receives one single-element chain with correct destination", func() {
		recorderClient := newRecorderHTTPClient()
		requestHandler := someRequestHandlerFactory.withRecorderHTTPClient(recorderClient)
		requestBody := fmt.Sprintf(
			`{ "request": [
			[ { "destination": "%s" } ]
		] }`, serviceHostname)
		resp := postRequest(requestHandler, requestBody)

		It("should return a success response", func() {
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})
		It("should forward no http requests", func() {
			Expect(recorderClient.recordedRequests).Should(HaveLen(0))
		})
	})

	Context("when it receives one single-element chain with wrong destination", func() {
		recorderClient := newRecorderHTTPClient()
		requestHandler := someRequestHandlerFactory.withRecorderHTTPClient(recorderClient)
		requestBody := fmt.Sprintf(
			`{ "request": [
			[ { "destination": "%s" } ]
		] }`, someOtherHostname1)
		resp := postRequest(requestHandler, requestBody)

		It("should return a bad request response", func() {
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})
		It("should forward no http requests", func() {
			Expect(recorderClient.recordedRequests).Should(HaveLen(0))
		})
	})

	Context("when it receives a request with a single chain", func() {
		recorderClient := newRecorderHTTPClient()
		requestHandler := someRequestHandlerFactory.withRecorderHTTPClient(recorderClient)
		requestBody := fmt.Sprintf(
			`{ "request": [
			[ { "destination": "%s" },
			  { "destination": "%s" } ]
		] }`, serviceHostname, someOtherHostname1)
		resp := postRequest(requestHandler, requestBody)

		It("should return a success response", func() {
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})
		// CHECK TARGETS!
		It("should forward one http request", func() {
			Expect(recorderClient.recordedRequests).Should(HaveLen(1))
		})
		It("should forward a request to the correct host, with shortened request body", func() {
			expectedBody := fmt.Sprintf(
				`{ "request": [
					[ { "destination": "%s" }]
				] }`, someOtherHostname1)
			expectedHeader := someOtherHostname1
			req := recorderClient.recordedRequests[0]
			Expect(req.host).Should(Equal(expectedHeader))
			Expect(req.requestBody).Should(MatchJSON(expectedBody))
		})
	})

	Context("when it receives a request with two chains", func() {
		recorderClient := newRecorderHTTPClient()
		requestHandler := someRequestHandlerFactory.withRecorderHTTPClient(recorderClient)
		requestBody := fmt.Sprintf(
			`{ "request": [
			[ { "destination": "%s" },
			  { "destination": "%s" },
			  { "destination": "%s" } ],
			[ { "destination": "%s" },
			  { "destination": "%s" } ]
		] }`, serviceHostname, someOtherHostname1, someOtherHostname2, serviceHostname, someOtherHostname2)
		resp := postRequest(requestHandler, requestBody)

		It("should return a success response", func() {
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})

		It("should forward two http requests", func() {
			Expect(recorderClient.recordedRequests).Should(HaveLen(2))
		})
		It("should forward the first request chain to the first receiver", func() {
			expectedBody := fmt.Sprintf(
				`{ "request": [
					[ { "destination": "%s" },
					  { "destination": "%s" } ]
				] }`, someOtherHostname1, someOtherHostname2)
			expectedHeader := someOtherHostname1
			req := recorderClient.recordedRequests[0]
			Expect(req.host).Should(Equal(expectedHeader))
			Expect(req.requestBody).Should(MatchJSON(expectedBody))
		})
		It("should forward the second request chain to the second receiver", func() {
			expectedBody := fmt.Sprintf(
				`{ "request": [
					[ { "destination": "%s" }]
				] }`, someOtherHostname2)
			expectedHeader := someOtherHostname2
			req := recorderClient.recordedRequests[1]
			Expect(req.host).Should(Equal(expectedHeader))
			Expect(req.requestBody).Should(MatchJSON(expectedBody))
		})
	})

	Context("when it receives a request chain with invalid first destination", func() {
		recorderClient := newRecorderHTTPClient()
		requestHandler := someRequestHandlerFactory.withRecorderHTTPClient(recorderClient)
		requestBody := fmt.Sprintf(
			`{ "request": [
			[ { "destination": "%s" },
			  { "destination": "%s" } ]
		] }`, someOtherHostname1, someOtherHostname2)
		resp := postRequest(requestHandler, requestBody)

		It("should return a bad request response", func() {
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})

		It("should forward no http requests", func() {
			Expect(recorderClient.recordedRequests).Should(HaveLen(0))
		})
	})

	Context("when it receives two separate requests with one valid chain with multiple elements each ", func() {
		recorderClient := newRecorderHTTPClient()
		requestHandler := someRequestHandlerFactory.withRecorderHTTPClient(recorderClient)
		requestBody1 := fmt.Sprintf(
			`{ "request": [
			[ { "destination": "%s" },
			  { "destination": "%s" } ]
		] }`, serviceHostname, someOtherHostname1)
		resp1 := postRequest(requestHandler, requestBody1)

		requestBody2 := fmt.Sprintf(
			`{ "request": [
			[ { "destination": "%s" },
			  { "destination": "%s" } ]
		] }`, serviceHostname, someOtherHostname2)
		resp2 := postRequest(requestHandler, requestBody2)

		It("should return a success response for first request", func() {
			Expect(resp1.StatusCode).Should(Equal(http.StatusOK))
		})
		It("should return a success response for second request", func() {
			Expect(resp2.StatusCode).Should(Equal(http.StatusOK))
		})
		It("should forward a total of two http requests", func() {
			Expect(recorderClient.recordedRequests).Should(HaveLen(2))
		})
	})
})
