package build

import (
	"strings"
	"testing"
)

func TestGetImageName1(t *testing.T) {
	const registry = "containerregistryelvia.azurecr.io"
	const systemName = "core"
	const imageName = "demo-api"

	expectedImageName := registry + "/" + systemName + "-" + imageName

	actualImageName := getImageName(registry, systemName, imageName)

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageName2(t *testing.T) {
	const registry = "containerregistryelvia.azurecr.o"
	const systemName = "core"
	const imageName = "demo-api"

	expectedImageName := registry + "/" + systemName + "/" + imageName

	actualImageName := getImageName(registry, systemName, imageName)

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageName3(t *testing.T) {
	const registry = "ghcr.io"
	const systemName = "core"
	const imageName = "demo-api"

	expectedImageName := registry + "/" + systemName + "/" + imageName

	actualImageName := getImageName(registry, systemName, imageName)

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageName4(t *testing.T) {
	const registry = "quay.io"
	const systemName = "core"
	const imageName = "demo-api"

	expectedImageName := registry + "/" + systemName + "/" + imageName

	actualImageName := getImageName(registry, systemName, imageName)

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestConstructBuildCommandArguments1(t *testing.T) {
	const dockerfilePath = "build/Dockerfile"
	const buildContext = "src/app"
	const imageName = "containerregistryelvia.azurecr.io/test-image"
	const cacheTag = "latest"

	imageNameWithCacheTag := imageName + ":" + cacheTag

	expectedArguments := "buildx build -f " + dockerfilePath + " --load --cache-to type=inline --cache-from " + imageNameWithCacheTag + " -t " + imageNameWithCacheTag + " " + buildContext

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

	expectedArguments := "buildx build -f " + dockerfilePath + " --load --cache-to type=inline --cache-from " + imageNameWithCacheTag + " -t " + imageNameWithCacheTag + " " + buildContext

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

	expectedArguments := "buildx build -f " + dockerfilePath + " --load --cache-to type=inline --cache-from " + imageNameWithCacheTag + " -t " + imageName + ":" + additionalTags[0] + " -t " + imageName + ":" + additionalTags[1] + " -t " + imageName + ":" + additionalTags[2] + " -t " + imageNameWithCacheTag + " " + buildContext

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
