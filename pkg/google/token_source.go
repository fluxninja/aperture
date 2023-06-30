package google

import (
	"context"

	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// SourceIn stores fields and parameters provided to the Token Source constructor.
type SourceIn struct {
	fx.In
	Unmarshaller config.Unmarshaller
}

// Config stores configuration for the Google Token Source
type Config struct {
	Enabled bool     `json:"enabled" default:"false"`
	Scopes  []string `json:"scopes" default:"[]"`
}

// Module provides a singleton pointer to oauth2.TokenSource via FX.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideGoogleTokenSource),
	)
}

const tokenSourceConfigKey = "token_source"

func provideGoogleTokenSource(in SourceIn) (*oauth2.TokenSource, error) {
	var config Config
	if err := in.Unmarshaller.UnmarshalKey(tokenSourceConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to unmarshal configuration")
	}

	if !config.Enabled {
		return nil, nil
	}

	tokenSource, err := google.DefaultTokenSource(context.Background(), config.Scopes...)
	if err != nil {
		return nil, err
	}

	return &tokenSource, nil
}
