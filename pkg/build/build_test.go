package build

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestGetImageName(t *testing.T) {
	t.Parallel()

	const (
		registry   = "containerregistryelvia.azurecr.io"
		systemName = "core"
		imageName  = "demo-api"
	)

	expectedImageName := registry + "/" + systemName + "-" + imageName

	actualImageName, err := GetImageName(registry, systemName, imageName)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageNameGHCR(t *testing.T) {
	t.Parallel()

	const (
		registry   = "ghcr.io"
		systemName = "core"
		imageName  = "demo-api"
	)

	expectedImageName := registry + "/" + systemName + "/" + imageName

	actualImageName, err := GetImageName(registry, systemName, imageName)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestGetImageNameOtherRegistry(t *testing.T) {
	t.Parallel()

	const (
		registry   = "quay.io"
		systemName = "core"
		imageName  = "demo-api"
	)

	expectedImageName := registry + "/" + systemName + "/" + imageName

	actualImageName, err := GetImageName(registry, systemName, imageName)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if actualImageName != expectedImageName {
		t.Errorf("Expected %s, got %s", expectedImageName, actualImageName)
	}
}

func TestBuildCommand1(t *testing.T) {
	t.Parallel()

	const (
		dockerfilePath = "build/Dockerfile"
		buildContext   = "src/app"
		imageName      = "containerregistryelvia.azurecr.io/test-image"
		cacheTag       = "latest"
	)

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
	t.Parallel()

	const (
		dockerfilePath = "Dockerfile"
		buildContext   = "."
		imageName      = "ghcr.io/test-image"
	)

	imageNameWithCacheTag := imageName + ":" + DefaultCacheTag

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
		DefaultCacheTag,
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
	t.Parallel()

	const (
		dockerfilePath = "Dockerfile"
		buildContext   = "."
		imageName      = "ghcr.io/test-image"
	)

	imageNameWithCacheTag := imageName + ":" + DefaultCacheTag
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
		DefaultCacheTag,
		additionalTags,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
