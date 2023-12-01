module github.com/fluxninja/aperture/v2

go 1.21.4

require (
	cloud.google.com/go/secretmanager v1.11.4
	github.com/BurntSushi/toml v1.3.2
	github.com/Henry-Sarabia/sliceconv v1.0.2
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/buger/jsonparser v1.1.1
	github.com/buraksezer/olric v0.0.0-00010101000000-000000000000
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/charmbracelet/bubbletea v0.24.2
	github.com/clarketm/json v1.17.1
	github.com/containerd/cgroups v1.1.0
	github.com/eapache/queue v1.1.0
	github.com/elastic/gmux v0.3.1
	github.com/elastic/gosigar v0.14.2
	github.com/emicklei/dot v1.6.0
	github.com/envoyproxy/go-control-plane v0.11.1
	github.com/fluxninja/aperture/api/v2 v2.0.0
	github.com/fluxninja/datasketches-go v0.0.0-20220916235224-7501a2d28551
	github.com/fluxninja/lumberjack v0.0.0-20220729045908-655029e4d814
	github.com/getsentry/sentry-go v0.24.1
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/zerologr v1.2.3
	github.com/go-openapi/runtime v0.26.0
	github.com/go-openapi/strfmt v0.21.7
	github.com/go-playground/validator/v10 v10.14.1
	github.com/gofrs/flock v0.8.1
	github.com/golang/mock v1.6.0
	github.com/google/go-jsonnet v0.20.0
	github.com/google/uuid v1.4.0
	github.com/gorilla/mux v1.8.1
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/memberlist v0.5.0
	github.com/integralist/go-findroot v0.0.0-20160518114804-ac90681525dc
	github.com/jonboulle/clockwork v0.3.0
	github.com/json-iterator/go v1.1.12
	github.com/jsonnet-bundler/jsonnet-bundler v0.5.1
	github.com/knadh/koanf/parsers/json v0.1.0
	github.com/knadh/koanf/parsers/yaml v0.1.0
	github.com/knadh/koanf/providers/confmap v0.1.0
	github.com/knadh/koanf/providers/posflag v0.1.0
	github.com/knadh/koanf/providers/rawbytes v0.1.0
	github.com/knadh/koanf/v2 v2.0.1
	github.com/lithammer/dedent v1.1.0
	github.com/looplab/tarjan v0.1.0
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/mapstructure v1.5.1-0.20220423185008-bf980b35cac4
	github.com/natefinch/atomic v1.0.1
	github.com/onsi/ginkgo/v2 v2.13.0
	github.com/onsi/gomega v1.29.0
	github.com/open-policy-agent/opa v0.58.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/servicegraphprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureeventhubreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/bigipreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/expvarreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/haproxyreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/httpcheckreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/iisreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sobjectsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nsxtreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/oracledbreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/pulsarreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/saphanareceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/signalfxreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmpreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/wavefrontreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.90.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.90.0
	github.com/prometheus/alertmanager v0.26.0
	github.com/prometheus/client_golang v1.17.0
	github.com/prometheus/common v0.45.0
	github.com/prometheus/prometheus v0.48.0
	github.com/reugn/go-quartz v0.7.0
	github.com/rs/zerolog v1.31.0
	github.com/sourcegraph/conc v0.3.0
	github.com/spf13/cast v1.5.1
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	github.com/technosophos/moniker v0.0.0-20210218184952-3ea787d3943b
	github.com/xeipuuv/gojsonschema v1.2.0
	go.etcd.io/etcd/api/v3 v3.5.10
	go.etcd.io/etcd/client/v3 v3.5.10
	go.opentelemetry.io/collector/component v0.90.0
	go.opentelemetry.io/collector/confmap v0.90.0
	go.opentelemetry.io/collector/connector v0.90.0
	go.opentelemetry.io/collector/consumer v0.90.0
	go.opentelemetry.io/collector/exporter v0.90.0
	go.opentelemetry.io/collector/exporter/loggingexporter v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/collector/exporter/otlpexporter v0.90.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/collector/extension v0.90.0
	go.opentelemetry.io/collector/extension/ballastextension v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/collector/extension/zpagesextension v0.90.0
	go.opentelemetry.io/collector/otelcol v0.90.0
	go.opentelemetry.io/collector/pdata v1.0.0
	go.opentelemetry.io/collector/processor v0.90.0
	go.opentelemetry.io/collector/processor/batchprocessor v0.90.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.90.0
	go.opentelemetry.io/collector/receiver v0.90.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.90.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.1
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/proto/otlp v1.0.0
	go.uber.org/automaxprocs v1.5.3
	go.uber.org/fx v1.20.1
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.26.0
	golang.org/x/exp v0.0.0-20231127185646-65229373498e
	golang.org/x/net v0.19.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231127180814-3a041ad873d4
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/yaml.v3 v3.0.1
	helm.sh/helm/v3 v3.13.1
	k8s.io/api v0.28.4
	k8s.io/apimachinery v0.28.4
	k8s.io/client-go v0.28.4
	k8s.io/klog/v2 v2.110.1
	sigs.k8s.io/controller-runtime v0.16.3
)

