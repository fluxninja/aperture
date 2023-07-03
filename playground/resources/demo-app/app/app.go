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
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/castai/promwrite"
	backoff "github.com/cenkalti/backoff/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

const (
	// Identifies the piece of code for purposes of otel's tracing.
	libraryName = "demo-app/app"
)

var (
	counter          = map[string]*Counter{}
	prometheusClient = &promwrite.Client{}
	mutex            = sync.RWMutex{}
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
	// If it is not set then value is -1 and we do not configure proxy.
	// Istio proxy should handle requests without additional config
	// if it is injected.
	envoyPort         int
	rabbitMQURL       string
	concurrency       int           // Maximum number of concurrent clients
	latency           time.Duration // Simulated latency for each request
	rejectRatio       float64       // Ratio of requests to be rejected
	cpuLoadPercentage int           // Percentage of CPU to be loaded
}

// ResponseBody is a response body for returning a response to all requests made on api endpoints
type ResponseBody struct {
	Message string `json:"message"`
}

// NewSimpleService creates a SimpleService instance.
func NewSimpleService(
	hostname string,
	port int,
	envoyPort int,
	rabbitMQURL string,
	concurrency int,
	latency time.Duration,
	rejectRatio float64,
	cpuLoadPercentage int,
) *SimpleService {
	return &SimpleService{
		hostname:          hostname,
		port:              port,
		envoyPort:         envoyPort,
		rabbitMQURL:       rabbitMQURL,
		concurrency:       concurrency,
		latency:           latency,
		rejectRatio:       rejectRatio,
		cpuLoadPercentage: cpuLoadPercentage,
	}
}

// Run starts listening for requests on given port.
func (ss SimpleService) Run() error {
	handler := &RequestHandler{
		hostname:          ss.hostname,
		latency:           ss.latency,
		rejectRatio:       ss.rejectRatio,
		concurrency:       ss.concurrency,
		limitClients:      make(chan struct{}, ss.concurrency),
		cpuLoadPercentage: ss.cpuLoadPercentage,
	}

	if ss.rabbitMQURL != "" {
		var conn *amqp.Connection
		operation := func() error {
			var err error
			conn, err = amqp.Dial(ss.rabbitMQURL)
			if err != nil {
				return err
			}
			return nil
		}
		err := backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewConstantBackOff(2*time.Second), 3))
		if err != nil {
			return err
		}
		defer conn.Close()

		pCh, err := conn.Channel()
		if err != nil {
			return err
		}
		defer pCh.Close()
		handler.rabbitMQChan = pCh

		cCh, err := conn.Channel()
		if err != nil {
			return err
		}
		defer cCh.Close()

		err = cCh.Qos(
			handler.concurrency*5, // prefetch count
			0,                     // prefetch size
			false,                 // global
		)
		if err != nil {
			return err
		}

		// Declare a queue for the service to consume from.
		queueArgs := amqp.Table{
			"x-max-length": int32(2000),
			"x-overflow":   string("reject-publish"),
		}
		q, err := cCh.QueueDeclare(
			handler.hostname, // name
			false,            // durable
			false,            // delete when unused
			false,            // exclusive
			false,            // no-wait
			queueArgs,        // arguments
		)
		if err != nil {
			return err
		}

		go consumeFromQueue(q, cCh, handler)
	}

	if ss.envoyPort == -1 {
		handler.httpClient = &http.Client{}
	} else {
		proxyURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", ss.envoyPort))
		if err != nil {
			log.Panic().Err(err).Msgf("Failed to parse url: %v", err)
		}
		handler.httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	}

	prometheusAddress := os.Getenv("PROMETHEUS_ADDRESS")
	if prometheusAddress != "" {

		if prometheusAddress[len(prometheusAddress)-1] == '/' {
			prometheusAddress = prometheusAddress[:len(prometheusAddress)-1]
		}

		prometheusClient = promwrite.NewClient(fmt.Sprintf("%s/api/v1/write", prometheusAddress))
		http.HandleFunc("/prometheus", prometheusHandler)
	}

	http.HandleFunc("/api/rate-limit", apiEndpointHandler)
	http.HandleFunc("/api/load-ramp", apiEndpointHandler)
	http.Handle("/request", handlerFunc(handler))

	address := fmt.Sprintf(":%d", ss.port)

	server := &http.Server{Addr: address}

	return server.ListenAndServe()
}

