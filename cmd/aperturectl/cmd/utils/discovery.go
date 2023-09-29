package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
)

// ParseEntities parses the entities from the cluster.
func ParseEntities(client IntrospectionClient, findBy string) error {
	toPrint := ""

	if findBy != "" {
		find := strings.Split(findBy, "=")
		if len(find) != 2 {
			return fmt.Errorf("invalid findBy argument: %s", findBy)
		}

		switch find[0] {
		case "name":
			if find[1] == "" {
				return fmt.Errorf("must provide a name to find by: %s", findBy)
			}
			resp, err := client.ListDiscoveryEntity(context.Background(), &cmdv1.ListDiscoveryEntityRequest{By: &cmdv1.ListDiscoveryEntityRequest_Name{Name: find[1]}})
			if err != nil {
				return err
			}
			entity, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(resp.Entity)
			if err != nil {
				return err
			}
			toPrint = string(entity)
		case "ip":
			if find[1] == "" {
				return fmt.Errorf("invalid findBy argument: %s", findBy)
			}
			resp, err := client.ListDiscoveryEntity(context.Background(), &cmdv1.ListDiscoveryEntityRequest{By: &cmdv1.ListDiscoveryEntityRequest_IpAddress{IpAddress: find[1]}})
			if err != nil {
				return err
			}
			entity, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(resp.Entity)
			if err != nil {
				return err
			}
			toPrint = string(entity)
		default:
			return fmt.Errorf("invalid findBy argument: %s", findBy)
		}
	} else {
		resp, err := client.ListDiscoveryEntities(context.Background(), &cmdv1.ListDiscoveryEntitiesRequest{})
		if err != nil {
			return err
		}
		if resp.ErrorsCount != 0 {
			fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
		}
		entities, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(resp.Entities)
		if err != nil {
			return err
		}
		toPrint = string(entities)
	}

	fmt.Printf("%s\n", toPrint)

	return nil
}
