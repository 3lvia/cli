package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
)

func Deploy(c *cli.Context) error {
	if c.NArg() <= 0 {
		return cli.Exit("No input provided", 1)
	}

	applicationName := c.Args().First()
	if applicationName == "" {
		return cli.Exit("Application name not provided", 1)
	}

	systemName := c.String("system-name")
	helmValuesFile := c.String("helm-values-file")

	environment := strings.ToLower(c.String("environment"))
	allowedEnvironments := []string{"dev", "test", "prod"}
	if !slices.Contains(allowedEnvironments, environment) {
		return cli.Exit(fmt.Sprintf("Invalid environment provided: must be one of %v", allowedEnvironments), 1)
	}

	workloadType := strings.ToLower(c.String("workload-type"))
	allowedWorkloadTypes := []string{"deployment", "statefulset"}
	if !slices.Contains(allowedWorkloadTypes, workloadType) {
		return cli.Exit(fmt.Sprintf("Invalid workload type provided: must be one of %v", allowedWorkloadTypes), 1)
	}

	runtimeCloudProvider := strings.ToLower(c.String("runtime-cloud-provider"))
	allowedRuntimeCloudProviders := []string{"aks", "gke"}
	if !slices.Contains(allowedRuntimeCloudProviders, runtimeCloudProvider) {
		return cli.Exit(fmt.Sprintf("Invalid runtime cloud provider provided: must be one of %v", allowedRuntimeCloudProviders), 1)
	}

	if runtimeCloudProvider == "aks" {
		authOptions := DeployToAKSOptions{
			Environment:          environment,
			AKSSubscriptionID:    c.String("aks-subscription-id"),
			AKSClusterName:       c.String("aks-cluster-name"),
			AKSResourceGroupName: c.String("aks-resource-group-name"),
		}
		if err := authenticateAKS(authOptions); err != nil {
			return cli.Exit(err, 1)
		}

	} else if runtimeCloudProvider == "gke" {
		if err := authenticateGKE(); err != nil {
			return cli.Exit(err, 1)
		}
	}

	helmOptions := HelmDeployOptions{
		ApplicationName: applicationName,
		SystemName:      systemName,
		HelmValuesFile:  helmValuesFile,
		Environment:     environment,
		WorkloadType:    workloadType,
	}
	if err := helmDeploy(helmOptions); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

type DeployToAKSOptions struct {
	Environment          string // required
	AKSSubscriptionID    string
	AKSClusterName       string
	AKSResourceGroupName string
}

func authenticateAKS(options DeployToAKSOptions) error {
	// const azureTenantID = "2186a6ec-c227-4291-9806-d95340bf439d"

	aksSubscriptionID, err := func() (string, error) {
		if options.AKSSubscriptionID == "" {
			switch options.Environment {
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
			return "akscluster" + options.Environment
		}
		return options.AKSClusterName
	}()

	aksResourceGroupName := func() string {
		if options.AKSResourceGroupName == "" {
			return "RUNTIMESERVICE-RG" + options.Environment
		}
		return options.AKSResourceGroupName
	}()

	fmt.Printf("Authenticating to AKS cluster %s in resource group %s in subscription %s\n", aksClusterName, aksResourceGroupName, aksSubscriptionID)

	return nil
}

func authenticateGKE() error {
	return fmt.Errorf("Not implemented")
}

type HelmDeployOptions struct {
	ApplicationName string
	SystemName      string
	HelmValuesFile  string
	Environment     string
	WorkloadType    string
}

func helmDeploy(options HelmDeployOptions) error {
	const (
		chartsNamespace     = "elvia-charts"
		chartsRepositoryURL = "https://raw.githubusercontent.com/3lvia/kubernetes-charts/master"
	)

	helmRepoAddCmd := exec.Command(
		"helm",
		"repo",
		"add",
		chartsNamespace,
		chartsRepositoryURL,
	)
	helmRepoAddCmd.Stdout = os.Stdout
	helmRepoAddCmd.Stderr = os.Stderr

	if err := helmRepoAddCmd.Run(); err != nil {
		return fmt.Errorf("Failed to add Helm repository: %w", err)
	}

	helmRepoUpdateCmd := exec.Command(
		"helm",
		"repo",
		"update",
	)
	helmRepoUpdateCmd.Stdout = os.Stdout
	helmRepoUpdateCmd.Stderr = os.Stderr

	if err := helmRepoUpdateCmd.Run(); err != nil {
		return fmt.Errorf("Failed to update Helm repository: %w", err)
	}

	helmDeployCmd := exec.Command(
		"helm",
		"upgrade",
		"--debug",
		"--install",
		"-n",
		options.SystemName,
		"-f",
		options.HelmValuesFile,
		options.ApplicationName,
		"elvia-charts/"+options.WorkloadType,
		"--set='environment="+options.Environment+"'",
	)
	helmDeployCmd.Stdout = os.Stdout
	helmDeployCmd.Stderr = os.Stderr

	if err := helmDeployCmd.Run(); err != nil {
		return fmt.Errorf("Failed to deploy Helm chart: %w", err)
	}

	return nil
}
