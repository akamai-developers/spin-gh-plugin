package spinapp

import "testing"

func TestGetVarSafeAppName(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		want    string
	}{
		{"simple name is upper-cased", "myapp", "MYAPP"},
		{"dashes become underscores", "my-cool-app", "MY_COOL_APP"},
		{"spaces and dots are sanitized", "my app.v2", "MY_APP_V2"},
		{"leading digit is prefixed with underscore", "1app", "_1APP"},
		{"leading and trailing separators are trimmed", "-my-app-", "MY_APP"},
		{"empty name yields empty string", "", ""},
		{"whitespace-only name yields empty string", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{name: tt.appName}
			if got := app.GetVarSafeAppName(); got != tt.want {
				t.Errorf("GetVarSafeAppName(%q) = %q, want %q", tt.appName, got, tt.want)
			}
		})
	}
}
