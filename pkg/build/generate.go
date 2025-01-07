package build

import (
	"embed"
	"fmt"
	"os"
	"path"
	"strings"
)

//go:embed *.tmpl*
var dockerfileTemplates embed.FS

type GenerateDockerfileOptions struct {
	BuildContext           string
	GoMainPackageDirectory string
}

func generateDockerfile(
	projectFileRelativePath string,
	options GenerateDockerfileOptions,
) (string, string, error) {
	directory, err := os.MkdirTemp("", "3lv-build-*")
	if err != nil {
		return "", "", fmt.Errorf("Failed to create temporary directory: %w", err)
	}

	projectFile := getProjectFilePathRelativeToBuildContext(projectFileRelativePath, options.BuildContext)
	buildContext := getBuildContextFromProjectFile(projectFileRelativePath, options.BuildContext)

	projectFileBase := path.Base(projectFileRelativePath) // could have used `projectFile` here also, doesn't matter

	if strings.HasSuffix(projectFileBase, ".csproj") {
		dockerfile, err := generateDockerfileForDotNet(
			projectFile,
			buildContext,
			directory,
		)
		if err != nil {
			return "", "", fmt.Errorf("Failed to generate Dockerfile for .NET project: %w", err)
		}

		return dockerfile, buildContext, nil
	} else if projectFileBase == "go.mod" {
		dockerfile, err := generateDockerfileForGo(
			buildContext,
			directory,
			options.GoMainPackageDirectory,
		)
		if err != nil {
			return "", "", fmt.Errorf("Failed to generate Dockerfile for Go project: %w", err)
		}

		return dockerfile, buildContext, nil
	} else if projectFileBase == "pyproject.toml" {
		dockerfile, err := generateDockerfileForPython(
			projectFile,
			buildContext,
			directory,
		)
		if err != nil {
			return "", "", fmt.Errorf("Failed to generate Dockerfile for Python project: %w", err)
		}

		return dockerfile, buildContext, nil
	} else if strings.HasPrefix(projectFileBase, "Dockerfile") ||
		strings.HasSuffix(projectFileBase, "Dockerfile") ||
		strings.Contains(projectFileBase, "Dockerfile") {
		return projectFileRelativePath, buildContext, nil
	}

	return "", "", fmt.Errorf(
		"Unsupported project file: %s. If you want to use a Dockerfile directly,"+
			" ensure the name of the Dockerfile contains the string 'Dockerfile'",
		projectFileRelativePath,
	)
}

func getProjectFilePathRelativeToBuildContext(
	projectFileRelativePath string,
	buildContextRelativePath string,
) string {
	if buildContextRelativePath == "" {
		return path.Base(projectFileRelativePath)
	}

	buildContextRelativePath = strings.TrimSuffix(buildContextRelativePath, "/")

	return strings.TrimPrefix(
		projectFileRelativePath,
		buildContextRelativePath+"/",
	)
}

func getBuildContextFromProjectFile(
	projectFileRelativePath string,
	buildContextRelativePath string,
) string {
	if buildContextRelativePath == "" {
		projectFileDirectory := path.Dir(projectFileRelativePath)

		if projectFileDirectory == "" {
			return "."
		}

		return projectFileDirectory
	}

	return strings.TrimSuffix(buildContextRelativePath, "/")
}
