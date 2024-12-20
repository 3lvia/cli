package scan

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestScanImageCommandNormal(t *testing.T) {
	t.Parallel()

	const (
		imageName    = "test-image:latest"
		severity     = "CRITICAL,HIGH"
		disableError = false
	)

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
			"--scanners",
			"vuln",
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
	t.Parallel()

	const (
		imageName    = "test-image:latest"
		severity     = "CRITICAL,HIGH,MEDIUM"
		disableError = true
	)

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
			"--scanners",
			"vuln",
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
	t.Parallel()

	const (
		imageName    = "test-image:latest"
		severity     = "CRITICAL"
		disableError = true
	)

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
			"--scanners",
			"vuln",
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
	t.Parallel()

	const (
		imageName    = "test-image:latest"
		severity     = "CRITICAL,HIGH,MEDIUM,LOW"
		disableError = true
	)

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
			"--scanners",
			"vuln",
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
	t.Parallel()

	const (
		imageName    = "test-image:v42"
		severity     = "CRITICAL,HIGH,MEDIUM,LOW,UNKNOWN"
		disableError = false
	)

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
			"--scanners",
			"vuln",
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
