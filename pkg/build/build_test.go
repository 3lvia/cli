package build

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestGetImageName1(t *testing.T) {
	const registry = "containerregistryelvia.azurecr.io"
	const systemName = "core"
	const imageName = "demo-api"

	expectedImageName := registry + "/" + systemName + "-" + imageName

	actualImageName, err := getImageName(registry, systemName, imageName)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageName2(t *testing.T) {
	const registry = "containerregistryelvia.azurecr.o"
	const systemName = "core"
	const imageName = "demo-api"

	expectedImageName := registry + "/" + systemName + "/" + imageName

	actualImageName, err := getImageName(registry, systemName, imageName)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageName3(t *testing.T) {
	const registry = "ghcr.io"
	const systemName = "core"
	const imageName = "demo-api"

	expectedImageName := registry + "/" + systemName + "/" + imageName

	actualImageName, err := getImageName(registry, systemName, imageName)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageName4(t *testing.T) {
	const registry = "quay.io"
	const systemName = "core"
	const imageName = "demo-api"

	expectedImageName := registry + "/" + systemName + "/" + imageName

	actualImageName, err := getImageName(registry, systemName, imageName)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestBuildCommand1(t *testing.T) {
	const dockerfilePath = "build/Dockerfile"
	const buildContext = "src/app"
	const imageName = "containerregistryelvia.azurecr.io/test-image"
	const cacheTag = "latest"

	imageNameWithCacheTag := imageName + ":" + cacheTag

	expectedCommandString := strings.Join(
		[]string{
			"docker",
			"buildx",
			"build",
			"-f",
			dockerfilePath,
			"--load",
			"--cache-to",
			"type=inline",
			"--cache-from",
			imageNameWithCacheTag,
			"-t",
			imageNameWithCacheTag,
			buildContext,
		},
		" ",
	)

	actualCommand := buildImageCommand(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		[]string{},
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestBuildCommand2(t *testing.T) {
	const dockerfilePath = "Dockerfile"
	const buildContext = "."
	const imageName = "ghcr.io/test-image"
	const cacheTag = "latest-cache"

	imageNameWithCacheTag := imageName + ":" + cacheTag

	expectedCommandString := strings.Join(
		[]string{
			"docker",
			"buildx",
			"build",
			"-f",
			dockerfilePath,
			"--load",
			"--cache-to",
			"type=inline",
			"--cache-from",
			imageNameWithCacheTag,
			"-t",
			imageNameWithCacheTag,
			buildContext,
		},
		" ",
	)

	actualCommand := buildImageCommand(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		[]string{},
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestBuildCommand3(t *testing.T) {
	const dockerfilePath = "Dockerfile"
	const buildContext = "."
	const imageName = "ghcr.io/test-image"
	const cacheTag = "latest-cache"

	imageNameWithCacheTag := imageName + ":" + cacheTag
	additionalTags := []string{"latest", "v42.0.1", "v420alpha"}

	expectedCommandString := strings.Join(
		[]string{
			"docker",
			"buildx",
			"build",
			"-f",
			dockerfilePath,
			"--load",
			"--cache-to",
			"type=inline",
			"--cache-from",
			imageNameWithCacheTag,
			"-t",
			imageName + ":" + additionalTags[0],
			"-t",
			imageName + ":" + additionalTags[1],
			"-t",
			imageName + ":" + additionalTags[2],
			"-t",
			imageNameWithCacheTag,
			buildContext,
		},
		" ",
	)

	actualCommand := buildImageCommand(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		additionalTags,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
