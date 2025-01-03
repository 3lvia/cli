package shared

import "testing"

func TestNameToEnvVar(t *testing.T) {
	t.Parallel()

	expectedAndInput := map[string]string{
		"3LV_PROJECT_FILE":                      "project-file",
		"3LV_RUNTIME_CLOUD_PROVIDER":            "runtime-cloud-provider",
		"3LV_SYSTEM_NAME":                       "system-name",
		"3LV_APPLICATION_NAME":                  "application-name",
		"3LV_WHAT_IS_THIS_ENVIRONMENT_VARIABLE": "what-is-this-environment-variable",
	}

	for expected, input := range expectedAndInput {
		actual := nameToEnvVar(input)
		if actual != expected {
			t.Errorf("expected %s but got %s", expected, actual)
		}
	}
}
