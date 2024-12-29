package main

import (
	"flag"
	"fmt"
	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/api"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	apiPath = "repos/github/gitignore/contents"
)

type GitHubFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
}

func main() {
	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	listFlag := flag.Bool("list", false, "List available templates")
	flag.Parse()

	if *listFlag {
		listTemplates(*client)
		return
	}

	template := flag.Arg(0)
	if template == "" {
		fmt.Println("No template specified")
		os.Exit(1)
	}

	generateGitignore(*client, template)
}

func listTemplates(client api.RESTClient) {
	var files []GitHubFile
	err := client.Get(apiPath, &files)
	if err != nil {
		fmt.Printf("Error getting templates: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Available templates:")
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".gitignore") {
			name := strings.TrimSuffix(file.Name, ".gitignore")
			fmt.Printf("- %s\n", name)
		}
	}
}

func generateGitignore(client api.RESTClient, template string) {
	var files []GitHubFile
	err := client.Get(apiPath, &files)
	if err != nil {
		fmt.Printf("Error getting templates: %s\n", err)
		os.Exit(1)
	}

	templateName := fmt.Sprintf("%s.gitignore", template)
	templateName = strings.ToLower(templateName)
	var templateFile *GitHubFile
	for _, file := range files {
		if strings.EqualFold(strings.ToLower(file.Name), templateName) {
			templateFile = &file
			break
		}
	}

	if templateFile == nil {
		fmt.Printf("Template %s not found\n", template)
		os.Exit(1)
	}

	// Get raw content
	stdout, _, err := gh.Exec("api", templateFile.DownloadURL, "-H", "Accept: application/vnd.github.v3.raw")
	if err != nil {
		fmt.Printf("Error getting raw content: %s\n", err)
		os.Exit(1)
	}
	rawContent := stdout.String()

	// Write to file
	filePath := filepath.Join(".", ".gitignore")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s\n", err)
			os.Exit(1)
		}
	}(file)

	_, err = io.WriteString(file, rawContent)
	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Template %s written to %s\n", template, filePath)
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
