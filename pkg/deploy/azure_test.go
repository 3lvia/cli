package deploy

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestAzLoginTenantCommand1(t *testing.T) {
	const tenantID = "test-tenant-id"

	expectedCommandString := strings.Join(
		[]string{
			"az",
			"login",
			"--tenant",
			tenantID,
		},
		" ",
	)

	actualCommand := azLoginTenantCommand(
		tenantID,
		AzLoginTenantCommandOptions{
			RunOptions: &command.RunOptions{DryRun: true},
		},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestAzLoginTenantCommand2(t *testing.T) {
	const tenantID = "test-tenant-id"
	const federatedToken = "test-federated-token"

	expectedCommandString := strings.Join(
		[]string{
			"az",
			"login",
			"--tenant",
			tenantID,
			"--federated-token",
			federatedToken,
		},
		" ",
	)

	actualCommand := azLoginTenantCommand(
		tenantID,
		AzLoginTenantCommandOptions{
			FederatedToken: federatedToken,
			RunOptions:     &command.RunOptions{DryRun: true},
		},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestAzAccountShowCommand(t *testing.T) {
	expectedCommandString := strings.Join(
		[]string{
			"az",
			"account",
			"show",
			"--query",
			"tenantId",
			"--output",
			"tsv",
		},
		" ",
	)

	actualCommand := azAccountShowCommand(&command.RunOptions{DryRun: true})

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestCheckKubeloginInstalledCommand(t *testing.T) {
	expectedCommandString := strings.Join(
		[]string{
			"kubelogin",
			"--version",
		},
		" ",
	)

	actualCommand := checkKubeloginInstalledCommand(&command.RunOptions{DryRun: true})

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestKubeloginConvertCommand(t *testing.T) {
	expectedCommandString := strings.Join(
		[]string{
			"kubelogin",
			"convert-kubeconfig",
			"-l",
			"azurecli",
		},
		" ",
	)

	actualCommand := kubeloginConvertCommand(&command.RunOptions{DryRun: true})

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestAzGetCredentialsCommand(t *testing.T) {
	const aksResourceGroupName = "test-aks-resource-group-name"
	const aksClusterName = "test-aks-cluster-name"
	const aksSubscriptionID = "1234-5678-9012-3456"
	const contextName = "test-context-name"

	expectedCommandString := strings.Join(
		[]string{
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
		},
		" ",
	)

	actualCommand := azGetCredentialsCommand(
		aksResourceGroupName,
		aksClusterName,
		aksSubscriptionID,
		contextName,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