require (
	cloud.google.com/go v0.110.10 // indirect
	cloud.google.com/go/compute/metadata v0.2.4-0.20230617002413-005d2dfb6b68 // indirect
	cloud.google.com/go/iam v1.1.5 // indirect
	cloud.google.com/go/pubsub v1.33.0 // indirect
	cloud.google.com/go/spanner v1.53.0 // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20211115184647-b584dd5df32c // indirect
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible // indirect
	code.cloudfoundry.org/rfc5424 v0.0.0-20201103192249-000122071b78 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/AthenZ/athenz v1.10.39 // indirect
	github.com/Azure/azure-amqp-common-go/v4 v4.2.0 // indirect
	github.com/Azure/azure-event-hubs-go/v3 v3.6.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.8.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2 v2.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.2.0 // indirect
	github.com/Azure/go-amqp v1.0.2 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.1.1 // indirect
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.21.0 // indirect
	github.com/IBM/sarama v1.42.1 // indirect
	github.com/JohnCGriffin/overflow v0.0.0-20211019200055-46fa312c352c // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/ReneKroon/ttlcache/v2 v2.11.0 // indirect
	github.com/RoaringBitmap/roaring v1.2.1 // indirect
	github.com/SAP/go-hdb v1.6.2 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/aerospike/aerospike-client-go/v6 v6.13.0 // indirect
	github.com/alecthomas/participle/v2 v2.1.0 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/antonmedv/expr v1.15.5 // indirect
	github.com/apache/arrow/go/v12 v12.0.1 // indirect
	github.com/apache/pulsar-client-go v0.8.1 // indirect
	github.com/apache/pulsar-client-go/oauth2 v0.0.0-20220120090717-25e59572242e // indirect
	github.com/apache/thrift v0.19.0 // indirect
	github.com/ardielle/ardielle-go v1.5.2 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aws/aws-sdk-go-v2 v1.22.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.71 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.26 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.29 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.36.0 // indirect
	github.com/aws/smithy-go v1.16.0 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/checkpoint-restore/go-criu/v5 v5.3.0 // indirect
	github.com/cilium/ebpf v0.9.1 // indirect
	github.com/cloudfoundry-incubator/uaago v0.0.0-20190307164349-8136b7bbe76e // indirect
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/containerd/ttrpc v1.2.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/danieljoos/wincred v1.2.0 // indirect
	github.com/denisenkom/go-mssqldb v0.12.3 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/eapache/go-resiliency v1.4.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/efficientgo/core v1.0.0-rc.2 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/facebook/time v0.0.0-20231121165353-cb922d512a84 // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/facebookgo/testname v0.0.0-20150612200628-5443337c3a12 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/zapr v1.2.4 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.4 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-openapi/validate v0.22.1 // indirect
	github.com/go-redis/redis/v7 v7.4.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/cadvisor v0.48.1 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/pprof v0.0.0-20230926050212-f7f687d19a98 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gosnmp/gosnmp v1.37.0 // indirect
	github.com/grobie/gomemcache v0.0.0-20230213081705-239240bbc445 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/cronexpr v1.1.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20230721134942-515895c7690c // indirect
	github.com/hetznercloud/hcloud-go/v2 v2.4.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/influxdata/go-syslog/v3 v3.0.1-0.20230911200830-875f5bc594a4 // indirect
	github.com/influxdata/influxdb-observability/common v0.5.8 // indirect
	github.com/influxdata/influxdb-observability/influx2otel v0.5.8 // indirect
	github.com/influxdata/line-protocol/v2 v2.2.1 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/karrick/godirwalk v1.17.0 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20181214104525-299bdde78165 // indirect
	github.com/leoluk/perflib_exporter v0.2.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/linkedin/goavro/v2 v2.9.8 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/metalmatze/signal v0.0.0-20210307161603-1c9aa721a97a // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/mongodb-forks/digest v1.0.5 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/mrunalp/fileutils v0.5.0 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/muesli/ansi v0.0.0-20211018074035-2e021307bc4b // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.1 // indirect
	github.com/nginxinc/nginx-prometheus-exporter v0.8.1-0.20201110005315-f5a5f8086c19 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/collectd v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/docker v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8stest v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kafka v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/azure v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/signalfx v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/skywalking v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.90.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/winperfcounters v0.90.0 // indirect
	github.com/opencontainers/runc v1.1.9 // indirect
	github.com/opencontainers/selinux v1.11.0 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/openshift/client-go v0.0.0-20230120202327-72f107311084 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.2 // indirect
	github.com/ovh/go-ovh v1.4.3 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/relvacode/iso8601 v1.3.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/seccomp/libseccomp-golang v0.9.2-0.20220502022130-f33da4d89646 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.3 // indirect
	github.com/signalfx/sapm-proto v0.13.0 // indirect
	github.com/sijms/go-ora/v2 v2.7.22 // indirect
	github.com/snowflakedb/gosnowflake v1.7.0 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/tidwall/btree v1.1.0 // indirect
	github.com/tidwall/redcon v1.6.2 // indirect
	github.com/tilinna/clock v1.1.0 // indirect
	github.com/tinylib/msgp v1.1.9 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/vishvananda/netlink v1.2.1-beta.2 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/vmware/go-vmware-nsxt v0.0.0-20230223012718-d31b8a1ca05e // indirect
	github.com/vmware/govmomi v0.33.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.mongodb.org/atlas v0.35.0 // indirect
	go.mongodb.org/mongo-driver v1.13.0 // indirect
	go.opentelemetry.io/collector v0.90.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.90.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v0.90.0 // indirect
	go.opentelemetry.io/collector/config/configgrpc v0.90.0 // indirect
	go.opentelemetry.io/collector/config/confighttp v0.90.0 // indirect
	go.opentelemetry.io/collector/config/confignet v0.90.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v0.90.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.90.0 // indirect
	go.opentelemetry.io/collector/config/configtls v0.90.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.90.0 // indirect
	go.opentelemetry.io/collector/extension/auth v0.90.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.0.0 // indirect
	go.opentelemetry.io/collector/service v0.90.0 // indirect
	go.opentelemetry.io/contrib/config v0.1.1 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.21.1 // indirect
	go.opentelemetry.io/otel/bridge/opencensus v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.21.0 // indirect
	golang.org/x/tools v0.16.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	gonum.org/v1/gonum v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20231120223509-83a465c0220f // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231127180814-3a041ad873d4 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools/v3 v3.4.0 // indirect
	k8s.io/apiextensions-apiserver v0.28.3 // indirect
	k8s.io/component-base v0.28.4 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kubelet v0.28.4 // indirect
	skywalking.apache.org/repo/goapi v0.0.0-20231026090926-09378dd56587 // indirect
)

