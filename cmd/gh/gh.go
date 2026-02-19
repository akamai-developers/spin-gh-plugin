package gh

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"
var rootCmd = &cobra.Command{
	Use:   "gh",
	Short: "Generates GitHub Actions for Spin Apps",
	Long: `Use the sub-commands provided to create tailored GitHub Action Workflows

This plugin will recursively examine the current folder and all its sub-folders to detect all Spin applications
and their components. 

It will create a GitHub Action workflow file that has all necessary language-specific tooling installed and will compile all Spin apps.

Optionally, it could setup continuous deployment for Akamai Functions as well.
`,
	Version: Version,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func ExecuteRootCommand() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
