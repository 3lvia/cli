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
	const versionOlderThan0_57_1 = false

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
		versionOlderThan0_57_1,
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
	const versionOlderThan0_57_1 = false

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
		versionOlderThan0_57_1,
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
	const versionOlderThan0_57_1 = false

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
		versionOlderThan0_57_1,
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
	const versionOlderThan0_57_1 = false

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
		versionOlderThan0_57_1,
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
	const versionOlderThan0_57_1 = false

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
		versionOlderThan0_57_1,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommandVersionOlderThan0_57_1(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH"
	const disableError = false
	const versionOlderThan0_57_1 = true

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
			"--ignore-unfixed",
			"--exit-code",
			"1",
			"--scanners",
			"vuln",
			"--db-repository",
			"mirror.gcr.io/aquasec/trivy-db:2",
			"--java-db-repository",
			"mirror.gcr.io/aquasec/trivy-java-db:1",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		versionOlderThan0_57_1,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestScanImageCommandAllSeveritiesAndVersionTagAndVersionOlderThan0_57_1(t *testing.T) {
	const imageName = "test-image:v42"
	const severity = "CRITICAL,HIGH,MEDIUM,LOW,UNKNOWN"
	const disableError = false
	const versionOlderThan0_57_1 = true

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
			"--ignore-unfixed",
			"--exit-code",
			"1",
			"--scanners",
			"vuln",
			"--db-repository",
			"mirror.gcr.io/aquasec/trivy-db:2",
			"--java-db-repository",
			"mirror.gcr.io/aquasec/trivy-java-db:1",
			imageName,
		},
		" ",
	)

	actualCommand := scanImageCommand(
		imageName,
		severity,
		disableError,
		versionOlderThan0_57_1,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
