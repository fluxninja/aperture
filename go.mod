module github.com/fluxninja/aperture

go 1.20

require (
	github.com/Henry-Sarabia/sliceconv v1.0.2
	github.com/Masterminds/semver/v3 v3.2.0
	github.com/benlaurie/objecthash v0.0.0-20180202135721-d1e3d6079fc1
	github.com/buger/jsonparser v1.1.1
	github.com/buraksezer/olric v0.4.9
	github.com/cenkalti/backoff/v4 v4.2.0
	github.com/charmbracelet/bubbletea v0.23.2
	github.com/clarketm/json v1.17.1
	github.com/containerd/cgroups v1.1.0
	github.com/eapache/queue v1.1.0
	github.com/elastic/gmux v0.2.0
	github.com/elastic/gosigar v0.14.2
	github.com/emicklei/dot v1.3.1
	github.com/envoyproxy/go-control-plane v0.11.0
	github.com/fluxninja/datasketches-go v0.0.0-20220916235224-7501a2d28551
	github.com/fluxninja/lumberjack v0.0.0-20220729045908-655029e4d814
	github.com/getsentry/sentry-go v0.18.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-git/go-git/v5 v5.6.0
	github.com/go-logr/zerologr v1.2.3
	github.com/go-openapi/runtime v0.25.0
	github.com/go-openapi/strfmt v0.21.3
	github.com/go-playground/validator/v10 v10.11.2
	github.com/gofrs/flock v0.8.1
	github.com/golang/mock v1.6.0
	github.com/google/go-jsonnet v0.19.1
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/memberlist v0.5.0
	github.com/integralist/go-findroot v0.0.0-20160518114804-ac90681525dc
	github.com/jonboulle/clockwork v0.3.0
	github.com/json-iterator/go v1.1.12
	github.com/jsonnet-bundler/jsonnet-bundler v0.5.1
	github.com/knadh/koanf v1.5.0
	github.com/lithammer/dedent v1.1.0
	github.com/looplab/tarjan v0.1.0
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/natefinch/atomic v1.0.1
	github.com/onsi/ginkgo/v2 v2.8.4
	github.com/onsi/gomega v1.27.2
	github.com/open-policy-agent/opa v0.49.2
	github.com/open-policy-agent/opa-envoy-plugin v0.49.2-envoy
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/datadogprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/servicegraphprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureeventhubreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dotnetdiagnosticsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/expvarreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/haproxyreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/httpcheckreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/iisreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sobjectsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nsxtreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/oracledbreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusexecreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/pulsarreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/saphanareceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/signalfxreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmpreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/wavefrontreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.73.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.73.0
	github.com/prometheus/alertmanager v0.25.0
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/common v0.41.0
	github.com/reugn/go-quartz v0.6.0
	github.com/rs/zerolog v1.29.0
	github.com/sourcegraph/conc v0.3.0
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.2
	github.com/technosophos/moniker v0.0.0-20210218184952-3ea787d3943b
	github.com/xeipuuv/gojsonschema v1.2.0
	go.etcd.io/etcd/api/v3 v3.5.7
	go.etcd.io/etcd/client/v3 v3.5.7
	go.opentelemetry.io/collector v0.73.0
	go.opentelemetry.io/collector/component v0.73.0
	go.opentelemetry.io/collector/confmap v0.73.0
	go.opentelemetry.io/collector/consumer v0.73.0
	go.opentelemetry.io/collector/exporter v0.73.0
	go.opentelemetry.io/collector/exporter/loggingexporter v0.73.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.73.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.73.0
	go.opentelemetry.io/collector/extension/ballastextension v0.73.0
	go.opentelemetry.io/collector/extension/zpagesextension v0.73.0
	go.opentelemetry.io/collector/pdata v1.0.0-rc7
	go.opentelemetry.io/collector/processor/batchprocessor v0.73.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.73.0
	go.opentelemetry.io/collector/receiver v0.73.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.73.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.40.0
	go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/proto/otlp v0.19.0
	go.uber.org/automaxprocs v1.5.1
	go.uber.org/fx v1.19.2
	go.uber.org/goleak v1.2.1
	go.uber.org/multierr v1.9.0
	go.uber.org/zap v1.24.0
	golang.org/x/exp v0.0.0-20230224173230-c95f2b4c22f2
	golang.org/x/net v0.8.0
	google.golang.org/genproto v0.0.0-20230223222841-637eb2293923
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/yaml.v3 v3.0.1
	helm.sh/helm/v3 v3.11.1
	k8s.io/api v0.26.2
	k8s.io/apimachinery v0.26.2
	k8s.io/client-go v0.26.2
	k8s.io/klog/v2 v2.90.0
	sigs.k8s.io/controller-runtime v0.14.4
)