require (
	cloud.google.com/go/compute v1.23.3 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.2 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.29 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.23 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/aws/aws-sdk-go v1.48.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buraksezer/consistent v0.10.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cncf/xds/go v0.0.0-20230607035331-e9ce68804cb4 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/digitalocean/godo v1.104.1 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v24.0.7+incompatible
	github.com/docker/go-connections v0.4.1-0.20210727194412-58542c764a11
	github.com/docker/go-units v0.5.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.2 // indirect
	github.com/facebookgo/symwalk v0.0.0-20150726040526-42004b9f3222
	github.com/fatih/color v1.15.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.3.0
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/go-cmp v0.6.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gophercloud/gophercloud v1.7.0 // indirect
	github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd // indirect
	github.com/hashicorp/consul/api v1.26.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.4 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/imdario/mergo v0.3.16
	github.com/ionos-cloud/sdk-go/v6 v6.1.9 // indirect
	github.com/jaegertracing/jaeger v1.48.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.17.3 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/linode/linodego v1.23.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/lukejoshuapark/infchan v1.0.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/miekg/dns v1.1.56 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.2.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc5 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/prometheus-community/prom-label-proxy v0.7.0
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/prometheus/statsd_exporter v0.22.8 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/cors v1.10.1 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.21 // indirect
	github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 // indirect
	github.com/shirou/gopsutil/v3 v3.23.10 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.2.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.10 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector/semconv v0.90.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.46.1 // indirect
	go.opentelemetry.io/contrib/zpages v0.46.1 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/mod v0.14.0
	golang.org/x/oauth2 v0.14.0
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/term v0.15.0 // indirect
	golang.org/x/text v0.14.0
	golang.org/x/time v0.4.0 // indirect
	google.golang.org/api v0.151.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9 // indirect
	k8s.io/utils v0.0.0-20230711102312-30195339c3c7
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.3.0 // indirect
	sigs.k8s.io/yaml v1.4.0
)