func apiEndpointHandler(w http.ResponseWriter, r *http.Request) {
	// Extract baggage and trace context from headers
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	// Start a new span (not needed if we want "just passthrough")
	ctx, span := otel.Tracer(libraryName).Start(ctx, "apiEndpointHandler")
	defer span.End()

	responseBody := ResponseBody{
		Message: "Request accepted",
	}

	err := json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		log.Autosample().Error().Err(err).Msg("Error encoding JSON")
		span.SetStatus(codes.Error, "rejected")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
	}

	r = r.WithContext(ctx)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

	span.SetStatus(codes.Ok, "accepted")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getOrCreateCounter(userID, userType string) *Counter {
	mutex.Lock()
	defer mutex.Unlock()

	key := fmt.Sprintf("%s-%s", userID, userType)
	c := counter[key]
	if c == nil {
		c = &Counter{}
		counter[key] = c
	}
	return c
}

func prometheusHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// Increment the counter
	userID := r.Header.Get("User-Id")
	userType := r.Header.Get("User-Type")

	// Increment the counter
	c := getOrCreateCounter(userID, userType)
	c.Increment()

	prometheusClient.Write(context.Background(), &promwrite.WriteRequest{
		TimeSeries: []promwrite.TimeSeries{
			{
				Labels: []promwrite.Label{
					{
						Name:  "__name__",
						Value: "demo_app_requests_total",
					},
					{
						Name:  "request_uri",
						Value: r.RequestURI,
					},
					{
						Name:  "request_user_id",
						Value: userID,
					},
					{
						Name:  "request_user_type",
						Value: userType,
					},
				},
				Sample: promwrite.Sample{
					Time:  time.Now(),
					Value: float64(c.Value()),
				},
			},
		},
	})
	end := time.Now()
	prometheusClient.Write(context.Background(), &promwrite.WriteRequest{
		TimeSeries: []promwrite.TimeSeries{
			{
				Labels: []promwrite.Label{
					{
						Name:  "__name__",
						Value: "demo_app_request_duration_micro_seconds",
					},
					{
						Name:  "request_uri",
						Value: r.RequestURI,
					},
					{
						Name:  "request_user_id",
						Value: userID,
					},
					{
						Name:  "request_user_type",
						Value: userType,
					},
				},
				Sample: promwrite.Sample{
					Time:  time.Now(),
					Value: float64(end.Sub(start).Milliseconds()),
				},
			},
		},
	})
	// Send a response
	w.Write([]byte("OK"))
}

func consumeFromQueue(q amqp.Queue, cCh *amqp.Channel, handler *RequestHandler) {
	// Consume messages from the queue.
	delivery, err := cCh.Consume(
		q.Name,           // queue
		handler.hostname, // consumer
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to consume messages")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for d := range delivery {
		var p Request
		err := json.Unmarshal(d.Body, &p)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal message")
			continue
		}
		code := http.StatusOK
		for _, chain := range p.Chains {
			if len(chain.subrequests) == 0 {
				log.Error().Err(err).Msg("Empty chain")
				break
			}
			requestDestination := chain.subrequests[0].Destination
			if !strings.HasPrefix(requestDestination, "http") {
				requestDestination = "http://" + requestDestination
			}
			requestDomain, parseErr := url.Parse(requestDestination)
			if requestDomain.Hostname() != handler.hostname {
				log.Error().Err(parseErr).Msg("Invalid destination")
				break
			}
			code, err = handler.processChain(ctx, chain, nil)
			if err != nil {
				log.Error().Err(err).Msg("Failed to process chain")
				break
			}
		}
		if code == http.StatusOK {
			err = d.Ack(false)
			if err != nil {
				log.Error().Err(err).Msg("Failed to ack message")
			}
		} else {
			err = d.Nack(false, false)
			if err != nil {
				log.Error().Err(err).Msg("Failed to nack message")
			}
		}
	}
}

