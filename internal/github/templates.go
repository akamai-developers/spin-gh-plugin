package github

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/fermyon/spin-gh-plugin/internal/spinapp"
)

//go:embed default_workflow.yaml.tmpl
var actionsTemplate string

func GetDefaultTemplate() string {
	return actionsTemplate
}

func getTemplateContents(customTemplateFilePath string) (string, error) {
	if len(customTemplateFilePath) == 0 {
		return actionsTemplate, nil
	}
	content, err := os.ReadFile(customTemplateFilePath)
	if err != nil {
		return "", fmt.Errorf("custom template not found: %w", err)
	}
	return string(content), nil
}

type templateData struct {
	ActionName              string
	DeployToAkamaiFunctions bool
	EnvironmentVariables    []*EnvVar
	Go                      bool
	JavaScript              bool
	OperatingSystem         string
	Python                  bool
	Rust                    bool
	SpinApps                []spinAppTemplateData
	SpinPlugins             string
	TargetBranch            string
	Tools
	ActionTriggers
}

type spinAppTemplateData struct {
	Components     []componentTemplateData
	Name           string
	DeploymentName string
	Path           string
	Setup          string
	Teardown       string
}

type componentTemplateData struct {
	Language                   string
	InstallDependenciesCommand string
	Path                       string
}

func newComponentMetadata(lang string, path string) componentTemplateData {

	return componentTemplateData{
		InstallDependenciesCommand: spinapp.GetInstallDependenciesByLang(lang),
		Language:                   lang,
		Path:                       path,
	}

}

func newSpinAppMetadata(apps []*spinapp.App, setupCmds []string, teardownCmds []string) []spinAppTemplateData {
	res := make([]spinAppTemplateData, len(apps))
	for idx, app := range apps {
		am := spinAppTemplateData{
			Name:           app.GetName(),
			DeploymentName: app.GetDeploymentName(),
			Path:           app.GetLocation(),
			Components:     []componentTemplateData{},
			Setup:          strings.Join(setupCmds, " && "),
			Teardown:       strings.Join(teardownCmds, " && "),
		}
		for _, comp := range app.GetComponents() {
			cm := newComponentMetadata(comp.Language, comp.Location)
			am.Components = append(am.Components, cm)
		}
		res[idx] = am
	}
	return res
}

func buildTemplateData(options RenderActionOptions) templateData {
	allPlugins := options.Plugins
	appSetup := []string{}
	appTeardown := []string{}
	var rust, golang, js, py bool
	for _, app := range options.SpinApps {
		for _, comp := range app.GetComponents() {
			langPlugin := spinapp.GetPluginByLang(comp.Language)
			langSetup := spinapp.GetSetupByLang(comp.Language)
			langTeardown := spinapp.GetTeardownByLang(comp.Language)

			if len(langPlugin) > 0 {
				allPlugins = append(allPlugins, langPlugin)
			}
			if len(langSetup) > 0 {
				appSetup = append(appSetup, langSetup)
			}
			if len(langTeardown) > 0 {
				appTeardown = append(appTeardown, langTeardown)
			}
			switch comp.Language {
			case spinapp.Rust.Name:
				rust = true
			case spinapp.GoLang.Name:
				golang = true
			case spinapp.JavaScript.Name:
				js = true
			case spinapp.Python.Name:
				py = true
			}
		}
	}
	allPlugins = distinct(allPlugins)
	appSetup = distinct(appSetup)
	appTeardown = distinct(appTeardown)

	return templateData{
		ActionName:              options.Name,
		Go:                      golang,
		JavaScript:              js,
		OperatingSystem:         options.OperatingSystem,
		Python:                  py,
		Rust:                    rust,
		SpinApps:                newSpinAppMetadata(options.SpinApps, appSetup, appTeardown),
		SpinPlugins:             strings.Join(allPlugins, ","),
		ActionTriggers:          options.ActionTriggers,
		Tools:                   options.Tools,
		EnvironmentVariables:    options.EnvironmentVariables,
		DeployToAkamaiFunctions: options.DeployToAkamaiFunctions,
	}
}

func distinct(slice []string) []string {
	unique := make([]string, 0, len(slice))
	seen := make(map[string]struct{})
	for _, s := range slice {
		if _, ok := seen[s]; !ok && s != "" {
			seen[s] = struct{}{}
			unique = append(unique, s)
		}
	}
	return unique
}
