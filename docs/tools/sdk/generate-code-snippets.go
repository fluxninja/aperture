package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/integralist/go-findroot/find"
)

// Snippet represents a code snippet.
type Snippet struct {
	Language string
	Name     string
	Code     string
}

// main is the entry point of the program.
func main() {
	repoRoot, err := find.Repo()
	if err != nil {
		log.Fatal(err)
	}
	sdksRoot := filepath.Join(repoRoot.Path, "sdks")

	snippets := make(map[string]map[string]string)

	err = filepath.Walk(sdksRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".py", ".ts", ".java", ".go":
			fileSnippets, err := processFile(path)
			if err != nil {
				return err
			}
			for _, snippet := range fileSnippets {
				if _, ok := snippets[snippet.Language]; !ok {
					snippets[snippet.Language] = make(map[string]string)
				}
				snippets[snippet.Language][snippet.Name] = snippet.Code
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonBytes, err := json.MarshalIndent(snippets, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	codeSnippetsFile := filepath.Join(repoRoot.Path, "docs/content/code-snippets.json")
	err = os.WriteFile(codeSnippetsFile, jsonBytes, 0o600)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

// processFile processes a single file and extracts snippets.
func processFile(filePath string) ([]Snippet, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var startPattern, endPattern string
	switch filepath.Ext(filePath) {
	case ".py":
		startPattern, endPattern = `# START: (\w+)`, `# END`
	case ".ts", ".java":
		startPattern, endPattern = `// START: (\w+)`, `// END`
	case ".go":
		startPattern, endPattern = `// START: (\w+)`, `// END`
	}

	return extractSnippets(string(content), startPattern, endPattern, filepath.Ext(filePath)[1:]), nil
}

// extractSnippets extracts snippets from the given content based on start and end patterns.
func extractSnippets(content, startPattern, endPattern, lang string) []Snippet {
	var snippets []Snippet
	startRegexp := regexp.MustCompile(startPattern)
	endRegexp := regexp.MustCompile(endPattern)

	startIndexes := startRegexp.FindAllStringSubmatchIndex(content, -1)
	endIndexes := endRegexp.FindAllIndex([]byte(content), -1)

	for i, start := range startIndexes {
		if i < len(endIndexes) {
			snippetName := content[start[2]:start[3]]
			snippetCode := content[start[1]:endIndexes[i][0]]
			snippets = append(snippets, Snippet{
				Language: lang,
				Name:     snippetName,
				Code:     strings.TrimSpace(snippetCode),
			})
		}
	}

	return snippets
}
