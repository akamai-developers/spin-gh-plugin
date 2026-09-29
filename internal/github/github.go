package github

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/akamai-developers/spin-gh-plugin/internal/spinapp"
)

type RenderActionOptions struct {
	CustomTemplatePath      string
	DeployToAkamaiFunctions bool
	PushOciArtifacts        bool
	GenerateSbom            bool
	GenerateSignature       bool
	DryRun                  bool
	Name                    string
	OperatingSystem         string
	Output                  string
	Overwrite               bool
	Plugins                 []string
	EnvironmentVariables    []*EnvVar
	SpinApps                []*spinapp.App
	SpinVersion             string
	Tools
	ActionTriggers
}

func RenderAction(options RenderActionOptions) error {
	templateContent, err := getTemplateContents(options.CustomTemplatePath)
	if err != nil {
		return err
	}
	templ, err := template.New("ci.yaml").Parse(templateContent)
	if err != nil {
		return err
	}
	target, err := getTarget(options)
	if err != nil {
		return err
	}
	defer target.Close()
	data := buildTemplateData(options)

	err = templ.Execute(target, data)
	if err != nil {
		return err
	}

	// Secrets summary always goes to stderr so that stdout stays a clean YAML
	// document (important when the user redirects --dry-run output to a file).
	writeRequiredSecrets(os.Stderr, data)
	return nil
}

type requiredSecret struct {
	Name        string
	Description string
}

// collectRequiredSecrets returns the repository secrets the generated workflow
// references and that the user has to create themselves. The built-in
// GITHUB_TOKEN is intentionally excluded because GitHub provides it
// automatically.
func collectRequiredSecrets(data templateData) []requiredSecret {
	var secrets []requiredSecret
	seen := map[string]bool{}
	add := func(name, description string) {
		if seen[name] {
			return
		}
		seen[name] = true
		secrets = append(secrets, requiredSecret{Name: name, Description: description})
	}

	if data.PushOciArtifacts {
		for _, app := range data.SpinApps {
			// ghcr.io logs in with the built-in GITHUB_TOKEN, and apps without a
			// login server don't authenticate at all.
			if app.OciLoginServer == "" || app.OciUseGitHubToken {
				continue
			}
			add(
				fmt.Sprintf("OCI_REGISTRY_PASSWORD_%s", app.VarSafeAppName),
				fmt.Sprintf("Password or token for %q to authenticate as %q against %s.", app.Name, app.OciUser, app.OciLoginServer),
			)
		}
	}

	if data.DeployToAkamaiFunctions {
		add("AKAMAI_FUNCTIONS_TOKEN", "Authentication token used by `spin aka login`.")
		add("AKAMAI_FUNCTIONS_ACCOUNT_ID", "Akamai Functions account ID used to list, link and deploy your apps.")
	}

	return secrets
}

func writeRequiredSecrets(w io.Writer, data templateData) {
	secrets := collectRequiredSecrets(data)
	if len(secrets) == 0 {
		return
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "The generated workflow references the following repository secrets.")
	fmt.Fprintln(w, "Add them under Settings → Secrets and variables → Actions before the workflow runs:")
	fmt.Fprintln(w)
	for _, s := range secrets {
		fmt.Fprintf(w, "  • %s\n", s.Name)
		fmt.Fprintf(w, "      %s\n", s.Description)
	}
	fmt.Fprintln(w)
}

func getTarget(options RenderActionOptions) (io.WriteCloser, error) {
	if options.DryRun {
		return os.Stdout, nil
	}
	if !options.Overwrite {
		if _, err := os.Stat(options.Output); err == nil {
			return nil, fmt.Errorf("pass --overwrite to overwrite an existing GitHub Action file")
		}
	}
	if _, err := os.Stat(options.Output); err == nil {
		if err := os.Remove(options.Output); err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(filepath.Dir(options.Output), os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.Create(options.Output)
	if err != nil {
		return nil, err
	}
	return file, nil
}
