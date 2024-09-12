package scan

import "testing"

func TestConstructScanImageArguments1(t *testing.T) {
	const imageName = "test-image:latest"
	const severity = "CRITICAL,HIGH"
	const format = "table"
	const disableError = false

	expectedCommand := []string{
		"image",
		"--severity",
		severity,
		"--exit-code",
		"1",
		"--format",
		format,
		"--timeout",
		"15m0s",
		imageName,
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		format,
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
	const format = "json"
	const disableError = true

	expectedCommand := []string{
		"image",
		"--severity",
		severity,
		"--exit-code",
		"0",
		"--format",
		format,
		"--timeout",
		"15m0s",
		"--output",
		"trivy.json",
		imageName,
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		format,
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
	const format = "sarif"
	const disableError = true

	expectedCommand := []string{
		"image",
		"--severity",
		severity,
		"--exit-code",
		"0",
		"--format",
		format,
		"--timeout",
		"15m0s",
		"--output",
		"trivy." + format,
		imageName,
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		format,
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
	const format = "markdown"
	const disableError = true

	expectedCommand := []string{
		"image",
		"--severity",
		severity,
		"--exit-code",
		"0",
		"--format",
		"json",
		"--timeout",
		"15m0s",
		"--output",
		"trivy.json",
		imageName,
	}

	actualCommand := constructScanImageArguments(
		imageName,
		severity,
		format,
		disableError,
	)

	for i, arg := range expectedCommand {
		if arg != actualCommand[i] {
			t.Errorf("Expected %s, got %s", arg, actualCommand[i])
		}
	}
}
