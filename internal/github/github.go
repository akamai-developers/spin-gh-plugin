package github

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/fermyon/spin-gh-plugin/internal/spinapp"
)

type RenderActionOptions struct {
	CustomTemplatePath      string
	DeployToAkamaiFunctions bool
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
	return nil
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
