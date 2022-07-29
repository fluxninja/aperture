package config

import (
	"flag"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/fx"

	"github.com/FluxNinja/aperture/pkg/info"
	"github.com/FluxNinja/aperture/pkg/log"
)

const (
	// ConfigPathFlag is the name of the flag for the configuration path.
	ConfigPathFlag = "config_path"
)

/*

	CommandLine

*/

// FlagSetBuilder is a function that helps users to build a Flagset.
type FlagSetBuilder func(*pflag.FlagSet) error

// FlagSetBuilderOut wraps the group of FlagSetBuilder and makes it handy to
// provide FlagSetBuilder via Fx.
type FlagSetBuilderOut struct {
	fx.Out
	Builder FlagSetBuilder `group:"FlagSetBuilders"`
}

// CommandLineIn holds parameters for NewCommandLine.
type CommandLineIn struct {
	fx.In
	// Builders help platform users set flags
	Builders []FlagSetBuilder `group:"FlagSetBuilders"`
	// pFlag error handling mode
	ErrorHandling pflag.ErrorHandling `optional:"true"`
}

// CommandLineOut holds the output of NewCommandLine, set of defined flags.
type CommandLineOut struct {
	fx.Out
	FlagSet *pflag.FlagSet
}

// CommandLineConfig holds configuration for CommandLine.
type CommandLineConfig struct {
	UnknownFlags bool
	ExitOnHelp   bool
}

// NewCommandLine returns a new CommandLineOut with new FlagSet.
func (config CommandLineConfig) NewCommandLine(cl CommandLineIn) (CommandLineOut, error) {
	fs := pflag.NewFlagSet(info.Service, cl.ErrorHandling)

	fs.ParseErrorsWhitelist.UnknownFlags = config.UnknownFlags

	fs.SetNormalizeFunc(wordSepNormalizeFunc)

	fs.String(ConfigPathFlag, DefaultConfigDirectory, "path to configuration file")

	// Add flags from flag.CommandLine
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	for _, builder := range cl.Builders {
		if err := builder(fs); err != nil {
			return CommandLineOut{}, err
		}
	}

	arguments := os.Args[1:]

	if err := fs.Parse(arguments); err != nil {
		if err == pflag.ErrHelp {
			if config.ExitOnHelp {
				// quietly exit
				os.Exit(0)
			}
		} else {
			log.Error().Err(err).Msg("Unable to parse command line!")
			return CommandLineOut{}, err
		}
	}

	return CommandLineOut{
		FlagSet: fs,
	}, nil
}

func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-"}
	to := "_"
	for _, sep := range from {
		name = strings.ReplaceAll(name, sep, to)
	}
	return pflag.NormalizedName(name)
}
