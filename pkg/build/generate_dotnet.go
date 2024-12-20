package build

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/3lvia/cli/pkg/utils"
)

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

	dockerfilePath, err := utils.WriteFileWithTemplate(
		directory,
		"Dockerfile",
		templateFile,
		dockerfileTemplates,
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

type CSharpProjectFile struct {
	XMLName       xml.Name      `xml:"Project"`
	SDK           string        `xml:"Sdk,attr"`
	PropertyGroup PropertyGroup `xml:"PropertyGroup"`
}

type PropertyGroup struct {
	AssemblyName    string `xml:"AssemblyName"`
	TargetFramework string `xml:"TargetFramework"`
}

func getXMLFromFile(fileName string) (*CSharpProjectFile, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("getXMLFromFile: Failed to open file: %w", err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("getXMLFromFile: Failed to read file: %w", err)
	}

	var project CSharpProjectFile

	err = xml.Unmarshal(bytes, &project)
	if err != nil {
		return nil, fmt.Errorf("getXMLFromFile: Failed to unmarshal file: %w", err)
	}

	return &project, nil
}

func findAssemblyName(
	csprojFileRelativePath string,
	csprojFileName string,
) (string, error) {
	var assemblyName string

	csprojXML, err := getXMLFromFile(csprojFileRelativePath)
	if err != nil {
		return "", err
	}

	assemblyName = csprojXML.PropertyGroup.AssemblyName

	if len(assemblyName) == 0 {
		basename := filepath.Base(csprojFileName)
		withoutExtension := strings.TrimSuffix(basename, filepath.Ext(basename))

		return withoutExtension + ".dll", nil
	}

	return assemblyName + ".dll", nil
}

func findBaseImageTag(csprojFileRelativePath string) (string, error) {
	csprojXML, err := getXMLFromFile(csprojFileRelativePath)
	if err != nil {
		return "", err
	}

	targetFramework := csprojXML.PropertyGroup.TargetFramework
	if len(targetFramework) == 0 {
		return "", fmt.Errorf(
			"findBaseImageTag: TargetFramework not found in csproj file: %s",
			csprojFileRelativePath,
		)
	}

	return targetFramework[3:] + "-alpine", nil
}

func findRuntimeBaseImage(csprojFileRelativePath string) (string, error) {
	csprojXML, err := getXMLFromFile(csprojFileRelativePath)
	if err != nil {
		return "", err
	}

	sdk := csprojXML.SDK
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
