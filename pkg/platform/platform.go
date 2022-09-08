package platform

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/spf13/pflag"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/wrapperspb"

	infov1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/info/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	fswatcher "github.com/fluxninja/aperture/pkg/filesystem/watcher"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	_ "github.com/fluxninja/aperture/pkg/net" // needed for docs
	"github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/pkg/net/health"
	"github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/plugins"
	"github.com/fluxninja/aperture/pkg/profilers"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/fluxninja/aperture/pkg/watchdog"
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(log.Printf))
}

// Config is the configuration for the platform.
type Config struct {
	// Additional config to be merged (used for unit tests etc)
	MergeConfig map[string]interface{}
	// Plugin symbols to look for
	PluginSymbols []string
}

var platform = initPlatform()

// Platform holds the state of the platform.
type Platform struct {
	statusRegistry status.Registry
	unmarshaller   config.Unmarshaller
	dotgraph       fx.DotGraph
}

func initPlatform() *Platform {
	return &Platform{}
}

// optionGroup enables to append more options to application.
type optionGroup []fx.Option

// New returns a new fx.App with the provided options.
func New(opts ...fx.Option) *fx.App {
	options := optionGroup(opts)
	defer func() {
		if v := recover(); v != nil {
			panichandler.Crash(v)
		}
	}()
	panichandler.RegisterPanicHandler(OnCrash)
	return fx.New(options...)
}

const (
	dotFileKey = "dot_file"
)

func setFlags(fs *pflag.FlagSet) error {
	fs.String(dotFileKey, "", "create fx dot file")
	return nil
}

func provideFlagSetBuilder() config.FlagSetBuilderOut {
	return config.FlagSetBuilderOut{Builder: setFlags}
}

// Module returns the platform module.
func (cfg Config) Module() fx.Option {
	// purge previous temp
	_ = os.RemoveAll(config.DefaultTempBase)
	// mkdir temp
	_ = os.MkdirAll(config.DefaultTempDirectory, os.ModePerm)

	// Create a temporary Fx App to load plugins
	var pluginOptions fx.Option
	var registry *plugins.PluginRegistry
	_ = fx.New(
		config.ModuleConfig{MergeConfig: cfg.MergeConfig, UnknownFlags: true}.Module(),
		plugins.ModuleConfig{PluginSymbols: cfg.PluginSymbols}.Module(),
		fx.Populate(&registry),
		fx.Populate(&pluginOptions),
		fx.NopLogger,
	)

	options := fx.Options(
		fx.Provide(provideFlagSetBuilder),
		config.ModuleConfig{MergeConfig: cfg.MergeConfig, UnknownFlags: false, ExitOnHelp: true}.Module(),
		config.LogModule(),
		health.HealthModule(),
		http.ProxyModule(),
		metrics.Module(),
		watchdog.Module(),
		fswatcher.Module(),
		profilers.Module(),
		ServerModule(false),
		etcdclient.Module(),
		jobs.Module(),
		status.Module(),
		fx.Populate(&platform.statusRegistry),
		platformStatusModule(),
		fx.Supply(registry),
		fx.Populate(&platform.unmarshaller),
		fx.Populate(&platform.dotgraph),
	)
	if pluginOptions != nil {
		options = fx.Options(
			options,
			pluginOptions,
		)
	}
	return options
}

// ServerModule returns the platform server module.
func ServerModule(testMode bool) fx.Option {
	var options fx.Option
	if testMode {
		options = fx.Options(
			fx.Provide(listener.ProvideTestListener),
		)
	} else {
		options = fx.Options(
			listener.Module(),
		)
	}
	options = fx.Options(options,
		tlsconfig.Module(),
		http.ServerModule(),
		grpc.GMuxServerModule(),
		grpcgateway.Module(),
		grpcgateway.RegisterHandler{Handler: infov1.RegisterInfoServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(grpc.RegisterInfoService),
	)
	return options
}

// Run is an fx helper function to gracefully start and stop an app container.
func Run(app *fx.App) {
	platform.statusRegistry = platform.statusRegistry.Child(platformStatusPath)
	// Check for dotflag
	if platform.unmarshaller != nil {
		dotfile := config.GetStringValue(platform.unmarshaller, dotFileKey, "")
		if dotfile != "" {
			bytes := []byte(platform.dotgraph)
			err := os.WriteFile(dotfile, bytes, 0o600)
			if err != nil {
				log.Error().Err(err).Str("file", dotfile).Msg("unable to write to file")
			}
		}
	}

	log.Info().Msg("Starting application")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to start application")
		return
	}

	defer stop(app)

	platform.statusRegistry.SetStatus(status.NewStatus(wrapperspb.String("platform running"), nil))

	// Wait for os.Signal
	<-app.Done()
	platform.statusRegistry.SetStatus(status.NewStatus(nil, errors.New("platform stopping")))
}

func stop(app *fx.App) {
	log.Info().Msg("Stopping application")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to stop application")
	}
	log.WaitFlush()
	// cleanup temp
	_ = os.RemoveAll(config.DefaultTempBase)
	platform.statusRegistry.Detach()
	os.Exit(0)
}

// OnCrash is the panic handler.
// TODO: Crash Report will be handled by Sentry plugin.
// Need to implement Panic Handler for the platform.
func OnCrash(interface{}, panichandler.Callstack) {}
