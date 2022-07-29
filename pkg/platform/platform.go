package platform

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/wrapperspb"

	infov1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/common/info/v1"
	"github.com/FluxNinja/aperture/pkg/config"
	etcdclient "github.com/FluxNinja/aperture/pkg/etcd/client"
	fswatcher "github.com/FluxNinja/aperture/pkg/filesystem/watcher"
	"github.com/FluxNinja/aperture/pkg/info"
	"github.com/FluxNinja/aperture/pkg/jobs"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/metrics"
	_ "github.com/FluxNinja/aperture/pkg/net" // needed for docs
	"github.com/FluxNinja/aperture/pkg/net/grpc"
	"github.com/FluxNinja/aperture/pkg/net/grpcgateway"
	"github.com/FluxNinja/aperture/pkg/net/health"
	"github.com/FluxNinja/aperture/pkg/net/http"
	"github.com/FluxNinja/aperture/pkg/net/listener"
	"github.com/FluxNinja/aperture/pkg/net/tlsconfig"
	"github.com/FluxNinja/aperture/pkg/panic"
	"github.com/FluxNinja/aperture/pkg/peers"
	"github.com/FluxNinja/aperture/pkg/plugins"
	"github.com/FluxNinja/aperture/pkg/profilers"
	"github.com/FluxNinja/aperture/pkg/status"
	"github.com/FluxNinja/aperture/pkg/watchdog"
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
	statusRegistry *status.Registry
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
	panic.RegisterPanicHandler(OnCrash)
	defer panic.Recover()
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
		peers.Constructor{}.Module(),
		profilers.Module(),
		ServerModule(false),
		etcdclient.Module(),
		jobs.Module(),
		status.Module(),
		fx.Invoke(grpc.RegisterStatusService),
		fx.Populate(&platform.statusRegistry),
		platformStatusModule(),
		plugins.ModuleConfig{OnlyCommandLineFlags: true}.Module(),
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to start application")
		return
	}

	defer stop(app)

	s := status.NewStatus(wrapperspb.String("platform running"), nil)
	err := platform.statusRegistry.Push(platformReadinessStatusName, s)
	if err != nil {
		log.Error().Err(err).Msg("Failed to push platform readiness status")
		return
	}

	// Wait for os.Signal
	<-app.Done()
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
	os.Exit(0)
}

var defaultDiagnosticDir = path.Join(config.DefaultAssetsDirectory, "diagnostic")

// OnCrash is the panic handler.
func OnCrash(e interface{}, s panic.Callstack) {
	log.Debug().Msg("Crash Reporter Registered")
	_ = os.MkdirAll(defaultDiagnosticDir, os.ModePerm)
	diagnosticDir := path.Join(defaultDiagnosticDir, time.Now().Format(time.RFC3339))

	// Crash Log writer
	fName := "/crash.log"
	crashlogger := panic.NewCrashFileWriter(filepath.Join(diagnosticDir, fName))
	crashLogWriter := panic.GetCrashWriter()
	crashLogWriter.Flush(crashlogger)
	panic.CloseCrashFileWriter(crashlogger)

	// Dump Status Registry
	groupStatus := platform.statusRegistry.Get("")
	if groupStatus != nil {
		gs, err := json.MarshalIndent(groupStatus, "", " ")
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal group status")
		}
		fName = "/status.json"
		_ = ioutil.WriteFile(filepath.Join(diagnosticDir, fName), gs, 0o600)
	} else {
		log.Info().Msg("No status information collected yet")
	}

	// Service version information
	versionInfo := info.GetVersionInfo()
	if versionInfo != nil {
		vInfo, err := json.MarshalIndent(versionInfo, "", " ")
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal version information")
		}
		fName = "/version-info.json"
		_ = ioutil.WriteFile(filepath.Join(diagnosticDir, fName), vInfo, 0o600)
	}
}
