package build

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetProjectFileAndBuildContext1(t *testing.T) {
	const expectedCsprojFileName = "demo-api.csproj"
	const expectedBuildContext = "."

	csprojFileName, buildContext := getProjectFileAndBuildContext(
		"demo-api.csproj",
		"",
	)

	if expectedCsprojFileName != csprojFileName {
		t.Errorf("Csproj file name mismatch: expected %s, got %s", expectedCsprojFileName, csprojFileName)
	}

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileAndBuildContext2(t *testing.T) {
	const expectedCsprojFileName = "demo-api.csproj"
	const expectedBuildContext = "src/Things/DemoApi"

	csprojFileName, buildContext := getProjectFileAndBuildContext(
		"demo-api.csproj",
		"src/Things/DemoApi",
	)

	if expectedCsprojFileName != csprojFileName {
		t.Errorf("Csproj file name mismatch: expected %s, got %s", expectedCsprojFileName, csprojFileName)
	}

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileAndBuildContext3(t *testing.T) {
	const expectedCsprojFileName = "demo-api.csproj"
	const expectedBuildContext = "src/Things/DemoApi"

	csprojFileName, buildContext := getProjectFileAndBuildContext(
		"src/Things/DemoApi/demo-api.csproj",
		"",
	)

	if expectedCsprojFileName != csprojFileName {
		t.Errorf("Csproj file name mismatch: expected %s, got %s", expectedCsprojFileName, csprojFileName)
	}

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileAndBuildContext4(t *testing.T) {
	const expectedCsprojFileName = "DemoApi/demo-api.csproj"
	const expectedBuildContext = "src/Things"

	csprojFileName, buildContext := getProjectFileAndBuildContext(
		"src/Things/DemoApi/demo-api.csproj",
		"src/Things",
	)

	if expectedCsprojFileName != csprojFileName {
		t.Errorf("Csproj file name mismatch: expected %s, got %s", expectedCsprojFileName, csprojFileName)
	}

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileAndBuildContext5(t *testing.T) {
	const expectedCsprojFileName = "Things/DemoApi/demo-api.csproj"
	const expectedBuildContext = "src"

	csprojFileName, buildContext := getProjectFileAndBuildContext(
		"src/Things/DemoApi/demo-api.csproj",
		"src",
	)

	if expectedCsprojFileName != csprojFileName {
		t.Errorf("Csproj file name mismatch: expected %s, got %s", expectedCsprojFileName, csprojFileName)
	}

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestFindAssemblyName1(t *testing.T) {
	const projectFile = "_test/no-assembly-name.csproj"
	const expectedAssemblyName = "no-assembly-name.dll"

	actualAssemblyName, err := findAssemblyName(
		projectFile,
		filepath.Base(projectFile),
	)
	if err != nil {
		t.Errorf("Error finding assembly name: %v", err)
	}

	if expectedAssemblyName != actualAssemblyName {
		t.Errorf("Assembly name mismatch: expected %s, got %s", expectedAssemblyName, actualAssemblyName)
	}
}

func TestFindAssemblyName2(t *testing.T) {
	const projectFile = "_test/assembly-name.csproj"
	const expectedAssemblyName = "SelfDefinedAssemblyName.dll"

	actualAssemblyName, err := findAssemblyName(
		projectFile,
		filepath.Base(projectFile),
	)
	if err != nil {
		t.Errorf("Error finding assembly name: %v", err)
	}

	if expectedAssemblyName != actualAssemblyName {
		t.Errorf("Assembly name mismatch: expected %s, got %s", expectedAssemblyName, actualAssemblyName)
	}
}

func TestFindRuntimeBaseImage1(t *testing.T) {
	const projectFile = "_test/no-sdk.csproj"

	_, err := findRuntimeBaseImage(projectFile)
	if err == nil {
		t.Errorf("Expected error finding runtime base image")
	}
}

func TestFindRuntimeBaseImage2(t *testing.T) {
	const projectFile = "_test/dotnet-runtime.csproj"
	const expectedBaseImage = "mcr.microsoft.com/dotnet/runtime"

	actualBaseImage, err := findRuntimeBaseImage(projectFile)
	if err != nil {
		t.Errorf("Error finding runtime base image: %v", err)
	}

	if expectedBaseImage != actualBaseImage {
		t.Errorf("Base image mismatch: expected %s, got %s", expectedBaseImage, actualBaseImage)
	}
}

func TestFindRuntimeBaseImage3(t *testing.T) {
	const expectedBaseImage = "mcr.microsoft.com/dotnet/aspnet"
	projectFiles := []string{
		"_test/dotnet-aspnet-web.csproj",
		"_test/dotnet-aspnet-blazor-web-assembly.csproj",
		"_test/dotnet-aspnet-razor.csproj",
		"_test/dotnet-aspnet-worker.csproj",
	}

	for _, projectFile := range projectFiles {
		actualBaseImage, err := findRuntimeBaseImage(projectFile)
		if err != nil {
			t.Errorf("Error finding runtime base image: %v", err)
		}

		if expectedBaseImage != actualBaseImage {
			t.Errorf("Base image mismatch: expected %s, got %s", expectedBaseImage, actualBaseImage)
		}
	}
}

func TestFindBaseImageTag1(t *testing.T) {
	const projectFile = "_test/this-does-not-exists.csproj"

	_, err := findBaseImageTag(projectFile)
	if err == nil {
		t.Errorf("Expected error finding base image tag")
	}
}

func TestFindBaseImageTag2(t *testing.T) {
	const projectFile = "_test/no-target-framework.csproj"

	_, err := findBaseImageTag(projectFile)
	if err == nil {
		t.Errorf("Expected error finding base image tag")
	}
}

func TestFindBaseImageTag3(t *testing.T) {
	dotnetFrameworkVersions := []string{
		"6.0",
		"7.0",
		"8.0",
		"9.0",
	}

	for _, frameworkVersion := range dotnetFrameworkVersions {
		projectFile := "_test/dotnet-" + frameworkVersion + ".csproj"
		expectedBaseImageTag := frameworkVersion + "-alpine"

		actualBaseImageTag, err := findBaseImageTag(projectFile)
		if err != nil {
			t.Errorf("Error finding base image tag: %v", err)
		}

		if expectedBaseImageTag != actualBaseImageTag {
			t.Errorf("Base image tag mismatch: expected %s, got %s", expectedBaseImageTag, actualBaseImageTag)
		}
	}
}

func TestGenerateGoDockerfile1(t *testing.T) {
	expectedDockerfile, err := os.ReadFile("_test/Dockerfile.test1")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	expectedBuildContext := "."

	const projectFile = "go.mod"
	const applicationName = "demo-api"

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
		applicationName,
		GenerateDockerfileOptions{},
	)
	if err != nil {
		t.Errorf("Error generating Dockerfile: %v", err)
	}

	actualDockerfile, err := os.ReadFile(actualDockerfilePath)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	if string(expectedDockerfile) != string(actualDockerfile) {
		t.Errorf("Dockerfile mismatch: expected %s, got %s", expectedDockerfile, actualDockerfile)
	}

	if expectedBuildContext != actualBuildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, actualBuildContext)
	}
}

func TestGenerateDockerfileWithDockerfile1(t *testing.T) {
	const projectFile = "Dockerfile.test" // doesn't need to exist
	const expectedBuildContext = "."

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
		"",
		GenerateDockerfileOptions{},
	)
	if err != nil {
		t.Errorf("Error generating Dockerfile: %v", err)
	}

	if projectFile != actualDockerfilePath {
		t.Errorf("Dockerfile path mismatch: expected %s, got %s", projectFile, actualDockerfilePath)
	}

	if expectedBuildContext != actualBuildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, actualBuildContext)
	}
}

