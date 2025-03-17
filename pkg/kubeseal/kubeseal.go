package kubeseal

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/3lvia/cli/pkg/auth"
	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/deploy"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/style"
	"github.com/urfave/cli/v3"
)

const (
	commandName = "kubeseal"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      commandName,
		Aliases:   []string{"k"},
		Usage:     "Create a sealed secret for a Kubernetes cluster.",
		UsageText: "3lv kubeseal [options] <secret-value>",
		Flags: []cli.Flag{
			shared.RuntimeCloudProviderFlag(),
			shared.EnvironmentFlag(),
			&cli.StringFlag{
				Name: "azure-client-id",
				Usage: "The client ID to use when authenticating with the registry." +
					" Must be combined with --azure-federated-token.",
				Sources: cli.EnvVars("3LV_AZURE_CLIENT_ID"),
			},
			&cli.StringFlag{
				Name: "azure-federated-token",
				Usage: "The federated token to use when authenticating with the Azure Container Registry." +
					" Must be combined with --client-id.",
				Sources: cli.EnvVars("3LV_AZURE_FEDERATED_TOKEN"),
			},
			&cli.StringFlag{
				Name:     "secret-name",
				Usage:    "The name of the secret to create.",
				Aliases:  []string{"n"},
				Required: true,
			},
			&cli.StringFlag{
				Name:  "output-file",
				Usage: "The name of the output file to write the sealed secret to. Defaults to standard output.",
			},
            &cli.BoolFlag{
                Name: "skip-authentication",
                Usage: "Skip authenticating with the Kubernetes cluster.",
            },
		},
		Action: Kubeseal,
	}
}

func Kubeseal(_ context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	runtimeCloudProvider := strings.ToLower(c.String("runtime-cloud-provider"))
	environment := c.String("environment")

	secretName := c.String("secret-name")
	if secretName == "" {
		return cli.Exit("secret-name is required", 1)
	}

	secretValue := c.Args().First()
	if secretValue == "" {
		return cli.Exit("secret-value is required", 1)
	}

    if !c.Bool("skip-authentication") {
        err := deploy.SetupKubernetes(
            &deploy.SetupAKSOptions{ // nolint:exhaustruct
                AzLoginOptions: &auth.AzLoginCommandOptions{
                    ClientID:       c.String("azure-client-id"),
                    FederatedToken: c.String("azure-federated-token"),
                },
            },
            &deploy.SetupGKEOptions{}, // nolint:exhaustruct
            "",                        // default to Elvia tenant
            runtimeCloudProvider,
            environment,
        )
        if err != nil {
            return cli.Exit(err, 1)
        }
    }

	checkKubesealInstalledOutput := checkKubesealInstalledCommand(nil)
	if command.IsError(checkKubesealInstalledOutput) {
		return checkKubesealInstalledOutput.Error
	}

	kubectlCreateSecretDryRunOutput := kubectlCreateSecretDryRunCommand(
		secretName,
		secretValue,
		nil,
	)
	if command.IsError(kubectlCreateSecretDryRunOutput) {
		return kubectlCreateSecretDryRunOutput.Error
	}

	secretJSON, err := os.CreateTemp("", "secret-*.json")
	if err != nil {
		return err
	}

	defer os.Remove(secretJSON.Name())

	_, err = secretJSON.WriteString(kubectlCreateSecretDryRunOutput.Output)
	if err != nil {
		return err
	}

	outputFile := c.String("output-file")

	kubesealOutput := kubesealCommand(
		secretJSON,
		outputFile,
		nil,
	)
	if command.IsError(kubesealOutput) {
		return kubesealOutput.Error
	}

	if outputFile != "" {
		style.PrintSuccess("Sealed secret written to " + outputFile + ".")
	}

	return nil
}

func checkKubesealInstalledCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command("kubeseal", "--version"),
		runOptions,
	)
}

func kubectlCreateSecretDryRunCommand(
	secretName string,
	secretValue string,
	runOptions *command.RunOptions,
) command.Output {
	if runOptions == nil {
		runOptions = &command.RunOptions{}
	}

	return command.Run(
		*exec.Command(
			"kubectl",
			"create",
			"secret",
			"generic",
			secretName,
			"--dry-run=client",
			"--from-literal",
			secretName+"="+secretValue,
			"-o",
			"json",
		),
		&command.RunOptions{
			PrintCommand: false, // don't print command with secret
			Silent:       true,  // don't print output since we output secret to stdout
			DryRun:       runOptions.DryRun,
		},
	)
}

func kubesealCommand(
	secretJSON *os.File,
	outputFile string,
	runOptions *command.RunOptions,
) command.Output {
	if runOptions == nil {
		runOptions = &command.RunOptions{}
	}

	cmd := *exec.Command(
		"kubeseal",
		"--controller-namespace",
		"sealed-secrets",
		"--controller-name",
		"sealed-secrets",
		"-f",
		secretJSON.Name(),
		"-o",
		"yaml",
	)

	if outputFile != "" {
		cmd.Args = append(cmd.Args, "-w", outputFile)
	}

	return command.Run(
		cmd,
		&command.RunOptions{
			PrintCommand: false, // don't print command with secret
			Silent:       runOptions.Silent,
			DryRun:       runOptions.DryRun,
		},
	)
}
