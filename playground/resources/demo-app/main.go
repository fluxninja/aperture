package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/fluxninja/aperture/playground/resources/demo-app/app"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

type (
	rabbitMQEnvVar      string
	elasticsearchEnvVar string
	pgsqlEnvVar         string
)

const (
	rabbitMQEnabledEnvVar  rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_ENABLED"
	rabbitMQHostEnvVar     rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_HOST"
	rabbitMQPortEnvVar     rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_PORT"
	rabbitMQUsernameEnvVar rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_USERNAME"
	rabbitMQPasswordEnvVar rabbitMQEnvVar = "SIMPLE_SERVICE_RABBITMQ_PASSWORD"

	elasticsearchEnabledEnvVar  elasticsearchEnvVar = "SIMPLE_SERVICE_ELASTICSEARCH_ENABLED"
	elasticsearchHostEnvVar     elasticsearchEnvVar = "SIMPLE_SERVICE_ELASTICSEARCH_HOST"
	elasticsearchPortEnvVar     elasticsearchEnvVar = "SIMPLE_SERVICE_ELASTICSEARCH_PORT"
	elasticsearchUsernameEnvVar elasticsearchEnvVar = "SIMPLE_SERVICE_ELASTICSEARCH_USERNAME"
	elasticsearchPasswordEnvVar elasticsearchEnvVar = "SIMPLE_SERVICE_ELASTICSEARCH_PASSWORD"

	pgsqlEnabledEnvVar  pgsqlEnvVar = "SIMPLE_SERVICE_PGSQL_ENABLED"
	pgsqlHostEnvVar     pgsqlEnvVar = "SIMPLE_SERVICE_PGSQL_HOST"
	pgsqlPortEnvVar     pgsqlEnvVar = "SIMPLE_SERVICE_PGSQL_PORT"
	pgsqlUserEnvVar     pgsqlEnvVar = "SIMPLE_SERVICE_PGSQL_USER"
	pgsqlPasswordEnvVar pgsqlEnvVar = "SIMPLE_SERVICE_PGSQL_PASSWORD"
)

func main() {
	hostname := hostnameFromEnv()
	port := portFromEnv()
	envoyPort := envoyPortFromEnv()
	concurrency := concurrencyFromEnv()
	latency := latencyFromEnv()
	rejectRatio := rejectRatioFromEnv()
	cpuLoadPercentage := cpuLoadPercentageFromEnv()

	// RabbitMQ related setup
	rabbitMQURL := ""
	if rabbitMQFromEnv(rabbitMQEnabledEnvVar) == "true" {
		rabbitMQHost := rabbitMQFromEnv(rabbitMQHostEnvVar)
		rabbitMQPort := rabbitMQFromEnv(rabbitMQPortEnvVar)
		rabbitMQUsername := rabbitMQFromEnv(rabbitMQUsernameEnvVar)
		rabbitMQPassword := rabbitMQFromEnv(rabbitMQPasswordEnvVar)
		rabbitMQURL = "amqp://" + rabbitMQUsername + ":" + rabbitMQPassword + "@" + rabbitMQHost + ":" + rabbitMQPort + "/"
	}

	// Elasticsearch related setup
	elaticsearchConfig := elasticsearch.Config{}
	if elasticsearchFromEnv(elasticsearchEnabledEnvVar) == "true" {
		elasticsearchHost := elasticsearchFromEnv(elasticsearchHostEnvVar)
		elasticsearchPort := elasticsearchFromEnv(elasticsearchPortEnvVar)
		elasticsearchUserName := elasticsearchFromEnv(elasticsearchUsernameEnvVar)
		elasticsearchPassword := elasticsearchFromEnv(elasticsearchPasswordEnvVar)

		elaticsearchConfig = elasticsearch.Config{
			Addresses: []string{
				"http://" + elasticsearchHost + ":" + elasticsearchPort,
			},
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 1,
				MaxIdleConns:        1,
				MaxConnsPerHost:     1000,
				DisableKeepAlives:   true,
			},
			DiscoverNodesOnStart:  true,
			DiscoverNodesInterval: 60 * time.Second,
			Username:              elasticsearchUserName,
			Password:              elasticsearchPassword,
		}
	}

	// PostgreSQL related setup
	pgsqlURL := ""
	if pgsqlFromEnv(pgsqlEnabledEnvVar) == "true" {
		pgsqlHost := pgsqlFromEnv(pgsqlHostEnvVar)
		pgsqlPort := pgsqlFromEnv(pgsqlPortEnvVar)
		pgsqlUser := pgsqlFromEnv(pgsqlUserEnvVar)
		pgsqlPassword := pgsqlFromEnv(pgsqlPasswordEnvVar)
		pgsqlURL = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", pgsqlUser, pgsqlPassword, pgsqlHost, pgsqlPort, pgsqlUser)
	}

	// We do not necessarily need tracing providers (just propagators), but lets
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

	service := app.NewSimpleService(hostname, port, envoyPort, rabbitMQURL, elaticsearchConfig, pgsqlURL, concurrency, latency, rejectRatio, cpuLoadPercentage)
	err := service.Run()
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func rabbitMQFromEnv(envVar rabbitMQEnvVar) string {
	value := os.Getenv(string(envVar))
	if value == "" {
		switch envVar {
		case rabbitMQEnabledEnvVar:
			return "false"
		case rabbitMQHostEnvVar:
			return "localhost"
		case rabbitMQPortEnvVar:
			return "5672"
		case rabbitMQUsernameEnvVar:
			return "user"
		case rabbitMQPasswordEnvVar:
			return ""
		default:
			return ""
		}
	}
	return value
}

