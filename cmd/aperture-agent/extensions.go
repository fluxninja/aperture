package main

// WARNING: This is a placeholder file and should not be edited normally.

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/info"
)

func init() {
	info.Extensions = GetExtensions()
}

func Module() fx.Option {
	return fx.Options()
}

func GetExtensions() []string {
	return []string{}
}
