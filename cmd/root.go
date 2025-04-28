package cmd

import (
	"fmt"
	"os"
	"strings" // For string manipulation (trimming, splitting)

	"github.com/OmSingh2003/git-util/pkg/gitops" // Importing the gitops package for helper functions
	"github.com/spf13/cobra"
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
		// Use the branch name provided via the --main flag.
		targetMainBranch := mainBranchName
		// If the flag was not set, try to auto-detect 'main' or 'master'.
		if targetMainBranch == "" {
			var err error
			targetMainBranch, err = gitops.DetectDefaultMainBranch()
			// Handle error if default branch couldn't be detected.
			if err != nil {
				return fmt.Errorf("could not detect default main branch: %w", err)
			}
			// fmt.Printf("Auto-detected main branch: %s\n", targetMainBranch) // Optional logging
		} else {
			// fmt.Printf("Using specified main branch: %s\n", targetMainBranch) // Optional logging
		}

		// --- Step 2: Run `git branch --merged <target>` ---
		// Execute the git command to list branches merged into the target branch.
		mergedBranchesOutput, err := runGitCommand("branch", "--merged", targetMainBranch)
		// Handle errors, specifically checking if the target branch exists.
		if err != nil {
			if strings.Contains(err.Error(), "warn: no such ref") || strings.Contains(err.Error(), "error: malformed object name") {
				return fmt.Errorf("specified main branch '%s' not found", targetMainBranch)
			}
			// Handle other errors from the git command.
			return fmt.Errorf("failed to list merged branches: %w", err)
		}

		// --- Step 3: Parse the output ---
		// Split the output into individual lines based on the newline character.
		lines := strings.Split(mergedBranchesOutput, "\n")

		// --- Step 4: Filter the branches ---
		// Create a slice to hold the names of branches identified for processing.
		var branchesToProcess []string
		// Iterate through each line of the output.
		for _, line := range lines {
			// Remove leading/trailing whitespace from the branch name.
			branchName := strings.TrimSpace(line)
			// Skip any empty lines.
			if branchName == "" {
				continue
			}
			// Skip the current branch (indicated by '* ').
			if strings.HasPrefix(branchName, "* ") {
				continue
			}
			// Skip the main branch itself, as we don't want to delete it.
			if branchName == targetMainBranch {
				continue
			}
			// If the branch passes checks, add it to the list.
			branchesToProcess = append(branchesToProcess, branchName)
		}

		// --- Step 5: Perform action ---
		// Check if any branches were found after filtering.
		if len(branchesToProcess) == 0 {
			fmt.Printf("No local branches found that are merged into %s (excluding the current branch).\n", targetMainBranch)
			return nil // Exit successfully, nothing to do.
		}

		// Check if the --delete flag was NOT provided (default action is listing).
		if !deleteBranches {
			// Print the list of branches that can be deleted.
			fmt.Printf("The following local branches are merged into %s and can potentially be deleted:\n", targetMainBranch)
			for _, branch := range branchesToProcess {
				fmt.Printf("  - %s\n", branch)
			}
			fmt.Println("\nRun with --delete flag (or -d) to remove them.")
		} else {
			// The --delete flag was provided.
			fmt.Printf("Processing deletion for branches merged into %s...\n", targetMainBranch)
			successCount := 0
			failCount := 0

			// Iterate through the branches identified for deletion.
			for _, branch := range branchesToProcess {
				// Check if the --dry-run flag was provided.
				if dryRun {
					// If dry run, just print what would happen.
					fmt.Printf("[Dry Run] Would attempt to delete branch: %s\n", branch)
					successCount++ // Count as success for dry run summary.
				} else {
					// Actual deletion attempt.
					fmt.Printf("Attempting to delete branch: %s...", branch)
					// Execute 'git branch -d <branch>' command. Using -d is safer than -D.
					_, err := gitops.RunGitCommand("branch", "-d", branch)
					// Check if the deletion command resulted in an error.
					if err != nil {
						// Print failure message including the error.
						fmt.Printf(" Failed (%v)\n", err)
						failCount++
					} else {
						// Print success message.
						fmt.Println(" Deleted.")
						successCount++
					}
				}
			}

			// Print a summary of the deletion process.
			fmt.Printf("\nSummary:\n")
			if dryRun {
				fmt.Printf("  Dry run complete. %d branches would have been targeted for deletion.\n", successCount)
			} else {
				fmt.Printf("  Successfully deleted: %d\n", successCount)
				fmt.Printf("  Failed to delete:   %d\n", failCount)
				// Add a hint about potential reasons for failure.
				if failCount > 0 {
					fmt.Println("  (Failures might occur if a branch has unmerged changes specific to it; use 'git branch -D' manually if needed.)")
				}
			}
		}

		return nil // Indicate successful execution of the root command's logic.
	},
}

// --- Standard Cobra Functions ---

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	// If the command execution returns an error, print it and exit with status 1.
	if err != nil {
		// Cobra automatically prints errors, but we exit non-zero.
		os.Exit(1)
	}
}

// init is run by Go automatically when the package is initialized.
// We define flags and configuration settings here.
func init() {
	// Define flags specific to the root command (branch cleaner).
	// Use StringVarP to define a string flag with a short version ('m').
	rootCmd.Flags().StringVarP(&mainBranchName, "main", "m", "", "Specify the main branch (e.g., main, master, develop)")
	// Use BoolVarP to define boolean flags with short versions ('d', 'n').
	rootCmd.Flags().BoolVarP(&deleteBranches, "delete", "d", false, "Actually delete the merged branches")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show what branches would be deleted without actually deleting")

	// NOTE: Subcommands like 'status' and 'sync' register themselves
	// with rootCmd using rootCmd.AddCommand(subCmd) in their respective init() functions.
}
