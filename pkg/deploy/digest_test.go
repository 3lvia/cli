package deploy

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

const renderedChart = `---
# Source: elvia-deployment/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: demo-api
spec:
  ports:
  - port: 80
---
# Source: elvia-deployment/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo-api
spec:
  template:
    spec:
      containers:
      - name: otel-sidecar
        image: ghcr.io/3lvia/core/otel:1.0
      - name: demo-api
        image: containerregistryelvia.azurecr.io/core/demo-api:v12
`

func TestMainContainerImage(t *testing.T) {
	t.Parallel()

	image, err := mainContainerImage(renderedChart)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if image != "containerregistryelvia.azurecr.io/core/demo-api:v12" {
		t.Errorf("Expected the main container image, got %s", image)
	}
}

func TestMainContainerImageWithoutDeployment(t *testing.T) {
	t.Parallel()

	_, err := mainContainerImage("apiVersion: v1\nkind: Service\nmetadata:\n  name: demo-api\n")
	if err == nil || !strings.Contains(err.Error(), "no Deployment") {
		t.Errorf("Expected a missing Deployment error, got %v", err)
	}
}

func TestMainContainerImageWithoutMainContainer(t *testing.T) {
	t.Parallel()

	rendered := strings.Replace(renderedChart, "- name: demo-api", "- name: something-else", 1)

	_, err := mainContainerImage(rendered)
	if err == nil {
		t.Error("Expected an error when no container is named after the Deployment")
	}
}

func TestSplitACRImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		image             string
		registryName      string
		repositoryWithTag string
		ok                bool
	}{
		{
			image:             "containerregistryelvia.azurecr.io/core/demo-api:v12",
			registryName:      "containerregistryelvia",
			repositoryWithTag: "core/demo-api:v12",
			ok:                true,
		},
		{
			image:             "containerregistryelvia.azurecr.io/core/demo-api:v12@sha256:abc",
			registryName:      "",
			repositoryWithTag: "",
			ok:                false,
		},
		{image: "containerregistryelvia.azurecr.io/core/demo-api", registryName: "", repositoryWithTag: "", ok: false},
		{image: "containerregistryelvia.azurecr.io/core/demo-api:", registryName: "", repositoryWithTag: "", ok: false},
		{image: "ghcr.io/3lvia/core/demo-api:v12", registryName: "", repositoryWithTag: "", ok: false},
		{image: "demo-api:v12", registryName: "", repositoryWithTag: "", ok: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.image, func(t *testing.T) {
			t.Parallel()

			registryName, repositoryWithTag, ok := splitACRImage(testCase.image)
			if registryName != testCase.registryName || repositoryWithTag != testCase.repositoryWithTag || ok != testCase.ok {
				t.Errorf(
					"Expected (%q, %q, %v), got (%q, %q, %v)",
					testCase.registryName, testCase.repositoryWithTag, testCase.ok,
					registryName, repositoryWithTag, ok,
				)
			}
		})
	}
}

func TestAzAcrImageDigestCommand(t *testing.T) {
	t.Parallel()

	expectedCommandString := strings.Join(
		[]string{
			"az",
			"acr",
			"repository",
			"show",
			"--name",
			"containerregistryelvia",
			"--image",
			"core/demo-api:v12",
			"--query",
			"digest",
			"--output",
			"tsv",
		},
		" ",
	)

	actualCommand := azAcrImageDigestCommand(
		"containerregistryelvia",
		"core/demo-api:v12",
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
