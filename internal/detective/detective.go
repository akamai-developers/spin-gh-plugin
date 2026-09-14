package detective

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/akamai-developers/spin-gh-plugin/internal/spinapp"
)

var ignoreFolders = []string{"node_modules", "target"}

func FindAllSpinApps() []*spinapp.App {
	apps := make([]*spinapp.App, 0)
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && slices.Contains(ignoreFolders, d.Name()) {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.EqualFold(d.Name(), spinapp.ManifestName) {
			relativePath, err := getRelativePath(filepath.Dir(path))
			if err != nil {
				return err
			}
			app, err := spinapp.NewApp(relativePath)
			if err != nil {
				log.Printf("%v", err)
				return err
			}
			log.Printf("Discovered Spin App: %s\n", app.ToString())
			findAllAppComponents(app)
			apps = append(apps, app)
			return filepath.SkipDir
		}
		return nil
	})
	log.Println("")
	if err != nil {
		return make([]*spinapp.App, 0)
	}
	return apps
}

func findAllAppComponents(app *spinapp.App) {
	appFolder := app.GetLocation()
	filepath.WalkDir(appFolder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && slices.Contains(ignoreFolders, d.Name()) {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			var language string
			switch strings.ToLower(d.Name()) {
			case spinapp.Rust.Hint:
				language = spinapp.Rust.Name
			case spinapp.GoLang.Hint:
				language = spinapp.GoLang.Name
			case spinapp.JavaScript.Hint:
				language = spinapp.JavaScript.Name
			case spinapp.Python.Hint:
				language = spinapp.Python.Name
			default:
				return nil
			}

			relativePath, err := getRelativePath(filepath.Dir(path))
			if err != nil {
				return err
			}

			component := spinapp.Component{
				Language: language,
				Location: relativePath,
			}
			app.AddComponent(component) // Assuming there's a method to add components to the app

			return filepath.SkipDir
		}
		return nil
	})

	for _, c := range app.GetComponents() {
		log.Printf(" - %s Component discovered at %s\n", c.Language, c.Location)
	}
}

func getRelativePath(p string) (string, error) {
	relativePath, err := filepath.Rel(".", p)
	if err != nil {
		return "", err
	}
	if relativePath != "." {
		relativePath = fmt.Sprintf("./%s", relativePath)
	}
	return relativePath, nil
}
