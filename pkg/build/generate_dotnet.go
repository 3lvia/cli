package build

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/3lvia/cli/pkg/utils"
)

type DockerfileVariablesDotnet struct {
	CsprojFile       string
	AssemblyName     string
	BaseImageTag     string
	RuntimeBaseImage string
}

func generateDockerfileForDotNet(
	projectFile string,
	buildContext string,
	directory string,
) (string, error) {
	csprojXML, err := getXMLFromFile(path.Join(buildContext, projectFile))
	if err != nil {
		return "", err
	}

	assemblyName, err := findAssemblyName(projectFile, csprojXML)
	if err != nil {
		return "", err
	}

	baseImageTag, err := findBaseImageTag(csprojXML)
	if err != nil {
		return "", err
	}

	runtimeBaseImage, err := findRuntimeBaseImage(csprojXML)
	if err != nil {
		return "", err
	}

	dockerfilePath, err := utils.WriteFileWithTemplate(
		directory,
		"Dockerfile",
		"Dockerfile.dotnet.tmpl",
		dockerfileTemplates,
		DockerfileVariablesDotnet{
			CsprojFile:       projectFile,
			AssemblyName:     assemblyName,
			BaseImageTag:     baseImageTag,
			RuntimeBaseImage: runtimeBaseImage,
		},
	)
	if err != nil {
		return "", err
	}

	return dockerfilePath, nil
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
	var project CSharpProjectFile

	file, err := os.Open(fileName)
	if err != nil {
		return &project, fmt.Errorf("Failed to open file: %w", err)
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return &project, fmt.Errorf("Failed to read file: %w", err)
	}

	err = xml.Unmarshal(bytes, &project)
	if err != nil {
		return &project, fmt.Errorf("Failed to unmarshal file: %w", err)
	}

	return &project, nil
}

func findAssemblyName(csprojFileName string, csprojXML *CSharpProjectFile) (string, error) {
	assemblyName := csprojXML.PropertyGroup.AssemblyName

	if len(assemblyName) == 0 {
		withoutExtension := strings.TrimSuffix(path.Base(csprojFileName), filepath.Ext(csprojFileName))

		return withoutExtension + ".dll", nil
	}

	return assemblyName + ".dll", nil
}

func findBaseImageTag(csprojXML *CSharpProjectFile) (string, error) {
	targetFramework := csprojXML.PropertyGroup.TargetFramework
	if len(targetFramework) == 0 {
		return "", errors.New("TargetFramework not found in .csproj-file.")
	}

	return targetFramework[3:] + "-alpine", nil
}

func findRuntimeBaseImage(csprojXML *CSharpProjectFile) (string, error) {
	sdk := csprojXML.SDK
	if len(sdk) == 0 {
		return "", errors.New("SDK not found in .csproj-file.")
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
