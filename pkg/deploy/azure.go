package deploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type SetupAKSOptions struct {
	AKSTenantID          string
	AKSSubscriptionID    string
	AKSClusterName       string
	AKSResourceGroupName string
}

func setupAKS(
	environment string,
	skipAuthentication bool,
	options SetupAKSOptions,
) error {
	if !skipAuthentication {
		if err := authenticateAKS(
			environment,
			AuthenticateAKSOptions{AKSTenantID: options.AKSTenantID},
		); err != nil {
			return fmt.Errorf("Failed to authenticate to AKS: %w", err)
		}
	}

	if err := checkKubeloginInstalled(); err != nil {
		return err
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

	contextName := "aks" + environment

	aksResourceGroupName := func() string {
		if options.AKSResourceGroupName == "" {
			return "RUNTIMESERVICE-RG" + environment
		}
		return options.AKSResourceGroupName
	}()

	if err := getAKSCredentials(
		aksResourceGroupName,
		aksClusterName,
		aksSubscriptionID,
		contextName,
	); err != nil {
		return fmt.Errorf("Failed to get AKS credentials: %w", err)
	}

	return nil
}

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

func checkKubeloginInstalled() error {
	if err := exec.Command("kubelogin", "--version").Run(); err != nil {
		return fmt.Errorf("kubelogin is not installed")
	}

	return nil
}

type AuthenticateAKSOptions struct {
	AKSTenantID string
}

func authenticateAKS(
	environment string,
	options AuthenticateAKSOptions,
) error {
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

	return nil
}

func getAKSCredentials(
	aksResourceGroupName string,
	aksClusterName string,
	aksSubscriptionID string,
	contextName string,
) error {
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

	return nil
}
