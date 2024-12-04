package scan

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v3"
)

var Command *cli.Command = &cli.Command{
	Name:    "scan",
	Aliases: []string{"s"},
	Usage:   "Scan image using Trivy.",
	Flags: []cli.Flag{
		shared.SeverityFlag("severity"),
		shared.FormatsFlag("formats"),
		shared.DisableErrorFlag("disable-error"),
	},
	Action: Scan,
}

func Scan(ctx context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		return cli.ShowAppHelp(c)
	}

	// Required args
	imageName := c.Args().First()
	if imageName == "" {
		style.Print(
			"Image name not provided.",
			&style.PrintOptions{Color: "red"},
		)
		return cli.ShowAppHelp(c)
	}

	// Optional args
	severity := c.String("severity")
	formats := utils.RemoveZeroValues(c.StringSlice("formats"))
	disableError := c.Bool("disable-error")

	err := ScanImage(imageName, severity, formats, disableError)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func scanImageCommand(
	imageName string,
	severity string,
	disableError bool,
	runOptions *command.RunOptions,
) command.Output {
	exitCode := func() string {
		if disableError {
			return "0"
		}

		return "1"
	}()

	cmd := exec.Command(
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
		exitCode,
		"--scanners",
		"vuln",
	)

	cmd.Args = append(cmd.Args, imageName)

	return command.Run(*cmd, runOptions)
}

func convertCommand(
	format string,
	runOptions *command.RunOptions,
) command.Output {
	if format == "table" {
		return command.Run(
			*exec.Command(
				"trivy",
				"convert",
				"--format",
				"table",
				"trivy.json",
			),
			runOptions,
		)
	}

	if format == "sarif" {
		return command.Run(
			*exec.Command(
				"trivy",
				"convert",
				"--format",
				"sarif",
				"--output",
				"trivy.sarif",
				"trivy.json",
			),
			runOptions,
		)
	}

	return command.Error(fmt.Errorf("Invalid format %s", format))
}

func ScanImage(
	imageName string,
	severity string,
	formats []string,
	disableError bool,
) error {
	scanImageOutput := scanImageCommand(
		imageName,
		severity,
		disableError,
		nil,
	)

	if _, err := os.Stat("trivy.json"); errors.Is(err, os.ErrNotExist) {
		if disableError {
			style.Print(
				"Trivy did not produce any output.",
				&style.PrintOptions{Color: "yellow"},
			)
			return nil
		}

		return fmt.Errorf("Trivy did not produce any output")
	}

	if slices.Contains(formats, "table") {
		style.Print(
			"Converted results to table format.",
			nil,
		)

		convertOutput := convertCommand(
			"table",
			nil,
		)
		if command.IsError(convertOutput) {
			return convertOutput.Error
		}
	}

	if slices.Contains(formats, "sarif") {
		style.Print(
			"Converted results to SARIF format.",
			nil,
		)

		convertOutput := convertCommand(
			"sarif",
			nil,
		)
		if command.IsError(convertOutput) {
			return convertOutput.Error
		}
	}

	if slices.Contains(formats, "markdown") {
		style.Print(
			"Converting results to markdown format.",
			nil,
		)

		result, err := parseJSONOutput()
		if err != nil {
			return err
		}

		markdown, err := toMarkdown(result)
		if err != nil {
			return err
		}
		if len(markdown) == 0 {
			style.Print(
				"Markdown output is empty, will write to empty file",
				&style.PrintOptions{Color: "yellow"},
			)
		}

		err = os.WriteFile("trivy.md", markdown, 0644)
		if err != nil {
			return err
		}
	}

	if !slices.Contains(formats, "json") {
		err := os.Remove("trivy.json")
		if err != nil {
			return err
		}
	} else {
		style.Print(
			"Keeping pre-existing JSON output.",
			nil,
		)
	}

	if command.IsError(scanImageOutput) {
		return scanImageOutput.Error
	}

	return nil
}
