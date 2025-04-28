package cmd

import (
	// Required by runGitCommand (if called from here or package)
	"fmt"
	"os" // Required by runGitCommand (if called from here or package)
	"path/filepath"
	"strings"

	// We don't need 'strconv' in sync.go

	"github.com/spf13/cobra"
	// NOTE: This code now relies on findGitRepos and runGitCommand
	// being defined elsewhere in the 'cmd' package (e.g., in root.go)
	// or being imported from a shared package if you refactor later.
)

// Variables to hold the flag values for the sync command
var (
	syncDirectory string
	syncAction    string
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize multiple Git repositories (fetch or pull).",
	Long: `Scans a directory for Git repositories and runs 'git fetch --prune' (default)
or 'git pull --ff-only' to synchronize them with their remotes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --- Determine Target Directory ---
		targetDir := syncDirectory
		if targetDir == "" {
			var err error
			targetDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
		}
		// Get absolute path for consistency
		var err error
		targetDir, err = filepath.Abs(targetDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for target directory: %w", err)
		}

		// --- Validate Action ---
		action := strings.ToLower(syncAction)
		if action != "fetch" && action != "pull" {
			return fmt.Errorf("invalid action '%s': must be 'fetch' or 'pull'", syncAction)
		}
		fmt.Printf("Scanning directory: %s (Action: %s)\n", targetDir, action)

		// --- Find Repositories ---
		// *** This call was missing in the code you pasted ***
		// Call the helper function (defined elsewhere in cmd package, e.g., root.go)
		repos, err := findGitRepos(targetDir)
		if err != nil {
			return fmt.Errorf("error finding repositories: %w", err)
		}

		if len(repos) == 0 {
			fmt.Println("No Git repositories found in the specified directory.")
			return nil
		}

		fmt.Printf("\n--- Synchronizing Repositories ---\n")

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
		successCount := 0
		failCount := 0
		for _, repoPath := range repos {
			relPath, _ := filepath.Rel(targetDir, repoPath)
			if relPath == "." {
				relPath = filepath.Base(targetDir)
			}

			fmt.Printf("%-*s : Syncing (%s)... ", maxLen, relPath, action)

			var gitArgs []string
			if action == "fetch" {
				gitArgs = []string{"fetch", "--prune"}
			} else { // action == "pull"
				gitArgs = []string{"pull", "--ff-only"}
			}

			// Prepend -C <path> to run in the correct directory
			gitFullArgs := append([]string{"-C", repoPath}, gitArgs...)

			// *** This call was missing in the code you pasted ***
			// Execute the command using the helper (defined elsewhere in cmd package, e.g., root.go)
			output, err := runGitCommand(gitFullArgs...) // output is now assigned

			// Check for errors after executing the command
			if err != nil {
				fmt.Printf("FAILED\n")
				// Print concise error, including output from the command
				fmt.Fprintf(os.Stderr, "  Error for %s: %v\n  Output: %s\n", relPath, err, output) // Using output here
				failCount++
			} else {
				fmt.Printf("OK\n")
				// Optional: Print output even on success if needed (can be noisy)
				// if output != "" { fmt.Printf("  Output: %s\n", output) }
				successCount++
			}
		} // End loop

		// Print summary
		fmt.Printf("\n--- Summary ---\n")
		fmt.Printf("Action '%s' completed.\n", action)
		fmt.Printf("  Successfully synced: %d\n", successCount)
		fmt.Printf("  Failed to sync:    %d\n", failCount)

		return nil
	},
}

func init() {
	// Register syncCmd with the root command
	rootCmd.AddCommand(syncCmd)

	// Define flags specific to the sync command
	syncCmd.Flags().StringVarP(&syncDirectory, "directory", "D", "", "Directory to scan for Git repositories (defaults to current directory)")
	syncCmd.Flags().StringVarP(&syncAction, "action", "a", "fetch", "Sync action to perform: 'fetch' (default) or 'pull'")
}

// --- Helper Functions ---
// *** DELETE THE DEFINITIONS of findGitRepos and runGitCommand from this file! ***
// They should only be defined ONCE in your cmd package (e.g., in cmd/root.go)
// Leaving them defined here WILL cause the 'redeclared' build errors again.

/* DELETE THIS FUNCTION DEFINITION from sync.go:
func findGitRepos(rootDir string) ([]string, error) { ... }
*/

/* DELETE THIS FUNCTION DEFINITION from sync.go:
func runGitCommand(args ...string) (string, error) { ... }
*/
