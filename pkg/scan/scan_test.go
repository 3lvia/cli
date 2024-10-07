package scan

import "testing"

func TestConstructScanImageArguments1(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH"
	const disableError = false
	const skipUpdate = false

	expectedCommand := []string{
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
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
		skipUpdate,
	)

	for i, arg := range expectedCommand {
		if arg != actualCommand[i] {
			t.Errorf("Expected %s, got %s", arg, actualCommand[i])
		}
	}
}

func TestConstructScanImageArguments2(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH,MEDIUM"
	const disableError = true
	const skipUpdate = false

	expectedCommand := []string{
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
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
		skipUpdate,
	)

	for i, arg := range expectedCommand {
		if arg != actualCommand[i] {
			t.Errorf("Expected %s, got %s", arg, actualCommand[i])
		}
	}
}

func TestConstructScanImageArguments3(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL"
	const disableError = true
	const skipUpdate = true

	expectedCommand := []string{
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
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
		skipUpdate,
	)

	for i, arg := range expectedCommand {
		if arg != actualCommand[i] {
			t.Errorf("Expected %s, got %s", arg, actualCommand[i])
		}
	}
}

func TestConstructScanImageArguments4(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH,MEDIUM,LOW"
	const disableError = true
	const skipUpdate = false

	expectedCommand := []string{
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
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
		skipUpdate,
	)

	for i, arg := range expectedCommand {
		if arg != actualCommand[i] {
			t.Errorf("Expected %s, got %s", arg, actualCommand[i])
		}
	}
}

func TestConstructScanImageArguments5(t *testing.T) {
	const imageName = "test-image:v42"
	const severity = "CRITICAL,HIGH,MEDIUM,LOW,UNKNOWN"
	const disableError = false
	const skipUpdate = true

	expectedCommand := []string{
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
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
		skipUpdate,
	)

	for i, arg := range expectedCommand {
		if arg != actualCommand[i] {
			t.Errorf("Expected %s, got %s", arg, actualCommand[i])
		}
	}
}
