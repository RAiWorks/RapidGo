package cli

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const starterRepo = "https://github.com/RAiWorks/RapidGo-starter/archive/refs/heads/main.zip"

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new RapidGo project from the starter template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Validate project name
		if strings.ContainsAny(name, `/\:*?"<>|`) {
			return fmt.Errorf("invalid project name: %s", name)
		}

		// Check target directory doesn't exist
		if _, err := os.Stat(name); err == nil {
			return fmt.Errorf("directory %q already exists", name)
		}

		fmt.Printf("Creating new RapidGo project: %s\n", name)

		// 1. Download starter template
		fmt.Println("  Downloading starter template...")
		zipPath, err := downloadStarter()
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		defer os.Remove(zipPath)

		// 2. Extract to project directory
		fmt.Println("  Extracting template...")
		if err := extractZip(zipPath, name); err != nil {
			return fmt.Errorf("extract failed: %w", err)
		}

		// 3. Replace module name
		fmt.Println("  Configuring module name...")
		if err := replaceModuleName(name); err != nil {
			return fmt.Errorf("module rename failed: %w", err)
		}

		// 4. Run go mod tidy
		fmt.Println("  Running go mod tidy...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = name
		tidyCmd.Stdout = os.Stdout
		tidyCmd.Stderr = os.Stderr
		if err := tidyCmd.Run(); err != nil {
			return fmt.Errorf("go mod tidy failed: %w", err)
		}

		fmt.Println("=================================")
		fmt.Printf("  Project %s created!\n", name)
		fmt.Println("=================================")
		fmt.Printf("  cd %s\n", name)
		fmt.Println("  cp .env.example .env")
		fmt.Println("  go run cmd/main.go serve")
		fmt.Println("=================================")

		return nil
	},
}

// downloadStarter downloads the starter template zip from GitHub.
func downloadStarter() (string, error) {
	resp, err := http.Get(starterRepo) //nolint:gosec // Fixed trusted URL, not user input
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from GitHub", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "rapidgo-starter-*.zip")
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}
	tmpFile.Close()
	return tmpFile.Name(), nil
}

// extractZip extracts a GitHub archive zip to the target directory.
// GitHub archives contain a top-level directory (e.g., "RapidGo-starter-main/")
// which is stripped during extraction.
func extractZip(zipPath, targetDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Find the common prefix (GitHub adds "RepoName-branch/")
	prefix := ""
	if len(r.File) > 0 {
		prefix = strings.SplitN(r.File[0].Name, "/", 2)[0] + "/"
	}

	for _, f := range r.File {
		// Strip the GitHub prefix directory
		relPath := strings.TrimPrefix(f.Name, prefix)
		if relPath == "" {
			continue
		}

		destPath := filepath.Join(targetDir, relPath)

		// Zip slip protection: ensure path doesn't escape target directory
		cleanDest := filepath.Clean(destPath)
		cleanTarget := filepath.Clean(targetDir) + string(os.PathSeparator)
		if !strings.HasPrefix(cleanDest, cleanTarget) && cleanDest != filepath.Clean(targetDir) {
			return fmt.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// replaceModuleName replaces "github.com/RAiWorks/RapidGo-starter" with
// the project name in go.mod and all .go files.
func replaceModuleName(projectDir string) error {
	oldModule := "github.com/RAiWorks/RapidGo-starter"
	newModule := projectDir

	return filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		name := filepath.Base(path)

		// Only process .go files and go.mod
		if ext != ".go" && name != "go.mod" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newContent := strings.ReplaceAll(string(content), oldModule, newModule)
		if newContent != string(content) {
			return os.WriteFile(path, []byte(newContent), info.Mode())
		}

		return nil
	})
}
