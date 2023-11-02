package secretmanager

import (
	"context"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/utils"
	"go.uber.org/fx"
	"golang.org/x/oauth2/google"
)

// Module is a fx module that provides secret manager instance.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			provideSecretManagerClient,
		),
	)
}

// SecretManagerClientIn holds parameters for provideSecretManagerClient.
type SecretManagerClientIn struct {
	fx.In
	Shutdowner             fx.Shutdowner
	Lifecycle              fx.Lifecycle
	InstallationModeConfig *agentinfo.InstallationModeConfig `optional:"true"`
}

// SecretManagerClient is a wrapper around secretmanager.Client.
type SecretManagerClient struct {
	Client       *secretmanager.Client
	GCPProjectID string
	ProjectID    string
}

func provideSecretManagerClient(in SecretManagerClientIn) (*SecretManagerClient, error) {
	if in.InstallationModeConfig.InstallationMode != utils.InstallationModeCloudAgent {
		return nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		defer cancel()
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil || creds.ProjectID == "" {
		defer cancel()
		return nil, fmt.Errorf("failed to get default credentials: %v", err)
	}

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		defer cancel()
		return nil, fmt.Errorf("PROJECT_ID environment variable is not set")
	}

	secretManagerClient := &SecretManagerClient{
		Client:       client,
		GCPProjectID: creds.ProjectID,
		ProjectID:    projectID,
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			return client.Close()
		},
	})

	return secretManagerClient, nil
}
