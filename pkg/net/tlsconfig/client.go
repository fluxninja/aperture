package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientTLSConfig is the configuration for client TLS.
// swagger:model
// +kubebuilder:object:generate=true
type ClientTLSConfig struct {
	CertFile           string `json:"cert_file"`
	KeyFile            string `json:"key_file"`
	CAFile             string `json:"ca_file"`
	KeyLogWriter       string `json:"key_log_file"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
}

// GetTLSConfig initializes tls.Config from config options.
func (c *ClientTLSConfig) GetTLSConfig() (*tls.Config, error) {
	var keyLogWriter io.Writer
	var err error
	if c.KeyLogWriter != "" {
		keyLogWriter, err = os.OpenFile(c.KeyLogWriter, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			return nil, err
		}
	}

	config := &tls.Config{
		InsecureSkipVerify: c.InsecureSkipVerify,
		KeyLogWriter:       keyLogWriter,
		MinVersion:         tls.VersionTLS12,
	}

	if c.CAFile != "" {
		var caCertPool *x509.CertPool
		caCert, err := os.ReadFile(c.CAFile)
		if err != nil {
			return nil, err
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		config.RootCAs = caCertPool
	}

	if c.CertFile != "" || c.KeyFile != "" {
		if c.CertFile == "" {
			return nil, errors.New("certificate file path missing")
		}
		if c.KeyFile == "" {
			return nil, errors.New("key file path missing")
		}
		clientCert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			return nil, err
		}
		config.Certificates = []tls.Certificate{clientCert}
	}

	return config, nil
}

// GetGRPCDialOptions creates GRPC DialOptions for TLS.
func (c *ClientTLSConfig) GetGRPCDialOptions(insecureEnabled bool) ([]grpc.DialOption, error) {
	if insecureEnabled {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}

	tlsConfig, err := c.GetTLSConfig()
	if err != nil {
		return nil, err
	}

	return []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))}, nil
}
