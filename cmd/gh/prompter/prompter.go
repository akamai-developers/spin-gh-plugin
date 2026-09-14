package prompter

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/akamai-developers/spin-gh-plugin/internal/spinapp"
)

var (
	deploymentNameRegex     = regexp.MustCompile(`^[a-z0-9]+([._-][a-z0-9]+)*$`)
	deploymentNameMaxLength = 128
)

func validateDeploymentName(name string) error {
	// 1. Check length
	if len(name) > deploymentNameMaxLength {
		return fmt.Errorf("Deployment name exceeds maximum length of %d", deploymentNameMaxLength)
	}

	// 2. Check pattern
	if !deploymentNameRegex.MatchString(name) {
		return fmt.Errorf("Deployment name contains invalid characters or format")
	}
	return nil
}

type ociConfig struct {
	LoginServer    string
	User           string
	References     []string
	UseGitHubToken bool
}

// PromptForAkamaiFunctionDeploymentDetails collects one deployment name per
// discovered Spin app. All apps are gathered on a single screen so the user can
// review them together, which matters when several apps were discovered.
func PromptForAkamaiFunctionDeploymentDetails(apps []*spinapp.App) {
	intro := huh.NewNote().
		Title("Deploy to Akamai Functions").
		Description(akamaiIntro(len(apps)))

	fields := []huh.Field{intro}
	for i := range apps {
		app := apps[i]
		fields = append(fields, huh.NewInput().
			Title(fmt.Sprintf("Deployment name for %q", app.GetName())).
			Description(fmt.Sprintf("Name on your Akamai Functions account · source: %s", app.GetLocation())).
			Placeholder(app.GetName()).
			Value(&app.DeploymentName).
			Validate(validateDeploymentName))
	}

	form := huh.NewForm(
		huh.NewGroup(fields...),
	).WithTheme(huh.ThemeCatppuccin()).WithOutput(os.Stderr)
	if err := form.Run(); err != nil {
		log.Fatalf("Error while gathering deployment names")
	}
}

func akamaiIntro(count int) string {
	secrets := "The generated workflow authenticates with the repository secrets\n" +
		"AKAMAI_FUNCTIONS_TOKEN and AKAMAI_FUNCTIONS_ACCOUNT_ID — add those in\n" +
		"GitHub before the workflow runs."
	if count == 1 {
		return "Choose the name your app should use on Akamai Functions.\n\n" + secrets
	}
	return fmt.Sprintf(
		"Found %d Spin apps. Choose the name each should use on Akamai\n"+
			"Functions (defaults match the app name).\n\n%s",
		count, secrets)
}

// PromptForOciDetails walks the user through the OCI publish settings for each
// discovered Spin app, one app at a time. Each app's form is self-identifying so
// progress stays clear when multiple apps were discovered.
func PromptForOciDetails(apps []*spinapp.App) {
	total := len(apps)
	for i := range apps {
		app := apps[i]
		ociCfg := promptForOciConfig(app, i+1, total)
		app.OciLoginServer = ociCfg.LoginServer
		app.OciReferences = ociCfg.References
		app.OciUser = ociCfg.User
		app.OciUseGitHubToken = ociCfg.UseGitHubToken
	}
}

const (
	defaultOciHost      = "ghcr.io"
	defaultOciNamespace = "my-org"
)

func generateArtifactPlaceholder(appName string) string {
	invalidOciCharRegex := regexp.MustCompile(`[^a-z0-9._-]+`)
	consecutiveDashes := regexp.MustCompile(`[-_.]+`)
	name := strings.ToLower(strings.TrimSpace(appName))
	name = invalidOciCharRegex.ReplaceAllString(name, "-")
	name = consecutiveDashes.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-_.")
	if name == "" {
		name = "the-unnamed"
	}

	return fmt.Sprintf("%s/%s/%s", defaultOciHost, defaultOciNamespace, name)
}

