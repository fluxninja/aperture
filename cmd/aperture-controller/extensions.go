package main

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/info"
)

// WARNING: This is a placeholder file and should not be edited normally.

func init() {
	info.Extensions = GetExtensions()
}

func Module() fx.Option {
	return fx.Options()
}

func GetExtensions() []string {
	return []string{}
}