require (
	cloud.google.com/go v0.107.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.11.0 // indirect
	cloud.google.com/go/pubsub v1.28.0 // indirect
	cloud.google.com/go/spanner v1.44.0 // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20211115184647-b584dd5df32c // indirect
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible // indirect
	code.cloudfoundry.org/rfc5424 v0.0.0-20201103192249-000122071b78 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/AthenZ/athenz v1.10.39 // indirect
	github.com/Azure/azure-amqp-common-go/v4 v4.0.0 // indirect
	github.com/Azure/azure-event-hubs-go/v3 v3.4.0 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.22.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.9.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.3.0 // indirect
	github.com/Azure/azure-storage-blob-go v0.15.0 // indirect
	github.com/Azure/go-amqp v0.18.1 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.11 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-agent/pkg/trace v0.43.0-rc.3.0.20230227093656-ca3af9ce4c3c // indirect
	github.com/DataDog/datadog-agent/pkg/util/cgroups v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-agent/pkg/util/log v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-agent/pkg/util/pointer v0.43.0-rc.3 // indirect
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.43.0-rc.3.0.20230206114529-17c7dfde736c // indirect
	github.com/DataDog/datadog-go/v5 v5.1.1 // indirect
	github.com/DataDog/go-tuf v0.3.0--fix-localmeta-fork // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.1.2 // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metrics v0.1.1 // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile v0.1.2 // indirect
	github.com/DataDog/sketches-go v1.4.1 // indirect
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.11.2 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20221026131551-cf6655e29de4 // indirect
	github.com/ReneKroon/ttlcache/v2 v2.11.0 // indirect
	github.com/SAP/go-hdb v1.1.1 // indirect
	github.com/Shopify/sarama v1.38.1 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/aerospike/aerospike-client-go/v6 v6.10.0 // indirect
	github.com/alecthomas/participle/v2 v2.0.0-beta.5 // indirect
	github.com/antonmedv/expr v1.12.1 // indirect
	github.com/apache/arrow/go/arrow v0.0.0-20211112161151-bc219186db40 // indirect
	github.com/apache/pulsar-client-go v0.8.1 // indirect
	github.com/apache/pulsar-client-go/oauth2 v0.0.0-20220120090717-25e59572242e // indirect
	github.com/apache/thrift v0.18.1 // indirect
	github.com/ardielle/ardielle-go v1.5.2 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go-v2 v1.17.5 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.15 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.33 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.29 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.27.11 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
	github.com/aymanbagabas/go-osc52 v1.2.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/checkpoint-restore/go-criu/v5 v5.3.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/cilium/ebpf v0.9.1 // indirect
	github.com/cloudflare/circl v1.1.0 // indirect
	github.com/cloudfoundry-incubator/uaago v0.0.0-20190307164349-8136b7bbe76e // indirect
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/containerd/ttrpc v1.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/denisenkom/go-mssqldb v0.12.2 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230111030713-bf00bc1b83b6 // indirect
	github.com/emicklei/go-restful/v3 v3.10.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/facebook/time v0.0.0-20220713225404-f7a0d7702d50 // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/facebookgo/testname v0.0.0-20150612200628-5443337c3a12 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/gabriel-vasile/mimetype v1.4.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.4.0 // indirect
	github.com/go-logr/zapr v1.2.3 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.3 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/spec v0.20.7 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-openapi/validate v0.22.0 // indirect
	github.com/go-redis/redis/v7 v7.4.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/google/cadvisor v0.47.1 // indirect
	github.com/google/pprof v0.0.0-20230111200839-76d1ae5aea2b // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gosnmp/gosnmp v1.35.0 // indirect
	github.com/grobie/gomemcache v0.0.0-20180201122607-1f779c573665 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20230124213148-69fd1a0e4bf7 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/influxdata/go-syslog/v3 v3.0.1-0.20210608084020-ac565dc76ba6 // indirect
	github.com/influxdata/influxdb-observability/common v0.3.1 // indirect
	github.com/influxdata/influxdb-observability/influx2otel v0.3.1 // indirect
	github.com/influxdata/line-protocol/v2 v2.2.1 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.3 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/karrick/godirwalk v1.17.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20181214104525-299bdde78165 // indirect
	github.com/leoluk/perflib_exporter v0.2.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/linkedin/goavro/v2 v2.9.8 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-ieproxy v0.0.9 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/mongodb-forks/digest v1.0.4 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/mrunalp/fileutils v0.5.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/muesli/ansi v0.0.0-20211018074035-2e021307bc4b // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.14.0 // indirect
	github.com/nginxinc/nginx-prometheus-exporter v0.8.1-0.20201110005315-f5a5f8086c19 // indirect
	github.com/observiq/ctimefmt v1.0.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/docker v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/signalfx v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/winperfcounters v0.73.0 // indirect
	github.com/opencontainers/runc v1.1.4 // indirect
	github.com/opencontainers/selinux v1.10.1 // indirect
	github.com/openlyinc/pointy v1.2.0 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/openshift/client-go v0.0.0-20230120202327-72f107311084 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	github.com/outcaste-io/ristretto v0.2.0 // indirect
	github.com/ovh/go-ovh v1.3.0 // indirect
	github.com/panta/machineid v1.0.2 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/prometheus v0.42.1-0.20230210113933-af1d9e01c7e4 // indirect
	github.com/relvacode/iso8601 v1.3.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/seccomp/libseccomp-golang v0.9.2-0.20220502022130-f33da4d89646 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.4.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.3 // indirect
	github.com/signalfx/sapm-proto v0.12.0 // indirect
	github.com/sijms/go-ora/v2 v2.5.33 // indirect
	github.com/skeema/knownhosts v1.1.0 // indirect
	github.com/snowflakedb/gosnowflake v1.6.18 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/tilinna/clock v1.1.0 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/vishvananda/netlink v1.1.1-0.20210330154013-f5de75959ad5 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/vmware/go-vmware-nsxt v0.0.0-20220328155605-f49a14c1ef5f // indirect
	github.com/vmware/govmomi v0.30.2 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	go.mongodb.org/atlas v0.23.0 // indirect
	go.mongodb.org/mongo-driver v1.11.2 // indirect
	go.opentelemetry.io/collector/featuregate v0.73.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.14.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	gonum.org/v1/gonum v0.12.0 // indirect
	google.golang.org/grpc/examples v0.0.0-20211119005141-f45e61797429 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools/v3 v3.3.0 // indirect
	k8s.io/apiextensions-apiserver v0.26.1 // indirect
	k8s.io/component-base v0.26.1 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kubelet v0.26.1 // indirect
	skywalking.apache.org/repo/goapi v0.0.0-20220121092418-9c455d0dda3f // indirect
)

