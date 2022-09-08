// +kubebuilder:validation:Optional
package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"path"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

const (
	defaultKey = "server.tls"
)

// Module is a fx module that constructs annotated instance of *tls.Config.
func Module() fx.Option {
	constructor := Constructor{ConfigKey: defaultKey}
	return fx.Options(
		constructor.Annotate(),
	)
}

// ServerTLSConfig holds configuration for setting up server TLS support.
// swagger:model
// +kubebuilder:object:generate=true
type ServerTLSConfig struct {
	// Path to credentials. This can be set via command line arguments as well.
	CertsPath string `json:"certs_path"`
	// Server Cert file
	ServerCert string `json:"server_cert" default:"ca.crt"`
	// Server Key file
	ServerKey string `json:"server_key" default:"ca.key"`
	// Client CA file
	ClientCA string `json:"client_ca" validate:"omitempty"`
	// Allowed CN
	AllowedCN string `json:"allowed_cn" validate:"omitempty,fqdn"`
	// Enabled TLS
	Enabled bool `json:"enabled" default:"false"`
}

// Constructor holds fields to create an annotated instance of *tls.Config.
type Constructor struct {
	Name          string
	ConfigKey     string
	DefaultConfig ServerTLSConfig
}

// Annotate creates an annotated instance of *tls.Config.
func (constructor Constructor) Annotate() fx.Option {
	name := config.NameTag(constructor.Name)
	return fx.Provide(
		fx.Annotate(
			constructor.provideTLSConfig,
			fx.ResultTags(name),
		),
	)
}

func (constructor Constructor) provideTLSConfig(unmarshaller config.Unmarshaller) (*tls.Config, error) {
	config := constructor.DefaultConfig
	if err := unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize tls configuration!")
		return nil, err
	}

	if config.Enabled {
		certPath := config.CertsPath
		serverCertKeyPair, err := tls.LoadX509KeyPair(
			path.Join(certPath, config.ServerCert),
			path.Join(certPath, config.ServerKey))
		if err != nil {
			log.Error().Err(err).Msg("failed to load server tls cert/key")
			return nil, err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverCertKeyPair},
			MinVersion:   tls.VersionTLS12,
		}

		if config.AllowedCN != "" {
			tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				for _, chains := range verifiedChains {
					if len(chains) != 0 {
						if config.AllowedCN == chains[0].Subject.CommonName {
							return nil
						}
					}
				}
				return errors.New("CommonName authentication failed")
			}
		}

		var clientCertPool *x509.CertPool
		if config.ClientCA != "" {

			caCert, err := os.ReadFile(path.Join(certPath, config.ClientCA))
			if err != nil {
				log.Error().Err(err).Msg("failed to load client CA")
				return nil, err
			}
			clientCertPool = x509.NewCertPool()
			clientCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.ClientCAs = clientCertPool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}

		return tlsConfig, nil
	}

	return nil, nil
}
