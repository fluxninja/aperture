package build

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	autils "github.com/fluxninja/aperture/pkg/utils"
)

const (
	exampleFmt = `# Build %s binary for Aperture

aperturectl --uri . build %s -c build-config.yaml -o /

Where build-config.yaml can be:
---
build:
  version: 1.0.0
  git_commit_hash: 1234567890
  git_branch: branch1
  ldflags:
    - -some-flag
    - -some-other-flag
  flags:
    - -some-flag
    - -some-other-flag
bundled_extensions: # remote extensions to be bundled
  - go_mod_name: github.com/org/name
    version: v1.0.0
    pkg_name: pkg
extensions: # built-in extensions to be enabled
  - fluxninja
  - sentry
replaces:
  - old: github.com/org/name
    new: github.com/org/name2
enable_core_extensions: false # default is true`
)

var (
	buildConfigFile   string
	outputDir         string
	apertureGoModName string
	cfg               Config
	serviceDir        string
)

// Config is the configuration for building the binary.
type Config struct {
	Build                BuildConfig       `json:"build"`
	BundledExtensions    []string          `json:"bundled_extensions"`
	Extensions           []ExtensionConfig `json:"extensions"`
	Replaces             []ReplaceConfig   `json:"replaces"`
	EnableCoreExtensions bool              `json:"enable_core_extensions" default:"true"`
}

// ExtensionConfig is the configuration for an extension.
type ExtensionConfig struct {
	// GoModName. e.g. github.com/fluxninja/aperture-extensions/extension/test
	GoModName string `json:"go_mod_name" validate:"required"`
	// Version e.g. v0.0.1
	Version string `json:"version" validate:"required"`
	// PkgName name of the extension. e.g. test
	PkgName string `json:"pkg_name"`
}

// ReplaceConfig is the configuration for a replace directive.
type ReplaceConfig struct {
	Old string `json:"old" validate:"required"`
	New string `json:"new" validate:"required"`
}

// BuildConfig picked up from environment variables.
type BuildConfig struct {
	Version       string   `json:"version" default:"unknown"`
	GitCommitHash string   `json:"git_commit_hash"`
	GitBranch     string   `json:"git_branch"`
	LdFlags       []string `json:"ldflags"`
	Flags         []string `json:"flags"`
}

const extensionsTpl = `package main

import (
	"go.uber.org/fx"
{{- range .Extensions }}
  "{{ .GoModName }}"
{{- end }}
{{- range .BundledExtensions }}
  "{{ .GoModName }}"
{{- end }}
)

func Module() fx.Option {
  return fx.Options(
    {{- range .Extensions }}
    {{ .PkgName }}.Module(),
    {{- end }}
    {{- range .BundledExtensions }}
    {{ .PkgName }}.Module(),
    {{- end }}
  )
}

func GetExtensions() []string {
  return []string{
	{{- range .Extensions }}
	"{{ .GoModName }}",
	{{- end }}
	{{- range .BundledExtensions }}
	"{{ .GoModName }}",
	{{- end }}
  }
}`