require (
	cloud.google.com/go/compute v1.18.0 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.2 // indirect
	github.com/Azure/azure-sdk-for-go v67.1.0+incompatible // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.22 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/aws/aws-sdk-go v1.44.216 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buraksezer/connpool v0.6.0 // indirect
	github.com/buraksezer/consistent v0.10.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cncf/xds/go v0.0.0-20230105202645-06c439db220b // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/digitalocean/godo v1.95.0 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v23.0.1+incompatible
	github.com/docker/go-connections v0.4.1-0.20210727194412-58542c764a11
	github.com/docker/go-units v0.5.0 // indirect
	github.com/elastic/go-licenser v0.4.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.9.1 // indirect
	github.com/facebookgo/symwalk v0.0.0-20150726040526-42004b9f3222
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.3 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/flatbuffers v22.11.23+incompatible // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gophercloud/gophercloud v1.1.1 // indirect
	github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd // indirect
	github.com/hashicorp/consul/api v1.19.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.4.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.6.0 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hetznercloud/hcloud-go v1.39.0 // indirect
	github.com/imdario/mergo v0.3.13
	github.com/ionos-cloud/sdk-go/v6 v6.1.3 // indirect
	github.com/jaegertracing/jaeger v1.41.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.16.0 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/linode/linodego v1.12.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/lukejoshuapark/infchan v1.0.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.1.17 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.73.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/bigipreceiver v0.73.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20220909204839-494a5a6aca78 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/prometheus/statsd_exporter v0.22.8 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/cors v1.8.3 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.12 // indirect
	github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 // indirect
	github.com/shirou/gopsutil/v3 v3.23.2 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.1.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.7 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector/semconv v0.73.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.40.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.40.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.37.0 // indirect
	go.opentelemetry.io/otel/metric v0.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.14.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/mod v0.9.0
	golang.org/x/oauth2 v0.6.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/text v0.8.0
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/api v0.111.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	k8s.io/kube-openapi v0.0.0-20221207184640-f3cff1453715 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/jsonnet-bundler/jsonnet-bundler => github.com/fluxninja/jsonnet-bundler v0.5.1-fn.patch.1
	go.opentelemetry.io/collector => github.com/fluxninja/opentelemetry-collector v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/component => github.com/fluxninja/opentelemetry-collector/component v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/confmap => github.com/fluxninja/opentelemetry-collector/confmap v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/consumer => github.com/fluxninja/opentelemetry-collector/consumer v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/exporter/loggingexporter => github.com/fluxninja/opentelemetry-collector/exporter/loggingexporter v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/exporter/otlpexporter => github.com/fluxninja/opentelemetry-collector/exporter/otlpexporter v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/exporter/otlphttpexporter => github.com/fluxninja/opentelemetry-collector/exporter/otlphttpexporter v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/extension/ballastextension => github.com/fluxninja/opentelemetry-collector/extension/ballastextension v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/extension/zpagesextension => github.com/fluxninja/opentelemetry-collector/extension/zpagesextension v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/pdata => github.com/fluxninja/opentelemetry-collector/pdata v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/processor/batchprocessor => github.com/fluxninja/opentelemetry-collector/processor/batchprocessor v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/processor/memorylimiterprocessor => github.com/fluxninja/opentelemetry-collector/processor/memorylimiterprocessor v0.73.0-fn.patch.1
	go.opentelemetry.io/collector/receiver/otlpreceiver => github.com/fluxninja/opentelemetry-collector/receiver/otlpreceiver v0.73.0-fn.patch.1
)

