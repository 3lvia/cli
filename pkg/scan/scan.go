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

func Command() *cli.Command {
	return &cli.Command{
		Name:      "scan",
		Aliases:   []string{"s"},
		Usage:     "Scan a container image using Trivy.",
		UsageText: "3lv scan [options] <image-name>",
		Flags: []cli.Flag{
			shared.SeverityFlag("severity"),
			shared.FormatsFlag("formats"),
			shared.DisableErrorFlag("disable-error"),
		},
		Action: Scan,
	}
}

func Scan(_ context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	// Required args
	imageName := c.Args().First()
	if imageName == "" {
		style.PrintError("Image name not provided.")

		cli.ShowSubcommandHelpAndExit(c, 1)
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
			style.PrintWarning("Trivy did not produce any output.")

			return nil
		}

		return errors.New("Trivy did not produce any output")
	}

	if slices.Contains(formats, "table") {
		style.PrintInfo("Converted results to table format.")

		if convertOutput := convertCommand(
			"table",
			nil,
		); command.IsError(convertOutput) {
			return convertOutput.Error
		}
	}

	if slices.Contains(formats, "sarif") {
		style.PrintInfo("Converted results to SARIF format.")

		if convertOutput := convertCommand(
			"sarif",
			nil,
		); command.IsError(convertOutput) {
			return convertOutput.Error
		}
	}

	if slices.Contains(formats, "markdown") {
		style.PrintInfo("Converting results to markdown format.")

		result, err := parseJSONOutput()
		if err != nil {
			return err
		}

		markdown, err := toMarkdown(result)
		if err != nil {
			return err
		}

		if len(markdown) == 0 {
			style.PrintWarning("Markdown output is empty, will write to empty file")
		}

		if err := os.WriteFile("trivy.md", markdown, 0o644); err != nil {
			return err
		}
	}

	if !slices.Contains(formats, "json") {
		if err := os.Remove("trivy.json"); err != nil {
			return err
		}
	} else {
		style.PrintInfo("Keeping pre-existing JSON output.")
	}

	if command.IsError(scanImageOutput) {
		return scanImageOutput.Error
	}

	return nil
}
