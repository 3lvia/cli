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
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	getRegistry("unknown")
}

func TestConstructBuildCommandArguments1(t *testing.T) {
	const dockerfilePath = "build/Dockerfile"
	const buildContext = "src/app"
	const imageName = "containerregistryelvia.azurecr.io/test-image"
	const cacheTag = "latest"

	imageNameWithCacheTag := imageName + ":" + cacheTag

	expectedArguments := "buildx build -f " + dockerfilePath + " --cache-to type=inline --cache-from " + imageNameWithCacheTag + " -t " + imageNameWithCacheTag + " " + buildContext

	actualArguments := constructBuildCommandArguments(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		[]string{},
	)

	if strings.Join(actualArguments, " ") != expectedArguments {
		t.Errorf("Expected %s, got %s", expectedArguments, actualArguments)
	}
}

func TestConstructBuildCommandArguments2(t *testing.T) {
	const dockerfilePath = "Dockerfile"
	const buildContext = "."
	const imageName = "ghcr.io/test-image"
	const cacheTag = "latest-cache"

	imageNameWithCacheTag := imageName + ":" + cacheTag

	expectedArguments := "buildx build -f " + dockerfilePath + " --cache-to type=inline --cache-from " + imageNameWithCacheTag + " -t " + imageNameWithCacheTag + " " + buildContext

	actualArguments := constructBuildCommandArguments(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		[]string{},
	)

	if strings.Join(actualArguments, " ") != expectedArguments {
		t.Errorf("Expected %s, got %s", expectedArguments, actualArguments)
	}
}

func TestConstructBuildCommandArguments3(t *testing.T) {
	const dockerfilePath = "Dockerfile"
	const buildContext = "."
	const imageName = "ghcr.io/test-image"
	const cacheTag = "latest-cache"

	imageNameWithCacheTag := imageName + ":" + cacheTag
	additionalTags := []string{"latest", "v42.0.1", "v420alpha"}

	expectedArguments := "buildx build -f " + dockerfilePath + " --cache-to type=inline --cache-from " + imageNameWithCacheTag + " -t " + imageName + ":" + additionalTags[0] + " -t " + imageName + ":" + additionalTags[1] + " -t " + imageName + ":" + additionalTags[2] + " -t " + imageNameWithCacheTag + " " + buildContext

	actualArguments := constructBuildCommandArguments(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		additionalTags,
	)

	if strings.Join(actualArguments, " ") != expectedArguments {
		t.Errorf("Expected %s, got %s", expectedArguments, actualArguments)
	}
}

func TestGetImageName1(t *testing.T) {
	const systemName = "test-system"
	const applicationName = "test-image"
	const registry = "ghcr.io"

	expectedImageName := registry + "/" + systemName + "-" + applicationName

	actualImageName := getImageName(systemName, applicationName, registry)

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageName2(t *testing.T) {
	const systemName = ""
	const applicationName = "test-image"
	const registry = "ghcr.io"

	expectedImageName := registry + "/" + applicationName

	actualImageName := getImageName(systemName, applicationName, registry)

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}
