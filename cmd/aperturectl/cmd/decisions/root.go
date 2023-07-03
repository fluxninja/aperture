package decisions

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
)

var decisionTypes = []string{
	"load_scheduler",
	"rate_limiter",
	"quota_scheduler",
	"pod_scaler",
	"sampler",
}

var (
	etcdHost     string
	etcdPort     string
	etcdClient   *clientv3.Client
	all          bool
	decisionType string
)

func init() {
	DecisionsCmd.PersistentFlags().StringVar(&etcdHost, "etcd-host", "localhost", "Etcd host")
	DecisionsCmd.PersistentFlags().StringVar(&etcdPort, "etcd-port", "2379", "Etcd port")
	DecisionsCmd.PersistentFlags().BoolVar(&all, "all", false, "Get all decisions")
	DecisionsCmd.PersistentFlags().StringVar(&decisionType, "decision-type", "", fmt.Sprintf("Type of the decision to get (%s)", strings.Join(decisionTypes, ", ")))
}

// DecisionsCmd is the command to apply a policy to the cluster.
var DecisionsCmd = &cobra.Command{
	Use:           "decisions",
	Short:         "Get Aperture Decisions",
	Long:          `Use this command to get the Aperture Decisions.`,
	SilenceErrors: true,
	Example: `
	aperturectl decisions --etcd-host="127.0.0.1" --etcd-port="2379" --all
	aperturectl decisions --etcd-host="127.0.0.1" --etcd-port="2379" --decision-type="load_scheduler"`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// build etcd address
		etcdAddr := etcdHost + ":" + etcdPort

		// validate etcd address
		_, err := url.Parse(etcdAddr)
		if err != nil {
			return err
		}

		// make sure etcd is reachable on the address using etcd go client
		etcdClient, err = clientv3.New(clientv3.Config{
			Endpoints:   []string{etcdAddr},
			DialTimeout: 2 * time.Second,
		})
		if err != nil {
			return err
		}

		if !all {
			if decisionType == "" {
				return errors.New("decision type is required or use --all to get all decisions")
			} else {
				var found bool
				for _, v := range decisionTypes {
					if v == decisionType {
						found = true
						break
					}
				}
				if !found {
					return errors.New("invalid decision type, use one of the valid types (" + strings.Join(decisionTypes, ", ") + ") or use --all to get all decisions")
				}
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if all {
			decisionType = ""
		}

		decisions, err := getDecisions(etcdClient, decisionType)
		if err != nil {
			return err
		}

		for k, v := range decisions {
			cmd.Printf("%s:\n%s\n\n", k, v)
		}

		etcdClient.Close()
		return nil
	},
}

func getDecisions(etcdClient *clientv3.Client, decisionType string) (map[string]string, error) {
	decisionsPath := "aperture/decisions/"
	if decisionType != "" {
		decisionsPath += decisionType + "/"
	}

	decisions := make(map[string]string)
	decisionsResp, err := etcdClient.Get(context.Background(), decisionsPath, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range decisionsResp.Kvs {
		decisionName, ok := strings.CutPrefix(string(kv.Key), decisionsPath)
		if !ok {
			continue
		}

		if all {
			decisionType = strings.Split(decisionName, "/")[0]
		}

		var m protoreflect.ProtoMessage
		switch decisionType {
		case decisionTypes[0]:
			m = &policysyncv1.LoadDecisionWrapper{}
		case decisionTypes[1], decisionTypes[2]:
			m = &policysyncv1.RateLimiterDecisionWrapper{}
		case decisionTypes[3]:
			m = &policysyncv1.ScaleDecisionWrapper{}
		case decisionTypes[4]:
			m = &policysyncv1.SamplerDecisionWrapper{}
		}

		err := proto.Unmarshal(kv.Value, m)
		if err != nil {
			return nil, err
		}
		mjson := protojson.Format(m)
		decisions[decisionName] = mjson
	}

	return decisions, nil
}
