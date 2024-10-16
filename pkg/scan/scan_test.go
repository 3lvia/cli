package scan

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestScanImageCommand1(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH"
	const disableError = false
	const skipUpdate = false

	expectedCommandString := strings.Join(
		[]string{
			"trivy",
			"image",
			"--severity",
			severity,
			"--exit-code",
			"1",
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
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		skipUpdate,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommand2(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH,MEDIUM"
	const disableError = true
	const skipUpdate = false

	expectedCommandString := strings.Join(
		[]string{
			"image",
			"--severity",
			severity,
			"--exit-code",
			"0",
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
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		skipUpdate,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommand3(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL"
	const disableError = true
	const skipUpdate = true

	expectedCommandString := strings.Join(
		[]string{
			"image",
			"--severity",
			severity,
			"--exit-code",
			"0",
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
			"--skip-db-update",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		skipUpdate,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommand4(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH,MEDIUM,LOW"
	const disableError = true
	const skipUpdate = false

	expectedCommandString := strings.Join(
		[]string{
			"image",
			"--severity",
			severity,
			"--exit-code",
			"0",
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
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		skipUpdate,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommand5(t *testing.T) {
	const imageName = "test-image:v42"
	const severity = "CRITICAL,HIGH,MEDIUM,LOW,UNKNOWN"
	const disableError = false
	const skipUpdate = true

	expectedCommandString := strings.Join(
		[]string{
			"image",
			"--severity",
			severity,
			"--exit-code",
			"1",
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
			"--skip-db-update",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		skipUpdate,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
