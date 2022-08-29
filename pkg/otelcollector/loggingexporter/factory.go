// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loggingexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

const (
	// The value of "type" key in configuration.
	typeStr                   = "aperturelogging"
	defaultSamplingInitial    = 2
	defaultSamplingThereafter = 500
)

// NewFactory creates a factory for Logging exporter.
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesExporter(createTracesExporter, component.StabilityLevelInDevelopment),
		component.WithMetricsExporter(createMetricsExporter, component.StabilityLevelInDevelopment),
		component.WithLogsExporter(createLogsExporter, component.StabilityLevelInDevelopment))
}

func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings:   config.NewExporterSettings(config.NewComponentID(typeStr)),
		SamplingInitial:    defaultSamplingInitial,
		SamplingThereafter: defaultSamplingThereafter,
	}
}

func createTracesExporter(ctx context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.TracesExporter, error) {
	return newTracesExporter(ctx, set, config)
}

func createMetricsExporter(ctx context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.MetricsExporter, error) {
	return newMetricsExporter(ctx, set, config)
}

func createLogsExporter(ctx context.Context, set component.ExporterCreateSettings, config config.Exporter) (component.LogsExporter, error) {
	return newLogsExporter(ctx, set, config)
}
