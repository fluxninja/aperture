package status_test

// import (
// 	"context"
// 	"time"

// 	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
// 	"github.com/fluxninja/aperture/cmd/aperture-agent/agent"
// 	"github.com/fluxninja/aperture/pkg/config"
// 	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
// 	"github.com/fluxninja/aperture/pkg/platform"
// 	"github.com/fluxninja/aperture/pkg/status"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// 	"go.uber.org/fx"
// )

// var (
// 	app *fx.App
// 	svc statusv1.StatusServiceServer
// )

// var _ = BeforeEach(func() {
// 	registry := status.NewRegistry()

// 	app = platform.New(
// 		config.ModuleConfig{
// 			MergeConfig: map[string]interface{}{
// 				"sentrywriter": map[string]interface{}{
// 					"disabled": true,
// 				},
// 			},
// 		}.Module(),
// 		fx.Supply(registry),
// 		fx.Provide(agent.ProvidePeersPrefix),
// 		grpcclient.ClientConstructor{Name: "status-grpc-client", ConfigKey: "peer_discovery.client.grpc"}.Annotate(),
// 		fx.Populate(&svc),
// 	)
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	err := app.Start(ctx)
// 	Expect(err).NotTo(HaveOccurred())
// })

// var _ = AfterEach(func() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	err := app.Stop(ctx)
// 	Expect(err).NotTo(HaveOccurred())
// })

// var _ = Describe("Status GetGroupStatus", func() {
// 	When("client request with existing keys comes in", func() {
// 		It("should return the status of the registry", func() {

// 		})
// 	})

// 	When("client request with non existing key", func() {
// 		It("should return empty status of the registry", func() {

// 		})
// 	})
// })
