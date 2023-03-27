package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

// GraphConfig represents the graph_config.yaml file.
type GraphConfig struct {
	PolicyFile string `json:"policy_file"`
	FileType   string `json:"file_type"`
	Depth      int    `json:"depth"`
}

func main() {
	gitRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	docsDir := filepath.Join(strings.TrimSpace(string(gitRoot)), "docs")
	contentDir := filepath.Join(docsDir, "content")

	err = filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Base(path) == "graph_config.yaml" {
			processGraphConfig(path, strings.TrimSpace(string(gitRoot)))
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func processGraphConfig(graphConfigFile, gitRoot string) {
	dir := filepath.Dir(graphConfigFile)

	data, err := os.ReadFile(graphConfigFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var graphConfigs []GraphConfig
	err = yaml.Unmarshal(data, &graphConfigs)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, gc := range graphConfigs {
		policyFilePath := filepath.Join(dir, gc.PolicyFile)
		mermaidFilePath := strings.TrimSuffix(policyFilePath, filepath.Ext(policyFilePath)) + ".mmd"

		if _, statErr := os.Stat(policyFilePath); os.IsNotExist(statErr) {
			fmt.Printf("Policy file not found: %s\n", policyFilePath)
			continue
		}

		fmt.Printf("Processing %s with depth: %d\n", policyFilePath, gc.Depth)

		aperturectl := gitRoot + "/cmd/aperturectl/main.go"
		var (
			cmd            *exec.Cmd
			stdout, stderr bytes.Buffer
		)
		switch gc.FileType {
		case "policy":
			// print the output of the command
			cmd = exec.Command("go", "run", aperturectl, "compile", "--policy", policyFilePath, "--mermaid", mermaidFilePath, "--depth", fmt.Sprintf("%d", gc.Depth))
		case "cr":
			cmd = exec.Command("go", "run", aperturectl, "compile", "--cr", policyFilePath, "--mermaid", mermaidFilePath, "--depth", fmt.Sprintf("%d", gc.Depth))
		default:
			fmt.Printf("Invalid file_type: %s\n", gc.FileType)
			continue
		}
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error while generating mermaid graph: %s\n", err)
			fmt.Println("Stderr output:", stderr.String())

		}
		fmt.Println("Stdout output:", stdout.String())
	}
}
