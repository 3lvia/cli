package deploy

import (
	"fmt"
	"os/exec"

	"github.com/3lvia/cli/pkg/auth"
	"github.com/3lvia/cli/pkg/command"
)

type SetupAKSOptions struct {
	SubscriptionID    string
	ClusterName       string
	ResourceGroupName string
	AzLoginOptions    *auth.AzLoginCommandOptions
}

func setupAKS(
	tenantID string,
	environment string,
	skipAuthentication bool,
	options *SetupAKSOptions,
) error {
	if options == nil {
		options = &SetupAKSOptions{}
	}
	if options.AzLoginOptions == nil {
		options.AzLoginOptions = &auth.AzLoginCommandOptions{}
	}

	subscriptionID, err := auth.GetElviaDefaultRuntimeSubscriptionID(
		environment,
		options.SubscriptionID,
	)
	if err != nil {
		return err
	}

	if !skipAuthentication {
		err = auth.AuthenticateAzure(
			tenantID,
			subscriptionID,
			options.AzLoginOptions,
		)
		if err != nil {
			return fmt.Errorf("Failed to authenticate with AKS: %w", err)
		}
	}

	checkKubeLoginInstalledOutput := checkKubeloginInstalledCommand(nil)
	if command.IsError(checkKubeLoginInstalledOutput) {
		return checkKubeLoginInstalledOutput.Error
	}

	clusterName := func() string {
		if options.ClusterName == "" {
			return "akscluster" + environment
		}
		return options.ClusterName
	}()

	contextName := "aks" + environment

	resourceGroupName := func() string {
		if options.ResourceGroupName == "" {
			return "RUNTIMESERVICE-RG" + environment
		}
		return options.ResourceGroupName
	}()

	runKubeloginConvert := options.AzLoginOptions.FederatedToken != "" && options.AzLoginOptions.ClientID != ""
	if err := getAKSCredentials(
		resourceGroupName,
		clusterName,
		subscriptionID,
		contextName,
		runKubeloginConvert,
	); err != nil {
		return fmt.Errorf("Failed to get AKS credentials: %w", err)
	}

	return nil
}

func getAKSCredentials(
	resourceGroupName string,
	clusterName string,
	subscriptionID string,
	contextName string,
	runKubeloginConvert bool,
) error {
	azGetCredentialsOutput := azGetCredentialsCommand(
		resourceGroupName,
		clusterName,
		subscriptionID,
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
