package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
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
	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// ControllerConfig is the config file structure for Aperture Cloud Controller.
type ControllerConfig struct {
	// When changing fields, remember to update docs/content/reference/configuration/aperturectl.md.
	URL         string `toml:"url"`
	APIKey      string `toml:"api_key"`
	ProjectName string `toml:"project_name"`
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
	projectName    string

	forwarderStopChan chan struct{}
	conn              *grpc.ClientConn
}

// InitFlags sets up flags for Cloud Controller.
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
		"Aperture Cloud User API Key to be used when using Cloud Controller",
	)
	flags.StringVar(
		&c.projectName,
		"project-name",
		"",
		"Aperture Cloud Project Name to be used when using Cloud Controller",
	)
	flags.StringVar(
		&c.config,
		"config",
		"",
		"Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG",
	)
}

// PreRunE verifies flags and should be run at PreRunE stage.
func (c *ControllerConn) PreRunE(_ *cobra.Command, _ []string) error {
	// Fetching config from environment variable
	if c.config == "" {
		c.config = os.Getenv(configEnv)
	}

	// Fetching config from default location
	if c.config == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			c.config = filepath.Join(homeDir, utils.AperturectlRootDir, "config")
			if _, err := os.Stat(c.config); err != nil {
				c.config = ""
			}
		}
	}

	if c.config == "" && (c.controllerAddr == "" || c.apiKey == "" || c.projectName == "") {
		return errors.New("missing required flag(s): --controller, --api-key, --project-name, --config")
	}

	if c.config != "" && (c.controllerAddr == "" || c.apiKey == "" || c.projectName == "") {
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

		if config.Controller.ProjectName == "" && c.projectName == "" {
			return fmt.Errorf("invalid config file '%s'. Missing key 'controller.project_name'", c.config)
		}

		if c.projectName == "" {
			c.projectName = config.Controller.ProjectName
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
func (c *ControllerConn) CloudPolicyClient() (utils.CloudPolicyClient, error) {
	if c.conn == nil {
		var err error
		c.conn, err = c.prepareGRPCClient()
		if err != nil {
			return nil, err
		}
	}

	return cloudv1.NewPolicyServiceClient(c.conn), nil
}

// CloudBlueprintsClient returns Cloud Controller BlueprintsClient, connecting to cloud controller if not yet connected.
func (c *ControllerConn) CloudBlueprintsClient() (utils.CloudBlueprintsClient, error) {
	if c.conn == nil {
		var err error
		c.conn, err = c.prepareGRPCClient()
		if err != nil {
			return nil, err
		}
	}

	return cloudv1.NewBlueprintsServiceClient(c.conn), nil
}

// IntrospectionClient returns Controller IntrospectionClient, connecting to controller if not yet connected.
func (c *ControllerConn) IntrospectionClient() (utils.IntrospectionClient, error) {
	return c.client()
}

// StatusClient returns Controller StatusClient, connecting to controller if not yet connected.
func (c *ControllerConn) StatusClient() (utils.StatusClient, error) {
	// StatusClient has no restrictions.
	return c.client()
}

// PolicyClient returns Controller PolicyClient, connecting to controller if not yet connected.
func (c *ControllerConn) PolicyClient() (utils.PolicyClient, error) {
	// PolicyClient has no restrictions.
	return c.client()
}

// client returns Controller Client, connecting to controller if not yet connected.
//
// This functions is not exposed to force callers to go through the check above.
func (c *ControllerConn) client() (cmdv1.ControllerClient, error) {
	if c.conn != nil {
		return cmdv1.NewControllerClient(c.conn), nil
	}

	var err error
	c.conn, err = c.prepareGRPCClient()
	if err != nil {
		return nil, err
	}

	return cmdv1.NewControllerClient(c.conn), nil
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

func (c *ControllerConn) prepareGRPCClient() (*grpc.ClientConn, error) {
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

	return grpc.Dial(addr, grpc.WithTransportCredentials(cred), grpc.WithUnaryInterceptor(c.cloudControllerInterceptor))
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
	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %s", c.apiKey), "projectName", c.projectName)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
