package grpcgateway

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/FluxNinja/aperture/pkg/config"
	"github.com/FluxNinja/aperture/pkg/info"
	"github.com/FluxNinja/aperture/pkg/log"
	commonhttp "github.com/FluxNinja/aperture/pkg/net/http"
)

const (
	defaultKey = "server.grpc_gateway"
)

// Module is an fx module that sets up grpc-http Gateway.
func Module() fx.Option {
	return fx.Options(
		Constructor{}.Annotate(),
	)
}

// Constructor holds fields to configure grpc-http gateway.
type Constructor struct {
	// Name of http server instance -- empty for main server. Gateway name will be the same as HTTP server name
	Name string
	// Config key
	Key string
	// Default gateway config
	DefaultConfig GRPCGatewayConfig
}

// GRPCGatewayConfig holds configuration for grpc-http gateway
// swagger:model
type GRPCGatewayConfig struct {
	// GRPC server address to connect to - By default it points to HTTP server port because FluxNinja stack runs GRPC and HTTP servers on the same port
	GRPCAddr string `json:"grpc_server_address" validate:"hostname_port" default:"0.0.0.0:1"`
}

// GRPCGateway holds fields required for grpc-http gateway.
type GRPCGateway struct {
	// Tracks gw context
	ctx context.Context
	// Invoking cancel will disconnect all gw client connections
	cancel   context.CancelFunc
	grpcAddr string
	gwmux    *runtime.ServeMux
	mux      *mux.Router
	dopts    []grpc.DialOption
}

func (gw *GRPCGateway) registerGW(fn RegisterHandlerFunc) error {
	return fn(gw.ctx, gw.gwmux, gw.grpcAddr, gw.dopts)
}

// Annotate returns fx options that set up grpc-http gateway.
func (constructor Constructor) Annotate() fx.Option {
	if constructor.Key == "" {
		constructor.Key = defaultKey
	}
	var group, name string

	if constructor.Name == "" {
		group = config.GroupTag(info.Service)
	} else {
		group = config.GroupTag(constructor.Name)
	}
	name = config.NameTag(constructor.Name)

	return fx.Options(
		fx.Invoke(
			fx.Annotate(
				constructor.setupGRPCGateway,
				fx.ParamTags(group, name),
			),
		),
	)
}

// RegisterHandlerFunc is a handler function that registers grpc-http gateway handlers.
type RegisterHandlerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

// RegisterHandler holds fields to register grpc-http gateway handlers.
type RegisterHandler struct {
	// Handler func to invoke on start
	Handler RegisterHandlerFunc
	// gateway instance to register with -- empty = default gateway
	Name string
}

// Annotate annotates handler function in rh and returns it.
// These handlers can be registered with grpc-http gateway.
func (rh RegisterHandler) Annotate() fx.Option {
	var group string
	if rh.Name == "" {
		group = config.GroupTag(info.Service)
	} else {
		group = config.GroupTag(rh.Name)
	}
	return fx.Provide(
		fx.Annotate(
			func() RegisterHandlerFunc { return rh.Handler },
			fx.ResultTags(group),
		),
	)
}

func (constructor Constructor) setupGRPCGateway(
	handlers []RegisterHandlerFunc,
	httpServer *commonhttp.Server,
	lifecycle fx.Lifecycle,
	unmarshaller config.Unmarshaller,
	// Client to health service
	healthClient grpc_health_v1.HealthClient,
) error {
	config := constructor.DefaultConfig

	if err := unmarshaller.UnmarshalKey(constructor.Key, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize grpc gateway configuration!")
		return err
	}

	dopts := []grpc.DialOption{}

	var tlsConfig *tls.Config
	if httpServer.TLSConfig != nil {
		tlsConfig = httpServer.TLSConfig.Clone()
		// trust local server
		tlsConfig.InsecureSkipVerify = true
		transportCredentials := credentials.NewTLS(tlsConfig)
		dopts = append(dopts, grpc.WithTransportCredentials(transportCredentials))
	} else {
		// connect in an insecure manner
		dopts = append(dopts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	gw := &GRPCGateway{}
	gw.gwmux = runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithIncomingHeaderMatcher(httpHeaderMatcher),
		runtime.WithHealthzEndpoint(healthClient),
	)

	gw.dopts = dopts
	gw.mux = httpServer.Mux
	// Set RootHandler for HTTP server
	httpServer.RootHandler = gw.gwmux

	// Actual address will be figured out on start because it may depend on HTTP listener address -- which is not known until it starts
	gw.grpcAddr = config.GRPCAddr
	gw.ctx, gw.cancel = context.WithCancel(context.Background())

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if gw.grpcAddr == "0.0.0.0:1" {
				gw.grpcAddr = httpServer.Listener.GetListener().Addr().String()
			}

			// register routes -- this starts making grpc connections to the grpc server
			for _, handler := range handlers {
				err := gw.registerGW(handler)
				if err != nil {
					log.Error().Err(err).Msg("Failed to register grpc gateway route!")
					return err
				}
			}
			return nil
		},
		OnStop: func(context.Context) error {
			// Cancel the gw context
			gw.cancel()
			return nil
		},
	})

	return nil
}

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}

	return nil
}

func httpHeaderMatcher(key string) (string, bool) {
	switch key {
	case "Apikey":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
