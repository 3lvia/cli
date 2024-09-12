package scan

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

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
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"F"},
			Usage:   "The format to use when outputting the scan results: can be table, json, sarif or markdown.",
			Value:   "table",
			Action: func(c *cli.Context, format string) error {
				if format != "table" && format != "json" && format != "sarif" && format != "markdown" {
					return cli.Exit("Invalid format provided", 1)
				}

				return nil
			},
			EnvVars: []string{"3LV_FORMAT"},
		},
		&cli.BoolFlag{
			Name:    "disable-error",
			Aliases: []string{"D"},
			Usage:   "Disable error exit code on vulnerabilities found",
			Value:   false,
			EnvVars: []string{"3LV_DISABLE_ERROR"},
		},
	},
	Action: Scan,
}

func Scan(c *cli.Context) error {
	if c.NArg() <= 0 {
		return cli.Exit("No input provided", 1)
	}

	// Required args
	imageName := c.Args().First()
	if imageName == "" {
		return cli.Exit("Image name not provided", 1)
	}

	// Optional args
	severity := c.String("severity")
	format := c.String("format")
	disableError := c.Bool("disable-error")

	err := ScanImage(imageName, severity, format, disableError)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func constructScanImageArguments(
	imageName string,
	severity string,
	format string,
	disableError bool,
) []string {
	exitCode := func() string {
		if disableError {
			return "0"
		}
		return "1"
	}()

	// Markdown output is created by parsing json output, not natively by Trivy
	trivyFormat := func() string {
		if format == "markdown" {
			return "json"
		}
		return format
	}()

	cmd := []string{
		"image",
		"--severity",
		severity,
		"--exit-code",
		exitCode,
		"--format",
		trivyFormat,
		"--timeout",
		"15m0s",
	}

	if trivyFormat == "json" || trivyFormat == "sarif" {
		return append(
			append(
				cmd,
				"--output",
				"trivy."+trivyFormat,
			),
			imageName,
		)
	}

	return append(cmd, imageName)
}

func ScanImage(
	imageName string,
	severity string,
	format string,
	disableError bool,
) error {

	scanCmd := exec.Command(
		"trivy",
		constructScanImageArguments(
			imageName,
			severity,
			format,
			disableError,
		)...,
	)
	scanCmd.Stdout = os.Stdout
	scanCmd.Stderr = os.Stderr

	log.Printf(scanCmd.String())

	_ = scanCmd.Run()

	if format == "markdown" {
		result, err := parseJSONOutput()
		if err != nil {
			return err
		}
		markdown, err := toMarkdown(result)
		if err != nil {
			return err
		}

		err = os.WriteFile("trivy.md", markdown, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