func buildRunE(service string) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		var configBytes []byte
		var err error

		serviceDir = filepath.Join(builderDir, "cmd", service)

		if buildConfigFile == "" {
			// look for this file under cmd/{service}/build-config.yaml
			buildConfigFile = filepath.Join(serviceDir, "build-config.yaml")
			if _, err = os.Stat(buildConfigFile); os.IsNotExist(err) {
				log.Info().Msgf("No build config file provided. Using default values.")
				buildConfigFile = ""
			} else {
				log.Info().Msgf("Using build config file: %s", buildConfigFile)
			}
		}

		if buildConfigFile != "" {
			// read config file and unmarshal into BuildConfig struct
			configBytes, err = os.ReadFile(buildConfigFile)
			if err != nil {
				log.Error().Err(err).Msg("failed to read build config file")
				return err
			}
		}

		unmarshaller, err := config.KoanfUnmarshallerConstructor{
			EnableEnv: true,
		}.NewKoanfUnmarshaller(configBytes)
		if err != nil {
			log.Error().Err(err).Msg("failed to create unmarshaller")
			return err
		}
		if err = unmarshaller.Unmarshal(&cfg); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal builder config")
			return err
		}

		if cfg.EnableCoreExtensions {
			// add fluxninja and sentry to the list of extensions if they are not already there
			if !autils.SliceContains(cfg.BundledExtensions, "fluxninja") {
				cfg.BundledExtensions = append(cfg.BundledExtensions, "fluxninja")
			}
			if !autils.SliceContains(cfg.BundledExtensions, "sentry") {
				cfg.BundledExtensions = append(cfg.BundledExtensions, "sentry")
			}
		}

		// outputDir
		if outputDir == "" {
			// use current directory
			outputDir, err = os.Getwd()
			if err != nil {
				return err
			}
			log.Info().Msgf("Using default output directory: %s", outputDir)
		} else {
			if err = os.MkdirAll(outputDir, 0o700); err != nil {
				log.Error().Err(err).Msg("failed to create output directory")
				return err
			}
		}

		// get lock file
		err = utils.WriterLock(builderURIRoot)
		if err != nil {
			log.Error().Err(err).Msg("failed to get lock file")
			return err
		}
		defer utils.Unlock(builderURIRoot)

		// for each extension module:
		// add replace directives to the final go.mod
		// generate code that calls the extension module's PlatformOptions(), AgentOptions(), ControllerOptions() functions
		// keep the code in cmd/aperture-agent/extensions.go and cmd/aperture-controller/extensions.go
		// Note: local extensions are in ./extension and we can use replace directives to point to their local paths
		goModPath := filepath.Join(builderDir, "go.mod")
		err = utils.BackupFile(goModPath)
		if err != nil {
			log.Error().Err(err).Msg("failed to backup go.mod file")
			return err
		}
		defer utils.RestoreFile(goModPath)

		goModBytes, err := os.ReadFile(goModPath)
		if err != nil {
			log.Error().Err(err).Msg("failed to read go.mod file")
			return err
		}

		// parse go.mod
		apertureGoModFile, err := modfile.Parse(goModPath, goModBytes, nil)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse go.mod file")
			return err
		}
		apertureGoModName = apertureGoModFile.Module.Mod.Path

		// add the extensions to the final go.mod
		for _, ext := range cfg.Extensions {
			if ext.PkgName == "" {
				ext.PkgName = getGoPkgName(ext.GoModName)
			}
			err = apertureGoModFile.AddRequire(ext.GoModName, ext.Version)
			if err != nil {
				log.Error().Err(err).Msg("failed to add require directive to go.mod file")
				return err
			}
			log.Info().Msgf("Adding extension %s", ext.GoModName)
		}
		// add the replace directive to the final go.mod
		for _, replace := range cfg.Replaces {
			err = apertureGoModFile.AddReplace(replace.Old, "", replace.New, "")
			if err != nil {
				log.Error().Err(err).Msg("failed to add replace directive to go.mod file")
				return err
			}
		}
		goModBytes, err = apertureGoModFile.Format()
		if err != nil {
			log.Error().Err(err).Msg("failed to format go.mod file")
			return err
		}
		err = os.WriteFile(goModPath, goModBytes, 0o600)
		if err != nil {
			log.Error().Err(err).Msg("failed to write go.mod file")
			return err
		}

		// build binaries
		err = buildBinary(service)
		if err != nil {
			log.Error().Err(err).Msgf("failed to build %s binary", service)
			return err
		}

		return nil
	}
}

func goModTidy() error {
	// execute go mod tidy
	goModTidyCmd := exec.Command("go", "mod", "tidy")
	goModTidyCmd.Dir = builderDir
	goModTidyCmd.Stdout = os.Stdout
	goModTidyCmd.Stderr = os.Stderr
	err := goModTidyCmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to execute go mod tidy")
		return err
	}
	return nil
}

func buildBinary(service string) error {
	defer func() {
		_ = goModTidy()
	}()

	extensionsFile := filepath.Join(serviceDir, "extensions.go")
	err := utils.BackupFile(extensionsFile)
	if err != nil {
		log.Error().Err(err).Msg("failed to backup agent extensions file")
		return err
	}
	defer utils.RestoreFile(extensionsFile)
	err = generateExtensionsCode(extensionsFile)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate agent extensions code")
		return err
	}

	err = goModTidy()
	if err != nil {
		return err
	}

	ldFlagsFinal, err := getLdFlags(service)
	if err != nil {
		log.Error().Err(err).Msg("failed to get ldflags")
		return err
	}

	flagsFinal := strings.Join(cfg.Build.Flags, " ")

	cmdString := fmt.Sprintf("go build -o %s %s %s", service, ldFlagsFinal, flagsFinal)

	buildCmd := exec.Command("bash", "-c", cmdString)
	buildCmd.Dir = serviceDir
	// print the command and directory
	log.Info().Msgf("Building binary: %s", buildCmd.String())
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	err = buildCmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to build binary")
		return err
	}

	// remove the existing binary if it exists
	err = os.Remove(filepath.Join(outputDir, service))
	if err != nil && !os.IsNotExist(err) {
		log.Error().Err(err).Msg("failed to remove existing binary")
		return err
	}
	err = os.Rename(filepath.Join(buildCmd.Dir, service), filepath.Join(outputDir, service))
	if err != nil {
		log.Error().Err(err).Msg("failed to move binary")
		return err
	}

	return nil
}

