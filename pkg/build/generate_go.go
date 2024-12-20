package build

import (
	"strings"

	"github.com/3lvia/cli/pkg/utils"
)

type DockerfileVariablesGo struct {
	GoModuleDirectory    string // required
	MainPackageDirectory string // required
}

func generateDockerfileForGo(
	projectFile string,
	applicationName string,
	directory string,
	options GenerateDockerfileOptions,
) (string, string, error) {
	goModuleDirectory, buildContext := getGoModuleDirectoryAndBuildContext(
		projectFile,
		options.BuildContext,
	)

	mainPackageDirectory := func() string {
		if options.GoMainPackageDirectory == "" {
			return "./cmd/" + applicationName
		}

		return options.GoMainPackageDirectory
	}()

	dockerfileVariables := DockerfileVariablesGo{
		GoModuleDirectory:    goModuleDirectory,
		MainPackageDirectory: mainPackageDirectory,
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
		return "", "", err
	}

	return dockerfilePath, buildContext, nil
}

func getGoModuleDirectoryAndBuildContext(
	projectFileRelativePath string,
	buildContextRelativePath string,
) (string, string) {
	projectFileName, buildContext := getProjectFileAndBuildContext(
		projectFileRelativePath,
		buildContextRelativePath,
	)

	return dotIfEmpty(
		strings.TrimSuffix(
			strings.TrimSuffix(
				projectFileName,
				"go.mod",
			),
			"/",
		),
	), buildContext
}

func dotIfEmpty(str string) string {
	if len(str) == 0 {
		return "."
	}

	return str
}
