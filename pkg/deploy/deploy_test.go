package deploy

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestDockerInspectCommand(t *testing.T) {
	t.Parallel()

	const imageName = "ghcr.io/3lvia/core/demo-api"

	expectedCommandString := strings.Join(
		[]string{
			"docker",
			"inspect",
			"--format",
			"{{index .RepoDigests 0}}",
			imageName + ":latest-cache",
		},
		" ",
	)

	actualCommand := dockerInspectCommand(
		imageName,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
