package deploy

import "testing"

func TestFormatDeploymentMessage1(t *testing.T) {
	repositoryName := "test-repo"
	commitMessage := "Did something very important"
	options := &FormatDeploymentMessageOptions{
		RunID: "123",
	}

	expected := "Deployed from GitHub Actions run " + options.RunID + " - " + commitMessage + " - " + repositoryName + " - <a href=\"https://github.com/" + repositoryName + "/actions/runs/" + options.RunID + "\">Link</a>"
	actual := formatDeploymentMessage(repositoryName, commitMessage, options)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatDeploymentMessage2(t *testing.T) {
	repositoryName := "core-very-important-repo"
	commitMessage := "Did something very important"
	options := &FormatDeploymentMessageOptions{
		RunID: "1237757570",
	}

	expected := "Deployed from GitHub Actions run " + options.RunID + " - " + commitMessage + " - " + repositoryName + " - <a href=\"https://github.com/" + repositoryName + "/actions/runs/" + options.RunID + "\">Link</a>"
	actual := formatDeploymentMessage(repositoryName, commitMessage, options)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatDeploymentMessage3(t *testing.T) {
	repositoryName := "core-very-important-repo"
	commitMessage := "Did something very important"
	options := &FormatDeploymentMessageOptions{
		RunID: "",
	}

	expected := "Manually deployed with CLI - " + commitMessage + " - " + repositoryName
	actual := formatDeploymentMessage(repositoryName, commitMessage, options)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatDeploymentMessage4(t *testing.T) {
	repositoryName := "core-not-important-repo"
	commitMessage := "Not very pretty"

	expected := "Manually deployed with CLI - " + commitMessage + " - " + repositoryName
	actual := formatDeploymentMessage(repositoryName, commitMessage, nil)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
