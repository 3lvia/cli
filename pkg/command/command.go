package command

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/style"
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
		CommandString: "",
		Error:         err,
		Output:        "",
	}
}

func ErrorString(err string) Output {
	return Output{
		CommandString: "",
		Error:         errors.New(err),
		Output:        "",
	}
}

type RunOptions struct {
	DryRun bool
	Silent bool
}

func Run(cmd exec.Cmd, options *RunOptions) Output {
	if options == nil {
		options = &RunOptions{}
	}

	if options.DryRun {
		return Output{
			CommandString: cmd.String(),
			Error:         nil,
			Output:        "",
		}
	}

	style.Print(
		cmd.String()+"\n",
		&style.PrintOptions{Color: "cyan"},
	)

	var errBuf, outBuf bytes.Buffer

	if options.Silent {
		cmd.Stdout = nil
		cmd.Stderr = nil
	} else {
		cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
		cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	}

	err := cmd.Run()
	if err != nil {
		return Output{
			CommandString: "",
			Error:         err,
			Output:        errBuf.String(),
		}
	}

	return Output{
		CommandString: cmd.String(),
		Error:         nil,
		Output:        outBuf.String(),
	}
}

func ExpectedCommandStringEqualsActualCommand(t *testing.T, expectedCommandString string, actualCommand Output) {
	t.Helper()

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
