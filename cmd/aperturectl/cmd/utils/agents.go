package utils

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

func ListAgents(client IntrospectionClient) error {
	agents, err := client.ListAgents(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	for _, agent := range agents.Agents {
		fmt.Println(agent)
	}

	return nil
}
