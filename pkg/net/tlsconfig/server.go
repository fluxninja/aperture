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
	constructor := Constructor{Key: defaultKey}
	return fx.Options(
		constructor.Annotate(),
	)
}

// ServerTLSConfig holds configuration for setting up server TLS support.
// swagger:model
// +kubebuilder:object:generate=true
type ServerTLSConfig struct {
	// Path to credentials. This can be set via command line arguments as well.
	//+kubebuilder:validation:Optional
	CertsPath string `json:"certs_path,omitempty"`
	// Server Cert file
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="ca.crt"
	ServerCert string `json:"server_cert" default:"ca.crt"`
	// Server Key file
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="ca.key"
	ServerKey string `json:"server_key" default:"ca.key"`
	// Client CA file
	//+kubebuilder:validation:Optional
	ClientCA string `json:"client_ca,omitempty" validate:"omitempty"`
	// Allowed CN
	//+kubebuilder:validation:Optional
	AllowedCN string `json:"allowed_cn,omitempty" validate:"omitempty,fqdn"`
	// Enabled TLS
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=false
	Enabled bool `json:"enabled" default:"false"`
}

// Constructor holds fields to create an annotated instance of *tls.Config.
type Constructor struct {
	Name          string
	Key           string
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
	if err := unmarshaller.UnmarshalKey(constructor.Key, &config); err != nil {
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
