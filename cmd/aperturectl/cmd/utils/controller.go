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
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

var controllerNs string

// ControllerConn manages flags for connecting to controller – either via
// address or kubeconfig.
type ControllerConn struct {
	controllerAddr string
	allowInsecure  bool
	isKube         bool
	kubeConfigPath string
	kubeConfig     *rest.Config

	forwarderStopChan chan struct{}
	conn              *grpc.ClientConn
}

// InitFlags sets up flags for kubeRestConfig.
func (c *ControllerConn) InitFlags(flags *flag.FlagSet) {
	flags.StringVar(
		&c.controllerAddr,
		"controller",
		"",
		"Address of Aperture controller",
	)
	flags.BoolVar(
		&c.allowInsecure,
		"insecure",
		false,
		"Allow insecure connection to controller",
	)
	flags.BoolVar(
		&c.isKube,
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
}

// PreRunE verifies flags (optionally loading kubeconfig) and should be run at PreRunE stage.
func (c *ControllerConn) PreRunE(_ *cobra.Command, _ []string) error {
	if c.controllerAddr == "" && !c.isKube {
		return errors.New("either --controller or --kube should be set")
	}

	if c.controllerAddr != "" && c.isKube {
		return errors.New("--controller cannot be used with --kube")
	}

	if c.kubeConfigPath != "" && !c.isKube {
		return errors.New("--kube-config can only be used with --kube")
	}

	if c.isKube {
		var err error
		c.kubeConfig, err = GetKubeConfig(c.kubeConfigPath)
		if err != nil {
			return err
		}
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
		cred = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // Requires enabling CLI option
		})
	}

	if !c.isKube {
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
	c.conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(cred))
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
			LabelSelector: labels.Set{"app.kubernetes.io/component": "aperture-controller"}.String(),
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
