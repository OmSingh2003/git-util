// pkg/gitops/helpers.go
package gitops 

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// --- Helper Functions ---
// and usable by other packages (like 'cmd').

// RunGitCommand executes a git command and returns its trimmed stdout output or an error
// including stderr content for better diagnostics.
func RunGitCommand(args ...string) (string, error) { // function to run git commands take array of strings as input and return the output or error uses valadic operator
	cmd := exec.Command("git", args...) //uses exec commnad to make an object and store upack args
	var stdout, stderr bytes.Buffer // variables stdout and stderr of type bytes.buffer to store output and error if any 
	cmd.Stdout = &stdout //  address of stdout (which is a bytes.Buffer) to cmd.Stdout.
	cmd.Stderr = &stderr //  address of stderr (which is a bytes.Buffer) to cmd.Stderr.
	err := cmd.Run() // returns error to err if any
	output := strings.TrimSpace(stdout.String()) // triming whitespaces or new lines 
	if err != nil { // conditional statement to check if err is null or not
		return output, fmt.Errorf("command 'git %s' failed: %w\nStderr: %s", strings.Join(args, " "), err, stderr.String()) // print error
	}
	return output, nil // return output and error as nill because there is no error if compiler reached here.
}

// DetectDefaultMainBranch tries to find 'main' or 'master' branch in the current repository
// by calling the exported RunGitCommand function.
func DetectDefaultMainBranch() (string, error) { // return type of string and error string: main or master
	_, errMain := RunGitCommand("show-ref", "--verify", "--quiet", "refs/heads/main") // Checks for main branch
	if errMain == nil { // is this is nill means no error main exits : lets gooooo
		return "main", nil // main branch exits !
	}
	_, errMaster := RunGitCommand("show-ref", "--verify", "--quiet", "refs/heads/master") // Checks for master branch same as main branch
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
