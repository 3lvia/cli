package build

import (
	"bytes"
	"embed"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed *.tmpl*
var templates embed.FS

type GenerateDockerfileOptions struct {
	GoMainPackageDirectory string
	BuildContext           string
	IncludeFiles           []string
	IncludeDirectories     []string
}

func generateDockerfile(
	projectFile string,
	applicationName string,
	options GenerateDockerfileOptions,
) (string, string, error) {
	directory, err := os.MkdirTemp("", "3lv-build")
	if err != nil {
		return "", "", fmt.Errorf("Failed to create temporary directory: %s", err)
	}

	if strings.HasSuffix(projectFile, ".csproj") {
		dockerfile, buildContext, err := generateDockerfileForDotNet(
			projectFile,
			directory,
			options,
		)
		if err != nil {
			return "", "", fmt.Errorf("Failed to generate Dockerfile for .NET project: %s", err)
		}

		return dockerfile, buildContext, nil
	} else if strings.HasSuffix(projectFile, ".mod") {
		dockerfile, buildContext, err := generateDockerfileForGo(
			projectFile,
			applicationName,
			directory,
			options,
		)
		if err != nil {
			return "", "", fmt.Errorf("Failed to generate Dockerfile for Go project: %s", err)
		}

		return dockerfile, buildContext, nil
	} else if strings.HasPrefix(projectFile, "Dockerfile") || strings.HasSuffix(projectFile, "Dockerfile") {
		return projectFile, path.Dir(projectFile), nil
	} else {
		return "", "", fmt.Errorf("Unsupported project file: %s", projectFile)
	}
}

type DockerfileVariablesDotnet struct {
	CsprojFile       string // required
	AssemblyName     string // required
	BaseImageTag     string // required
	RuntimeBaseImage string // required
}

func generateDockerfileForDotNet(
	projectFile string,
	directory string,
	options GenerateDockerfileOptions,
) (string, string, error) {
	csprojFileName, buildContext := getProjectFileAndBuildContext(
		projectFile,
		options.BuildContext,
	)

	assemblyName, err := findAssemblyName(
		projectFile,
		csprojFileName,
	)
	if err != nil {
		return "", "", err
	}

	baseImageTag, err := findBaseImageTag(projectFile)
	if err != nil {
		return "", "", err
	}

	runtimeBaseImage, err := findRuntimeBaseImage(projectFile)
	if err != nil {
		return "", "", err
	}

	const templateFile = "Dockerfile.dotnet.tmpl"

	dockerfilePath, err := writeDockerfile(
		directory,
		templateFile,
		DockerfileVariablesDotnet{
			CsprojFile:       csprojFileName,
			AssemblyName:     assemblyName,
			BaseImageTag:     baseImageTag,
			RuntimeBaseImage: runtimeBaseImage,
		},
	)
	if err != nil {
		return "", "", err
	}

	return dockerfilePath, buildContext, nil
}

type DockerfileVariablesGo struct {
	ModuleDirectory      string // required
	BuildContext         string // required
	MainPackageDirectory string // required
	IncludeFiles         []string
	IncludeDirectories   []string
}

func generateDockerfileForGo(
	projectFile string,
	applicationName string,
	dir string,
	options GenerateDockerfileOptions,
) (string, string, error) {
	moduleDirectory, buildContext := getModuleDirectoryAndBuildContext(
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
		ModuleDirectory:      moduleDirectory,
		BuildContext:         buildContext,
		MainPackageDirectory: mainPackageDirectory,
		IncludeFiles:         options.IncludeFiles,
		IncludeDirectories:   options.IncludeDirectories,
	}

	const templateFile = "Dockerfile.go.tmpl"
	dockerfilePath, err := writeDockerfile(dir, templateFile, dockerfileVariables)
	if err != nil {
		return "", "", err
	}

	return dockerfilePath, buildContext, nil
}

func writeDockerfile(
	dir string,
	templateFile string,
	dockerfileVariables any,
) (string, error) {
	dockerfilePath := path.Join(dir, "Dockerfile")
	dockerfile, err := os.Create(dockerfilePath)
	if err != nil {
		return "", fmt.Errorf("Failed to create Dockerfile: %s", err)
	}

	defer dockerfile.Close()

	dockerfileTemplate, err := template.New(templateFile).ParseFS(templates, templateFile)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Dockerfile template: %s", err)
	}

	var dockerfileBuffer bytes.Buffer
	err = dockerfileTemplate.Execute(&dockerfileBuffer, dockerfileVariables)
	if err != nil {
		return "", fmt.Errorf("Failed to execute Dockerfile template: %s", err)
	}

	if _, err := dockerfile.Write(dockerfileBuffer.Bytes()); err != nil {
		return "", fmt.Errorf("Failed to write Dockerfile: %s", err)
	}

	return dockerfilePath, nil
}

type CSharpProjectFile struct {
	XMLName       xml.Name `xml:"Project"`
	SDK           string   `xml:"Sdk,attr"`
	PropertyGroup PropertyGroup
}

type PropertyGroup struct {
	AssemblyName    string
	TargetFramework string
}

func getXMLFromFile(fileName string) (*CSharpProjectFile, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("getXMLFromFile: Failed to open file: %s", err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("getXMLFromFile: Failed to read file: %s", err)
	}

	var project CSharpProjectFile
	err = xml.Unmarshal(bytes, &project)
	if err != nil {
		return nil, fmt.Errorf("getXMLFromFile: Failed to unmarshal file: %s", err)
	}

	return &project, nil
}

func findAssemblyName(
	csprojFileRelativePath string,
	csprojFileName string,
) (string, error) {
	var assemblyName string
	csprojXml, err := getXMLFromFile(csprojFileRelativePath)
	if err != nil {
		return "", err
	}

	assemblyName = csprojXml.PropertyGroup.AssemblyName

	if len(assemblyName) == 0 {
		basename := filepath.Base(csprojFileName)
		withoutExtension := strings.TrimSuffix(basename, filepath.Ext(basename))

		return withoutExtension + ".dll", nil
	}

	return assemblyName + ".dll", nil
}

func findBaseImageTag(csprojFileRelativePath string) (string, error) {
	csprojXml, err := getXMLFromFile(csprojFileRelativePath)
	if err != nil {
		return "", err
	}

	targetFramework := csprojXml.PropertyGroup.TargetFramework
	if len(targetFramework) == 0 {
		return "", fmt.Errorf(
			"findBaseImageTag: TargetFramework not found in csproj file: %s",
			csprojFileRelativePath,
		)
	}

	return targetFramework[3:] + "-alpine", nil

}

func findRuntimeBaseImage(csprojFileRelativePath string) (string, error) {
	csprojXml, err := getXMLFromFile(csprojFileRelativePath)
	if err != nil {
		return "", err
	}

	sdk := csprojXml.SDK
	if len(sdk) == 0 {
		return "", fmt.Errorf(
			"SDK not found in csproj file: %s",
			csprojFileRelativePath,
		)
	}

	switch sdk {
	case "Microsoft.NET.Sdk":
		return "mcr.microsoft.com/dotnet/runtime", nil
	case "Microsoft.NET.Sdk.Web",
		"Microsoft.NET.Sdk.BlazorWebAssembly",
		"Microsoft.NET.Sdk.Razor",
		"Microsoft.NET.Sdk.Worker":
		return "mcr.microsoft.com/dotnet/aspnet", nil
	default:
		return "", fmt.Errorf("Unknown SDK: %s", sdk)
	}
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

func getModuleDirectoryAndBuildContext(
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