replace (
	github.com/buraksezer/olric => github.com/fluxninja/olric v0.5.4-fn.patch.12
	github.com/fluxninja/aperture/api/v2 => ./api
	github.com/jsonnet-bundler/jsonnet-bundler => github.com/fluxninja/jsonnet-bundler v0.5.1-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter => github.com/fluxninja/opentelemetry-collector-contrib/exporter/fileexporter v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter => github.com/fluxninja/opentelemetry-collector-contrib/exporter/kafkaexporter v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter => github.com/fluxninja/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter => github.com/fluxninja/opentelemetry-collector-contrib/exporter/signalfxexporter v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter => github.com/fluxninja/opentelemetry-collector-contrib/exporter/splunkhecexporter v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension => github.com/fluxninja/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension => github.com/fluxninja/opentelemetry-collector-contrib/extension/healthcheckextension v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer => github.com/fluxninja/opentelemetry-collector-contrib/extension/observer v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension => github.com/fluxninja/opentelemetry-collector-contrib/extension/pprofextension v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage => github.com/fluxninja/opentelemetry-collector-contrib/extension/storage v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil => github.com/fluxninja/opentelemetry-collector-contrib/internal/aws/awsutil v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight => github.com/fluxninja/opentelemetry-collector-contrib/internal/aws/containerinsight v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil => github.com/fluxninja/opentelemetry-collector-contrib/internal/aws/ecsutil v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s => github.com/fluxninja/opentelemetry-collector-contrib/internal/aws/k8s v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics => github.com/fluxninja/opentelemetry-collector-contrib/internal/aws/metrics v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy => github.com/fluxninja/opentelemetry-collector-contrib/internal/aws/proxy v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray => github.com/fluxninja/opentelemetry-collector-contrib/internal/aws/xray v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/collectd => github.com/fluxninja/opentelemetry-collector-contrib/internal/collectd v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => github.com/fluxninja/opentelemetry-collector-contrib/internal/common v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => github.com/fluxninja/opentelemetry-collector-contrib/internal/coreinternal v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/datadog => github.com/fluxninja/opentelemetry-collector-contrib/internal/datadog v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/docker => github.com/fluxninja/opentelemetry-collector-contrib/internal/docker v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter => github.com/fluxninja/opentelemetry-collector-contrib/internal/filter v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig => github.com/fluxninja/opentelemetry-collector-contrib/internal/k8sconfig v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8stest => github.com/fluxninja/opentelemetry-collector-contrib/internal/k8stest v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kafka => github.com/fluxninja/opentelemetry-collector-contrib/internal/kafka v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet => github.com/fluxninja/opentelemetry-collector-contrib/internal/kubelet v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders => github.com/fluxninja/opentelemetry-collector-contrib/internal/metadataproviders v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent => github.com/fluxninja/opentelemetry-collector-contrib/internal/sharedcomponent v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk => github.com/fluxninja/opentelemetry-collector-contrib/internal/splunk v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr => github.com/fluxninja/opentelemetry-collector-contrib/pkg/batchperresourceattr v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal => github.com/fluxninja/opentelemetry-collector-contrib/pkg/batchpersignal v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata => github.com/fluxninja/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl => github.com/fluxninja/opentelemetry-collector-contrib/pkg/ottl v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest => github.com/fluxninja/opentelemetry-collector-contrib/pkg/pdatatest v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil => github.com/fluxninja/opentelemetry-collector-contrib/pkg/pdatautil v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry => github.com/fluxninja/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza => github.com/fluxninja/opentelemetry-collector-contrib/pkg/stanza v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/azure => github.com/fluxninja/opentelemetry-collector-contrib/pkg/translator/azure v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger => github.com/fluxninja/opentelemetry-collector-contrib/pkg/translator/jaeger v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus => github.com/fluxninja/opentelemetry-collector-contrib/pkg/translator/opencensus v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus => github.com/fluxninja/opentelemetry-collector-contrib/pkg/translator/prometheus v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite => github.com/fluxninja/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/signalfx => github.com/fluxninja/opentelemetry-collector-contrib/pkg/translator/signalfx v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/skywalking => github.com/fluxninja/opentelemetry-collector-contrib/pkg/translator/skywalking v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin => github.com/fluxninja/opentelemetry-collector-contrib/pkg/translator/zipkin v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/winperfcounters => github.com/fluxninja/opentelemetry-collector-contrib/pkg/winperfcounters v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/attributesprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/filterprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/resourceprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/routingprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/servicegraphprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/servicegraphprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/spanprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/transformprocessor v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/aerospikereceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/apachereceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/azureblobreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureeventhubreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/azureeventhubreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/carbonreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/chronyreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/collectdreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/expvarreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/expvarreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/filelogreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/haproxyreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/haproxyreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/httpcheckreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/httpcheckreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/iisreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/iisreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/journaldreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sobjectsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/k8sobjectsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/kafkareceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/nginxreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nsxtreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/nsxtreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/oracledbreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/oracledbreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/podmanreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/pulsarreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/pulsarreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/purefareceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator => github.com/fluxninja/opentelemetry-collector-contrib/receiver/receivercreator v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/redisreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/riakreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/saphanareceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/saphanareceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/sapmreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/signalfxreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/signalfxreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/skywalkingreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmpreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/snmpreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/solacereceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/splunkhecreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/sqlqueryreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/sqlserverreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/statsdreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/syslogreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/udplogreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/wavefrontreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/wavefrontreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.90.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.90.0-fn.patch.1
	go.opentelemetry.io/collector => github.com/fluxninja/opentelemetry-collector v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/component => github.com/fluxninja/opentelemetry-collector/component v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/confmap => github.com/fluxninja/opentelemetry-collector/confmap v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/connector => github.com/fluxninja/opentelemetry-collector/connector v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/consumer => github.com/fluxninja/opentelemetry-collector/consumer v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/exporter => github.com/fluxninja/opentelemetry-collector/exporter v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/exporter/loggingexporter => github.com/fluxninja/opentelemetry-collector/exporter/loggingexporter v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/exporter/otlpexporter => github.com/fluxninja/opentelemetry-collector/exporter/otlpexporter v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/exporter/otlphttpexporter => github.com/fluxninja/opentelemetry-collector/exporter/otlphttpexporter v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/extension => github.com/fluxninja/opentelemetry-collector/extension v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/extension/ballastextension => github.com/fluxninja/opentelemetry-collector/extension/ballastextension v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/extension/zpagesextension => github.com/fluxninja/opentelemetry-collector/extension/zpagesextension v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/otelcol => github.com/fluxninja/opentelemetry-collector/otelcol v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/pdata => github.com/fluxninja/opentelemetry-collector/pdata v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/processor => github.com/fluxninja/opentelemetry-collector/processor v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/processor/batchprocessor => github.com/fluxninja/opentelemetry-collector/processor/batchprocessor v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/processor/memorylimiterprocessor => github.com/fluxninja/opentelemetry-collector/processor/memorylimiterprocessor v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/receiver => github.com/fluxninja/opentelemetry-collector/receiver v0.90.0-fn.patch.1
	go.opentelemetry.io/collector/receiver/otlpreceiver => github.com/fluxninja/opentelemetry-collector/receiver/otlpreceiver v0.90.0-fn.patch.1
)
