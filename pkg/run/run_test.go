package run

import (
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestDockerComposeUpCommand(t *testing.T) {
	t.Parallel()

	const composeFile = "test-compose-file"

	expectedCommandString := "docker compose -f " + composeFile + " up"

	actualCommand := dockerComposeUpCommand(
		composeFile,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
