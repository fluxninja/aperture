package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd"
	"github.com/integralist/go-findroot/find"
	"github.com/spf13/cobra/doc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GenerateDocs generated the reference docs for the CLI.
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	root, err := find.Repo()
	if err != nil {
		log.Fatal(err)
	}
	docsDir := filepath.Join(root.Path, "docs/content/reference/aperturectl")

	// remove all subdirectories within the docsDir except docsDir itself
	err = filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && path != docsDir {
			err = os.RemoveAll(path)
			if err != nil {
				return nil
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = doc.GenMarkdownTreeCustom(cmd.RootCmd, docsDir, filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
	// walk the generated docs and move them to subdirectories based on their filename where _ separates the subdirectory
	err = filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".md") {
			// inside this file, replace instances of "sub-command" with "sub command"
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content = []byte(strings.ReplaceAll(string(content), "sub-command", "sub command"))
			content = []byte(strings.ReplaceAll(string(content), "via your OS's", "using your operating system's"))
			// write the file back
			err = os.WriteFile(path, content, 0o600)
			if err != nil {
				return err
			}

			subdir := transform(path)
			subdir = filepath.Join(docsDir, subdir)
			// create the subdirectory and move the file to it.
			err = os.MkdirAll(filepath.Dir(subdir), 0o755)
			if err != nil {
				return err
			}
			err = os.Rename(path, subdir)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func transform(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	// remove aperturectl prefix and replace _ with /
	link := strings.ReplaceAll(strings.TrimPrefix(filename, "aperturectl"), "_", "/")
	// remove the extension
	link = strings.TrimSuffix(link, ext)

	// tokenize filename by _ and pick up the last token
	tokens := strings.Split(filename, "_")
	filename = tokens[len(tokens)-1]

	link = filepath.Join(link, filename)
	return link
}

func filePrepender(filename string) string {
	name := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filepath.Base(filename)))
	// sidebar label is the last token in name
	tokens := strings.Split(name, "_")
	sidebarLabel := tokens[len(tokens)-1]
	sidebarLabel = cases.Title(language.English).String(sidebarLabel)

	return fmt.Sprintf(`---
sidebar_label: %s
hide_title: true
keywords:
- aperturectl
- %s
---

`, sidebarLabel, name)
}

func linkHandler(name string) string {
	return filepath.Join("/reference/aperturectl", transform(name))
}