func TestGenerateDockerfileWithDockerfile2(t *testing.T) {
	const projectFile = "src/Project/Dockerfile.test" // doesn't need to exist
	const expectedBuildContext = "src/Project"

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
		"",
		GenerateDockerfileOptions{},
	)
	if err != nil {
		t.Errorf("Error generating Dockerfile: %v", err)
	}

	if projectFile != actualDockerfilePath {
		t.Errorf("Dockerfile path mismatch: expected %s, got %s", projectFile, actualDockerfilePath)
	}

	if expectedBuildContext != actualBuildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, actualBuildContext)
	}
}

func TestGenerateDockerfileWithDockerfile3(t *testing.T) {
	const projectFile = "src/Project/Dockerfile.test" // doesn't need to exist
	const expectedBuildContext = "src/Project"

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
		expectedBuildContext,
		GenerateDockerfileOptions{},
	)
	if err != nil {
		t.Errorf("Error generating Dockerfile: %v", err)
	}

	if projectFile != actualDockerfilePath {
		t.Errorf("Dockerfile path mismatch: expected %s, got %s", projectFile, actualDockerfilePath)
	}

	if expectedBuildContext != actualBuildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, actualBuildContext)
	}
}

func TestDotIfEmpty1(t *testing.T) {
	const expected = "value"
	actual := dotIfEmpty("value")

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestDotIfEmpty2(t *testing.T) {
	const expected = "."
	actual := dotIfEmpty("")

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestDotIfEmpty3(t *testing.T) {
	const expected = "."

	type testStruct struct {
		value string
	}
	test := testStruct{}

	actual := dotIfEmpty(test.value)

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
