package utils

import (
	"context"
	"fmt"

	statusv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// ParseStatus parses the status.
func ParseStatus(client StatusClient) error {
	getStatusReq := &statusv1.GroupStatusRequest{
		Path: "",
	}
	statusResp, err := client.GetStatus(
		context.Background(),
		getStatusReq,
	)
	if err != nil {
		return err
	}

	result, err := status.ParseGroupStatus(make(map[string]string), "", statusResp)
	if err != nil {
		return err
	}
	for k, v := range result {
		fmt.Printf("%s: %s\n", k, v)
	}

	return nil
}
