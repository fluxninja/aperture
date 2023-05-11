// +kubebuilder:validation:Optional
package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	guuid "github.com/google/uuid"
	"github.com/spf13/pflag"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/filesystem"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/panichandler"
)

var (
	// DefaultAssetsDirectory is path to default assets directory.
	DefaultAssetsDirectory = path.Join("/", "etc", info.Prefix, info.Service)
	// DefaultLogDirectory is path to default log directory.
	DefaultLogDirectory = path.Join("/", "var", "log", info.Prefix, info.Service)
	// DefaultTempBase is path to default temporary base.
	DefaultTempBase = path.Join(os.TempDir(), info.Prefix, info.Service)
	// DefaultTempDirectory is path to default temporary directory.
	DefaultTempDirectory = path.Join(DefaultTempBase, guuid.NewString())
	// DefaultConfigDirectory is path to default config directory.
	DefaultConfigDirectory = path.Join(DefaultAssetsDirectory, "config")
	// EnvPrefix is default environment prefix.
	EnvPrefix = strings.ReplaceAll(strings.ToUpper(info.Service), "-", "_") + "_"
)

const (
	// DefaultConfigFileExt is default config file extension.
	DefaultConfigFileExt = ".yaml"
	// DefaultKoanfDelim is default koanf delimiter.
	DefaultKoanfDelim = "."
)

// ConfigPath stores the path to config directory for other modules
// to look for their config files.
type ConfigPath struct {
	Path string
}

// ModuleConfig holds configuration for the config module.
type ModuleConfig struct {
	MergeConfig  map[string]interface{}
	UnknownFlags bool
	ExitOnHelp   bool
}

// Module is a fx Module that invokes CreateGlobalPanicHandlerRegistry and provides an annotated
// instance of CommandLine and FileUnmarshaller.
func (config ModuleConfig) Module() fx.Option {
	return fx.Options(
		fx.Invoke(
			panichandler.RegisterPanicHandlers,
		),
		fx.Provide(CommandLineConfig{UnknownFlags: config.UnknownFlags, ExitOnHelp: config.ExitOnHelp}.NewCommandLine),
		FileUnmarshallerConstructor{
			Name:        info.Service,
			PathFlag:    ConfigPathFlag,
			Path:        DefaultConfigDirectory,
			MergeConfig: config.MergeConfig,
			NoErrOnFile: true,
			EnableEnv:   true,
		}.Annotate(),
	)
}

// Annotate creates an annotated instance of FileUnmarshaller.
func (constructor FileUnmarshallerConstructor) Annotate() fx.Option {
	var returnName string
	if constructor.Name == info.Service {
		returnName = ``
	} else {
		returnName = NameTag(constructor.Name)
	}

	return fx.Options(
		fx.Provide(
			fx.Annotate(
				constructor.NewFileUnmarshaller,
				fx.ParamTags(`optional:"true"`),
				fx.ResultTags(returnName, returnName),
			),
		),
	)
}

// NameTag returns the name tag for the given name in the manner of fmt.Sprintf.
func NameTag(name string) string {
	return fmt.Sprintf("name:\"%s\"", name)
}

// GroupTag returns the group tag for the given group name in the manner of fmt.Sprintf.
func GroupTag(group string) string {
	return fmt.Sprintf("group:\"%s\"", group)
}

// FileUnmarshallerConstructor holds fields to create an annotated instance of FileUnmarshaller.
type FileUnmarshallerConstructor struct {
	// Optional Merge Config
	MergeConfig map[string]interface{}
	// Config Name -- config file name without the extension -- it is also the annotated name of koanf instance
	Name string
	// Command line flag for reading file path
	PathFlag string
	// If flag is empty or not provided on CL, fallback to explicit path
	Path string
	// Extension of file (exact) - empty = yaml
	FileExt string
	// Enable AutomaticEnv
	EnableEnv bool
	// NoErrOnFile
	NoErrOnFile bool
}

// NewFileUnmarshaller creates a new instance of FileUnmarshaller and ConfigPath that unmarshals the config file.
func (constructor FileUnmarshallerConstructor) NewFileUnmarshaller(flagSet *pflag.FlagSet) (Unmarshaller, ConfigPath, error) {
	var configPath string
	var err error

	if constructor.PathFlag != "" {
		configPath, err = flagSet.GetString(constructor.PathFlag)
		if err != nil {
			return nil, ConfigPath{}, err
		}
	}

	if configPath == "" {
		if constructor.Path == "" {
			configPath = DefaultConfigDirectory
		} else {
			configPath = constructor.Path
		}
	}

	if constructor.FileExt == "" {
		constructor.FileExt = DefaultConfigFileExt
	}
	fileInfo := filesystem.NewFileInfo(configPath, constructor.Name, constructor.FileExt)

	bytes, err := fileInfo.ReadAsByteBufferFromFile()
	if err != nil {
		if !constructor.NoErrOnFile {
			return nil, ConfigPath{}, err
		}
	}

	unmarshaller, err := KoanfUnmarshallerConstructor{
		FlagSet:     flagSet,
		EnableEnv:   constructor.EnableEnv,
		MergeConfig: constructor.MergeConfig,
	}.NewKoanfUnmarshaller(bytes)
	if err != nil {
		return nil, ConfigPath{}, err
	}

	return unmarshaller, ConfigPath{Path: configPath}, nil
}
