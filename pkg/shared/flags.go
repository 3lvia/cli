package shared

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

func nameToEnvVar(name string) string {
	return "3LV_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
}

func ProjectFileFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "project-file",
		Aliases: []string{"f"},
		Usage: "The project file to use. We currently support .NET (*.csproj), Go (go.mod)," +
			" Python with uv (uv.lock) or a generic Docker project (Dockerfile).",
	}
}

func RuntimeCloudProviderFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "runtime-cloud-provider",
		Aliases: []string{"r"},
		Usage:   "The runtime cloud provider to use (aks, gke, iss).",
		Value:   "aks",
		Action: func(_ context.Context, _ *cli.Command, runtimeCloudProvider string) error {
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

func SystemNameFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "system-name",
		Aliases: []string{"s"},
		Usage:   usage,
		Sources: cli.EnvVars("3LV_SYSTEM_NAME"),
	}
}

func ApplicationNameFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "application-name",
		Aliases: []string{"a"},
		Usage:   usage,
	}
}

func HelmValuesFileFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "helm-values-file",
		Aliases: []string{"F"},
		Usage:   "The Helm values file used for deploying.",
		Value:   ".github/deploy/values.yml",
	}
}

func SeverityFlag(name string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    name,
		Aliases: []string{"S"},
		Usage: "The severity to use when scanning the image: can be any combination of" +
			" CRITICAL, HIGH, MEDIUM, LOW, or UNKNOWN separated by commas",
		Value:   "CRITICAL,HIGH",
		Sources: cli.EnvVars(nameToEnvVar(name)),
	}
}

func FormatsFlag(name string) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:    name,
		Aliases: []string{"F"},
		Usage:   "The formats to use when outputting the Trivy scan results: can be table, json, sarif or markdown.",
		Value:   []string{"table"},
		Action: func(_ context.Context, _ *cli.Command, formats []string) error {
			for _, format := range formats {
				if format != "table" && format != "json" && format != "sarif" && format != "markdown" {
					return cli.Exit("Invalid format provided", 1)
				}
			}

			return nil
		},
		Sources: cli.EnvVars(nameToEnvVar(name)),
	}
}

func DisableErrorFlag(name string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    name,
		Aliases: []string{"D"},
		Usage:   "Disable error exit code on vulnerabilities found by Trivy.",
		Value:   false,
		Sources: cli.EnvVars(nameToEnvVar(name)),
	}
}

func RegistryFlag(usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "registry",
		Aliases: []string{"r"},
		Usage:   usage,
		Sources: cli.EnvVars("3LV_REGISTRY"),
	}
}
