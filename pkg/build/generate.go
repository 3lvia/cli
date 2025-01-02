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
	GoMainPackageDirectory string
	BuildContext           string
}

func generateDockerfile(
	projectFile string,
	applicationName string,
	options GenerateDockerfileOptions,
) (string, string, error) {
	directory, err := os.MkdirTemp("", "3lv-build-*")
	if err != nil {
		return "", "", fmt.Errorf("Failed to create temporary directory: %w", err)
	}

	projectFileBase := path.Base(projectFile)

	if strings.HasSuffix(projectFileBase, ".csproj") {
		dockerfile, buildContext, err := generateDockerfileForDotNet(
			projectFile,
			directory,
			options,
		)
		if err != nil {
			return "", "", fmt.Errorf("Failed to generate Dockerfile for .NET project: %w", err)
		}

		return dockerfile, buildContext, nil
	} else if projectFileBase == "go.mod" {
		dockerfile, buildContext, err := generateDockerfileForGo(
			projectFile,
			applicationName,
			directory,
			options,
		)
		if err != nil {
			return "", "", fmt.Errorf("Failed to generate Dockerfile for Go project: %w", err)
		}

		return dockerfile, buildContext, nil
	} else if projectFileBase == "pyproject.toml" {
		dockerfile, buildContext, err := generateDockerfileForPython(
			projectFile,
			directory,
			options,
		)
		if err != nil {
			return "", "", fmt.Errorf("Failed to generate Dockerfile for Python project: %w", err)
		}

		return dockerfile, buildContext, nil
	} else if strings.HasPrefix(projectFileBase, "Dockerfile") ||
		strings.HasSuffix(projectFileBase, "Dockerfile") ||
		strings.Contains(projectFileBase, "Dockerfile") {
		if options.BuildContext == "" {
			return projectFile, path.Dir(projectFile), nil
		}

		return projectFile, options.BuildContext, nil
	}

	return "", "", fmt.Errorf(
		"Unsupported project file: %s. If you want to use a Dockerfile directly,"+
			" ensure the name of the Dockerfile contains the string 'Dockerfile'",
		projectFileBase,
	)
}

func getProjectFileAndBuildContext(
	projectFileRelativePath string,
	buildContextRelativePath string,
) (string, string) {
	if len(buildContextRelativePath) == 0 {
		return path.Base(projectFileRelativePath), path.Dir(projectFileRelativePath)
	}

	if strings.HasSuffix(buildContextRelativePath, "/") {
		return strings.TrimPrefix(
			projectFileRelativePath,
			buildContextRelativePath,
		), buildContextRelativePath
	}

	return strings.TrimPrefix(
		projectFileRelativePath,
		buildContextRelativePath+"/",
	), buildContextRelativePath
}
