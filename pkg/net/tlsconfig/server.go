// +kubebuilder:validation:Optional
package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

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
	// Server Cert file path
	CertFile string `json:"cert_file" validate:"omitempty,file"`
	// Server Key file path
	KeyFile string `json:"key_file" validate:"omitempty,file"`
	// Client CA file
	ClientCAFile string `json:"client_ca_file" validate:"omitempty,file"`
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
		fx.Annotate(
			constructor.provideServerTLSConfig,
			fx.ResultTags(name),
		),
	)
}

func (constructor Constructor) provideServerTLSConfig(unmarshaller config.Unmarshaller) (ServerTLSConfig, error) {
	config := constructor.DefaultConfig
	if err := unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize tls configuration!")
		return ServerTLSConfig{}, err
	}
	return config, nil
}

func (constructor Constructor) provideTLSConfig(unmarshaller config.Unmarshaller) (*tls.Config, error) {
	config, err := constructor.provideServerTLSConfig(unmarshaller)
	if err != nil {
		return nil, err
	}

	if config.Enabled {
		serverCertKeyPair, err := tls.LoadX509KeyPair(
			config.CertFile,
			config.KeyFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to load server tls cert/key")
			return nil, err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverCertKeyPair},
			MinVersion:   tls.VersionTLS12,
		}

		if config.AllowedCN != "" {
			tlsConfig.VerifyPeerCertificate = func(_ [][]byte, verifiedChains [][]*x509.Certificate) error {
				if len(verifiedChains) == 0 || len(verifiedChains[0]) == 0 {
					return errors.New("no verified chains")
				}
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
		if config.ClientCAFile != "" {

			caCert, err := os.ReadFile(config.ClientCAFile)
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
