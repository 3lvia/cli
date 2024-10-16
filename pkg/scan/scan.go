package scan

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"

	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

const commandName = "scan"

var Command *cli.Command = &cli.Command{
	Name:    "scan",
	Aliases: []string{"s"},
	Usage:   "Scan image using Trivy",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "severity",
			Aliases: []string{"S"},
			Usage:   "The severity to use when scanning the image: can be any combination of CRITICAL, HIGH, MEDIUM, LOW, or UNKNOWN separated by commas",
			Value:   "CRITICAL,HIGH",
			EnvVars: []string{"3LV_SEVERITY"},
		},
		&cli.StringSliceFlag{
			Name:    "formats",
			Aliases: []string{"F"},
			Usage:   "The formats to use when outputting the scan results: can be table, json, sarif or markdown.",
			Value:   cli.NewStringSlice("table"),
			Action: func(c *cli.Context, formats []string) error {
				for _, format := range formats {
					if format != "table" && format != "json" && format != "sarif" && format != "markdown" {
						return cli.Exit("Invalid format provided", 1)
					}
				}

				return nil
			},
			EnvVars: []string{"3LV_FORMATS"},
		},
		&cli.BoolFlag{
			Name:    "disable-error",
			Aliases: []string{"D"},
			Usage:   "Disable error exit code on vulnerabilities found",
			Value:   false,
			EnvVars: []string{"3LV_DISABLE_ERROR"},
		},
		&cli.BoolFlag{
			Name:    "skip-db-update",
			Usage:   "Skip update Trivy vulnerability database",
			Value:   false,
			EnvVars: []string{"3LV_SKIP_DB_UPDATE"},
		},
	},
	Action: Scan,
}

func Scan(c *cli.Context) error {
	if c.NArg() <= 0 {
		return cli.ShowCommandHelp(c, commandName)
	}

	// Required args
	imageName := c.Args().First()
	if imageName == "" {
		log.Println("Image name not provided")
		return cli.ShowCommandHelp(c, commandName)
	}

	// Optional args
	severity := c.String("severity")
	formats := utils.RemoveZeroValues(c.StringSlice("formats"))
	disableError := c.Bool("disable-error")
	skipDBUpdate := c.Bool("skip-db-update")

	err := ScanImage(imageName, severity, formats, disableError, skipDBUpdate)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func scanImageCommand(
	imageName string,
	severity string,
	disableError bool,
	skipDBUpdate bool,
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
		"--exit-code",
		exitCode,
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
	)

	if skipDBUpdate {
		cmd.Args = append(cmd.Args, "--skip-db-update")
	}

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
	skipDBUpdate bool,
) error {
	scanImageOutput := scanImageCommand(
		imageName,
		severity,
		disableError,
		skipDBUpdate,
		nil,
	)

	if _, err := os.Stat("trivy.json"); errors.Is(err, os.ErrNotExist) {
		if disableError {
			log.Println("Trivy did not produce any output")
			return nil
		}

		return fmt.Errorf("Trivy did not produce any output")
	}

	if slices.Contains(formats, "table") {
		log.Println("Converting results to table format")

		convertOutput := convertCommand(
			"table",
			nil,
		)
		if command.IsError(convertOutput) {
			return convertOutput.Error
		}
	}

	if slices.Contains(formats, "sarif") {
		log.Println("Converting results to SARIF format")

		convertOutput := convertCommand(
			"sarif",
			nil,
		)
		if command.IsError(convertOutput) {
			return convertOutput.Error
		}
	}

	if slices.Contains(formats, "markdown") {
		log.Println("Converting results to Markdown format")

		result, err := parseJSONOutput()
		if err != nil {
			return err
		}

		markdown, err := toMarkdown(result)
		if err != nil {
			return err
		}
		if len(markdown) == 0 {
			log.Println("Markdown output is empty, will write to empty file")
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
		log.Println("Keeping pre-existing JSON output")
	}

	if command.IsError(scanImageOutput) {
		return scanImageOutput.Error
	}

	return nil
}
