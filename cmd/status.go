package cmd

import (
	"fmt"
	"os"            // Needed for os.Getwd
	"path/filepath" // To walk directories

	// For string checking
	"github.com/spf13/cobra"
)

// Variable to store the directory flag value
var statusDirectory string

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of multiple Git repositories within a directory.",
	Long: `Scans a directory for Git repositories and reports their status,
including uncommitted changes, untracked files, and ahead/behind status
compared to the upstream branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the target directory
		targetDir := statusDirectory
		if targetDir == "" {
			// Default to current directory if flag not set
			var err error
			targetDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
		}
		fmt.Printf("Scanning directory: %s\n\n", targetDir)

		// --- TODO: Implement core logic ---
		// 1. Find all directories containing a '.git' folder within targetDir.
		// 2. For each found repository path:
		//    a. Run 'git -C <path> status --porcelain=v1'
		//    b. Run 'git -C <path> rev-list --left-right --count HEAD...@{u}' (handle upstream error)
		//    c. Parse outputs to determine status (Clean/Dirty, Ahead/Behind/Diverged)
		//    d. Print the status for the repo.

		fmt.Println("\n(Core status logic not yet implemented)") // Placeholder

		// Example of finding repositories (we will refine this)
		repos, err := findGitRepos(targetDir)
		if err != nil {
			return fmt.Errorf("error finding repositories: %w", err)
		}

		if len(repos) == 0 {
			fmt.Println("No Git repositories found in the specified directory.")
			return nil
		}

		fmt.Printf("Found %d Git repositories:\n", len(repos))
		for _, repo := range repos {
			// Just print the path for now
			relPath, _ := filepath.Rel(targetDir, repo) // Show relative path
			if relPath == "." {
				relPath = filepath.Base(targetDir) // Use directory name if it's the target dir itself
			}
			fmt.Printf(" - %s\n", relPath)
			// TODO: Add status check logic here inside the loop
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd) // Add statusCmd to the root command

	// Define flags for the status command
	// Use P for short flag 'd'
	statusCmd.Flags().StringVarP(&statusDirectory, "directory", "D", "", "Directory to scan for Git repositories (defaults to current directory)")
	// Note: Using -D to avoid conflict with root command's -d for delete. Choose another if you prefer.
}

// --- Helper Function (Example - can be moved/improved later) ---

// findGitRepos walks the directory tree and finds paths containing a .git subdirectory.
func findGitRepos(rootDir string) ([]string, error) {
	var repos []string
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Skip directories we can't read, but report other errors
			fmt.Fprintf(os.Stderr, "Warning: Error accessing path %q: %v\n", path, err)
			return filepath.SkipDir // Stop walking this directory branch
		}

		// Check if the entry is a directory named .git
		if d.IsDir() && d.Name() == ".git" {
			// The parent directory is the Git repository root
			repoPath := filepath.Dir(path)
			repos = append(repos, repoPath)
			// Skip walking further down into the .git directory itself
			return filepath.SkipDir
		}

		// Skip vendor/node_modules directories to speed things up (optional)
		if d.IsDir() && (d.Name() == "vendor" || d.Name() == "node_modules") {
			return filepath.SkipDir
		}

		// If we are checking the root directory itself, don't skip it if it doesn't contain .git immediately
		// Allow WalkDir to continue into subdirectories unless skipped above.
		return nil
	})

	if err != nil {
		return nil, err
	}
	return repos, nil
}
