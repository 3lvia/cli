package create

import (
	"fmt"
	"path"

	"github.com/orsinium-labs/enum"
)

type Language enum.Member[string]

var (
	Dotnet = Language{"dotnet"} // not technically a language, but it fits our naming scheme
	Go     = Language{"go"}     //nolint:varnamelen
	Python = Language{"python"}

	Languages = enum.New(
		Dotnet,
		Go,
		Python,
	)
)

type Template enum.Member[string]

var (
	Dotnet8WebAPI = Template{"dotnet8-webapi"}
	// Dotnet8WebApp = Template{"dotnet8-webapp"}.
	Dotnet8Worker = Template{"dotnet8-worker"}
	GoWebAPI      = Template{"go-webapi"}
	PythonWebAPI  = Template{"python-webapi"}
	PythonWorker  = Template{"python-worker"}

	Templates = enum.New(
		Dotnet8WebAPI,
		//	Dotnet8WebApp,
		Dotnet8Worker,
		GoWebAPI,
		PythonWebAPI,
		PythonWorker,
	)
)

func (template Template) getLanguage() Language {
	switch template {
	case Dotnet8WebAPI /*Dotnet8WebApp,*/, Dotnet8Worker:
		return Dotnet
	case GoWebAPI:
		return Go
	case PythonWebAPI, PythonWorker:
		return Python
	default: // should never happen
		return Language{"error"}
	}
}

func (template Template) getProjectDirectory(
	outputDirectory string,
	applicationName string,
) string {
	switch template.getLanguage() {
	case Dotnet:
		return path.Join(
			outputDirectory,
			toPascalCaseWithoutHyphens(applicationName),
		)
	default:
		return path.Join(outputDirectory, applicationName)
	}
}

func (template Template) getProjectFile(
	applicationName string,
) (string, error) {
	switch template.getLanguage() {
	case Dotnet:
		return toPascalCaseWithoutHyphens(applicationName) + ".csproj", nil
	case Go:
		return "go.mod", nil
	case Python:
		return "pyproject.toml", nil
	default:
		return "", fmt.Errorf("Could not find project file for template '%s'", template)
	}
}
