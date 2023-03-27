package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/playground/resources/demo-app/app"
)

type rabbitMQEnvVar string

const (
	rabbitMQEnable         rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_ENABLE"
	rabbitMQHostEnvVar     rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_HOST"
	rabbitMQPortEnvVar     rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_PORT"
	rabbitMQUserEnvVar     rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_USER"
	rabbitMQPasswordEnvVar rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_PASSWORD"
)

func main() {
	hostname := hostnameFromEnv()
	port := portFromEnv()
	envoyPort := envoyPortFromEnv()
	concurrency := concurrencyFromEnv()
	latency := latencyFromEnv()
	rejectRatio := rejectRatioFromEnv()

	// RabbitMQ related setup
	rabbitMQURL := ""
	if rabbitMQFromEnv(rabbitMQEnable) == "true" {
		rabbitMQHost := rabbitMQFromEnv(rabbitMQHostEnvVar)
		rabbitMQPort := rabbitMQFromEnv(rabbitMQPortEnvVar)
		rabbitMQUser := rabbitMQFromEnv(rabbitMQUserEnvVar)
		rabbitMQPassword := rabbitMQFromEnv(rabbitMQPasswordEnvVar)
		rabbitMQURL = "amqp://" + rabbitMQUser + ":" + rabbitMQPassword + "@" + rabbitMQHost + ":" + rabbitMQPort + "/"
	}

	// We don't necessarily need tracing providers (just propagators), but lets
	// do them anyway to have a "more realistic" otel usage
	// exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	// if err != nil {
	// 	log.Panic().Err(err).Msgf("Failed to set up exporter: %v", err)
	// }
	tp := trace.NewTracerProvider(
		// trace.WithBatcher(exporter),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Failed to shutdown tracer")
		}
	}()
	otel.SetTracerProvider(tp)
	// Setup Propagators
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	service := app.NewSimpleService(hostname, port, envoyPort, rabbitMQURL, concurrency, latency, rejectRatio)
	err := service.Run()
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func rabbitMQFromEnv(envVar rabbitMQEnvVar) string {
	value := os.Getenv(string(envVar))
	if value == "" {
		switch envVar {
		case rabbitMQEnable:
			return "false"
		case rabbitMQHostEnvVar:
			return "localhost"
		case rabbitMQPortEnvVar:
			return "5672"
		case rabbitMQUserEnvVar:
			return "guest"
		case rabbitMQPasswordEnvVar:
			return "guest"
		default:
			return ""
		}
	}
	return value
}

func envoyPortFromEnv() int {
	portValue, exists := os.LookupEnv("ENVOY_EGRESS_PORT")
	if !exists {
		// We do not use manually configured envoy proxy
		return -1
	} else {
		envoyPort, err := strconv.Atoi(portValue)
		if err != nil {
			log.Panic().Err(err).Msgf("Failed converting ENVOY_EGRESS_PORT value: %v", err)
		}
		return envoyPort
	}
}

func portFromEnv() int {
	port, err := strconv.Atoi(os.Getenv("SIMPLE_SERVICE_PORT"))
	if err != nil {
		log.Panic().Err(err).Msgf("Failed converting SIMPLE_SERVICE_PORT: %v", err)
	}
	return port
}

func hostnameFromEnv() string {
	return os.Getenv("HOSTNAME")
}

func concurrencyFromEnv() int {
	concurrencyValue, exists := os.LookupEnv("SIMPLE_SERVICE_CONCURRENCY")
	if !exists {
		return 10
	}
	concurrency, err := strconv.Atoi(concurrencyValue)
	if err != nil {
		log.Panic().Err(err).Msgf("Failed converting SIMPLE_SERVICE_CONCURRENCY: %v", err)
	}
	return concurrency
}

func latencyFromEnv() time.Duration {
	latencyValue, exists := os.LookupEnv("SIMPLE_SERVICE_LATENCY")
	if !exists {
		return time.Millisecond * 50
	}
	latency, err := time.ParseDuration(latencyValue)
	if err != nil {
		log.Panic().Err(err).Msgf("Failed converting SIMPLE_SERVICE_LATENCY: %v", err)
	}
	return latency
}

func rejectRatioFromEnv() float64 {
	rejectRatioValue, exists := os.LookupEnv("SIMPLE_SERVICE_REJECT_RATIO")
	if !exists {
		return 0.05
	}
	rejectRatio, err := strconv.ParseFloat(rejectRatioValue, 64)
	if err != nil {
		log.Panic().Err(err).Msgf("Failed converting SIMPLE_SERVICE_REJECT_RATIO: %v", err)
	}
	return rejectRatio
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("demoapp"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
