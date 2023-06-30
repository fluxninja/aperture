package google

import (
	"context"

	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/fluxninja/aperture/v2/pkg/config"
	tokenconfig "github.com/fluxninja/aperture/v2/pkg/google/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// SourceIn stores fields and parameters provided to the Token Source constructor.
type SourceIn struct {
	fx.In
	Unmarshaller config.Unmarshaller
}

// Module provides a singleton pointer to oauth2.TokenSource via FX.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideGoogleTokenSource),
	)
}

const tokenSourceConfigKey = "token_source"

func provideGoogleTokenSource(in SourceIn) (*oauth2.TokenSource, error) {
	log.Info().Msg("Initializing Google Token Source")
	var config tokenconfig.Config
	if err := in.Unmarshaller.UnmarshalKey(tokenSourceConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to unmarshal configuration")
	}

	if !config.Enabled {
		log.Info().Msg("Google Token Source disabled. Will return nil.")
		return nil, nil
	}

	log.Info().Msg("Google Token Source enabled. Will create and return it.")

	tokenSource, err := google.DefaultTokenSource(context.Background(), config.Scopes...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Google Token Source")
		return nil, err
	}

	log.Info().Any("source", tokenSource).Msg("Created Token Source for GCP")

	return &tokenSource, nil
}
