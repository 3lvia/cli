package deploy

import (
	"fmt"
	"testing"
)

const repositoryOwner = "3lvia"

func TestFormatDeploymentMessage1(t *testing.T) {
	const repositoryName = "test-repo"
	const commitMessage = "Did something very important"
	options := &FormatDeploymentMessageOptions{
		RunID: "123",
	}

	expected := fmt.Sprintf(
		"Deployed from GitHub Actions run %s - %s - %s - <a href=\"https://github.com/%s/%s/actions/runs/%s\">Link</a>",
		options.RunID,
		commitMessage,
		repositoryName,
		repositoryOwner,
		repositoryName,
		options.RunID,
	)
	actual := formatDeploymentMessage(repositoryName, commitMessage, options)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatDeploymentMessage2(t *testing.T) {
	const repositoryName = "core-very-important-repo"
	const commitMessage = "Did something very important"
	options := &FormatDeploymentMessageOptions{
		RunID: "1237757570",
	}

	expected := fmt.Sprintf(
		"Deployed from GitHub Actions run %s - %s - %s - <a href=\"https://github.com/%s/%s/actions/runs/%s\">Link</a>",
		options.RunID,
		commitMessage,
		repositoryName,
		repositoryOwner,
		repositoryName,
		options.RunID,
	)
	actual := formatDeploymentMessage(repositoryName, commitMessage, options)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatDeploymentMessage3(t *testing.T) {
	const repositoryName = "core-very-important-repo"
	const commitMessage = "Did something very important"
	options := &FormatDeploymentMessageOptions{
		RunID: "",
	}

	expected := fmt.Sprintf(
		"Manually deployed with CLI - %s - %s",
		commitMessage,
		repositoryName,
	)
	actual := formatDeploymentMessage(repositoryName, commitMessage, options)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatDeploymentMessage4(t *testing.T) {
	const repositoryName = "core-not-important-repo"
	const commitMessage = "Not very pretty"

	expected := fmt.Sprintf(
		"Manually deployed with CLI - %s - %s",
		commitMessage,
		repositoryName,
	)
	actual := formatDeploymentMessage(repositoryName, commitMessage, nil)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
