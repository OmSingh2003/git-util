package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Variables to hold the flag values
var (
	mainBranchName string
	deleteBranches bool
	dryRun         bool
)

// rootCmd represents the base command when called without any subcommands
// For now, this will be our main 'git-cleaner' functionality command.
var rootCmd = &cobra.Command{
	// Use: is the one-line usage message.
	// Change this to 'git-util' to reflect our project name.
	Use:   "git-util",
	Short: "A utility tool for common Git operations.",
	Long: `git-util helps automate and simplify various Git tasks.
The first feature implemented is cleaning up merged local branches.
More features might be added later.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		// This is where the main logic for our branch cleaner will go.
		fmt.Println("Executing Git Branch Cleaner functionality...")
		fmt.Printf("Flags:\n")
		fmt.Printf("  --main: %s\n", mainBranchName)
		fmt.Printf("  --delete: %t\n", deleteBranches)
		fmt.Printf("  --dry-run: %t\n", dryRun)

		// --- TODO: Implement the core logic here ---
		// 1. Determine the target main branch (use flag or auto-detect).
		// 2. Run `git branch --merged <target_main_branch>`.
		// 3. Parse the output.
		// 4. Filter the branches (remove main branch, current branch '*').
		// 5. Perform action: Print list OR (if --delete) print/execute delete commands.

		fmt.Println("\n(Core logic not yet implemented)") // Placeholder message

		return nil // Return nil for success, or an error if something goes wrong
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// init() runs before main(). We define our flags here.

	// rootCmd.Flags() are flags specific to *this* command.
	// Use rootCmd.PersistentFlags() if you want flags to be available to subcommands too (useful later).
	rootCmd.Flags().StringVar(&mainBranchName, "main", "", "Specify the main branch (e.g., main, master, develop)")
	rootCmd.Flags().BoolVar(&deleteBranches, "delete", false, "Actually delete the merged branches")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what branches would be deleted without actually deleting")

}
