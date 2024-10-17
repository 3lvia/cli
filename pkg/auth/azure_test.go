package auth

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

	actualCommand := azLoginCommand(
		tenantID,
		&AzLoginCommandOptions{
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
	const clientID = "test-client-id"
	const federatedToken = "test-federated-token"

	expectedCommandString := strings.Join(
		[]string{
			"az",
			"login",
			"--tenant",
			tenantID,
			"--service-principal",
			"--username",
			clientID,
			"--federated-token",
			federatedToken,
		},
		" ",
	)

	actualCommand := azLoginCommand(
		tenantID,
		&AzLoginCommandOptions{
			FederatedToken: federatedToken,
			ClientID:       clientID,
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
