package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/rs/zerolog"
)

const (
	// Identifies the piece of code for purposes of otel's tracing.
	libraryName = "demo_app/app"
)

// HTTPClient interface for making http calls that can be mocked in tests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// SimpleService can be used to mimic read services in test environments.
// It will listen for request on defined port, and send requests to envoyPort.
type SimpleService struct {
	hostname string
	port     int
	// If the ENVOY_EGRESS_PORT environment variable is set then we use
	// envoy port to configure address to envoy proxy.
	// If it's not set then value is -1 and we do not configure proxy.
	// Istio proxy should handle requests without additional config
	// if it's injected.
	envoyPort   int
	concurrency int
	latency     time.Duration
}

// NewSimpleService creates a SimpleService instance.
func NewSimpleService(hostname string, port, envoyPort int, concurrency int, latency time.Duration) *SimpleService {
	return &SimpleService{
		hostname:    hostname,
		port:        port,
		envoyPort:   envoyPort,
		concurrency: concurrency,
		latency:     latency,
	}
}

// Run starts listening for requests on given port.
func (simpleService SimpleService) Run() error {
	var handler *RequestHandler
	if simpleService.envoyPort == -1 {
		handler = &RequestHandler{
			hostname:   simpleService.hostname,
			httpClient: &http.Client{},
		}
	} else {
		proxyURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", simpleService.envoyPort))
		if err != nil {
			log.Panic().Err(err).Msgf("Failed to parse url: %v", err)
		}
		handler = &RequestHandler{
			hostname:   simpleService.hostname,
			httpClient: &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}},
		}
	}

	http.Handle("/request", limitClients(handler, simpleService.concurrency, simpleService.latency))
	address := fmt.Sprintf(":%d", simpleService.port)

	server := &http.Server{Addr: address}

	return server.ListenAndServe()
}

func limitClients(h *RequestHandler, n int, l time.Duration) http.Handler {
	logger := log.Sample(&zerolog.BasicSampler{N: 100})

	sem := make(chan struct{}, n)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info().Msgf("Received request: %s", r.URL.Path)
		sem <- struct{}{}
		defer func() {
			<-sem
		}()
		time.Sleep(l)
		h.ServeHTTP(w, r)
	})
}

// Request contains a collection of destination chains.
// It is the format of incoming requests.
type Request struct {
	Chains []SubrequestChain `json:"request"`
}

// SubrequestChain is a single chain of destinations.
type SubrequestChain struct {
	subrequests []Subrequest
}

// MarshalJSON writes the subrequest chain as a JSON list and not a JSON object.
func (sc SubrequestChain) MarshalJSON() ([]byte, error) {
	return json.Marshal(sc.subrequests)
}

// UnmarshalJSON creates the subrequest chain, reading a JSON list and not a JSON object.
func (sc *SubrequestChain) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &sc.subrequests)
	if err != nil {
		return err
	}
	return nil
}

// Subrequest contains information on a single destination
// to which it should be routed.
type Subrequest struct {
	Destination string `json:"destination"`
}

// Response is returned by the service after successfully processing a Request.
type Response struct{}

// RequestHandler handles processing of incoming requests.
type RequestHandler struct {
	httpClient HTTPClient
	hostname   string
}

func (h RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var p Request

	// Extract baggage and trace context from headers
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	// Start a new span (not needed if we want "just passthrough")
	ctx, span := otel.Tracer(libraryName).Start(ctx, "ServeHTTP")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subresponses := make([]Response, len(p.Chains))
	for i, chain := range p.Chains {
		if len(chain.subrequests) == 0 {
			http.Error(w, "Received empty subrequest chain", http.StatusBadRequest)
			return
		}
		requestDestination := chain.subrequests[0].Destination
		if requestDestination != h.hostname {
			http.Error(w, "Invalid message destination", http.StatusBadRequest)
			return
		}
		subresponses[i], err = h.processChain(ctx, chain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	response := generateResponse(subresponses)
	jsonString, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(w, string(jsonString))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h RequestHandler) processChain(ctx context.Context, chain SubrequestChain) (Response, error) {
	log.Info().Msg("processChain()")
	if len(chain.subrequests) == 1 {
		return h.processRequest(chain.subrequests[0])
	}

	requestForwardingDestination := chain.subrequests[1].Destination
	trimmedSubrequestChain := SubrequestChain{
		subrequests: chain.subrequests[1:],
	}
	trimmedRequest := Request{
		Chains: []SubrequestChain{trimmedSubrequestChain},
	}
	return h.forwardRequest(ctx, requestForwardingDestination, trimmedRequest)
}

func (h RequestHandler) processRequest(s Subrequest) (Response, error) {
	return Response{}, nil
}

func (h RequestHandler) forwardRequest(ctx context.Context, destinationHostname string, requestBody Request) (Response, error) {
	address := fmt.Sprintf("http://%s/request", destinationHostname)

	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		return Response{}, err
	}

	request, err := http.NewRequest("POST", address, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return Response{}, err
	}

	request = request.WithContext(ctx)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	request.Header.Set("Content-Type", "application/json")

	response, err := h.httpClient.Do(request)
	if err != nil {
		return Response{}, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}

	err = response.Body.Close()
	if err != nil {
		return Response{}, err
	}

	var rsp Response

	err = json.Unmarshal(responseBody, &rsp)
	if err != nil {
		return Response{}, err
	}

	return rsp, nil
}

func generateResponse(subresponses []Response) Response {
	return Response{}
}
