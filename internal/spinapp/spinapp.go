package spinapp

import (
	"fmt"
	"regexp"
	"strings"
)

type App struct {
	location          string
	languages         []Language
	name              string
	DeploymentName    string
	OciReferences     []string
	OciUser           string
	OciLoginServer    string
	OciUseGitHubToken bool
	components        []Component
}

type Component struct {
	Language string
	Location string
}

func NewApp(location string) (*App, error) {
	appName, err := getAppNameFromManifest(location)
	if err != nil {
		return nil, err
	}
	return &App{
		name:           appName,
		DeploymentName: appName,
		location:       location,
	}, nil
}

func (app *App) AddComponent(c Component) {
	app.components = append(app.components, c)
}

func (app *App) GetComponents() []Component {
	return app.components
}
func (app *App) GetName() string {
	return app.name
}

func (app *App) GetVarSafeAppName() string {
	var (
		invalidCharRegex  = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
		leadingDigitRegex = regexp.MustCompile(`^[0-9]`)
	)

	if strings.TrimSpace(app.name) == "" {
		return ""
	}
	varSafeAppName := strings.ToUpper(app.name)
	varSafeAppName = invalidCharRegex.ReplaceAllString(varSafeAppName, "_")
	varSafeAppName = strings.Trim(varSafeAppName, "_")
	if leadingDigitRegex.MatchString(varSafeAppName) {
		varSafeAppName = "_" + varSafeAppName
	}
	return varSafeAppName
}

func (app *App) GetLanguages() []Language {
	return app.languages
}

func (app *App) GetLocation() string {
	return app.location
}

func (app *App) ToString() string {
	return fmt.Sprintf("%s at %s", app.GetName(), app.GetLocation())
}
