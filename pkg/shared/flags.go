package shared

import (
	"fmt"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
)

func nameToEnvVar(name string) string {
	return fmt.Sprintf("3LV_%s", strings.ToUpper(strings.ReplaceAll(name, "-", "_")))
}

func ProjectFileFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "project-file",
		Aliases:  []string{"f"},
		Usage:    "The project file to use. We currently support .NET (*.csproj), Go (go.mod) or a generic project (Dockerfile).",
		Required: true,
	}
}

func RuntimeCloudProviderFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "runtime-cloud-provider",
		Aliases: []string{"r"},
		Usage:   "The runtime cloud provider to use (aks, gke, iss).",
		Value:   "aks",
		Action: func(c *cli.Context, runtimeCloudProvider string) error {
			allowedRuntimeCloudProviders := []string{"aks", "gke", "iss"}
			if !slices.Contains(allowedRuntimeCloudProviders, strings.ToLower(runtimeCloudProvider)) {
				return cli.Exit(
					fmt.Sprintf(
						"Invalid runtime cloud provider '%s' provided: must be one of %v (ignoring case)",
						runtimeCloudProvider,
						allowedRuntimeCloudProviders),
					1,
				)
			}

			return nil
		},
	}
}

func SystemNameFlag(usage string, required bool) *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "system-name",
		Aliases:  []string{"s"},
		Usage:    usage,
		Required: required,
		EnvVars:  []string{"3LV_SYSTEM_NAME"},
	}
}

func ApplicationNameFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "application-name",
		Aliases:  []string{"a"},
		Usage:    usage,
		Required: true,
	}
}

func HelmValuesFileFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "helm-values-file",
		Aliases: []string{"f"},
		Usage:   "The Helm values file used for deploying.",
		Value:   ".github/deploy/values.yml",
	}
}

func SeverityFlag(name string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    name,
		Aliases: []string{"S"},
		Usage:   "The severity to use when scanning the image: can be any combination of CRITICAL, HIGH, MEDIUM, LOW, or UNKNOWN separated by commas",
		Value:   "CRITICAL,HIGH",
		EnvVars: []string{nameToEnvVar(name)},
	}
}

func FormatsFlag(name string) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:    name,
		Aliases: []string{"F"},
		Usage:   "The formats to use when outputting the Trivy scan results: can be table, json, sarif or markdown.",
		Value:   cli.NewStringSlice("table"),
		Action: func(c *cli.Context, formats []string) error {
			for _, format := range formats {
				if format != "table" && format != "json" && format != "sarif" && format != "markdown" {
					return cli.Exit("Invalid format provided", 1)
				}
			}

			return nil
		},
		EnvVars: []string{nameToEnvVar(name)},
	}
}

func DisableErrorFlag(name string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    name,
		Aliases: []string{"D"},
		Usage:   "Disable error exit code on vulnerabilities found by Trivy.",
		Value:   false,
		EnvVars: []string{nameToEnvVar(name)},
	}
}

func RegistryFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "registry",
		Aliases: []string{"r"},
		Usage:   usage,
		EnvVars: []string{"3LV_REGISTRY"},
	}
}
