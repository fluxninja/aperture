package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var controllerNs string

// ControllerConfig is the config file structure for ARC Controller.
type ControllerConfig struct {
	URL    string `toml:"url"`
	APIKey string `toml:"api_key"`
}

// Config is the config file structure for Aperture.
type Config struct {
	Controller *ControllerConfig `toml:"controller"`
}

// ControllerConn manages flags for connecting to controller – either via
// address or kubeconfig.
type ControllerConn struct {
	// kube is true if controller should be found in Kubernetes cluster.
	kube bool

	controllerAddr string
	allowInsecure  bool
	skipVerify     bool
	kubeConfigPath string
	kubeConfig     *rest.Config
	apiKey         string
	config         string

	forwarderStopChan chan struct{}
	conn              *grpc.ClientConn
}

// InitFlags sets up flags for kubeRestConfig.
func (c *ControllerConn) InitFlags(flags *flag.FlagSet) {
	flags.StringVar(
		&c.controllerAddr,
		"controller",
		"",
		"Address of Aperture Controller",
	)
	flags.BoolVar(
		&c.allowInsecure,
		"insecure",
		false,
		"Allow connection to controller running without TLS",
	)
	flags.BoolVar(
		&c.skipVerify,
		"skip-verify",
		false,
		"Skip TLS certificate verification while connecting to controller",
	)
	flags.BoolVar(
		&c.kube,
		"kube",
		false,
		"Find controller in Kubernetes cluster, instead of connecting directly",
	)
	flags.StringVar(
		&c.kubeConfigPath,
		"kube-config",
		"",
		"Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG",
	)
	flags.StringVar(
		&controllerNs,
		"controller-ns",
		"",
		"Namespace in which the Aperture Controller is running",
	)
	flags.StringVar(
		&c.apiKey,
		"api-key",
		"",
		"FluxNinja ARC API Key to be used when using Cloud Controller",
	)
	flags.StringVar(
		&c.config,
		"config",
		"",
		"Path to the Aperture config file",
	)
}

// PreRunE verifies flags (optionally loading kubeconfig) and should be run at PreRunE stage.
func (c *ControllerConn) PreRunE(_ *cobra.Command, _ []string) error {
	if c.config != "" && c.apiKey != "" {
		return errors.New("--api-key cannot be used with --config")
	}

	// Fetching config from environment variable
	if c.config == "" {
		c.config = os.Getenv(configEnv)
	}

	// Fetching config from default location
	if c.config == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			c.config = filepath.Join(homeDir, ".aperturectl", "config")
		}
	}

	if c.controllerAddr == "" && !c.kube && c.config == "" {
		log.Info().Msg("Neither --controller nor --kube nor --config flag is set. Assuming --kube=true.")
		c.kube = true
	}

	if c.controllerAddr != "" && c.kube {
		return errors.New("--controller cannot be used with --kube")
	}

	if c.kubeConfigPath != "" && !c.kube {
		return errors.New("--kube-config can only be used with --kube")
	}

	if c.kube {
		var err error
		c.kubeConfig, err = GetKubeConfig(c.kubeConfigPath)
		if err != nil {
			return err
		}
	} else if c.config != "" {
		var err error
		c.config, err = filepath.Abs(c.config)
		if err != nil {
			return fmt.Errorf("failed to resolve config file '%s' path: %w", c.config, err)
		}

		config := &Config{}
		_, err = toml.DecodeFile(c.config, config)
		if err != nil {
			return fmt.Errorf("failed to read config file '%s': %w", c.config, err)
		}

		if config.Controller == nil {
			return fmt.Errorf("invalid config file '%s'. Missing key 'controller'", c.config)
		}

		if config.Controller.URL == "" {
			return fmt.Errorf("invalid config file '%s'. Missing key 'controller.url'", c.config)
		}

		if config.Controller.APIKey == "" {
			return fmt.Errorf("invalid config file '%s'. Missing key 'controller.api_key'", c.config)
		}

		c.controllerAddr = config.Controller.URL
		c.apiKey = config.Controller.APIKey
	}

	return nil
}

