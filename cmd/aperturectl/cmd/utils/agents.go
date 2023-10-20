package utils

import (
	"context"
	"fmt"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
)

func ListAgents(client IntrospectionClient, agentGroup string) error {
	agents, err := client.ListAgents(context.Background(), &cmdv1.ListAgentsRequest{
		AgentGroup: agentGroup,
	})
	if err != nil {
		return err
	}

	for _, agent := range agents.Agents {
		fmt.Println(agent)
	}

	return nil
}
