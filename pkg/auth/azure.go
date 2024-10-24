package auth

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/3lvia/cli/pkg/command"
)

const ElviaTenantID = "2186a6ec-c227-4291-9806-d95340bf439d"
const ElviaDefaultRuntimeSubscriptionID = "9edbf217-b7c1-4f6a-ae76-d046cf932ff0"
const ElviaDefaultRuntimeDevTestSubscriptionID = "ceb9518c-528f-4c91-9b5a-c051d383e7a8"

func GetElviaDefaultRuntimeSubscriptionID(
	environment string,
	subscriptionID string,
) (string, error) {
	if subscriptionID != "" {
		return subscriptionID, nil
	}

	switch environment {
	case "dev", "test", "sandbox":
		return ElviaDefaultRuntimeDevTestSubscriptionID, nil
	case "prod":
		return ElviaDefaultRuntimeSubscriptionID, nil
	default:
		return "", fmt.Errorf("Unknown environment: %s", environment)
	}
}

func AuthenticateAzure(
	tenantID string,
	subscriptionID string,
	options *AzLoginCommandOptions,
) error {
	azAccountShowCommandOutput := azAccountShowCommand(nil)
	if command.IsError(azAccountShowCommandOutput) {
		azLoginCommandOutput := azLoginCommand(
			tenantID,
			options,
		)
		if command.IsError(azLoginCommandOutput) {
			return fmt.Errorf("Failed to authenticate to Azure: %w", azLoginCommandOutput.Error)
		}
		return nil
	}

	if tenantID := azAccountShowCommandOutput.Output; tenantID != "" {
		azAccountShowCmdOutputString := strings.TrimSpace(string(tenantID))
		if azAccountShowCmdOutputString != tenantID {
			azLoginTenantCommandOutput := azLoginCommand(
				tenantID,
				options,
			)
			if command.IsError(azLoginTenantCommandOutput) {
				return fmt.Errorf("Failed to authenticate to Azure: %w", azLoginTenantCommandOutput.Error)
			}
		}
	}

	azSetSubscriptionCommandOutput := azSetSubscriptionCommand(
		subscriptionID,
		nil,
	)
	if command.IsError(azSetSubscriptionCommandOutput) {
		return fmt.Errorf("Failed to set subscription: %w", azSetSubscriptionCommandOutput.Error)
	}

	return nil
}

func azSetSubscriptionCommand(
	subscriptionID string,
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"az",
			"account",
			"set",
			"--subscription",
			subscriptionID,
		),
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

type AzLoginCommandOptions struct {
	FederatedToken string
	ClientID       string
	RunOptions     *command.RunOptions
}

func azLoginCommand(
	tenantID string,
	options *AzLoginCommandOptions,
) command.Output {
	if options == nil {
		options = &AzLoginCommandOptions{}
	}

	if tenantID == "" {
		return command.Error(fmt.Errorf("Tenant ID is required"))
	}

	cmd := exec.Command(
		"az",
		"login",
		"--tenant",
		tenantID,
	)

	if options.FederatedToken != "" {
		if options.ClientID == "" {
			return command.Error(fmt.Errorf("Client ID is required when federated token is provided"))
		}

		cmd.Args = append(cmd.Args, "--service-principal")

		cmd.Args = append(cmd.Args, "--username")
		cmd.Args = append(cmd.Args, options.ClientID)

		cmd.Args = append(cmd.Args, "--federated-token")
		cmd.Args = append(cmd.Args, options.FederatedToken)
	}

	result := command.Run(*cmd, options.RunOptions)
	if command.IsError(result) {
		return command.Error(
			fmt.Errorf(
				"Failed to login to Azure tenant %s: %w",
				tenantID,
				result.Error,
			),
		)
	}

	return result
}
