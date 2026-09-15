package spinapp

type Language struct {
	Name                       string
	Hint                       string
	Plugin                     string
	Setup                      string
	Teardown                   string
	InstallDependenciesCommand string
}

var (
	Rust       = Language{Name: "Rust", Hint: "cargo.toml"}
	GoLang     = Language{Name: "Go", Hint: "go.mod"}
	JavaScript = Language{
		Name:                       "JavaScript",
		Hint:                       "package.json",
		Plugin:                     "",
		InstallDependenciesCommand: "",
	}
)
var Python = Language{
	Name:                       "Python",
	Hint:                       "requirements.txt",
	Plugin:                     "",
	Setup:                      "python3 -m venv venv && source venv/bin/activate",
	InstallDependenciesCommand: "pip install -r requirements.txt",
	Teardown:                   "",
}

var all = map[string]Language{
	"Rust":       Rust,
	"Go":         GoLang,
	"JavaScript": JavaScript,
	"Python":     Python,
}

func GetPluginByLang(lang string) string {
	return all[lang].Plugin
}

func GetInstallDependenciesByLang(lang string) string {
	return all[lang].InstallDependenciesCommand
}

func GetSetupByLang(lang string) string {
	return all[lang].Setup
}

func GetTeardownByLang(lang string) string {
	return all[lang].Teardown
}
