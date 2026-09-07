package build

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestGetProjectFileRelativeToBuildContextAndGetBuildContextFromProjectFile(t *testing.T) {
	t.Parallel()

	const (
		expectedProjectFile  = "demo-api.csproj"
		expectedBuildContext = "."

		projectFileInput  = "demo-api.csproj"
		buildContextInput = ""
	)

	projectFile := getProjectFilePathRelativeToBuildContext(
		projectFileInput,
		buildContextInput,
	)

	if expectedProjectFile != projectFile {
		t.Errorf("Project file mismatch: expected %s, got %s", expectedProjectFile, projectFile)
	}

	buildContext := getBuildContextFromProjectFile(
		projectFileInput,
		buildContextInput,
	)

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileRelativeToBuildContextAndGetBuildContextFromProjectFile2(t *testing.T) {
	t.Parallel()

	const (
		expectedProjectFile  = "demo-api.csproj"
		expectedBuildContext = "src/Things/DemoApi"

		projectFileInput  = "demo-api.csproj"
		buildContextInput = "src/Things/DemoApi"
	)

	projectFile := getProjectFilePathRelativeToBuildContext(
		projectFileInput,
		buildContextInput,
	)

	if expectedProjectFile != projectFile {
		t.Errorf("Project file mismatch: expected %s, got %s", expectedProjectFile, projectFile)
	}

	buildContext := getBuildContextFromProjectFile(
		projectFileInput,
		buildContextInput,
	)

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileRelativeToBuildContextAndGetBuildContextFromProjectFile3(t *testing.T) {
	t.Parallel()

	const (
		expectedProjectFile  = "demo-api.csproj"
		expectedBuildContext = "src/Things/DemoApi"

		projectFileInput  = "src/Things/DemoApi/demo-api.csproj"
		buildContextInput = ""
	)

	projectFile := getProjectFilePathRelativeToBuildContext(
		projectFileInput,
		buildContextInput,
	)

	if expectedProjectFile != projectFile {
		t.Errorf("Project file mismatch: expected %s, got %s", expectedProjectFile, projectFile)
	}

	buildContext := getBuildContextFromProjectFile(
		projectFileInput,
		buildContextInput,
	)

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileRelativeToBuildContextAndGetBuildContextFromProjectFile4(t *testing.T) {
	t.Parallel()

	const (
		expectedProjectFile  = "DemoApi/demo-api.csproj"
		expectedBuildContext = "src/Things"

		projectFileInput  = "src/Things/DemoApi/demo-api.csproj"
		buildContextInput = "src/Things"
	)

	projectFile := getProjectFilePathRelativeToBuildContext(
		projectFileInput,
		buildContextInput,
	)

	if expectedProjectFile != projectFile {
		t.Errorf("Project file mismatch: expected %s, got %s", expectedProjectFile, projectFile)
	}

	buildContext := getBuildContextFromProjectFile(
		projectFileInput,
		buildContextInput,
	)

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileRelativeToBuildContextAndGetBuildContextFromProjectFile5(t *testing.T) {
	t.Parallel()

	const (
		expectedProjectFile  = "DemoApi/demo-api.csproj"
		expectedBuildContext = "src/Things"

		projectFileInput  = "src/Things/DemoApi/demo-api.csproj"
		buildContextInput = "src/Things/"
	)

	projectFile := getProjectFilePathRelativeToBuildContext(
		projectFileInput,
		buildContextInput,
	)

	if expectedProjectFile != projectFile {
		t.Errorf("Project file mismatch: expected %s, got %s", expectedProjectFile, projectFile)
	}

	buildContext := getBuildContextFromProjectFile(
		projectFileInput,
		buildContextInput,
	)

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestGetProjectFileAndBuildContext6(t *testing.T) {
	t.Parallel()

	const (
		expectedProjectFile  = "Things/DemoApi/demo-api.csproj"
		expectedBuildContext = "src"

		projectFileInput  = "src/Things/DemoApi/demo-api.csproj"
		buildContextInput = "src"
	)

	projectFile := getProjectFilePathRelativeToBuildContext(
		projectFileInput,
		buildContextInput,
	)

	if expectedProjectFile != projectFile {
		t.Errorf("Project file mismatch: expected %s, got %s", expectedProjectFile, projectFile)
	}

	buildContext := getBuildContextFromProjectFile(
		projectFileInput,
		buildContextInput,
	)

	if expectedBuildContext != buildContext {
		t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, buildContext)
	}
}

func TestFindAssemblyName1(t *testing.T) {
	t.Parallel()

	const (
		projectFile          = "_test/no-assembly-name.csproj"
		expectedAssemblyName = "no-assembly-name.dll"
	)

	csprojXML, err := getXMLFromFile(projectFile)
	if err != nil {
		t.Errorf("Error reading XML from file: %v", err)
	}

	actualAssemblyName, err := findAssemblyName(
		projectFile,
		csprojXML,
	)
	if err != nil {
		t.Errorf("Error finding assembly name: %v", err)
	}

	if expectedAssemblyName != actualAssemblyName {
		t.Errorf("Assembly name mismatch: expected %s, got %s", expectedAssemblyName, actualAssemblyName)
	}
}

func TestFindAssemblyName2(t *testing.T) {
	t.Parallel()

	const (
		projectFile          = "_test/assembly-name.csproj"
		expectedAssemblyName = "SelfDefinedAssemblyName.dll"
	)

	csprojXML, err := getXMLFromFile(projectFile)
	if err != nil {
		t.Errorf("Error reading XML from file: %v", err)
	}

	actualAssemblyName, err := findAssemblyName(
		projectFile,
		csprojXML,
	)

	if expectedAssemblyName != actualAssemblyName {
		t.Errorf("Assembly name mismatch: expected %s, got %s", expectedAssemblyName, actualAssemblyName)
	}
}

func TestFindRuntimeBaseImage1(t *testing.T) {
	t.Parallel()

	const projectFile = "_test/no-sdk.csproj"

	csprojXML, err := getXMLFromFile(projectFile)
	if err != nil {
		t.Errorf("Error reading XML from file: %v", err)
	}

	_, err = findRuntimeBaseImage(csprojXML)
	if err == nil {
		t.Errorf("Expected error finding runtime base image")
	}
}

func TestFindRuntimeBaseImage2(t *testing.T) {
	t.Parallel()

	const (
		projectFile       = "_test/dotnet-runtime.csproj"
		expectedBaseImage = "mcr.microsoft.com/dotnet/runtime"
	)

	csprojXML, err := getXMLFromFile(projectFile)
	if err != nil {
		t.Errorf("Error reading XML from file: %v", err)
	}

	actualBaseImage, err := findRuntimeBaseImage(csprojXML)
	if err != nil {
		t.Errorf("Error finding runtime base image: %v", err)
	}

	if expectedBaseImage != actualBaseImage {
		t.Errorf("Base image mismatch: expected %s, got %s", expectedBaseImage, actualBaseImage)
	}
}

func TestFindRuntimeBaseImage3(t *testing.T) {
	t.Parallel()

	const expectedBaseImage = "mcr.microsoft.com/dotnet/aspnet"

	projectFiles := []string{
		"_test/dotnet-aspnet-web.csproj",
		"_test/dotnet-aspnet-blazor-web-assembly.csproj",
		"_test/dotnet-aspnet-razor.csproj",
		"_test/dotnet-aspnet-worker.csproj",
	}

	for _, projectFile := range projectFiles {
		csprojXML, err := getXMLFromFile(projectFile)
		if err != nil {
			t.Errorf("Error reading XML from file: %v", err)
		}

		actualBaseImage, err := findRuntimeBaseImage(csprojXML)
		if err != nil {
			t.Errorf("Error finding runtime base image: %v", err)
		}

		if expectedBaseImage != actualBaseImage {
			t.Errorf("Base image mismatch: expected %s, got %s", expectedBaseImage, actualBaseImage)
		}
	}
}

func TestFindBaseImageTag1(t *testing.T) {
	t.Parallel()

	const projectFile = "_test/this-does-not-exists.csproj"

	_, err := getXMLFromFile(projectFile)
	if err == nil {
		t.Errorf("Expected error reading XML from file")
	}
}

func TestFindBaseImageTag2(t *testing.T) {
	t.Parallel()

	const projectFile = "_test/no-target-framework.csproj"

	csprojXML, err := getXMLFromFile(projectFile)
	if err == nil {
		t.Errorf("Expected error reading XML from file")
	}

	_, err = findBaseImageTag(csprojXML)
	if err == nil {
		t.Errorf("Expected error finding base image tag")
	}
}

func TestFindBaseImageTag3(t *testing.T) {
	t.Parallel()

	dotnetFrameworkVersions := []string{
		"6.0",
		"7.0",
		"8.0",
		"9.0",
	}

	for _, frameworkVersion := range dotnetFrameworkVersions {
		projectFile := "_test/dotnet-" + frameworkVersion + ".csproj"
		expectedBaseImageTag := frameworkVersion + "-alpine"

		csprojXML, err := getXMLFromFile(projectFile)
		if err != nil {
			t.Errorf("Error reading XML from file: %v", err)
		}

		actualBaseImageTag, err := findBaseImageTag(csprojXML)
		if err != nil {
			t.Errorf("Error finding base image tag: %v", err)
		}

		if expectedBaseImageTag != actualBaseImageTag {
			t.Errorf("Base image tag mismatch: expected %s, got %s", expectedBaseImageTag, actualBaseImageTag)
		}
	}
}

func TestGenerateDotnetDockerfile(t *testing.T) {
	t.Parallel()

	expectedDockerfile, err := os.ReadFile("_test/Dockerfile.dotnet.test")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	const (
		projectFile          = "_test/dotnet-8.0.csproj"
		expectedBuildContext = "_test"
	)

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
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

func TestGenerateGoDockerfile(t *testing.T) {
	t.Parallel()

	expectedDockerfile, err := os.ReadFile("_test/Dockerfile.go.test")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	const (
		projectFile          = "go.mod"
		expectedBuildContext = "."
	)

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
		GenerateDockerfileOptions{
			GoMainPackageDirectory: "./cmd/demo-api-go",
		},
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

func TestGeneratePythonDockerfile1(t *testing.T) {
	t.Parallel()

	expectedDockerfile, err := os.ReadFile("_test/Dockerfile.python.test1")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	const (
		projectFile          = "pyproject.toml"
		expectedBuildContext = "."
		applicationName      = "demo-api-python"
	)

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
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

func TestGeneratePythonDockerfile2(t *testing.T) {
	t.Parallel()

	for i, version := range []string{"3.12", "3.10"} {
		expectedDockerfile, err := os.ReadFile("_test/Dockerfile.python.test" + strconv.FormatInt(int64(i+2), 10))
		if err != nil {
			t.Errorf("Error reading expected Dockerfile: %v", err)
		}

		tempDir := t.TempDir()
		projectFile := filepath.Join(tempDir, "pyproject.toml")
		expectedBuildContext := tempDir

		file, err := os.Create(filepath.Join(tempDir, ".python-version"))
		if err != nil {
			t.Errorf("Error creating file: %v", err)
		}

		if _, err := file.WriteString(version + "\n"); err != nil {
			t.Errorf("Error writing to file: %v", err)
		}

		actualDockerfilePath, actualBuildContext, err := generateDockerfile(
			projectFile,
			GenerateDockerfileOptions{},
		)
		if err != nil {
			t.Errorf("Error generating Dockerfile: %v", err)
		}

		actualDockerfile, err := os.ReadFile(actualDockerfilePath)
		if err != nil {
			t.Errorf("Error reading generated Dockerfile: %v", err)
		}

		if string(expectedDockerfile) != string(actualDockerfile) {
			t.Errorf("Dockerfile mismatch: expected %s, got %s", expectedDockerfile, actualDockerfile)
		}

		if expectedBuildContext != actualBuildContext {
			t.Errorf("Build context mismatch: expected %s, got %s", expectedBuildContext, actualBuildContext)
		}
	}
}

func TestGenerateDockerfileWithDockerfile1(t *testing.T) {
	t.Parallel()

	const (
		projectFile          = "Dockerfile.test" // doesn't need to exist
		expectedBuildContext = "."
	)

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
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
	t.Parallel()

	const (
		projectFile          = "src/Project/Dockerfile.test" // doesn't need to exist
		expectedBuildContext = "src/Project"
	)

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
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
	t.Parallel()

	const (
		projectFile          = "src/Project/Dockerfile.test" // doesn't need to exist
		expectedBuildContext = "src/Project"
	)

	actualDockerfilePath, actualBuildContext, err := generateDockerfile(
		projectFile,
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
