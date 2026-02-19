package spinapp

import (
	"fmt"
)

type App struct {
	location       string
	languages      []Language
	name           string
	deploymentName string
	components     []Component
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
		deploymentName: appName,
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

func (app *App) GetDeploymentName() string {
	return app.deploymentName
}

func (app *App) SetDeploymentName(d string) {
	app.deploymentName = d
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
