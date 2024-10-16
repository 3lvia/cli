package command

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type Output struct {
	CommandString string
	Error         error
	Output        string
}

func IsError(output Output) bool {
	return output.Error != nil
}

func Error(err error) Output {
	return Output{
		Error: err,
	}
}

func ErrorString(err string) Output {
	return Output{
		Error: errors.New(err),
	}
}

type RunOptions struct {
	DryRun bool
}

func Run(cmd exec.Cmd, options *RunOptions) Output {
	if options == nil {
		options = &RunOptions{}
	}

	if options.DryRun {
		return Output{
			CommandString: cmd.String(),
		}
	}

	log.Print(cmd.String())

	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)

	err := cmd.Run()
	if err != nil {
		return Output{
			Error:  err,
			Output: errBuf.String(),
		}
	}

	return Output{
		CommandString: cmd.String(),
		Output:        outBuf.String(),
	}
}

func ExpectedCommandStringEqualsActualCommand(t *testing.T, expectedCommandString string, actualCommand Output) {
	if IsError(actualCommand) {
		t.Errorf("Expected no error, got %s", actualCommand.Error)
	}

	if len(actualCommand.CommandString) == 0 {
		t.Errorf("Expected a command string, got an empty string")
	}

	if !strings.HasSuffix(actualCommand.CommandString, expectedCommandString) {
		t.Errorf("Expected command string to end with %s, got %s", expectedCommandString, actualCommand.CommandString)
	}
}
