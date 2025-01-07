package build

import (
	"github.com/3lvia/cli/pkg/utils"
)

type DockerfileVariablesGo struct {
	GoModuleDirectory    string
	MainPackageDirectory string
}

func generateDockerfileForGo(
	buildContext string,
	directory string,
	goMainPackageDirectory string,
) (string, error) {
	dockerfileVariables := DockerfileVariablesGo{
		GoModuleDirectory:    buildContext,
		MainPackageDirectory: goMainPackageDirectory,
	}

	const templateFile = "Dockerfile.go.tmpl"

	dockerfilePath, err := utils.WriteFileWithTemplate(
		directory,
		"Dockerfile",
		templateFile,
		dockerfileTemplates,
		dockerfileVariables,
	)
	if err != nil {
		return "", err
	}

	return dockerfilePath, nil
}
