package create

import "testing"

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

			actual := template.getProjectDirectory(
				outputDirectory,
				applicationName,
			)

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

			actual := template.getProjectDirectory(
				outputDirectory,
				applicationName,
			)

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

		actual, err := template.getProjectFile(
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

	actual, err := GoWebAPI.getProjectFile(
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

	actual, err := PythonWebAPI.getProjectFile(
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

func TestGetLanguage(t *testing.T) {
	t.Parallel()

	for _, template := range Templates.Members() {
		actual := template.getLanguage()

		switch template {
		case Dotnet8WebAPI, Dotnet8Worker:
			if actual != Dotnet {
				t.Errorf(
					"expected %q, got %q",
					Dotnet,
					actual,
				)
			}
		case GoWebAPI:
			if actual != Go {
				t.Errorf(
					"expected %q, got %q",
					Go,
					actual,
				)
			}
		case PythonWebAPI, PythonWorker:
			if actual != Python {
				t.Errorf(
					"expected %q, got %q",
					Python,
					actual,
				)
			}
		default:
			t.Errorf("unexpected template: %s", template)
		}
	}
}
