package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// ControllerConfig is the config file structure for Aperture Cloud Controller.
type ControllerConfig struct {
	// When changing fields, remember to update docs/content/reference/configuration/aperturectl.md.
	URL    string `toml:"url"`
	APIKey string `toml:"api_key"`
}

// Config is the config file structure for Aperture.
type Config struct {
	// When changing fields, remember to update docs/content/reference/configuration/aperturectl.md.
	Controller *ControllerConfig `toml:"controller"`
}

// ControllerConn manages flags for connecting to controller – either via
// address or kubeconfig.
type ControllerConn struct {
	controllerAddr string
	allowInsecure  bool
	skipVerify     bool
	apiKey         string
	config         string

	forwarderStopChan chan struct{}
	conn              *grpc.ClientConn
}

// CloudInitFlags sets up flags for Cloud Controller.
func (c *ControllerConn) InitFlags(flags *flag.FlagSet) {
	flags.StringVar(
		&c.controllerAddr,
		"controller",
		"",
		"Address of Aperture Cloud Controller",
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
	flags.StringVar(
		&c.apiKey,
		"api-key",
		"",
		"Aperture Cloud API Key to be used when using Cloud Controller",
	)
	flags.StringVar(
		&c.config,
		"config",
		"",
		"Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG",
	)
}

// PreRunE verifies flags (optionally loading kubeconfig) and should be run at PreRunE stage.
func (c *ControllerConn) PreRunE(_ *cobra.Command, _ []string) error {
	// Fetching config from environment variable
	if c.config == "" {
		c.config = os.Getenv(configEnv)
	}

	// Fetching config from default location
	if c.config == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			c.config = filepath.Join(homeDir, ".aperturectl", "config")
			if _, err := os.Stat(c.config); err != nil {
				c.config = ""
			}
		}
	}

	if c.config != "" && (c.controllerAddr == "" || c.apiKey == "") {
		var err error
		c.config, err = filepath.Abs(c.config)
		if err != nil {
			return fmt.Errorf("failed to resolve config file '%s' path: %w", c.config, err)
		}

		log.Info().Msgf("Using config file '%s'", c.config)
		config := &Config{}
		_, err = toml.DecodeFile(c.config, config)
		if err != nil {
			return fmt.Errorf("failed to read config file '%s': %w", c.config, err)
		}

		if config.Controller == nil {
			return fmt.Errorf("invalid config file '%s'. Missing key 'controller'", c.config)
		}

		if config.Controller.URL == "" && c.controllerAddr == "" {
			return fmt.Errorf("invalid config file '%s'. Missing key 'controller.url'", c.config)
		}

		if config.Controller.APIKey == "" && c.apiKey == "" {
			return fmt.Errorf("invalid config file '%s'. Missing key 'controller.api_key'", c.config)
		}

		if c.controllerAddr == "" {
			c.controllerAddr = config.Controller.URL
		}

		if c.apiKey == "" {
			c.apiKey = config.Controller.APIKey
		}
	}

	return nil
}

// CloudPolicyClient returns Cloud Controller PolicyClient, connecting to cloud controller if not yet connected.
func (c *ControllerConn) CloudPolicyClient() (CloudPolicyClient, error) {
	// PolicyClient has no restrictions.
	return c.policyServiceClient()
}

func (c *ControllerConn) prepareCred() credentials.TransportCredentials {
	var cred credentials.TransportCredentials
	if c.allowInsecure {
		cred = insecure.NewCredentials()
	} else if c.skipVerify {
		cred = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // Requires enabling CLI option
		})
	}

	return cred
}

// client returns Cloud Controller Client, connecting to controller if not yet connected.
//
// This functions is not exposed to force callers to go through the check above.
func (c *ControllerConn) policyServiceClient() (cloudv1.PolicyServiceClient, error) {
	if c.conn != nil {
		return cloudv1.NewPolicyServiceClient(c.conn), nil
	}

	var addr string
	cred := c.prepareCred()
	addr = c.controllerAddr

	if cred == nil {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred = credentials.NewClientTLSFromCert(certPool, "")
	}

	var err error
	c.conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(cred), grpc.WithUnaryInterceptor(c.cloudControllerInterceptor))
	if err != nil {
		return nil, err
	}

	return cloudv1.NewPolicyServiceClient(c.conn), nil
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
