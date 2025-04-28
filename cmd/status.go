package cmd

import (
	// Imports needed by the RunE logic:
	"fmt"
	"os" // For os.Getwd, os.Stderr, etc.
	"path/filepath"
	"strconv" // To convert string counts to int
	"strings"

	// Import the new gitops package
	"github.com/OmSingh2003/git-util/pkg/gitops" // <-- Import Added

	// Cobra import
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
		// --- Determine Target Directory ---
		targetDir := statusDirectory
		if targetDir == "" {
			var err error
			targetDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
		}
		var err error
		targetDir, err = filepath.Abs(targetDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for target directory: %w", err)
		}

		fmt.Printf("Scanning directory: %s\n", targetDir)

		// --- Find Git Repositories ---
		// Call the helper function from the gitops package
		repos, err := gitops.FindGitRepos(targetDir) // <-- Updated Call
		if err != nil {
			return fmt.Errorf("error finding repositories: %w", err)
		}

		if len(repos) == 0 {
			fmt.Println("No Git repositories found in the specified directory.")
			return nil
		}

		fmt.Printf("\n--- Repository Status ---\n")

		// --- Calculate Max Path Length for Formatting ---
		maxLen := 0
		for _, repoPath := range repos {
			relPath, _ := filepath.Rel(targetDir, repoPath)
			if relPath == "." {
				relPath = filepath.Base(targetDir)
			}
			if len(relPath) > maxLen {
				maxLen = len(relPath)
			}
		}

		// --- Process Each Repository ---
		for _, repoPath := range repos {
			relPath, _ := filepath.Rel(targetDir, repoPath)
			if relPath == "." {
				relPath = filepath.Base(targetDir)
			}

			// --- Check Working Directory Status ---
			isDirty := false
			// Call the helper function from the gitops package
			statusOutput, err := gitops.RunGitCommand("-C", repoPath, "status", "--porcelain=v1") // <-- Updated Call
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to get status for %s: %v\n", relPath, err)
				isDirty = true
			} else if statusOutput != "" {
				isDirty = true
			}

			// --- Check Ahead/Behind Status ---
			ahead := 0
			behind := 0
			upstreamStatus := ""
			// Call the helper function from the gitops package
			revOutput, err := gitops.RunGitCommand("-C", repoPath, "rev-list", "--left-right", "--count", "HEAD...@{u}") // <-- Updated Call

			if err != nil {
				if strings.Contains(err.Error(), "no upstream configured") || strings.Contains(err.Error(), "unknown revision") {
					upstreamStatus = " [No Upstream]"
				} else {
					fmt.Fprintf(os.Stderr, "Warning: failed to get ahead/behind count for %s: %v\n", relPath, err)
					upstreamStatus = " [Error]"
				}
			} else {
				parts := strings.Split(strings.TrimSpace(revOutput), "\t")
				if len(parts) == 2 {
					ahead, _ = strconv.Atoi(parts[0])
					behind, _ = strconv.Atoi(parts[1])

					if ahead > 0 && behind > 0 {
						upstreamStatus = fmt.Sprintf(" [Ahead %d, Behind %d]", ahead, behind)
					} else if ahead > 0 {
						upstreamStatus = fmt.Sprintf(" [Ahead %d]", ahead)
					} else if behind > 0 {
						upstreamStatus = fmt.Sprintf(" [Behind %d]", behind)
					}
				} else {
					upstreamStatus = " [Error Parsing Revs]"
				}
			}

			// --- Combine and Print Status ---
			finalStatus := ""
			if isDirty {
				finalStatus = "Dirty" + upstreamStatus
			} else if upstreamStatus != "" && !strings.HasPrefix(upstreamStatus, " [Error") && upstreamStatus != " [No Upstream]" {
				finalStatus = "Clean" + upstreamStatus
			} else if upstreamStatus != "" {
				finalStatus = "Clean" + upstreamStatus
			} else {
				finalStatus = "Clean"
			}

			fmt.Printf("%-*s : %s\n", maxLen, relPath, finalStatus)

		} // End loop through repos

		return nil
	},
}

// --- init() function ---
func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVarP(&statusDirectory, "directory", "D", "", "Directory to scan for Git repositories (defaults to current directory)")
}

// --- Helper Functions ---
// *** DO NOT DEFINE findGitRepos or runGitCommand here! ***
// They should be defined only ONCE (e.g., in pkg/gitops/helpers.go)
