package cmd

import (
	"bytes"  // To capture command output
	"errors" // For creating custom errors
	"fmt"
	"os"
	"os/exec" // To run external commands (like git)
	"strings" // For string manipulation (trimming, splitting)

	"github.com/spf13/cobra"
)

// Variables to hold the flag values
var (
	mainBranchName string
	deleteBranches bool
	dryRun         bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-util",
	Short: "A utility tool for common Git operations.",
	Long: `git-util helps automate and simplify various Git tasks.
The first feature implemented is cleaning up merged local branches.
More features might be added later.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Executing Git Branch Cleaner...")

		// --- Step 1: Determine the target main branch ---
		targetMainBranch := mainBranchName // Use flag value if provided
		if targetMainBranch == "" {
			var err error
			targetMainBranch, err = detectDefaultMainBranch()
			if err != nil {
				return fmt.Errorf("could not detect default main branch: %w", err)
			}
			fmt.Printf("Auto-detected main branch: %s\n", targetMainBranch)
		} else {
			fmt.Printf("Using specified main branch: %s\n", targetMainBranch)
		}

		// --- Step 2: Run `git branch --merged <target>` ---
		fmt.Printf("Finding local branches merged into %s...\n", targetMainBranch)
		mergedBranchesOutput, err := runGitCommand("branch", "--merged", targetMainBranch)
		if err != nil {
			if strings.Contains(err.Error(), "warn: no such ref") || strings.Contains(err.Error(), "error: malformed object name") {
				return fmt.Errorf("specified main branch '%s' not found", targetMainBranch)
			}
			return fmt.Errorf("failed to list merged branches: %w", err)
		}

		// --- Step 3: Parse the output ---
		lines := strings.Split(mergedBranchesOutput, "\n")

		// --- Step 4: Filter the branches ---
		var branchesToProcess []string // Slice to hold branches we might delete/list
		// currentBranch := "" // We don't strictly need to store the current branch name globally

		fmt.Println("\nProcessing branches:")
		for _, line := range lines {
			// Trim leading/trailing whitespace
			branchName := strings.TrimSpace(line)

			// Skip empty lines that might result from splitting
			if branchName == "" {
				continue
			}

			// Check for current branch indication '*' and skip it
			if strings.HasPrefix(branchName, "* ") {
				// currentBranch = strings.TrimPrefix(branchName, "* ") // Extract if needed later
				fmt.Printf(" - Found current branch: %s (skipping)\n", strings.TrimPrefix(branchName, "* "))
				continue // Skip the current branch - safer not to delete it automatically
			}

			// Skip the main branch itself
			if branchName == targetMainBranch {
				fmt.Printf(" - Found target main branch: %s (skipping)\n", branchName)
				continue
			}

			// If it passes all checks, add it to the list of branches to potentially delete/list
			fmt.Printf(" + Found merged branch: %s (candidate for action)\n", branchName)
			branchesToProcess = append(branchesToProcess, branchName)
		}

		// --- Step 5: Perform action (Next Step) ---
		fmt.Printf("\nBranches identified for processing (excluding current and main): %v\n", branchesToProcess)

		fmt.Println("\n(Action logic not yet implemented)")

		return nil // Return nil for success
	},
}

// --- Helper Functions ---

// runGitCommand executes a git command and returns its trimmed output or an error.
func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run() // Use Run instead of Output to capture stderr

	if err != nil {
		// Combine stderr with the error message for more context
		return "", fmt.Errorf("command 'git %s' failed: %w\nStderr: %s", strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// detectDefaultMainBranch tries to find 'main' or 'master' branch.
func detectDefaultMainBranch() (string, error) {
	// Check if 'main' exists
	_, errMain := runGitCommand("show-ref", "--verify", "--quiet", "refs/heads/main")
	if errMain == nil {
		return "main", nil // 'main' found
	}

	// Check if 'master' exists
	_, errMaster := runGitCommand("show-ref", "--verify", "--quiet", "refs/heads/master")
	if errMaster == nil {
		return "master", nil // 'master' found
	}

	// Neither found
	return "", errors.New("neither 'main' nor 'master' branch found. Please specify with --main flag")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init defines flags and configuration settings.
func init() {
	rootCmd.Flags().StringVar(&mainBranchName, "main", "", "Specify the main branch (e.g., main, master, develop)")
	rootCmd.Flags().BoolVar(&deleteBranches, "delete", false, "Actually delete the merged branches")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what branches would be deleted without actually deleting")
}
