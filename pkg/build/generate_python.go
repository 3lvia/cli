package build

import (
	"github.com/3lvia/cli/pkg/utils"
)

type DockerfileVariablesPython struct{}

func generateDockerfileForPython(
	projectFile string,
	directory string,
	options GenerateDockerfileOptions,
) (string, string, error) {
	_, buildContext := getProjectFileAndBuildContext(
		projectFile,
		options.BuildContext,
	)

	dockerfileVariables := DockerfileVariablesPython{}

	const templateFile = "Dockerfile.python.tmpl"

	dockerfilePath, err := utils.WriteFileWithTemplate(
		directory,
		"Dockerfile",
		templateFile,
		dockerfileTemplates,
		dockerfileVariables,
	)
	if err != nil {
		return "", "", err
	}

	return dockerfilePath, buildContext, nil
}
