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

	lang := filepath.Ext(filePath)[1:]
	startPattern := `(?m)^\s*// START: (\w+)`
	endPattern := `(?m)^\s*// END: (\w+)`

	if lang == "py" {
		startPattern = `(?m)^\s*# START: (\w+)`
		endPattern = `(?m)^\s*# END: (\w+)`
	}

	return extractSnippets(string(content), startPattern, endPattern, lang), nil
}

func extractSnippets(content, startPattern, endPattern, lang string) []Snippet {
	var snippets []Snippet
	startRegexp := regexp.MustCompile(startPattern)
	endRegexp := regexp.MustCompile(endPattern)

	type snippetMarker struct {
		name  string
		index int
	}

	snippetStack := make([]snippetMarker, 0)
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if startMatches := startRegexp.FindStringSubmatch(line); startMatches != nil {
			// Push new snippet start onto the stack
			snippetName := startMatches[1]
			snippetStack = append(snippetStack, snippetMarker{index: i, name: snippetName})
		} else if endMatches := endRegexp.FindStringSubmatch(line); endMatches != nil {
			// Process snippet end
			endName := endMatches[1]
			if len(snippetStack) == 0 || snippetStack[len(snippetStack)-1].name != endName {
				fmt.Printf("Error: Unmatched end marker '%s' at line %d", endName, i+1)
				os.Exit(1)
			}

			startMarker := snippetStack[len(snippetStack)-1]
			snippetStack = snippetStack[:len(snippetStack)-1]

			// Extract snippet from startMarker.index to current line
			snippetCode := strings.Join(lines[startMarker.index+1:i], "\n")
			snippetCode = startRegexp.ReplaceAllString(snippetCode, "") // Remove start comments
			snippetCode = endRegexp.ReplaceAllString(snippetCode, "")   // Remove end comments

			snippets = append(snippets, Snippet{
				Language: lang,
				Name:     startMarker.name,
				Code:     strings.TrimSpace(snippetCode),
			})
		}
	}

	// Check if there are any unmatched Start markers left in the stack
	if len(snippetStack) > 0 {
		fmt.Printf("Error: Unmatched start marker %s\n", snippetStack[len(snippetStack)-1].name)
		os.Exit(1)
	}

	return snippets
}