// Client returns Controller Client, connecting to controller if not yet connected.
func (c *ControllerConn) Client() (cmdv1.ControllerClient, error) {
	if c.conn != nil {
		return cmdv1.NewControllerClient(c.conn), nil
	}

	var addr string
	var cred credentials.TransportCredentials
	if c.allowInsecure {
		cred = insecure.NewCredentials()
	} else if c.skipVerify {
		cred = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // Requires enabling CLI option
		})
	}

	if !c.kube {
		addr = c.controllerAddr

		if cred == nil {
			certPool, err := x509.SystemCertPool()
			if err != nil {
				return nil, err
			}
			cred = credentials.NewClientTLSFromCert(certPool, "")
		}
	} else {
		deployment, err := GetControllerDeployment(c.kubeConfig, controllerNs)
		if err != nil {
			return nil, err
		}
		controllerNs = deployment.GetNamespace()
		port, cert, err := c.startPortForward()
		if err != nil {
			return nil, fmt.Errorf("failed to start port forward for Aperture Controller: %w", err)
		}

		addr = fmt.Sprintf("localhost:%d", port)

		if cred == nil {
			if cert == nil {
				return nil, errors.New("cannot find controller cert and --insecure is off")
			}

			certPool := x509.NewCertPool()
			ok := certPool.AppendCertsFromPEM(cert)
			if !ok {
				return nil, fmt.Errorf("cannot apply controller cert")
			}
			cred = credentials.NewClientTLSFromCert(certPool, fmt.Sprintf("%s.%s", deployment.GetName(), deployment.GetNamespace()))
		}
	}

	var err error
	c.conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(cred), grpc.WithUnaryInterceptor(c.cloudControllerInterceptor))
	if err != nil {
		return nil, err
	}

	return cmdv1.NewControllerClient(c.conn), nil
}

// PostRun cleans up ControllerConn's resources, and should be run at PostRun stage.
func (c *ControllerConn) PostRun(_ *cobra.Command, _ []string) {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close controller connection")
		}
	}

	if c.forwarderStopChan != nil {
		close(c.forwarderStopChan)
	}
}

func (c *ControllerConn) startPortForward() (localPort uint16, cert []byte, err error) {
	clientset, err := kubernetes.NewForConfig(c.kubeConfig)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// FIXME Forwarding to a service would be nicer solution, but could not make
	// it work for some reason, thus forwarding to pod directly.
	pods, err := clientset.CoreV1().Pods(controllerNs).List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.Set{"app.kubernetes.io/component": "aperture-controller"}.String(),
		FieldSelector: labels.Set{"status.phase": "Running"}.String(),
	})
	if err != nil {
		return 0, nil, fmt.Errorf("failed to list pods: %w", err)
	}
	if len(pods.Items) == 0 {
		return 0, nil, fmt.Errorf("no pods found")
	}

	pod := &pods.Items[0]
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", pod.Namespace, pod.Name)

	transport, upgrader, err := spdy.RoundTripperFor(c.kubeConfig)
	if err != nil {
		return 0, nil, err
	}

	hostIP := strings.TrimPrefix(c.kubeConfig.Host, "https://")
	targetURL := url.URL{Scheme: "https", Path: path, Host: hostIP}
	dialer := spdy.NewDialer(
		upgrader,
		&http.Client{Transport: transport},
		http.MethodPost,
		&targetURL,
	)

	c.forwarderStopChan = make(chan struct{})
	readyChan := make(chan struct{})
	fw, err := portforward.New(
		dialer,
		[]string{":8080"},
		c.forwarderStopChan,
		readyChan,
		io.Discard,
		io.Discard,
	)
	if err != nil {
		return 0, nil, err
	}

	fwErrChan := make(chan error, 1)
	go func() {
		fwErrChan <- fw.ForwardPorts()
	}()

	secrets, err := clientset.CoreV1().Secrets(controllerNs).List(
		context.Background(),
		metav1.ListOptions{
			LabelSelector: labels.Set{"app.kubernetes.io/name": "aperture"}.String(),
		},
	)
	if err != nil || len(secrets.Items) == 0 {
		return 0, nil, fmt.Errorf("no secrets found for Aperture Controller certificate")
	}

	for _, secret := range secrets.Items {
		if !strings.HasSuffix(secret.Name, "controller-cert") {
			continue
		}
		cert = secret.Data["crt.pem"]
	}

	select {
	case err = <-fwErrChan:
		return 0, nil, err
	case <-readyChan:
	}
	ports, err := fw.GetPorts()
	if err != nil {
		return 0, nil, err
	}

	return ports[0].Local, cert, nil
}

// IsKube returns true if controller should be found in Kubernetes cluster.
func (c *ControllerConn) IsKube() bool {
	return c.kube
}

// GetKubeRestConfig returns kubeRestConfig.
func (c *ControllerConn) GetKubeRestConfig() *rest.Config {
	return c.kubeConfig
}

// GetControllerNs returns namespace in which the Aperture Controller is running.
func GetControllerNs() string {
	return controllerNs
}

func (c *ControllerConn) cloudControllerInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if c.apiKey == "" {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	md := metadata.Pairs("apikey", c.apiKey)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
