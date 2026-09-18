package gh

import (
	"log"

	"github.com/akamai-developers/spin-gh-plugin/cmd/gh/prompter"
	"github.com/akamai-developers/spin-gh-plugin/internal/detective"
	gh "github.com/akamai-developers/spin-gh-plugin/internal/github"
	"github.com/spf13/cobra"
)

type CreateActionOptions struct {
	DryRun bool
	gh.ActionTriggers
	DeployToAkamaiFunctions bool
	PushOciArtifacts        bool
	GenerateSbom            bool
	GenerateSignature       bool
	Name                    string
	OperatingSystem         string
	Output                  string
	Overwrite               bool
	Plugins                 []string
	EnvironmentVariables    []string
	TemplatePath            string
	Tools                   gh.Tools
}

var options CreateActionOptions = CreateActionOptions{
	Tools: gh.DefaultTools(),
}

var createActionCmd = &cobra.Command{
	Use:     "create-action",
	Aliases: []string{"create", "generate"},
	Short:   "Examines your Spin App and creates a GitHub Action workflow file",
	Run: func(cmd *cobra.Command, args []string) {
		// although we could simply enable OCI publishing
		// we should not reconfigure behavior on-the-fly
		if options.GenerateSbom || options.GenerateSignature {
			if !options.PushOciArtifacts {
				log.Fatal("Generating SBOM or Signing artifacts requires OCI artifact publishing flag (--push-oci-artifacts)")
			}
		}
		apps := detective.FindAllSpinApps()
		if len(apps) == 0 {
			log.Fatal("Could not find Spin App(s) under the current directory")
		}
		if options.DeployToAkamaiFunctions {
			prompter.PromptForAkamaiFunctionDeploymentDetails(apps)
		}
		if options.PushOciArtifacts {
			prompter.PromptForOciDetails(apps)
		}

		envVars, err := gh.ParseEnvVars(options.EnvironmentVariables)
		if err != nil {
			log.Fatalf("Invalid Environment variable provided %v", err)
		}

		renderOptions := gh.RenderActionOptions{
			DeployToAkamaiFunctions: options.DeployToAkamaiFunctions,
			PushOciArtifacts:        options.PushOciArtifacts,
			CustomTemplatePath:      options.TemplatePath,
			GenerateSbom:            options.GenerateSbom,
			GenerateSignature:       options.GenerateSignature,
			DryRun:                  options.DryRun,
			Name:                    options.Name,
			OperatingSystem:         options.OperatingSystem,
			Output:                  options.Output,
			Overwrite:               options.Overwrite,
			Plugins:                 options.Plugins,
			SpinApps:                apps,
			ActionTriggers:          options.ActionTriggers,
			Tools:                   options.Tools,
			EnvironmentVariables:    envVars,
		}
		err = gh.RenderAction(renderOptions)
		if err != nil {
			log.Fatalf("Error while rendering template %v", err)
		}
	},
}

func init() {

	createActionCmd.Flags().BoolVarP(&options.ActionTriggers.ManualDispatch, "manual", "", false, "Trigger Action on workflow dispatch")
	createActionCmd.Flags().StringVarP(&options.ActionTriggers.Schedule, "cron", "", "", "Trigger Action on cron schedule")
	createActionCmd.Flags().StringVarP(&options.ActionTriggers.Push, "ci", "", "main", "Trigger Action for every push on the specified branch")
	createActionCmd.Flags().StringVarP(&options.ActionTriggers.PullRequest, "pr", "", "", "Trigger Action for every PR targeting the specified branch")

	createActionCmd.Flags().StringVarP(&options.Tools.Rust, "rust-version", "", options.Tools.Rust, "Set Rust version for GitHub Actions")
	createActionCmd.Flags().StringVarP(&options.Tools.RustTarget, "rust-target", "", options.Tools.RustTarget, "Set the desired Rust Target for GitHub Actions")
	createActionCmd.Flags().StringVarP(&options.Tools.Go, "go-version", "", options.Tools.Go, "Set Go version for GitHub Actions")
	createActionCmd.Flags().StringVarP(&options.Tools.Node, "node-version", "", options.Tools.Node, "Set Node.js version for GitHub Actions")
	createActionCmd.Flags().StringVarP(&options.Tools.Python, "python-version", "", options.Tools.Python, "Set Python version for GitHub Actions")
	createActionCmd.Flags().StringVarP(&options.Tools.Spin, "spin-version", "", "", "Set Spin version for GitHub Actions (default: current stable Spin release)")
	createActionCmd.Flags().StringSliceVarP(&options.Plugins, "spin-plugins", "p", []string{}, "Specify required Spin plugins")

	createActionCmd.Flags().BoolVarP(&options.Overwrite, "overwrite", "", false, "Overwrite existing output file")
	createActionCmd.Flags().BoolVarP(&options.DryRun, "dry-run", "", false, "Print GitHub Action to stdout instead of writing to a file")
	createActionCmd.Flags().StringVarP(&options.Output, "output", "o", ".github/workflows/ci.yaml", "Path where the GitHub Action will be created")

	createActionCmd.Flags().StringVarP(&options.Name, "name", "n", "CI", "Name for the GitHub Action")
	createActionCmd.Flags().StringVarP(&options.TemplatePath, "template", "t", "", "Specify the path to a custom template for creating the GitHub Action")
	createActionCmd.Flags().StringSliceVarP(&options.EnvironmentVariables, "env", "", []string{}, "Specify Environment Variables (format key=value)")
	createActionCmd.Flags().StringVarP(&options.OperatingSystem, "os", "", "ubuntu-latest", "Specify the desired operating system for the GitHub Action")

	createActionCmd.Flags().BoolVarP(&options.DeployToAkamaiFunctions, "deploy-to-akamai-functions", "", false, "Add steps for deploying your Spin App(s) to Akamai Functions")
	createActionCmd.Flags().BoolVarP(&options.PushOciArtifacts, "push-oci-artifacts", "", false, "Add steps to publish your Spin App(s) to an OCI registry")
	createActionCmd.Flags().BoolVarP(&options.GenerateSbom, "sbom", "", false, "Generate SBOM for OCI artifacts")
	createActionCmd.Flags().BoolVarP(&options.GenerateSignature, "sign", "", false, "Sign OCI artifacts")
	rootCmd.AddCommand(createActionCmd)
}
