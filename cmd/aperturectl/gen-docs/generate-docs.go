package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd"
	"github.com/integralist/go-findroot/find"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var position int

// GenerateDocs generated the reference docs for the CLI.
func main() {
	root, err := find.Repo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		os.Exit(1)
	}

	position = getTotalCommandsCount(cmd.RootCmd) + 2
	err = doc.GenMarkdownTreeCustom(cmd.RootCmd, filepath.Join(root.Path, "docs/content/get-started/aperture-cli"), filePrepender, linkHandler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		os.Exit(1)
	}
}

func getTotalCommandsCount(command *cobra.Command) int {
	count := len(command.Commands())
	for _, subCommand := range command.Commands() {
		if len(subCommand.Commands()) != 0 {
			count += getTotalCommandsCount(subCommand)
		}
	}
	return count
}

func filePrepender(filename string) string {
	name := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filepath.Base(filename)))
	sidebarPosition := position
	position--
	title := cases.Title(language.English).String(strings.ReplaceAll(name, "_", " "))
	return fmt.Sprintf(`---
title: %s
description: %s
keywords:
- aperturectl
- %s
sidebar_position: %d
---

`, title, title, name, sidebarPosition)
}

func linkHandler(name string) string {
	return name
}
