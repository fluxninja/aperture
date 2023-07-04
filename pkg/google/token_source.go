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

func provideGoogleTokenSource(in SourceIn) (oauth2.TokenSource, error) {
	log.Info().Msg("Initializing Google Token Source")
	var config tokenconfig.Config
	if err := in.Unmarshaller.UnmarshalKey(tokenSourceConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to unmarshal configuration")
		return nil, err
	}

	if !config.Enabled {
		return nil, nil
	}

	tokenSource, err := google.DefaultTokenSource(context.Background(), config.Scopes...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Google Token Source")
		return nil, err
	}

	return tokenSource, nil
}
