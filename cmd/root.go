package cmd

import (
	"fmt"
	"os"
	"strings" // Still needed for string manipulation in RunE

	"github.com/OmSingh2003/git-util/pkg/gitops" // Importing the gitops package
	"github.com/spf13/cobra"
	// Removed bytes, errors, os/exec, path/filepath as they are likely only used by helpers in gitops package now
)

// Variables to hold the flag values for the root command (branch cleaner)
var (
	mainBranchName string
	deleteBranches bool
	dryRun         bool
)

// rootCmd represents the base command when called without any subcommands
// Its primary function currently is the 'branch cleaner' logic.
var rootCmd = &cobra.Command{
	Use:   "git-util",
	Short: "A utility tool for common Git operations.",
	Long: `git-util helps automate and simplify various Git tasks.
The first feature implemented is cleaning up merged local branches.
More features might be added later via subcommands (e.g., status, sync).`,
	// RunE executes the logic for the root command (branch cleaner)
	RunE: func(cmd *cobra.Command, args []string) error {
		// --- Step 1: Determine the target main branch ---
		targetMainBranch := mainBranchName
		if targetMainBranch == "" {
			var err error
			// Call helper from gitops package
			targetMainBranch, err = gitops.DetectDefaultMainBranch()
			if err != nil {
				return fmt.Errorf("could not detect default main branch: %w", err)
			}
		}

		// --- Step 2: Run `git branch --merged <target>` ---
		// Call helper from gitops package
		mergedBranchesOutput, err := gitops.RunGitCommand("branch", "--merged", targetMainBranch) // <-- Updated Call
		if err != nil {
			if strings.Contains(err.Error(), "warn: no such ref") || strings.Contains(err.Error(), "error: malformed object name") {
				return fmt.Errorf("specified main branch '%s' not found", targetMainBranch)
			}
			return fmt.Errorf("failed to list merged branches: %w", err)
		}

		// --- Step 3: Parse the output ---
		lines := strings.Split(mergedBranchesOutput, "\n")

		// --- Step 4: Filter the branches ---
		var branchesToProcess []string
		for _, line := range lines {
			branchName := strings.TrimSpace(line)
			if branchName == "" {
				continue
			}
			if strings.HasPrefix(branchName, "* ") {
				continue
			}
			if branchName == targetMainBranch {
				continue
			}
			branchesToProcess = append(branchesToProcess, branchName)
		}

		// --- Step 5: Perform action ---
		if len(branchesToProcess) == 0 {
			fmt.Printf("No local branches found that are merged into %s (excluding the current branch).\n", targetMainBranch)
			return nil
		}

		if !deleteBranches {
			fmt.Printf("The following local branches are merged into %s and can potentially be deleted:\n", targetMainBranch)
			for _, branch := range branchesToProcess {
				fmt.Printf("  - %s\n", branch)
			}
			fmt.Println("\nRun with --delete flag (or -d) to remove them.")
		} else {
			fmt.Printf("Processing deletion for branches merged into %s...\n", targetMainBranch)
			successCount := 0
			failCount := 0
			for _, branch := range branchesToProcess {
				if dryRun {
					fmt.Printf("[Dry Run] Would attempt to delete branch: %s\n", branch)
					successCount++
				} else {
					fmt.Printf("Attempting to delete branch: %s...", branch)
					// Call helper from gitops package
					// We capture output in case the error message needs it, even if we don't print it on success.
					_, err := gitops.RunGitCommand("branch", "-d", branch) // <-- Updated Call
					if err != nil {
						fmt.Printf(" Failed (%v)\n", err) // Error from RunGitCommand includes stderr
						failCount++
					} else {
						fmt.Println(" Deleted.")
						successCount++
					}
				}
			}
			fmt.Printf("\nSummary:\n")
			if dryRun {
				fmt.Printf("  Dry run complete. %d branches would have been targeted for deletion.\n", successCount)
			} else {
				fmt.Printf("  Successfully deleted: %d\n", successCount)
				fmt.Printf("  Failed to delete:   %d\n", failCount)
				if failCount > 0 {
					fmt.Println("  (Failures might occur if a branch has unmerged changes specific to it; use 'git branch -D' manually if needed.)")
				}
			}
		}
		return nil
	},
}

// --- Standard Cobra Functions ---

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init is run by Go automatically when the package is initialized.
func init() {
	// Define flags specific to the root command (branch cleaner).
	rootCmd.Flags().StringVarP(&mainBranchName, "main", "m", "", "Specify the main branch (e.g., main, master, develop)")
	rootCmd.Flags().BoolVarP(&deleteBranches, "delete", "d", false, "Actually delete the merged branches")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show what branches would be deleted without actually deleting")
}

// --- Helper Function Definitions ---
// *** Helper function definitions (runGitCommand, detectDefaultMainBranch, findGitRepos) should be DELETED from this file ***
// They now live in pkg/gitops/helpers.go
