package githubactions

import "testing"

func TestGetExampleWorkflowFileURL(t *testing.T) {
	tests := []struct {
		language             string
		runtimeCloudProvider string
		want                 string
	}{
		{
			language:             "dotnet",
			runtimeCloudProvider: "aks",
			want:                 exampleWorkflowBaseURL + "/build-deploy-dotnet.yml",
		},
		{
			language:             "dotnet",
			runtimeCloudProvider: "gke",
			want:                 exampleWorkflowBaseURL + "/build-deploy-dotnet-google.yml",
		},
		{
			language:             "go",
			runtimeCloudProvider: "aks",
			want:                 exampleWorkflowBaseURL + "/build-deploy-go.yml",
		},
		{
			language:             "go",
			runtimeCloudProvider: "gke",
			want:                 exampleWorkflowBaseURL + "/build-deploy-go-google.yml",
		},
		{
			language:             "dockerfile",
			runtimeCloudProvider: "aks",
			want:                 exampleWorkflowBaseURL + "/build-deploy-dockerfile.yml",
		},
		{
			language:             "dockerfile",
			runtimeCloudProvider: "gke",
			want:                 exampleWorkflowBaseURL + "/build-deploy-dockerfile-google.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.language+"-"+tt.runtimeCloudProvider, func(t *testing.T) {
			got, err := getExampleWorkflowFileURL(tt.language, tt.runtimeCloudProvider)
			if err != nil {
				t.Errorf("getExampleWorkflowFileURL() error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("getExampleWorkflowFileURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetExampleWorkflowFileURLInvalidLanguage(t *testing.T) {
	_, err := getExampleWorkflowFileURL("invalid", "aks")
	if err == nil {
		t.Errorf("getExampleWorkflowFileURL() error = %v, want not nil", err)
	}
}

func TestGetExampleWorkflowFileURLInvalidRuntimeCloudProvider(t *testing.T) {
	_, err := getExampleWorkflowFileURL("dotnet", "invalid")
	if err == nil {
		t.Errorf("getExampleWorkflowFileURL() error = %v, want not nil", err)
	}
}

func TestGetExampleWorkflowFileURLInvalidLanguageAndRuntimeCloudProvider(t *testing.T) {
	_, err := getExampleWorkflowFileURL("invalid", "invalid")
	if err == nil {
		t.Errorf("getExampleWorkflowFileURL() error = %v, want not nil", err)
	}
}
