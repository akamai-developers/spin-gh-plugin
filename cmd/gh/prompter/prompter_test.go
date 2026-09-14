package prompter

import "testing"

func TestGenerateArtifactPlaceholder(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		want    string
	}{
		{"simple name is lower-cased", "MyApp", "ghcr.io/my-org/myapp"},
		{"spaces become dashes", "my cool app", "ghcr.io/my-org/my-cool-app"},
		{"consecutive separators collapse", "my___app", "ghcr.io/my-org/my-app"},
		{"leading and trailing separators are trimmed", "  -my-app-  ", "ghcr.io/my-org/my-app"},
		{"empty name falls back to placeholder", "", "ghcr.io/my-org/the-unnamed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateArtifactPlaceholder(tt.appName); got != tt.want {
				t.Errorf("generateArtifactPlaceholder(%q) = %q, want %q", tt.appName, got, tt.want)
			}
		})
	}
}
