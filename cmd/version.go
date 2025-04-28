package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These variables will be populated by GoReleaser during the build process using ldflags.
// They need to be package-level variables in the 'cmd' package (or wherever your version command is).
var (
	version = "dev"     // Default value if not built with ldflags
	commit  = "none"    // Default value
	date    = "unknown" // Default value
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number, commit hash, and build date of git-util",
	Long:  `All software has versions. This is git-util's.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Print the populated variables.
		fmt.Printf("git-util version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", date)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd) // Add versionCmd to the root command
	// No flags needed for the version command itself.
}
