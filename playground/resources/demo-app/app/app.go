package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"

	"github.com/fluxninja/aperture/pkg/log"
)

const (
	// Identifies the piece of code for purposes of otel's tracing.
	libraryName = "demo-app/app"
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
	rejectRatio float64
}

// NewSimpleService creates a SimpleService instance.
func NewSimpleService(hostname string, port, envoyPort int,
	concurrency int,
	latency time.Duration,
	rejectRatio float64,
) *SimpleService {
	return &SimpleService{
		hostname:    hostname,
		port:        port,
		envoyPort:   envoyPort,
		concurrency: concurrency,
		latency:     latency,
		rejectRatio: rejectRatio,
	}
}

// Run starts listening for requests on given port.
func (simpleService SimpleService) Run() error {
	handler := &RequestHandler{
		hostname:     simpleService.hostname,
		latency:      simpleService.latency,
		rejectRatio:  simpleService.rejectRatio,
		concurrency:  simpleService.concurrency,
		limitClients: make(chan struct{}, simpleService.concurrency),
	}
	if simpleService.envoyPort == -1 {
		handler.httpClient = &http.Client{}
	} else {
		proxyURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", simpleService.envoyPort))
		if err != nil {
			log.Panic().Err(err).Msgf("Failed to parse url: %v", err)
		}
		handler.httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	}

	http.Handle("/request", handlerFunc(handler))
	address := fmt.Sprintf(":%d", simpleService.port)

	server := &http.Server{Addr: address}

	return server.ListenAndServe()
}

func handlerFunc(h *RequestHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Autosample().Info().Msg("received request")
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

// RequestHandler handles processing of incoming requests.
type RequestHandler struct {
	httpClient   HTTPClient
	limitClients chan struct{}
	hostname     string
	concurrency  int
	latency      time.Duration
	rejectRatio  float64
}

func (h RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var p Request

	// Extract baggage and trace context from headers
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	// Start a new span (not needed if we want "just passthrough")
	ctx, span := otel.Tracer(libraryName).Start(ctx, "ServeHTTP")
	defer span.End()

	// randomly reject requests based on rejectRatio
	// nolint:gosec
	if h.rejectRatio > 0 && rand.Float64() < h.rejectRatio {
		span.SetStatus(codes.Error, "rejected")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Autosample().Warn().Err(err).Msg("Failed to read request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Autosample().Warn().Err(err).Msg("Failed to unmarshal request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code := http.StatusOK
	for _, chain := range p.Chains {
		if len(chain.subrequests) == 0 {
			log.Autosample().Warn().Err(err).Msg("Empty chain")
			http.Error(w, "Received empty subrequest chain", http.StatusBadRequest)
			return
		}
		requestDestination := chain.subrequests[0].Destination
		if !strings.HasPrefix(requestDestination, "http") {
			requestDestination = "http://" + requestDestination
		}
		requestDomain, err := url.Parse(requestDestination)
		if requestDomain.Hostname() != h.hostname {
			log.Autosample().Warn().Err(err).Msg("Invalid destination")
			http.Error(w, "Invalid message destination", http.StatusBadRequest)
			return
		}
		code, err = h.processChain(ctx, chain)
		if err != nil {
			log.Autosample().Warn().Err(err).Msg("Failed to process chain")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if code == http.StatusOK {
		span.SetStatus(codes.Ok, "OK")
	} else {
		span.SetStatus(codes.Error, "Error")
	}
}

func (h RequestHandler) processChain(ctx context.Context, chain SubrequestChain) (int, error) {
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

func (h RequestHandler) processRequest(s Subrequest) (int, error) {
	if h.concurrency > 0 {
		h.limitClients <- struct{}{}
		defer func() {
			<-h.limitClients
		}()
	}
	if h.latency > 0 {
		// Fake workload
		time.Sleep(h.latency)
	}
	return http.StatusOK, nil
}

func (h RequestHandler) forwardRequest(ctx context.Context, destinationHostname string, requestBody Request) (int, error) {
	address := fmt.Sprintf("http://%s", destinationHostname)

	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		log.Bug().Err(err).Msg("bug: Failed to marshal request")
		return http.StatusInternalServerError, err
	}

	request, err := http.NewRequest("POST", address, bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Autosample().Error().Err(err).Msg("Failed to create request")
		return http.StatusInternalServerError, err
	}

	request = request.WithContext(ctx)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	request.Header.Set("Content-Type", "application/json")

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Autosample().Error().Err(err).Msg("Failed to send request")
		return http.StatusInternalServerError, err
	}
	return response.StatusCode, nil
}
