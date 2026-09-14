package github

import (
	"testing"
)

func secretNames(secrets []requiredSecret) []string {
	names := make([]string, len(secrets))
	for i, s := range secrets {
		names[i] = s.Name
	}
	return names
}

func TestCollectRequiredSecrets(t *testing.T) {
	tests := []struct {
		name string
		data templateData
		want []string
	}{
		{
			name: "no secrets required",
			data: templateData{},
			want: nil,
		},
		{
			name: "oci app with custom registry needs a password secret",
			data: templateData{
				PushOciArtifacts: true,
				SpinApps: []spinAppTemplateData{
					{Name: "my-app", VarSafeAppName: "MY_APP", OciLoginServer: "registry.example.com", OciUser: "alice"},
				},
			},
			want: []string{"OCI_REGISTRY_PASSWORD_MY_APP"},
		},
		{
			name: "ghcr apps using the built-in token need no secret",
			data: templateData{
				PushOciArtifacts: true,
				SpinApps: []spinAppTemplateData{
					{Name: "gh-app", VarSafeAppName: "GH_APP", OciLoginServer: "ghcr.io", OciUseGitHubToken: true},
				},
			},
			want: nil,
		},
		{
			name: "apps without a login server need no secret",
			data: templateData{
				PushOciArtifacts: true,
				SpinApps: []spinAppTemplateData{
					{Name: "local-app", VarSafeAppName: "LOCAL_APP"},
				},
			},
			want: nil,
		},
		{
			name: "duplicate secret names are only reported once",
			data: templateData{
				PushOciArtifacts: true,
				SpinApps: []spinAppTemplateData{
					{Name: "app", VarSafeAppName: "APP", OciLoginServer: "registry.example.com", OciUser: "alice"},
					{Name: "app", VarSafeAppName: "APP", OciLoginServer: "registry.example.com", OciUser: "alice"},
				},
			},
			want: []string{"OCI_REGISTRY_PASSWORD_APP"},
		},
		{
			name: "akamai deployment adds token and account id secrets",
			data: templateData{
				DeployToAkamaiFunctions: true,
			},
			want: []string{"AKAMAI_FUNCTIONS_TOKEN", "AKAMAI_FUNCTIONS_ACCOUNT_ID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := secretNames(collectRequiredSecrets(tt.data))
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("secret[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
