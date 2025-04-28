package cmd

import (
	// Imports needed by the RunE logic:
	"fmt"
	"os" // For os.Getwd, os.Stderr, etc.
	"path/filepath"
	"strconv" // To convert string counts to int
	"strings"

	// Cobra import
	"github.com/spf13/cobra"
	// NOTE: We are assuming findGitRepos and runGitCommand are defined
	// elsewhere in the 'cmd' package (e.g., in root.go) and are therefore accessible here.
	// No 'bytes' or 'os/exec' import needed here if runGitCommand is defined elsewhere.
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
		// --- Determine Target Directory ---
		// Use the directory provided by the flag.
		targetDir := statusDirectory
		// If the flag was not provided, default to the current working directory.
		if targetDir == "" {
			var err error
			targetDir, err = os.Getwd()
			// Handle error getting the current directory.
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
		}
		// Convert the target directory path to an absolute path for consistency.
		var err error
		targetDir, err = filepath.Abs(targetDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for target directory: %w", err)
		}

		fmt.Printf("Scanning directory: %s\n", targetDir)

		// --- Find Git Repositories ---
		// Call the helper function (defined elsewhere in cmd package, e.g., root.go)
		repos, err := findGitRepos(targetDir)
		if err != nil {
			return fmt.Errorf("error finding repositories: %w", err)
		}

		// If no repositories were found, inform the user and exit.
		if len(repos) == 0 {
			fmt.Println("No Git repositories found in the specified directory.")
			return nil
		}

		fmt.Printf("\n--- Repository Status ---\n")

		// --- Calculate Max Path Length for Formatting ---
		// Find the longest relative repository path for nice column alignment.
		maxLen := 0
		for _, repoPath := range repos {
			// Calculate path relative to the target directory.
			relPath, _ := filepath.Rel(targetDir, repoPath)
			// If the repo is the target directory itself, use its base name.
			if relPath == "." {
				relPath = filepath.Base(targetDir)
			}
			// Update maxLen if the current path is longer.
			if len(relPath) > maxLen {
				maxLen = len(relPath)
			}
		}

		// --- Process Each Repository ---
		// Loop through the list of found repository paths.
		for _, repoPath := range repos {
			// --- Get Relative Path for Display ---
			// Calculate the relative path again for display purposes.
			relPath, _ := filepath.Rel(targetDir, repoPath)
			if relPath == "." {
				relPath = filepath.Base(targetDir)
			}

			// --- Check Working Directory Status ---
			// isDirty tracks if the working directory has changes or untracked files.
			isDirty := false
			// Run 'git status --porcelain=v1' which gives script-friendly output.
			// Use '-C <path>' to run the command within the specific repository's directory.
			// Call the helper function (defined elsewhere in cmd package, e.g., root.go)
			statusOutput, err := runGitCommand("-C", repoPath, "status", "--porcelain=v1")
			// If the git status command itself failed, log a warning and assume the repo is dirty.
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to get status for %s: %v\n", relPath, err)
				isDirty = true // Treat as dirty if status command fails
			} else if statusOutput != "" {
				// If the command succeeded AND produced any output, it means there are changes.
				isDirty = true
			}

			// --- Check Ahead/Behind Status ---
			// Initialize ahead/behind counts and the status string.
			ahead := 0
			behind := 0
			upstreamStatus := "" // Stores strings like "[Ahead 2]", "[No Upstream]", etc.
			// Run 'git rev-list' to count commits between local HEAD and its upstream (@{u}).
			// The output format is "ahead_count<tab>behind_count".
			// Call the helper function (defined elsewhere in cmd package, e.g., root.go)
			revOutput, err := runGitCommand("-C", repoPath, "rev-list", "--left-right", "--count", "HEAD...@{u}")

			// Handle errors from the rev-list command.
			if err != nil {
				// Check specifically for the error indicating no upstream branch is configured.
				if strings.Contains(err.Error(), "no upstream configured") || strings.Contains(err.Error(), "unknown revision") {
					upstreamStatus = " [No Upstream]"
				} else {
					// Log other errors encountered while checking ahead/behind status.
					fmt.Fprintf(os.Stderr, "Warning: failed to get ahead/behind count for %s: %v\n", relPath, err)
					upstreamStatus = " [Error]"
				}
			} else {
				// If the command succeeded, parse the "ahead\tbehind" output.
				parts := strings.Split(strings.TrimSpace(revOutput), "\t")
				// Ensure we got exactly two parts (ahead and behind counts).
				if len(parts) == 2 {
					// Convert the string counts to integers. Ignore potential errors for simplicity.
					ahead, _ = strconv.Atoi(parts[0])
					behind, _ = strconv.Atoi(parts[1])

					// Format the upstreamStatus string based on the counts.
					if ahead > 0 && behind > 0 {
						upstreamStatus = fmt.Sprintf(" [Ahead %d, Behind %d]", ahead, behind) // Diverged
					} else if ahead > 0 {
						upstreamStatus = fmt.Sprintf(" [Ahead %d]", ahead)
					} else if behind > 0 {
						upstreamStatus = fmt.Sprintf(" [Behind %d]", behind)
					}
					// If ahead=0 and behind=0, upstreamStatus remains empty ("Synced").
				} else {
					// Handle unexpected output format from rev-list.
					upstreamStatus = " [Error Parsing Revs]"
				}
			}

			// --- Combine and Print Status ---
			// Determine the final status string to display.
			finalStatus := ""
			if isDirty {
				// If dirty, show "Dirty" plus any upstream info.
				finalStatus = "Dirty" + upstreamStatus
			} else if upstreamStatus != "" && !strings.HasPrefix(upstreamStatus, " [Error") && upstreamStatus != " [No Upstream]" {
				// If clean AND there's a valid ahead/behind status, show "Clean" plus that status.
				finalStatus = "Clean" + upstreamStatus
			} else if upstreamStatus != "" {
				// If clean but there was an error/no upstream, show "Clean" plus that specific status.
				finalStatus = "Clean" + upstreamStatus
			} else {
				// If clean and synced (upstreamStatus is empty), just show "Clean".
				finalStatus = "Clean"
			}

			// Print the formatted status line, aligning the path using maxLen.
			fmt.Printf("%-*s : %s\n", maxLen, relPath, finalStatus)

		} // End loop through repos

		return nil // Indicate successful execution of the command
	},
}

// --- init() function ---
// This function runs when the package is initialized.
func init() {
	// Register the statusCmd as a subcommand of the rootCmd (git-util).
	rootCmd.AddCommand(statusCmd)
	// Define the '--directory' flag for the status command.
	// -P adds the short flag '-D'.
	// The flag value will be stored in the 'statusDirectory' variable.
	statusCmd.Flags().StringVarP(&statusDirectory, "directory", "D", "", "Directory to scan for Git repositories (defaults to current directory)")
}

// --- Helper Functions ---
// *** DO NOT DEFINE findGitRepos or runGitCommand here! ***
// They should be defined only ONCE in the cmd package (e.g., in cmd/root.go)