func promptForOciConfig(app *spinapp.App, index, total int) ociConfig {
	appName := app.GetName()

	var ociArtifactName string
	var ociAuthRequired bool
	var ociLoginServer string
	var ociUser string
	var ociTags string
	var ociCustomTags string

	isGitHubContainerRegistry := func(input string) bool {
		return strings.HasPrefix(input, "ghcr.io")
	}

	// The password/token is never entered here — it is read at runtime from a
	// per-app repository secret. Surface that name so the user can create it.
	passwordSecret := fmt.Sprintf("OCI_REGISTRY_PASSWORD_%s", app.GetVarSafeAppName())

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(fmt.Sprintf("Publish %q to an OCI registry%s", appName, appProgress(index, total))).
				Description("These values become the publish steps in your workflow.\n"+
					"ghcr.io is detected automatically and authenticated with the built-in GITHUB_TOKEN."),

			huh.NewInput().
				Title("Artifact name").
				Description("Full repository path, without a tag.").
				Placeholder(generateArtifactPlaceholder(appName)).
				Value(&ociArtifactName).
				Validate(func(str string) error {
					if len(strings.TrimSpace(str)) == 0 {
						return fmt.Errorf("artifact name cannot be empty")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Tags").
				Description("Which tag(s) to publish on every run.").
				Options(
					huh.NewOption("latest — standard production tag", "latest"),
					huh.NewOption("commit SHA — one immutable tag per commit", "git_sha"),
					huh.NewOption("latest + commit SHA", "both"),
					huh.NewOption("custom tag(s)…", "custom"),
				).
				Value(&ociTags),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Custom tag(s)").
				Description("Comma-separated, e.g. v1.0.0, staging").
				Placeholder("v1.0.0, staging").
				Value(&ociCustomTags).
				Validate(func(str string) error {
					if len(strings.TrimSpace(str)) == 0 {
						return fmt.Errorf("custom tag cannot be empty")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return ociTags != "custom"
		}),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Does this registry require authentication?").
				Description("Adds a login step to the workflow before pushing.").
				Affirmative("Yes").
				Negative("No").
				Value(&ociAuthRequired),
		).WithHideFunc(func() bool {
			// ghcr.io is always handled separately (GITHUB_TOKEN), so we never
			// ask about its authentication.
			return isGitHubContainerRegistry(ociArtifactName)
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Login server").
				Description("Registry host used for `spin registry login`.").
				PlaceholderFunc(func() string {
					s, _, ok := strings.Cut(ociArtifactName, "/")
					if !ok {
						return defaultOciHost
					}
					return s
				}, nil).
				Value(&ociLoginServer),

			huh.NewInput().
				Title("Username").
				Description(fmt.Sprintf("Login username. The matching password/token is read from the repository secret %s.", passwordSecret)).
				Value(&ociUser),
		).
			Title(fmt.Sprintf("Registry credentials%s", appProgress(index, total))).
			Description("Only the username is stored in the workflow — keep the password/token in GitHub secrets.").
			WithHideFunc(func() bool {
				return !ociAuthRequired || isGitHubContainerRegistry(ociArtifactName)
			}),
	).WithTheme(huh.ThemeCatppuccin()).WithOutput(os.Stderr)

	if err := form.Run(); err != nil {
		log.Fatalf("Could not collect OCI configuration for %s\n", appName)
	}

	ociReferences := []string{}
	switch ociTags {
	case "latest":
		ociReferences = append(ociReferences, buildOciReference(ociArtifactName, tagLatest))
	case "git_sha":
		ociReferences = append(ociReferences, buildOciReference(ociArtifactName, "${{ github.sha }}"))
	case "both":
		ociReferences = append(ociReferences, buildOciReference(ociArtifactName, tagLatest))
		ociReferences = append(ociReferences, buildOciReference(ociArtifactName, "${{ github.sha }}"))
	case "custom":
		customTags := strings.Split(ociCustomTags, ",")
		for _, t := range customTags {
			cleaned := strings.TrimSpace(t)
			if cleaned != "" {
				ociReferences = append(ociReferences, buildOciReference(ociArtifactName, cleaned))
			}
		}
	}

	// ghcr.io is authenticated with the built-in GITHUB_TOKEN, so we skip the
	// credential prompts above and log in against ghcr.io directly.
	if isGitHubContainerRegistry(ociArtifactName) {
		return ociConfig{
			LoginServer:    defaultOciHost,
			References:     ociReferences,
			UseGitHubToken: true,
		}
	}

	return ociConfig{
		LoginServer: ociLoginServer,
		User:        ociUser,
		References:  ociReferences,
	}
}

// appProgress renders a " (2/3)" suffix, but only when more than one app was
// discovered, so the single-app case stays uncluttered.
func appProgress(index, total int) string {
	if total <= 1 {
		return ""
	}
	return fmt.Sprintf(" (%d/%d)", index, total)
}

const tagLatest = "latest"

func buildOciReference(artifact string, tag string) string {
	return fmt.Sprintf("%s:%s", artifact, tag)
}
