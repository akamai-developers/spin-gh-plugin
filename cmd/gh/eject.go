package gh

import (
	"fmt"
	"log"
	"os"

	gh "github.com/akamai-developers/spin-gh-plugin/internal/github"
	"github.com/spf13/cobra"
)

var output string
var overwrite bool
var ejectCmd = &cobra.Command{
	Use:   "eject",
	Short: "Render the default GitHub Actions template",
	Run: func(cmd *cobra.Command, args []string) {
		template := gh.GetDefaultTemplate()
		if output != "" {
			if !overwrite {
				if _, err := os.Stat(output); err == nil {
					log.Fatalf("File %s already exists. Use --overwrite to overwrite it.\n", output)
					return
				}
			}
			err := os.WriteFile(output, []byte(template), 0644)
			if err != nil {
				log.Fatalf("Failed to write to file %s: %v\n", output, err)
			}
			log.Printf("Template written to %s\n", output)
		}
		// output was empty -> print to STDOUT
		fmt.Println(template)
	},
}

func init() {
	ejectCmd.Flags().StringVarP(&output, "output", "o", "", "Path where the default template should be rendered to")
	ejectCmd.Flags().BoolVarP(&overwrite, "overwrite", "", false, "Overwrite existing output file")
	rootCmd.AddCommand(ejectCmd)
}
