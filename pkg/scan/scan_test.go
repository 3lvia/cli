package scan

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestScanImageCommandNormal(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH"
	const disableError = false

	expectedCommandString := strings.Join(
		[]string{
			"trivy",
			"image",
			"--severity",
			severity,
			"--timeout",
			"15m0s",
			"--format",
			"json",
			"--output",
			"trivy.json",
			"--db-repository",
			"ghcr.io/3lvia/trivy-db",
			"--java-db-repository",
			"ghcr.io/3lvia/trivy-java-db",
			"--ignore-unfixed",
			"--exit-code",
			"1",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageDisableErrorAndMoreSeverities(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH,MEDIUM"
	const disableError = true

	expectedCommandString := strings.Join(
		[]string{
			"image",
			"--severity",
			severity,
			"--timeout",
			"15m0s",
			"--format",
			"json",
			"--output",
			"trivy.json",
			"--db-repository",
			"ghcr.io/3lvia/trivy-db",
			"--java-db-repository",
			"ghcr.io/3lvia/trivy-java-db",
			"--ignore-unfixed",
			"--exit-code",
			"0",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommandDisableErrorAndLessSeverities(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL"
	const disableError = true

	expectedCommandString := strings.Join(
		[]string{
			"image",
			"--severity",
			severity,
			"--timeout",
			"15m0s",
			"--format",
			"json",
			"--output",
			"trivy.json",
			"--db-repository",
			"ghcr.io/3lvia/trivy-db",
			"--java-db-repository",
			"ghcr.io/3lvia/trivy-java-db",
			"--ignore-unfixed",
			"--exit-code",
			"0",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommandEventMoreSeverities(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH,MEDIUM,LOW"
	const disableError = true

	expectedCommandString := strings.Join(
		[]string{
			"image",
			"--severity",
			severity,
			"--timeout",
			"15m0s",
			"--format",
			"json",
			"--output",
			"trivy.json",
			"--db-repository",
			"ghcr.io/3lvia/trivy-db",
			"--java-db-repository",
			"ghcr.io/3lvia/trivy-java-db",
			"--ignore-unfixed",
			"--exit-code",
			"0",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommandAllSeveritiesAndVersionTag(t *testing.T) {
	const imageName = "test-image:v42"
	const severity = "CRITICAL,HIGH,MEDIUM,LOW,UNKNOWN"
	const disableError = false

	expectedCommandString := strings.Join(
		[]string{
			"image",
			"--severity",
			severity,
			"--timeout",
			"15m0s",
			"--format",
			"json",
			"--output",
			"trivy.json",
			"--db-repository",
			"ghcr.io/3lvia/trivy-db",
			"--java-db-repository",
			"ghcr.io/3lvia/trivy-java-db",
			"--ignore-unfixed",
			"--exit-code",
			"1",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
