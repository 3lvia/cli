package command

import (
	"errors"
	"os/exec"
	"testing"
)

func TestIsError1(t *testing.T) {
	t.Parallel()

	output := Output{
		Error: nil,
	}

	actual := IsError(output)

	const expected = false

	if actual != expected {
		t.Errorf("Expected %t to be %t", actual, expected)
	}
}

func TestIsError2(t *testing.T) {
	t.Parallel()

	output := Output{
		Error: errors.New("error"),
	}

	actual := IsError(output)

	const expected = true

	if actual != expected {
		t.Errorf("Expected %t to be %t", actual, expected)
	}
}

func TestRunQuietCapturesOutput(t *testing.T) {
	t.Parallel()

	output := Run(*exec.Command("echo", "hello"), &RunOptions{Quiet: true})

	if IsError(output) {
		t.Fatalf("Expected no error, got %s", output.Error)
	}

	if output.Output != "hello\n" {
		t.Errorf("Expected captured output %q, got %q", "hello\n", output.Output)
	}
}
