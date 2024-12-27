package build

import (
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/utils"
)

type DockerfileVariablesPython struct {
	PythonVersion string
}

func generateDockerfileForPython(
	projectFile string,
	directory string,
	options GenerateDockerfileOptions,
) (string, string, error) {
	_, buildContext := getProjectFileAndBuildContext(
		projectFile,
		options.BuildContext,
	)

	pythonVersion := getPythonVersion(
		path.Dir(projectFile),
		buildContext,
	)

	dockerfileVariables := DockerfileVariablesPython{
		PythonVersion: pythonVersion,
	}

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

func getPythonVersion(directories ...string) string {
	const defaultPythonVersion = "3.13"

	// removes duplicates
	slices.Sort(directories)
	directories = slices.Compact(directories)

	style.PrintInfo(
		fmt.Sprintf(
			"Looking for .python-version file in directories: '%v'.",
			strings.Join(directories, ", "),
		),
	)

	for _, directory := range directories {
		versionFile := path.Join(directory, ".python-version")

		if _, err := os.Stat(versionFile); os.IsNotExist(err) {
			style.PrintWarning(
				fmt.Sprintf(
					"No .python-version file found in '%s', will try next directory.\n",
					directory,
				),
			)

			continue
		}

		file, err := os.Open(versionFile)
		if err != nil {
			style.PrintWarning(
				fmt.Sprintf(
					"Failed to open .python-version file in '%s', will try next directory.\n",
					directory,
				),
			)

			continue
		}

		defer file.Close()

		contents, err := io.ReadAll(file)
		if err != nil {
			style.PrintWarning(
				fmt.Sprintf(
					"Failed to read .python-version file in '%s', will try next directory.\n",
					directory,
				),
			)

			continue
		}

		pythonVersion := strings.TrimSpace(string(contents))

		style.PrintInfo(
			fmt.Sprintf(
				"Found .python-version file in '%s' with version %s.\n",
				directory,
				pythonVersion,
			),
		)

		return pythonVersion
	}

	style.PrintWarning(
		fmt.Sprintf(
			"Did not find any .python-version files, using default version %s.",
			defaultPythonVersion,
		),
	)

	return defaultPythonVersion
}