func handlerFunc(h *RequestHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Autosample().Trace().Msg("received request")
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
	httpClient        HTTPClient
	rabbitMQChan      *amqp.Channel
	limitClients      chan struct{}
	hostname          string
	concurrency       int
	latency           time.Duration
	rejectRatio       float64
	cpuLoadPercentage int
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

	headers := r.Header.Clone()

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
		code, err = h.processChain(ctx, chain, headers)
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

func (h RequestHandler) processChain(ctx context.Context, chain SubrequestChain, headers http.Header) (int, error) {
	if len(chain.subrequests) == 1 {
		return h.processRequest()
	}
	forwardingDestination := chain.subrequests[1].Destination
	trimmedSubrequestChain := SubrequestChain{
		subrequests: chain.subrequests[1:],
	}
	trimmedRequest := Request{
		Chains: []SubrequestChain{trimmedSubrequestChain},
	}
	return h.forwardRequest(ctx, forwardingDestination, trimmedRequest, headers)
}

func busyWait(duration time.Duration) {
	startTime := time.Now()
	for time.Since(startTime) < duration {
		// Just busy wait
	}
}

var ongoingRequests int32

func (h RequestHandler) processRequest() (int, error) {
	if h.concurrency > 0 {
		h.limitClients <- struct{}{}
		defer func() {
			<-h.limitClients
		}()
	}

	atomic.AddInt32(&ongoingRequests, 1)
	defer atomic.AddInt32(&ongoingRequests, -1)

	if h.latency > 0 {
		if h.cpuLoadPercentage > 0 {
			// Simulate CPU load by busy waiting
			numCores := runtime.NumCPU()

			// Calculate busy wait and sleep durations based on h.loadCPU and ongoing requests
			adjustedLoad := (float64(h.cpuLoadPercentage) / 100) * float64(atomic.LoadInt32(&ongoingRequests))
			if adjustedLoad > 100 {
				adjustedLoad = 100
			}
			totalDuration := h.latency.Seconds()
			busyWaitDuration := time.Duration(totalDuration * adjustedLoad / 100.0 * float64(time.Second))
			sleepDuration := time.Duration(totalDuration * (100.0 - adjustedLoad) / 100.0 * float64(time.Second))

			var wg sync.WaitGroup
			wg.Add(numCores)
			for i := 0; i < numCores; i++ {
				go func() {
					defer wg.Done()
					startTime := time.Now()
					for time.Since(startTime) < h.latency {
						busyWait(busyWaitDuration)
						time.Sleep(sleepDuration)
					}
				}()
			}
			wg.Wait()
		} else {
			time.Sleep(h.latency)
		}
	}

	return http.StatusOK, nil
}

func (h RequestHandler) forwardRequest(ctx context.Context, destination string, requestBody Request, headers http.Header) (int, error) {
	jsonRequest, err := json.Marshal(requestBody)
	if err != nil {
		log.Bug().Err(err).Msg("bug: Failed to marshal request")
		return http.StatusInternalServerError, err
	}

	if h.rabbitMQChan != nil {
		destinationHostname := strings.TrimSuffix(destination, "/request")
		err = h.rabbitMQChan.PublishWithContext(
			ctx,                 // context
			"",                  // exchange
			destinationHostname, // routing key
			false,               // mandatory
			false,               // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonRequest,
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send request")
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}

	address := fmt.Sprintf("http://%s", destination)
	request, err := http.NewRequest("POST", address, bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Autosample().Error().Err(err).Msg("Failed to create request")
		return http.StatusInternalServerError, err
	}

	request = request.WithContext(ctx)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	if headers != nil {
		request.Header = headers
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := h.httpClient.Do(request)
	if err != nil {
		log.Autosample().Error().Err(err).Msg("Failed to send request")
		return http.StatusInternalServerError, err
	}
	return response.StatusCode, nil
}