replace (
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter => github.com/fluxninja/opentelemetry-collector-contrib/exporter/fileexporter v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter => github.com/fluxninja/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension => github.com/fluxninja/opentelemetry-collector-contrib/extension/healthcheckextension v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension => github.com/fluxninja/opentelemetry-collector-contrib/extension/pprofextension v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/attributesprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/datadogprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/datadogprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/filterprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/resourceprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/routingprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/servicegraphprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/servicegraphprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/spanprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor => github.com/fluxninja/opentelemetry-collector-contrib/processor/transformprocessor v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/aerospikereceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/apachereceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/azureblobreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureeventhubreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/azureeventhubreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/carbonreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/chronyreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/collectdreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dotnetdiagnosticsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/dotnetdiagnosticsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/expvarreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/expvarreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/filelogreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/haproxyreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/haproxyreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/httpcheckreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/httpcheckreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/iisreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/iisreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/journaldreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sobjectsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/k8sobjectsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/kafkareceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/nginxreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nsxtreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/nsxtreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/oracledbreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/oracledbreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/podmanreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusexecreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/prometheusexecreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/pulsarreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/pulsarreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/purefareceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator => github.com/fluxninja/opentelemetry-collector-contrib/receiver/receivercreator v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/redisreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/riakreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/saphanareceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/saphanareceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/sapmreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/signalfxreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/signalfxreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/skywalkingreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmpreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/snmpreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/solacereceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/splunkhecreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/sqlqueryreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/sqlserverreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/statsdreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/syslogreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/udplogreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/wavefrontreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/wavefrontreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.73.0-fn.patch.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver => github.com/fluxninja/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.73.0-fn.patch.1
)