func elasticsearchFromEnv(envVar elasticsearchEnvVar) string {
	value := os.Getenv(string(envVar))
	if value == "" {
		switch envVar {
		case elasticsearchEnabledEnvVar:
			return "false"
		case elasticsearchHostEnvVar:
			return "localhost"
		case elasticsearchPortEnvVar:
			return "9200"
		case elasticsearchUsernameEnvVar:
			return "elastic"
		case elasticsearchPasswordEnvVar:
			return ""
		default:
			return ""
		}
	}
	return value
}

func pgsqlFromEnv(envVar pgsqlEnvVar) string {
	value := os.Getenv(string(envVar))
	if value == "" {
		switch envVar {
		case pgsqlEnabledEnvVar:
			return "false"
		case pgsqlHostEnvVar:
			return "localhost"
		case pgsqlPortEnvVar:
			return "5432"
		case pgsqlUserEnvVar:
			return "postgres"
		case pgsqlPasswordEnvVar:
			return ""
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
	}
	envoyPort, err := strconv.Atoi(portValue)
	if err != nil {
		log.Panic().Err(err).Msg("Failed converting ENVOY_EGRESS_PORT value")
	}
	return envoyPort
}

func portFromEnv() int {
	port, err := strconv.Atoi(os.Getenv("SIMPLE_SERVICE_PORT"))
	if err != nil {
		log.Panic().Err(err).Msg("Failed converting SIMPLE_SERVICE_PORT")
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
		log.Panic().Err(err).Msg("Failed converting SIMPLE_SERVICE_CONCURRENCY")
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
		log.Panic().Err(err).Msg("Failed converting SIMPLE_SERVICE_LATENCY")
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
		log.Panic().Err(err).Msg("Failed converting SIMPLE_SERVICE_REJECT_RATIO")
	}

	return rejectRatio
}

func cpuLoadPercentageFromEnv() int {
	loadCPUValue, exists := os.LookupEnv("SIMPLE_SERVICE_CPU_LOAD")
	if !exists {
		return 0
	}

	loadCPUI64, err := strconv.ParseInt(loadCPUValue, 10, 32)
	if err != nil {
		log.Panic().Err(err).Msg("Failed converting SIMPLE_SERVICE_CPU_LOAD")
	}

	loadCPU := int(loadCPUI64)

	if loadCPU < 0 || loadCPU > 100 {
		log.Panic().Msg("SIMPLE_SERVICE_CPU_LOAD must be between 0 and 100")
	}

	return loadCPU
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
