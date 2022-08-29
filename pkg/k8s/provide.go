package k8s

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/zerologr"
	"go.uber.org/fx"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
	"k8s.io/klog/v2"

	"github.com/fluxninja/aperture/pkg/log"
	commonhttp "github.com/fluxninja/aperture/pkg/net/http"
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
	K8sClient *http.Client `name:"k8s-http-client"`
	Logger    log.Logger
}

// Module provides a K8sClient.
func Module() fx.Option {
	return fx.Options(
		commonhttp.ClientConstructor{Name: "k8s-http-client", Key: httpConfigKey}.Annotate(),
		fx.Provide(Providek8sClient),
	)
}

// Providek8sDynamicClient provides a dynamic kubernetes client.
func Providek8sDynamicClient(in K8sClientConstructorIn) (dynamic.Interface, error) {
	dynamicClient, err := newk8sDynamicClient(in.K8sClient)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new dynamic client!")
		return nil, err
	}

	return dynamicClient, nil
}

func newk8sDynamicClient(client *http.Client) (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// client-go does not allow custom Transport when TLS is enabled.
	// We need to extract TLS settings from inClusterConfig and put them in our
	// custom transport.
	transportConfig, err := config.TransportConfig()
	if err != nil {
		return nil, fmt.Errorf("creating transport config: %v", err)
	}
	tlsConfig, err := transport.TLSConfigFor(transportConfig)
	if err != nil {
		return nil, fmt.Errorf("creating TLS Config: %v", err)
	}
	transport := client.Transport.(*http.Transport)
	transport.TLSClientConfig = tlsConfig

	// as TLS config is provided inside our custom transport, we need to erase
	// TLS settings here.
	config.TLSClientConfig = rest.TLSClientConfig{}
	config.Transport = transport
	config.Timeout = client.Timeout

	return dynamic.NewForConfig(config)
}

// K8sClient provides an interface for kubernetes client.
type K8sClient interface {
	GetClientSet() *kubernetes.Clientset
	GetErr() error
	GetErrNotInCluster() bool
}

// RealK8sClient implements kubernetes client set.
type RealK8sClient struct {
	clientSet *kubernetes.Clientset
	err       error
}

// NewK8sClient returns a new kubernetes client.
func NewK8sClient(clientSet *kubernetes.Clientset, err error) *RealK8sClient {
	return &RealK8sClient{
		clientSet: clientSet,
		err:       err,
	}
}

// GetClientSet returns the kubernetes client set.
func (r *RealK8sClient) GetClientSet() *kubernetes.Clientset {
	return r.clientSet
}

// GetErr returns the error of the client.
func (r *RealK8sClient) GetErr() error {
	return r.err
}

// GetErrNotInCluster returns true if client's error equals to ErrNotInCluster, unable to load in-cluster configuration.
func (r *RealK8sClient) GetErrNotInCluster() bool {
	return r.err == rest.ErrNotInCluster
}

// Providek8sClient provides a new kubernetes client and sets logger.
func Providek8sClient(in K8sClientConstructorIn) K8sClient {
	k8sClientSet, err := newk8sClientSetAndErr(in.K8sClient)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new client set!")
	}
	zerolog := in.Logger.Zerolog()
	klog.SetLogger(zerologr.New(zerolog))
	return NewK8sClient(k8sClientSet, err)
}

func newk8sClientSetAndErr(client *http.Client) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	transportConfig, err := config.TransportConfig()
	if err != nil {
		return nil, fmt.Errorf("creating transport config: %v", err)
	}
	tlsConfig, err := transport.TLSConfigFor(transportConfig)
	if err != nil {
		return nil, fmt.Errorf("creating TLS Config: %v", err)
	}
	transport := client.Transport.(*http.Transport)
	transport.TLSClientConfig = tlsConfig

	config.TLSClientConfig = rest.TLSClientConfig{}
	config.Transport = transport
	config.Timeout = client.Timeout
	return kubernetes.NewForConfig(config)
}
