package githubactions

import "testing"

func TestGetExampleWorkflowFileURL(t *testing.T) {
	t.Parallel()

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

	for _, testCase := range tests {
		t.Run(testCase.language+"-"+testCase.runtimeCloudProvider, func(t *testing.T) {
			t.Parallel()

			got, err := getExampleWorkflowFileURL(testCase.language, testCase.runtimeCloudProvider)
			if err != nil {
				t.Errorf("getExampleWorkflowFileURL() error = %v", err)

				return
			}

			if got != testCase.want {
				t.Errorf("getExampleWorkflowFileURL() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestGetExampleWorkflowFileURLInvalidLanguage(t *testing.T) {
	t.Parallel()

	_, err := getExampleWorkflowFileURL("invalid", "aks")
	if err == nil {
		t.Errorf("getExampleWorkflowFileURL() error = %v, want not nil", err)
	}
}

func TestGetExampleWorkflowFileURLInvalidRuntimeCloudProvider(t *testing.T) {
	t.Parallel()

	_, err := getExampleWorkflowFileURL("dotnet", "invalid")
	if err == nil {
		t.Errorf("getExampleWorkflowFileURL() error = %v, want not nil", err)
	}
}

func TestGetExampleWorkflowFileURLInvalidLanguageAndRuntimeCloudProvider(t *testing.T) {
	t.Parallel()

	_, err := getExampleWorkflowFileURL("invalid", "invalid")
	if err == nil {
		t.Errorf("getExampleWorkflowFileURL() error = %v, want not nil", err)
	}
}
