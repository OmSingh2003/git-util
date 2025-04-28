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
// These functions are defined here to be accessible by the root command and potentially
// other commands within the 'cmd' package (like status, sync).
// TODO: Ideally, move these helpers to a dedicated internal package (e.g., pkg/githelpers) later.

// runGitCommand executes a git command and returns its trimmed stdout output or an error
// including stderr content for better diagnostics.
func runGitCommand(args ...string) (string, error) {
	// Create the command object for 'git' with the provided arguments.
	cmd := exec.Command("git", args...)
	// Create buffers to capture standard output and standard error.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command.
	err := cmd.Run()

	// Capture stdout even if an error occurred, as it might contain partial info.
	output := strings.TrimSpace(stdout.String())

	// Check if the command execution resulted in an error.
	if err != nil {
		// Return the captured stdout along with a formatted error message
		// that includes the original error and the captured stderr.
		return output, fmt.Errorf("command 'git %s' failed: %w\nStderr: %s", strings.Join(args, " "), err, stderr.String())
	}

	// If no error, return the trimmed stdout and a nil error.
	return output, nil
}

// detectDefaultMainBranch tries to find 'main' or 'master' branch in the current repository.
func detectDefaultMainBranch() (string, error) {
	// Use 'git show-ref' to check if 'refs/heads/main' exists locally.
	_, errMain := runGitCommand("show-ref", "--verify", "--quiet", "refs/heads/main")
	// If the command succeeded (err is nil), 'main' exists.
	if errMain == nil {
		return "main", nil
	}

	// If 'main' wasn't found, check for 'master'.
	_, errMaster := runGitCommand("show-ref", "--verify", "--quiet", "refs/heads/master")
	// If the command succeeded (err is nil), 'master' exists.
	if errMaster == nil {
		return "master", nil
	}

	// If neither 'main' nor 'master' was found, return an error.
	return "", errors.New("neither 'main' nor 'master' branch found. Please specify with --main flag")
}

// findGitRepos walks the directory tree starting from rootDir and finds paths
// containing a .git subdirectory, indicating a Git repository root.
func findGitRepos(rootDir string) ([]string, error) {
	// Initialize a slice to store the paths of found repositories.
	var repos []string
	// Use filepath.WalkDir for efficient recursive directory traversal.
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		// Handle errors encountered while accessing files/directories.
		if err != nil {
			// Print a warning but continue the walk if possible.
			fmt.Fprintf(os.Stderr, "Warning: Error accessing path %q: %v\n", path, err)
			// Tell WalkDir to skip this directory if it's inaccessible.
			return filepath.SkipDir
		}

		// Check if the current item is a directory named exactly ".git".
		if d.IsDir() && d.Name() == ".git" {
			// The parent directory of ".git" is the root of the repository.
			repoPath := filepath.Dir(path)
			// Add the found repository path to the slice.
			repos = append(repos, repoPath)
			// Tell WalkDir to not descend into the ".git" directory itself.
			return filepath.SkipDir
		}

		// Optimization: Skip common large dependency/build directories.
		if d.IsDir() && (d.Name() == "vendor" || d.Name() == "node_modules" || d.Name() == "target" || d.Name() == "build") {
			// Tell WalkDir to skip descending into these directories.
			return filepath.SkipDir
		}

		// For all other files/directories, continue the walk normally.
		return nil
	})

	// Check for any overall error during the WalkDir process.
	if err != nil {
		return nil, err
	}
	// Return the slice of found repository paths.
	return repos, nil
}
