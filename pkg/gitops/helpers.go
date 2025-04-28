// pkg/gitops/helpers.go
package gitops // Declare package name

import (
	// Imports needed by the helper functions
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// --- Helper Functions ---
// These functions are defined with Uppercase names to be exported (public)
// and usable by other packages (like 'cmd').

// RunGitCommand executes a git command and returns its trimmed stdout output or an error
// including stderr content for better diagnostics.
func RunGitCommand(args ...string) (string, error) { // Renamed to Uppercase
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	output := strings.TrimSpace(stdout.String())
	if err != nil {
		return output, fmt.Errorf("command 'git %s' failed: %w\nStderr: %s", strings.Join(args, " "), err, stderr.String())
	}
	return output, nil
}

// DetectDefaultMainBranch tries to find 'main' or 'master' branch in the current repository
// by calling the exported RunGitCommand function.
func DetectDefaultMainBranch() (string, error) { // Renamed to Uppercase
	// Use 'RunGitCommand' (now exported) to check refs.
	_, errMain := RunGitCommand("show-ref", "--verify", "--quiet", "refs/heads/main") // Call Uppercase version
	if errMain == nil {
		return "main", nil
	}
	_, errMaster := RunGitCommand("show-ref", "--verify", "--quiet", "refs/heads/master") // Call Uppercase version
	if errMaster == nil {
		return "master", nil
	}
	return "", errors.New("neither 'main' nor 'master' branch found. Please specify with --main flag")
}

// FindGitRepos walks the directory tree starting from rootDir and finds paths
// containing a .git subdirectory, indicating a Git repository root.
func FindGitRepos(rootDir string) ([]string, error) { // Renamed to Uppercase
	var repos []string
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Error accessing path %q: %v\n", path, err)
			return filepath.SkipDir
		}
		if d.IsDir() && d.Name() == ".git" {
			repoPath := filepath.Dir(path)
			repos = append(repos, repoPath)
			return filepath.SkipDir
		}
		if d.IsDir() && (d.Name() == "vendor" || d.Name() == "node_modules" || d.Name() == "target" || d.Name() == "build") {
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return repos, nil
}
