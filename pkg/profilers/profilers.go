// +kubebuilder:validation:Optional
package profilers

import (
	"context"
	httppprof "net/http/pprof"
	"os"
	"path"
	"runtime/pprof"

	"github.com/fluxninja/lumberjack"
	"github.com/gorilla/mux"
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

// defaultPath is the default path that is used to store profiles.
var defaultPath = path.Join(config.DefaultLogDirectory, "profiles")

// Module is a fx module that provides the profilers.
func Module() fx.Option {
	constructor := Constructor{ConfigKey: defaultKey}

	return fx.Options(
		fx.Invoke(constructor.setupProfilers),
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
// +kubebuilder:object:generate=true
type ProfilersConfig struct {
	// Register routes. Profile types profile, symbol and cmdline will be registered at /debug/pprof/{profile,symbol,cmdline}.
	RegisterHTTPRoutes bool `json:"register_http_routes" default:"true"`
	// Path to save performance profiles. "default" path is `/var/log/aperture/<service>/profiles`.
	ProfilesPath string `json:"profiles_path" default:"default"`
	// Flag to enable cpu profiling on process start and save it to a file. HTTP interface will not work if this is enabled as CPU profile will always be running.
	CPUProfile bool `json:"cpu_profiler" default:"false"`
}

// Constructor holds fields to create an instance of profilers.
type Constructor struct {
	ConfigKey     string
	DefaultConfig ProfilersConfig
}

func (constructor Constructor) setupProfilers(unmarshaller config.Unmarshaller,
	router *mux.Router,
	lc fx.Lifecycle,
) error {
	config := constructor.DefaultConfig

	if err := unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Error().Err(err)
		return err
	}
	if config.ProfilesPath == "default" {
		config.ProfilesPath = defaultPath
	}

	var cpuProfileFile *lumberjack.Logger
	var err error

	if config.CPUProfile {
		filename := path.Join(config.ProfilesPath, defaultCPUFile)
		log.Debug().Str("filename", filename).Msg("opening cpu profile writer")
		cpuProfileFile = newProfileWriter(filename)
	}

	if config.RegisterHTTPRoutes {
		router.HandleFunc(httpPathPrefix, httppprof.Index)
		router.HandleFunc(path.Join(httpPathPrefix, "cmdline"), httppprof.Cmdline)
		router.HandleFunc(path.Join(httpPathPrefix, "profile"), httppprof.Profile)
		router.HandleFunc(path.Join(httpPathPrefix, "symbol"), httppprof.Symbol)
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
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
		OnStop: func(context.Context) error {
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
