package utils

import (
	"errors"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/tui"
	"github.com/ghodss/yaml"
)

// GetBlueprintsTUIModel prepares the TUI model for the blueprints command.
func GetBlueprintsTUIModel(blueprintsDir string, selectAll bool) ([]string, *tui.CheckBoxModel, error) {
	blueprints, err := GetBlueprints(blueprintsDir)
	if err != nil {
		return nil, nil, err
	}

	if len(blueprints) == 0 {
		return nil, nil, errors.New("no blueprints found")
	}

	model := tui.InitialCheckboxModel(blueprints, "Which blueprints would you like to apply?")
	if !selectAll {
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			return nil, nil, err
		}
	} else {
		for i := range blueprints {
			model.Selected[i] = struct{}{}
		}
	}

	return blueprints, model, nil
}

// GetBlueprints returns a list of files that look like blueprints.
func GetBlueprints(blueprintsDir string) ([]string, error) {
	blueprints := []string{}

	return blueprints, filepath.Walk(blueprintsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) == ".yaml" || filepath.Ext(info.Name()) == ".yml" {
			// read the blueprint and look for the blueprint and uri fields
			bytes, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			var blueprint map[string]interface{}
			if err := yaml.Unmarshal(bytes, &blueprint); err != nil {
				return nil
			}
			if blueprint["blueprint"] != nil && blueprint["uri"] != nil {
				// add this filename to the list of blueprints
				blueprints = append(blueprints, path)
			}
		}
		return nil
	})
}
