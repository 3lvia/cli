package deploy

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/3lvia/cli/pkg/command"
)

type SetupAKSOptions struct {
	AKSTenantID          string
	AKSSubscriptionID    string
	AKSClusterName       string
	AKSResourceGroupName string
	FederatedToken       string
}

func setupAKS(
	environment string,
	skipAuthentication bool,
	options SetupAKSOptions,
) error {
	if !skipAuthentication {
		if err := authenticateAKS(
			AuthenticateAKSOptions{
				AKSTenantID:    options.AKSTenantID,
				FederatedToken: options.FederatedToken,
			},
		); err != nil {
			return fmt.Errorf("Failed to authenticate to AKS: %w", err)
		}
	}

	checkKubeLoginInstalledOutput := checkKubeloginInstalledCommand(nil)
	if command.IsError(checkKubeLoginInstalledOutput) {
		return checkKubeLoginInstalledOutput.Error
	}

	aksSubscriptionID, err := func() (string, error) {
		if options.AKSSubscriptionID == "" {
			switch environment {
			case "sandbox", "dev", "test":
				return "ceb9518c-528f-4c91-9b5a-c051d383e7a8", nil
			case "prod":
				return "9edbf217-b7c1-4f6a-ae76-d046cf932ff0", nil
			default:
				// cannot happen, but have to do this since Go's type system is garbage
				return "", fmt.Errorf("Invalid environment provided")
			}
		}
		return options.AKSSubscriptionID, nil
	}()
	if err != nil {
		return err
	}

	aksClusterName := func() string {
		if options.AKSClusterName == "" {
			return "akscluster" + environment
		}
		return options.AKSClusterName
	}()

	contextName := "aks" + environment

	aksResourceGroupName := func() string {
		if options.AKSResourceGroupName == "" {
			return "RUNTIMESERVICE-RG" + environment
		}
		return options.AKSResourceGroupName
	}()

	runKubeloginConvert := skipAuthentication
	if err := getAKSCredentials(
		aksResourceGroupName,
		aksClusterName,
		aksSubscriptionID,
		contextName,
		runKubeloginConvert,
	); err != nil {
		return fmt.Errorf("Failed to get AKS credentials: %w", err)
	}

	return nil
}

type AuthenticateAKSOptions struct {
	AKSTenantID    string
	FederatedToken string
}

func authenticateAKS(
	options AuthenticateAKSOptions,
) error {
	aksTenantID := func() string {
		if options.AKSTenantID == "" {
			return "2186a6ec-c227-4291-9806-d95340bf439d"
		}
		return options.AKSTenantID
	}()

	azAccountShowCommandOutput := azAccountShowCommand(nil)
	if command.IsError(azAccountShowCommandOutput) {
		azLoginTenantCommandOutput := azLoginTenantCommand(
			aksTenantID,
			AzLoginTenantCommandOptions{FederatedToken: options.FederatedToken},
		)
		if command.IsError(azLoginTenantCommandOutput) {
			return fmt.Errorf("Failed to authenticate to AKS: %w", azLoginTenantCommandOutput.Error)
		}
	}

	if tenantID := azAccountShowCommandOutput.Output; tenantID != "" {
		azAccountShowCmdOutputString := strings.TrimSpace(string(tenantID))
		if azAccountShowCmdOutputString != aksTenantID {
			azLoginTenantCommandOutput := azLoginTenantCommand(
				aksTenantID,
				AzLoginTenantCommandOptions{FederatedToken: options.FederatedToken},
			)
			if command.IsError(azLoginTenantCommandOutput) {
				return fmt.Errorf("Failed to authenticate to AKS: %w", azLoginTenantCommandOutput.Error)
			}
		}
	}

	return nil
}

func getAKSCredentials(
	aksResourceGroupName string,
	aksClusterName string,
	aksSubscriptionID string,
	contextName string,
	runKubeloginConvert bool,
) error {
	azGetCredentialsOutput := azGetCredentialsCommand(
		aksResourceGroupName,
		aksClusterName,
		aksSubscriptionID,
		contextName,
		nil,
	)
	if command.IsError(azGetCredentialsOutput) {
		return fmt.Errorf("Failed to get AKS credentials: %w", azGetCredentialsOutput.Error)
	}

	if runKubeloginConvert {
		kubeloginConvertOutput := kubeloginConvertCommand(nil)
		if command.IsError(kubeloginConvertOutput) {
			return fmt.Errorf("Failed to convert AKS credentials: %w", kubeloginConvertOutput.Error)
		}
	}

	return nil
}

func azGetCredentialsCommand(
	aksResourceGroupName string,
	aksClusterName string,
	aksSubscriptionID string,
	contextName string,
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"az",
			"aks",
			"get-credentials",
			"--resource-group",
			aksResourceGroupName,
			"--name",
			aksClusterName,
			"--context",
			contextName,
			"--subscription",
			aksSubscriptionID,
			"--overwrite-existing",
		),
		runOptions,
	)
}

func kubeloginConvertCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"kubelogin",
			"convert-kubeconfig",
			"-l",
			"azurecli",
		),
		runOptions,
	)
}

func checkKubeloginInstalledCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command("kubelogin", "--version"),
		runOptions,
	)
}

func azAccountShowCommand(runOptions *command.RunOptions) command.Output {
	return command.Run(
		*exec.Command(
			"az",
			"account",
			"show",
			"--query",
			"tenantId",
			"--output",
			"tsv",
		),
		runOptions,
	)
}

type AzLoginTenantCommandOptions struct {
	FederatedToken string
	RunOptions     *command.RunOptions
}

func azLoginTenantCommand(
	tenantID string,
	options AzLoginTenantCommandOptions,
) command.Output {
	cmd := exec.Command(
		"az",
		"login",
		"--tenant",
		tenantID,
	)

	if options.FederatedToken != "" {
		cmd.Args = append(cmd.Args, "--federated-token")
		cmd.Args = append(cmd.Args, options.FederatedToken)
	}

	result := command.Run(*cmd, options.RunOptions)
	if command.IsError(result) {
		return command.Error(
			fmt.Errorf("Failed to login to Azure tenant %s", tenantID),
		)
	}

	return result
}
