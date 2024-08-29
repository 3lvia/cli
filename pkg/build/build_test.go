package build

import (
	"strings"
	"testing"
)

func TestGetRegistry1(t *testing.T) {
	const expectedRegistry = "containerregistryelvia.azurecr.io"

	actualRegistry := getRegistry("")

	if actualRegistry != expectedRegistry {
		t.Errorf("Expected %s, got %s", expectedRegistry, actualRegistry)
	}
}

func TestGetRegistry2(t *testing.T) {
	const expectedRegistry = "containerregistryelvia.azurecr.io"

	actualRegistry := getRegistry("acr")

	if actualRegistry != expectedRegistry {
		t.Errorf("Expected %s, got %s", expectedRegistry, actualRegistry)
	}
}

func TestGetRegistry3(t *testing.T) {
	const expectedRegistry = "ghcr.io/3lvia"

	actualRegistry := getRegistry("ghcr")

	if actualRegistry != expectedRegistry {
		t.Errorf("Expected %s, got %s", expectedRegistry, actualRegistry)
	}
}

func TestGetRegistry4(t *testing.T) {
	const expectedRegistry = "random"

	actualRegistry := getRegistry(expectedRegistry)

	if actualRegistry != expectedRegistry {
		t.Errorf("Expected %s, got %s", expectedRegistry, actualRegistry)
	}
}

func TestConstructBuildCommandArguments1(t *testing.T) {
	const dockerfilePath = "build/Dockerfile"
	const buildContext = "src/app"
	const imageName = "containerregistryelvia.azurecr.io/test-image"

	expectedArguments := "buildx build -f " + dockerfilePath + " -t " + imageName + ":latest " + buildContext

	actualArguments := constructBuildCommandArguments(
		dockerfilePath,
		buildContext,
		imageName,
		[]string{"latest"},
	)

	if strings.Join(actualArguments, " ") != expectedArguments {
		t.Errorf("Expected %s, got %s", expectedArguments, actualArguments)
	}
}

func TestConstructBuildCommandArguments2(t *testing.T) {
	const dockerfilePath = "Dockerfile"
	const buildContext = "."
	const imageName = "ghcr.io/test-image"

	expectedArguments := "buildx build -f " + dockerfilePath + " -t " + imageName + ":latest-cache " + buildContext

	actualArguments := constructBuildCommandArguments(
		dockerfilePath,
		buildContext,
		imageName,
		[]string{"latest-cache"},
	)

	if strings.Join(actualArguments, " ") != expectedArguments {
		t.Errorf("Expected %s, got %s", expectedArguments, actualArguments)
	}
}

func TestConstructBuildCommandArguments3(t *testing.T) {
	const dockerfilePath = "Dockerfile"
	const buildContext = "."
	const imageName = "ghcr.io/test-image"

	expectedArguments := "buildx build -f " + dockerfilePath + " -t " + imageName + ":latest-cache" + " -t " + imageName + ":latest" + " -t " + imageName + ":v42.0.1" + " -t " + imageName + ":v420alpha " + buildContext

	actualArguments := constructBuildCommandArguments(
		dockerfilePath,
		buildContext,
		imageName,
		[]string{"latest-cache", "latest", "v42.0.1", "v420alpha"},
	)

	if strings.Join(actualArguments, " ") != expectedArguments {
		t.Errorf("Expected %s, got %s", expectedArguments, actualArguments)
	}
}
