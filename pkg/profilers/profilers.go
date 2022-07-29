package profilers

import (
	"context"
	httppprof "net/http/pprof"
	"os"
	"path"
	"runtime/pprof"

	"github.com/fluxninja/lumberjack"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

const (
	defaultKey = "profilers"
	// HTTPPathPrefix is the path prefix to the profiler HTTP endpoints.
	httpPathPrefix = "/debug/pprof"
	// DefaultCPUFile is the default filename for the cpu profile.
	defaultCPUFile = "cpu.prof"
)

var (
	// DefaultPath is the default path that is used to store profiles.
	DefaultPath = path.Join(config.DefaultLogDirectory, "profiles")
	// DefaultPathFlag is the default path flag for the profiler.
	DefaultPathFlag = defaultKey + ".profiles_path"
)

func (constructor Constructor) setFlags(fs *pflag.FlagSet) error {
	fs.String(constructor.PathKey, DefaultPath, "path to performance profiles")
	return nil
}

func (constructor Constructor) provideFlagSetBuilder() config.FlagSetBuilderOut {
	return config.FlagSetBuilderOut{Builder: constructor.setFlags}
}

// Module is a fx module that provides the profilers.
func Module() fx.Option {
	constructor := Constructor{Key: defaultKey, PathKey: DefaultPathFlag}

	return fx.Options(
		fx.Invoke(constructor.setupProfilers),
		fx.Provide(constructor.provideFlagSetBuilder),
	)
}

// swagger:operation POST /profilers common-configuration Profilers
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/ProfilersConfig"

// ProfilersConfig holds configuration for profilers.
// swagger:model
type ProfilersConfig struct {
	// Path to save performance profiles. This can be set via command line arguments as well.
	ProfilesPath string `json:"profiles_path"`
	// Flag to enable cpu profiling
	CPUProfile bool `json:"cpu_profiler" default:"false"`
}

// Constructor holds fields to create an instance of profilers.
type Constructor struct {
	Key           string
	PathKey       string
	DefaultConfig ProfilersConfig
}

func (constructor Constructor) setupProfilers(unmarshaller config.Unmarshaller,
	router *mux.Router,
	lc fx.Lifecycle,
) error {
	profilesPath := cast.ToString(unmarshaller.Get(constructor.PathKey))

	config := constructor.DefaultConfig

	if err := unmarshaller.UnmarshalKey(constructor.Key, &config); err != nil {
		log.Error().Err(err)
		return err
	}

	var cpuProfileFile *lumberjack.Logger
	var err error
	registerRoutes := false

	if config.CPUProfile {
		filename := path.Join(profilesPath, defaultCPUFile)
		log.Debug().Str("filename", filename).Msg("opening cpu profile writer")
		cpuProfileFile = newProfileWriter(filename)
		registerRoutes = true
	}

	if registerRoutes {
		router.HandleFunc(httpPathPrefix, httppprof.Index)
		router.HandleFunc(path.Join(httpPathPrefix, "cmdline"), httppprof.Cmdline)
		router.HandleFunc(path.Join(httpPathPrefix, "profile"), httppprof.Profile)
		router.HandleFunc(path.Join(httpPathPrefix, "symbol"), httppprof.Symbol)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Interface("config", config).Msg("profiler config")

			if config.CPUProfile {
				log.Debug().Msg("starting cpu profile")
				err = pprof.StartCPUProfile(cpuProfileFile)
				if err != nil {
					log.Error().Err(err).Msg("Failed to start CPU profiling")
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if config.CPUProfile {
				log.Debug().Msg("stopping cpu profile")
				pprof.StopCPUProfile()
				closeProfileWriter(cpuProfileFile)
			}
			return nil
		},
	})

	return nil
}

func newProfileWriter(filename string) *lumberjack.Logger {
	writer := &lumberjack.Logger{
		Filename:   filename,
		MaxBackups: 10,
		MaxAge:     7,
	}
	return writer
}

func closeProfileWriter(lg *lumberjack.Logger) {
	filename := lg.Filename
	_ = lg.Rotate()
	_ = lg.Close()
	// remove the new file created by lumberjack after rotate is called
	_ = os.Remove(filename)
}
