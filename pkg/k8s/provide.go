package k8s

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/zerologr"
	"go.uber.org/fx"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"
	"k8s.io/client-go/transport"
	"k8s.io/klog/v2"

	"github.com/fluxninja/aperture/v2/pkg/log"
	commonhttp "github.com/fluxninja/aperture/v2/pkg/net/http"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

var (
	// swagger:operation POST /kubernetes_client common-configuration KubernetesClient
	// ---
	// x-fn-config-env: true
	// parameters:
	// - name: http_client
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/HTTPClientConfig"

	// KubernetesClientConfigKey is the key used to store the KubernetesClientConfig in the config.
	kubernetesClientConfigKey = "kubernetes_client"
	// HttpConfigKey is the key used to store the HTTPClientConfig in the config.
	httpConfigKey = strings.Join([]string{kubernetesClientConfigKey, "http_client"}, ".")
)

// K8sClientConstructorIn holds parameter for Providek8sClient and Providek8sDynamicClient.
type K8sClientConstructorIn struct {
	fx.In
	HTTPClient *http.Client `name:"k8s-http-client"`
	Shutdowner fx.Shutdowner
}

// Module provides a K8sClient.
func Module() fx.Option {
	return fx.Options(
		commonhttp.ClientConstructor{Name: "k8s-http-client", ConfigKey: httpConfigKey}.Annotate(),
		fx.Provide(Providek8sClient),
	)
}

// K8sClient provides an interface for kubernetes client.
type K8sClient interface {
	GetClientSet() *kubernetes.Clientset
	GetScaleClient() scale.ScalesGetter
	GetDynamicClient() dynamic.Interface
	GetRESTMapper() apimeta.RESTMapper
	ScaleForGroupKind(context.Context, string, string, schema.GroupKind) (*autoscalingv1.Scale, schema.GroupResource, error)
}

// RealK8sClient provides access to Kubernetes Clients.
type RealK8sClient struct {
	clientSet     *kubernetes.Clientset
	scaleClient   scale.ScalesGetter
	dynamicClient dynamic.Interface
	mapper        apimeta.RESTMapper
}

// RealK8sClient implements K8sClient.
var _ K8sClient = &RealK8sClient{}

// NewK8sClient returns a new kubernetes client, or nil if outside a Kubernetes cluster.
func NewK8sClient(httpClient *http.Client, shutdowner fx.Shutdowner) (*RealK8sClient, error) {
	clientSet, config, err := newK8sClientSet(httpClient, shutdowner)
	if err == rest.ErrNotInCluster {
		log.Info().Msg("Not in Kubernetes Cluster, creating nil client")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	zerolog := log.WithComponent("k8s-client").GetZerolog()
	klog.SetLogger(zerologr.New(zerolog))

	discoveryClient := clientSet.DiscoveryClient
	cachedDiscoveryClient := cacheddiscovery.NewMemCacheClient(discoveryClient)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(discoveryClient)
	scaleClient, err := scale.NewForConfig(config, mapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)
	if err != nil {
		log.Fatal().Err(err).Msg("Unexpected error, unable to create Kubernetes Scale Client")
		utils.Shutdown(shutdowner)
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Unexpected error, unable to create Kubernetes Dynamic Client")
		utils.Shutdown(shutdowner)
		return nil, err
	}

	return &RealK8sClient{
		clientSet:     clientSet,
		scaleClient:   scaleClient,
		dynamicClient: dynamicClient,
		mapper:        mapper,
	}, nil
}

// GetClientSet returns the kubernetes client set.
func (r *RealK8sClient) GetClientSet() *kubernetes.Clientset {
	return r.clientSet
}

// GetScaleClient returns the kubernetes scale client.
func (r *RealK8sClient) GetScaleClient() scale.ScalesGetter {
	return r.scaleClient
}

// GetDynamicClient returns the kubernetes dynamic client.
func (r *RealK8sClient) GetDynamicClient() dynamic.Interface {
	return r.dynamicClient
}

// GetRESTMapper returns the rest mapper.
func (r *RealK8sClient) GetRESTMapper() apimeta.RESTMapper {
	return r.mapper
}

// ScaleForGroupKind attempts to fetch the scale for the given Group and Kind.
// The possible Resources for the group and kind are retrieved. Scale is fetched
// for each Resource in the RESTMapping with the given name and namespace, until
// a working one is found.  If none work, the first error is returned.  It returns
// both the scale, as well as the group-resource from the working mapping.
func (r *RealK8sClient) ScaleForGroupKind(ctx context.Context, namespace, name string, groupKind schema.GroupKind) (*autoscalingv1.Scale, schema.GroupResource, error) {
	mappings, err := r.GetRESTMapper().RESTMappings(groupKind)
	if err != nil {
		return nil, schema.GroupResource{}, err
	}
	var firstErr error
	for i, mapping := range mappings {
		targetGR := mapping.Resource.GroupResource()
		scale, err := r.GetScaleClient().Scales(namespace).Get(ctx, targetGR, name, metav1.GetOptions{})
		if err == nil {
			return scale, targetGR, nil
		}

		// if this is the first error, remember it,
		// then go on and try other mappings until we find a good one
		if i == 0 {
			firstErr = err
		}
	}

	// make sure we handle an empty set of mappings
	if firstErr == nil {
		firstErr = fmt.Errorf("unrecognized resource")
	}

	return nil, schema.GroupResource{}, firstErr
}

// Providek8sClient provides a new kubernetes client and sets logger.
func Providek8sClient(in K8sClientConstructorIn) (K8sClient, error) {
	return NewK8sClient(in.HTTPClient, in.Shutdowner)
}

func newK8sClientSet(client *http.Client, shutdowner fx.Shutdowner) (*kubernetes.Clientset, *rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Info().Msg("Not in Kubernetes cluster, could not create Kubernetes client")
		return nil, nil, err
	}

	transportConfig, err := config.TransportConfig()
	if err != nil {
		// Call shutdowner to catch this issue early since Kubernetes Client is marked as optional downstream.
		log.Fatal().Err(err).Msg("Unexpected error creating Kubernetes Client's transport config")
		utils.Shutdown(shutdowner)
		return nil, nil, err
	}
	tlsConfig, err := transport.TLSConfigFor(transportConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Unexpected error creating Kubernetes Client's TLS Config")
		utils.Shutdown(shutdowner)
		return nil, nil, err
	}
	transport := client.Transport.(*http.Transport)
	transport.TLSClientConfig = tlsConfig

	config.TLSClientConfig = rest.TLSClientConfig{}
	config.Transport = transport
	config.Timeout = client.Timeout
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Unexpected error creating Kubernetes Client")
		utils.Shutdown(shutdowner)
		return nil, nil, err
	}
	return clientSet, config, nil
}
