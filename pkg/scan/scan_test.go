package scan

import "testing"

func TestConstructScanImageArguments1(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH"
	const disableError = false

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
		imageName,
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
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
		imageName,
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
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
		imageName,
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
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
		imageName,
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		disableError,
	)

	for i, arg := range expectedCommand {
		if arg != actualCommand[i] {
			t.Errorf("Expected %s, got %s", arg, actualCommand[i])
		}
	}
}
