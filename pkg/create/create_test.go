package create

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/build"
	"github.com/3lvia/cli/pkg/command"
)

func TestCookiecutterCommand(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		template      Template
		pythonVersion string
	}{
		{
			name:          "demo-api-python",
			template:      PythonWebAPI,
			pythonVersion: "3.11",
		},
		{
			name:          "demo-api-python",
			template:      PythonWebAPI,
			pythonVersion: "",
		},
		{
			name:          "demo-api-go",
			template:      GoWebAPI,
			pythonVersion: "",
		},
		{
			name:          "demo-api",
			template:      Dotnet8WebAPI,
			pythonVersion: "",
		},
	}

	const (
		outputDirectory = "test-output-directory"
		applicationName = "test-application-name"
		systemName      = "test-system-name"
	)

	for _, testCase := range testCases {
		expectedCommandString := strings.Join(
			[]string{
				"cookiecutter",
				"gh:3lvia/application-templates",
				"--directory",
				testCase.template.Value,
				"--output-dir",
				outputDirectory,
				"--no-input",
				"application_name=" + applicationName,
				"application_name_pascal_case=" + toPascalCaseWithoutHyphens(applicationName),
				"system_name=" + systemName,
			},
			" ",
		)

		if testCase.template == PythonWebAPI {
			if testCase.pythonVersion == "" {
				expectedCommandString += " python_version=" + build.DefaultPythonVersion
			} else {
				expectedCommandString += " python_version=" + testCase.pythonVersion
			}
		}

		actualCommand := cookiecutterCommand(
			testCase.template,
			outputDirectory,
			applicationName,
			systemName,
			testCase.pythonVersion,
			&command.RunOptions{DryRun: true},
		)

		command.ExpectedCommandStringEqualsActualCommand(
			t,
			expectedCommandString,
			actualCommand,
		)
	}
}

func TestCheckCookiecutterInstalledCommand(t *testing.T) {
	t.Parallel()

	expectedCommandString := "cookiecutter --version"

	actualCommand := checkCookiecutterInstalledCommand(
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestCheckPipxInstalledCommand(t *testing.T) {
	t.Parallel()

	expectedCommandString := "pipx --version"

	actualCommand := checkPipxInstalledCommand(
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestInstallCookiecutterCommand(t *testing.T) {
	t.Parallel()

	expectedCommandString := "sudo pipx install cookiecutter --global"

	actualCommand := installCookiecutterCommand(
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestCheckUvInstalledCommand(t *testing.T) {
	t.Parallel()

	expectedCommandString := "uv --version"

	actualCommand := checkUvInstalledCommand(
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestInstallUvCommand(t *testing.T) {
	t.Parallel()

	expectedCommandString := "pipx install uv"

	actualCommand := installUvCommand(
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestUvSyncCommand(t *testing.T) {
	t.Parallel()

	const projectDirectory = "test-project-directory"

	expectedCommandString := strings.Join(
		[]string{
			"uv",
			"sync",
			"--directory",
			projectDirectory,
		},
		" ",
	)

	actualCommand := uvSyncCommand(
		projectDirectory,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestToPascalCaseWithoutHyphens(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "test",
			expected: "Test",
		},
		{
			input:    "test-test",
			expected: "TestTest",
		},
		{
			input:    "test_test-test",
			expected: "Test_testTest",
		},
		{
			input:    "test-test-test",
			expected: "TestTestTest",
		},
	}

	for _, testCase := range testCases {
		actual := toPascalCaseWithoutHyphens(testCase.input)

		if actual != testCase.expected {
			t.Errorf(
				"expected %q, got %q",
				testCase.expected,
				actual,
			)
		}
	}
}

func TestGetProjectDirectoryForTemplateDotnet(t *testing.T) {
	t.Parallel()

	const applicationName = "demo-api"

	for _, outputDirectory := range []string{"", "applications", "src/applications"} {
		for _, template := range []Template{Dotnet8WebAPI, Dotnet8Worker} {
			expected := func() string {
				if outputDirectory == "" {
					return toPascalCaseWithoutHyphens(applicationName)
				}

				return outputDirectory + "/" + toPascalCaseWithoutHyphens(applicationName)
			}()

			actual, err := getProjectDirectoryForTemplate(
				template,
				outputDirectory,
				applicationName,
			)
			if err != nil {
				t.Error(err)
			}

			if actual != expected {
				t.Errorf(
					"expected %q, got %q",
					expected,
					actual,
				)
			}
		}
	}
}

func TestGetProjectDirectoryForTemplateGoPython(t *testing.T) {
	t.Parallel()

	const applicationName = "demo-api"

	for _, outputDirectory := range []string{"", "applications", "src/applications"} {
		for _, template := range []Template{PythonWebAPI, GoWebAPI} {
			expected := func() string {
				if outputDirectory == "" {
					return applicationName
				}

				return outputDirectory + "/" + applicationName
			}()

			actual, err := getProjectDirectoryForTemplate(
				template,
				outputDirectory,
				applicationName,
			)
			if err != nil {
				t.Error(err)
			}

			if actual != expected {
				t.Errorf(
					"expected %q, got %q",
					expected,
					actual,
				)
			}
		}
	}
}

func TestGetProjectFileForTemplateDotnet(t *testing.T) {
	t.Parallel()

	const applicationName = "demo-api"

	for _, template := range []Template{Dotnet8WebAPI, Dotnet8Worker} {
		expected := toPascalCaseWithoutHyphens(applicationName) + ".csproj"

		actual, err := getProjectFileForTemplate(
			template,
			applicationName,
		)
		if err != nil {
			t.Error(err)
		}

		if actual != expected {
			t.Errorf(
				"expected %q, got %q",
				expected,
				actual,
			)
		}
	}
}

func TestGetProjectFileForTemplateGo(t *testing.T) {
	t.Parallel()

	const (
		applicationName = "demo-api"
		expected        = "go.mod"
	)

	actual, err := getProjectFileForTemplate(
		GoWebAPI,
		applicationName,
	)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf(
			"expected %q, got %q",
			expected,
			actual,
		)
	}
}

func TestGetProjectFileForTemplatePython(t *testing.T) {
	t.Parallel()

	const (
		applicationName = "demo-api"
		expected        = "pyproject.toml"
	)

	actual, err := getProjectFileForTemplate(
		PythonWebAPI,
		applicationName,
	)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf(
			"expected %q, got %q",
			expected,
			actual,
		)
	}
}
