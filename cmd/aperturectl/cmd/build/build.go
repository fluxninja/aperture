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
)

var (
	buildConfigFile string
	outputDir       string
)

// BuildConfig is the configuration for building the binary.
type BuildConfig struct {
	BundledExtensions []string          `json:"bundled_extensions"`
	Extensions        []ExtensionConfig `json:"extensions"`
	Replaces          []ReplaceConfig   `json:"replaces"`
	LdFlags           []string          `json:"ldflags"`
	Flags             []string          `json:"flags"`
}

// ExtensionConfig is the configuration for an extension.
type ExtensionConfig struct {
	// GoModName. e.g. github.com/fluxninja/aperture-extensions/extension/test
	GoModName string `json:"go_mod_name" validate:"required"`
	// Import path e.g. github.com/fluxninja/aperture-extensions/extension/test v0.0.1.
	ImportPath string `json:"import_path" validate:"required"`
	// PkgName name of the extension. e.g. test
	PkgName string `json:"pkg_name"`
}

// ReplaceConfig is the configuration for a replace directive.
type ReplaceConfig struct {
	Old string `json:"old" validate:"required"`
	New string `json:"new" validate:"required"`
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

func PlatformModule() fx.Option {
  return fx.Options(
    {{- range .Extensions }}
    {{ .PkgName }}.PlatformModule(),
    {{- end }}
    {{- range .BundledExtensions }}
    {{ .Name }}.PlatformModule(),
    {{- end }}
  )
}

func {{ .ModuleName }}Module() fx.Option {
  return fx.Options(
    {{- range .Extensions }}
    {{ .PkgName }}.{{ $.ModuleName }}Module(),
    {{- end }}
    {{- range .BundledExtensions }}
    {{ .Name }}.{{ $.ModuleName }}Module(),
    {{- end }}
  )
}`

func buildRunE(cmd string) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		if buildConfigFile == "" {
			return fmt.Errorf("build config file is required")
		}
		// read config file and unmarshal into BuildConfig struct
		configBytes, err := os.ReadFile(buildConfigFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to read build config file")
			return err
		}
		var cfg BuildConfig
		if err = config.UnmarshalYAML(configBytes, &cfg); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal build config file")
			return err
		}

		// outputDir
		if outputDir == "" {
			return fmt.Errorf("output directory is required")
		}
		if err = os.MkdirAll(outputDir, 0o700); err != nil {
			log.Error().Err(err).Msg("failed to create output directory")
			return err
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
		apertureGoModName := apertureGoModFile.Module.Mod.Path

		// add the extensions to the final go.mod
		for _, ext := range cfg.Extensions {
			if ext.PkgName == "" {
				ext.PkgName = getGoPkgName(ext.GoModName)
			}
			if err != nil {
				log.Error().Err(err).Msg("failed to get go package name")
				return err
			}
			err = apertureGoModFile.AddModuleStmt(ext.ImportPath)
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

		// generate code using templates that calls the extension's module's
		// PlatformModule(), AgentModule(), ControllerModule() functions
		// keep the code in cmd/aperture-agent/extensions.go and cmd/aperture-controller/extensions.go
		// Note: local extensions are in ./extension and we can use replace directives to point to their local paths

		agentExtensionsFile := filepath.Join(builderDir, "cmd", "aperture-agent", "extensions.go")
		err = utils.BackupFile(agentExtensionsFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to backup agent extensions file")
			return err
		}
		defer utils.RestoreFile(agentExtensionsFile)
		err = generateExtensionsCode(cfg, "Agent", apertureGoModName, agentExtensionsFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to generate agent extensions code")
			return err
		}

		controllerExtensionsFile := filepath.Join(builderDir, "cmd", "aperture-controller", "extensions.go")
		err = utils.BackupFile(controllerExtensionsFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to backup controller extensions file")
			return err
		}
		defer utils.RestoreFile(controllerExtensionsFile)
		err = generateExtensionsCode(cfg, "Controller", apertureGoModName, controllerExtensionsFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to generate controller extensions code")
			return err
		}

		// execute go mod tidy
		goModTidyCmd := exec.Command("go", "mod", "tidy")
		goModTidyCmd.Stdout = os.Stdout
		goModTidyCmd.Stderr = os.Stderr
		err = goModTidyCmd.Run()
		if err != nil {
			log.Error().Err(err).Msg("failed to execute go mod tidy")
			return err
		}

		// build binaries
		err = buildBinary(cmd, cfg.LdFlags, cfg.Flags)
		if err != nil {
			log.Error().Err(err).Msg("failed to build agent binary")
			return err
		}

		return nil
	}
}

func buildBinary(service string, ldFlags, flags []string) error {
	ldFlagsFinal, err := getLdFlags(service, ldFlags)
	if err != nil {
		log.Error().Err(err).Msg("failed to get ldflags")
		return err
	}
	flagsFinal := strings.Join(flags, " ")

	cmdString := fmt.Sprintf("go build -o %s %s %s", service, ldFlagsFinal, flagsFinal)

	buildCmd := exec.Command("bash", "-c", cmdString)
	buildCmd.Dir = filepath.Join(builderDir, "cmd", service)
	// print the command and directory
	log.Info().Msgf("building binary: %s", buildCmd.String())
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

func getLdFlags(service string, ldFlags []string) (string, error) {
	// pick up version from env
	version := os.Getenv("APERTURE_VERSION")
	if version == "" {
		version = "0.0.1"
	}
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
	prefix := os.Getenv("APERTURE_PREFIX")
	if prefix == "" {
		prefix = "aperture"
	}
	gitBranch := os.Getenv("APERTURE_GIT_BRANCH")
	if gitBranch == "" {
		gitBranch = "undefined"
	}
	gitCommitHash := os.Getenv("APERTURE_GIT_COMMIT_HASH")
	if gitCommitHash == "" {
		gitCommitHash = "undefined"
	}
	// build time is 'date -Iseconds'
	buildTime := exec.Command("date", "-Iseconds")
	buildTimeOut, err := buildTime.Output()
	if err != nil {
		return "", err
	}
	// add all the ldflags
	ldFlagsFinal := "--ldflags \"-s -w -extldflags \"-Wl,--allow-multiple-definition\" "
	for _, flag := range ldFlags {
		ldFlagsFinal += fmt.Sprintf("%s ", flag)
	}
	ldFlagsFinal += fmt.Sprintf("-X 'github.com/fluxninja/aperture/pkg/info.Version=%s' ", version)
	ldFlagsFinal += fmt.Sprintf("-X 'github.com/fluxninja/aperture/pkg/info.BuildOS=%s/%s' ", strings.TrimSpace(string(goosOut)), strings.TrimSpace(string(goarchOut)))
	ldFlagsFinal += fmt.Sprintf("-X 'github.com/fluxninja/aperture/pkg/info.BuildHost=%s' ", strings.TrimSpace(string(hostnameOut)))
	ldFlagsFinal += fmt.Sprintf("-X 'github.com/fluxninja/aperture/pkg/info.BuildTime=%s' ", strings.TrimSpace(string(buildTimeOut)))
	ldFlagsFinal += fmt.Sprintf("-X 'github.com/fluxninja/aperture/pkg/info.GitBranch=%s' ", strings.TrimSpace(string(gitBranch)))
	ldFlagsFinal += fmt.Sprintf("-X 'github.com/fluxninja/aperture/pkg/info.GitCommitHash=%s' ", strings.TrimSpace(string(gitCommitHash)))
	ldFlagsFinal += fmt.Sprintf("-X 'github.com/fluxninja/aperture/pkg/info.Prefix=%s' ", prefix)
	ldFlagsFinal += fmt.Sprintf("-X 'github.com/fluxninja/aperture/pkg/info.Service=%s'", service)
	ldFlagsFinal += "\"" // close the ldflags
	return ldFlagsFinal, nil
}

func generateExtensionsCode(cfg BuildConfig, moduleName string, apertureModuleName string, dest string) error {
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
		ModuleName        string
		Extensions        []ExtensionConfig
		BundledExtensions []struct {
			Name      string
			GoModName string
		}
	}{
		ModuleName: moduleName,
		Extensions: cfg.Extensions,
	}
	// add bundled extensions which are local and their path is extensions/<name>
	for _, ext := range cfg.BundledExtensions {
		data.BundledExtensions = append(data.BundledExtensions, struct {
			Name      string
			GoModName string
		}{
			Name:      ext,
			GoModName: apertureModuleName + "/extensions/" + ext,
		})
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
