package deploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func azLoginTenant(tenantID string) error {
	azLoginCmd := exec.Command(
		"az",
		"login",
		"--tenant",
		tenantID,
	)
	azLoginCmd.Stdout = os.Stdout
	azLoginCmd.Stderr = os.Stderr

	log.Print(azLoginCmd.String())

	if err := azLoginCmd.Run(); err != nil {
		return fmt.Errorf("Failed to authenticate to Azure: %w", err)
	}

	return nil
}

type AuthenticateAKSOptions struct {
	AKSTenantID          string
	AKSSubscriptionID    string
	AKSClusterName       string
	AKSResourceGroupName string
}

func authenticateAKS(environment string, options AuthenticateAKSOptions) error {
	aksTenantID := func() string {
		if options.AKSTenantID == "" {
			return "2186a6ec-c227-4291-9806-d95340bf439d"
		}
		return options.AKSTenantID
	}()

	azAccountShowCmd := exec.Command(
		"az",
		"account",
		"show",
		"--query",
		"tenantId",
		"--output",
		"tsv",
	)
	log.Print(azAccountShowCmd.String())

	azAccountShowCmdOutput, err := azAccountShowCmd.Output()
	if err != nil {
		err := azLoginTenant(aksTenantID)
		if err != nil {
			return err
		}
	}

	if azAccountShowCmdOutput != nil {
		azAccountShowCmdOutputString := strings.TrimSpace(string(azAccountShowCmdOutput))
		if azAccountShowCmdOutputString != aksTenantID {
			err := azLoginTenant(aksTenantID)
			if err != nil {
				return err
			}
		}
	}

	aksSubscriptionID, err := func() (string, error) {
		if options.AKSSubscriptionID == "" {
			switch environment {
			case "dev", "test":
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

	contextName := "atlascluster" + environment

	aksResourceGroupName := func() string {
		if options.AKSResourceGroupName == "" {
			return "RUNTIMESERVICE-RG" + environment
		}
		return options.AKSResourceGroupName
	}()

	azGetCredentialsCmd := exec.Command(
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
	)
	azGetCredentialsCmd.Stdout = os.Stdout
	azGetCredentialsCmd.Stderr = os.Stderr

	log.Print(azGetCredentialsCmd.String())

	if err := azGetCredentialsCmd.Run(); err != nil {
		return fmt.Errorf("Failed to authenticate to AKS: %w", err)
	}

	setContextCmd := exec.Command(
		"kubectl",
		"config",
		"use-context",
		contextName,
	)
	setContextCmd.Stdout = os.Stdout
	setContextCmd.Stderr = os.Stderr

	log.Print(setContextCmd.String())

	if err := setContextCmd.Run(); err != nil {
		return fmt.Errorf("Failed to set kubectl context: %w", err)
	}

	return nil
}