func getLdFlags(service string) (string, error) {
	// goos is 'go env GOOS'
	goos := exec.Command("go", "env", "GOOS")
	goosOut, err := goos.Output()
	if err != nil {
		return "", err
	}
	// goarch is 'go env GOARCH'
	goarch := exec.Command("go", "env", "GOARCH")
	goarchOut, err := goarch.Output()
	if err != nil {
		return "", err
	}
	// hostname is 'hostname'
	hostname := exec.Command("hostname")
	hostnameOut, err := hostname.Output()
	if err != nil {
		return "", err
	}
	// build time is 'date -Iseconds'
	buildTime := exec.Command("date", "-Iseconds")
	buildTimeOut, err := buildTime.Output()
	if err != nil {
		return "", err
	}

	// if builderDir is a git project, find the git commit hash
	if cfg.Build.GitCommitHash == "" {
		// git commit is 'git rev-parse HEAD'
		gitCommit := exec.Command("git", "rev-parse", "HEAD")
		gitCommit.Dir = builderDir
		gitCommitOut, err := gitCommit.Output()
		if err != nil {
			cfg.Build.GitCommitHash = "unknown"
		} else {
			cfg.Build.GitCommitHash = strings.TrimSpace(string(gitCommitOut))
		}
	}

	// if builderDir is a git project, find the git branch
	if cfg.Build.GitBranch == "" {
		// git branch is 'git rev-parse --abbrev-ref HEAD'
		gitBranch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		gitBranch.Dir = builderDir
		gitBranchOut, err := gitBranch.Output()
		if err != nil {
			cfg.Build.GitBranch = "unknown"
		} else {
			cfg.Build.GitBranch = strings.TrimSpace(string(gitBranchOut))
		}
	}

	// add all the ldflags
	// check if we are on Linux
	ldFlagsFinal := "--ldflags \"-s -w "
	if strings.TrimSpace(string(goosOut)) == "linux" {
		ldFlagsFinal += "-extldflags \"-Wl,--allow-multiple-definition\" "
	}

	for _, flag := range cfg.Build.LdFlags {
		ldFlagsFinal += fmt.Sprintf("%s ", flag)
	}
	ldFlagsFinal += fmt.Sprintf("-X '%s/pkg/info.Version=%s' ", apertureGoModName, cfg.Build.Version)
	ldFlagsFinal += fmt.Sprintf("-X '%s/pkg/info.BuildOS=%s/%s' ", apertureGoModName, strings.TrimSpace(string(goosOut)), strings.TrimSpace(string(goarchOut)))
	ldFlagsFinal += fmt.Sprintf("-X '%s/pkg/info.BuildHost=%s' ", apertureGoModName, strings.TrimSpace(string(hostnameOut)))
	ldFlagsFinal += fmt.Sprintf("-X '%s/pkg/info.BuildTime=%s' ", apertureGoModName, strings.TrimSpace(string(buildTimeOut)))
	ldFlagsFinal += fmt.Sprintf("-X '%s/pkg/info.GitBranch=%s' ", apertureGoModName, strings.TrimSpace(string(cfg.Build.GitBranch)))
	ldFlagsFinal += fmt.Sprintf("-X '%s/pkg/info.GitCommitHash=%s' ", apertureGoModName, strings.TrimSpace(string(cfg.Build.GitCommitHash)))
	ldFlagsFinal += fmt.Sprintf("-X '%s/pkg/info.Prefix=aperture' ", apertureGoModName)
	ldFlagsFinal += fmt.Sprintf("-X '%s/pkg/info.Service=%s'", apertureGoModName, service)
	ldFlagsFinal += "\"" // close the ldflags
	return ldFlagsFinal, nil
}

func generateExtensionsCode(dest string) error {
	// create the destination file
	f, err := os.Create(dest)
	if err != nil {
		log.Error().Err(err).Msg("failed to create extensions file")
		return err
	}
	defer f.Close()
	// execute the template
	t := template.Must(template.New("extensions").Parse(extensionsTpl))

	data := struct {
		Extensions        []ExtensionConfig
		BundledExtensions []struct {
			Name      string
			GoModName string
			PkgName   string
		}
	}{
		Extensions: cfg.Extensions,
	}
	// add bundled extensions which are local and their path is extensions/<name>
	for _, ext := range cfg.BundledExtensions {
		goModName := apertureGoModName + "/extensions/" + ext
		pkgName := getGoPkgName(goModName)
		data.BundledExtensions = append(data.BundledExtensions, struct {
			Name      string
			GoModName string
			PkgName   string
		}{
			Name:      ext,
			GoModName: goModName,
			PkgName:   pkgName,
		})
		log.Info().Msgf("Adding bundled extension: %s", ext)
	}

	err = t.Execute(f, data)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute extensions template")
		return err
	}
	return nil
}

func getGoPkgName(goModName string) string {
	goModNameParts := strings.Split(goModName, "/")
	// if the last part is v{version}, use the second last part
	if strings.HasPrefix(goModNameParts[len(goModNameParts)-1], "v") {
		return goModNameParts[len(goModNameParts)-2]
	} else {
		return goModNameParts[len(goModNameParts)-1]
	}
}
