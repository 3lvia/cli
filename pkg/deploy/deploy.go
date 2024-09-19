package deploy

import (
	"fmt"
	"os/exec"
	"path"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
)

var Command *cli.Command = &cli.Command{
	Name:    "deploy",
	Aliases: []string{"d"},
	Usage:   "Deploy the project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "system-name",
			Aliases:  []string{"s"},
			Usage:    "The system name to use",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "helm-values-file",
			Aliases:  []string{"f"},
			Usage:    "The helm values file to use",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "image-tag",
			Aliases:  []string{"i"},
			Usage:    "The image tag to deploy",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "environment",
			Aliases: []string{"e"},
			Usage:   "The environment to deploy to",
			Value:   "dev",
			Action: func(c *cli.Context, environment string) error {
				allowedEnvironments := []string{"dev", "test", "prod"}
				if !slices.Contains(allowedEnvironments, environment) {
					return cli.Exit(fmt.Sprintf("Invalid environment provided: must be one of %v", allowedEnvironments), 1)
				}

				return nil
			},
		},
		&cli.StringFlag{
			Name:    "workload-type",
			Aliases: []string{"w"},
			Usage:   "The workload type to use",
			Value:   "deployment",
			Action: func(c *cli.Context, workloadType string) error {
				allowedWorkloadTypes := []string{"deployment", "statefulset"}
				if !slices.Contains(allowedWorkloadTypes, workloadType) {
					return cli.Exit(fmt.Sprintf("Invalid workload type provided: must be one of %v", allowedWorkloadTypes), 1)
				}

				return nil
			},
		},
		&cli.StringFlag{
			Name:    "runtime-cloud-provider",
			Aliases: []string{"r"},
			Usage:   "The runtime cloud provider to use",
			Value:   "aks",
			Action: func(c *cli.Context, runtimeCloudProvider string) error {
				allowedRuntimeCloudProviders := []string{"aks", "gke"}
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
		},
		&cli.StringFlag{
			Name:    "commit-hash",
			Aliases: []string{"c"},
			Usage:   "The commit hash to use",
		},
		&cli.StringFlag{
			Name:    "repository-name",
			Aliases: []string{"n"},
			Usage:   "The repository name to use",
		},
		&cli.BoolFlag{
			Name:    "skip-authentication",
			Aliases: []string{"A"},
			Usage:   "Skips authentication against the runtime cloud provider",
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"D"},
			Usage:   "Simulate the deployment without actually deploying.",
		},
		&cli.StringFlag{
			Name:   "aks-tenant-id",
			Usage:  "The AKS tenant ID to use",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "aks-subscription-id",
			Usage:  "The AKS subscription ID to use",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "aks-cluster-name",
			Usage:  "The AKS cluster name to use",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "aks-resource-group-name",
			Usage:  "The AKS resource group name to use",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "gke-project-id",
			Usage:  "The GKE project ID to use",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "gke-cluster-name",
			Usage:  "The GKE cluster name to use",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "gke-cluster-zone",
			Usage:  "The GKE cluster zone to use",
			Hidden: true,
		},
	},
	Action: Deploy,
}

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
	imageTag := c.String("image-tag")

	commitHash, err := resolveCommitHash(c.String("commit-hash"))
	if err != nil {
		return cli.Exit(err, 1)
	}

	repositoryName, err := resolveRepositoryName(c.String("repository-name"))
	if err != nil {
		return cli.Exit(err, 1)
	}

	environment := strings.ToLower(c.String("environment"))
	workloadType := strings.ToLower(c.String("workload-type"))
	runtimeCloudProvider := strings.ToLower(c.String("runtime-cloud-provider"))
	skipAuthentication := c.Bool("skip-authentication")
	dryRun := c.Bool("dry-run")

	if !skipAuthentication {
		if runtimeCloudProvider == "aks" {
			authOptions := AuthenticateAKSOptions{
				AKSTenantID:          c.String("aks-tenant-id"),
				AKSSubscriptionID:    c.String("aks-subscription-id"),
				AKSClusterName:       c.String("aks-cluster-name"),
				AKSResourceGroupName: c.String("aks-resource-group-name"),
			}
			if err := authenticateAKS(environment, authOptions); err != nil {
				return cli.Exit(err, 1)
			}

		} else if runtimeCloudProvider == "gke" {
			authOptions := AuthenticateGKEOptions{
				GKEProjectID:       c.String("gke-project-id"),
				GKEClusterName:     c.String("gke-cluster-name"),
				GKEClusterLocation: c.String("gke-cluster-zone"),
			}
			if err := authenticateGKE(environment, authOptions); err != nil {
				return cli.Exit(err, 1)
			}
		}
	}

	const (
		chartsNamespace     = "elvia-charts"
		chartsRepositoryURL = "https://raw.githubusercontent.com/3lvia/kubernetes-charts/master"
	)

	if err := helmRepoAdd(chartsNamespace, chartsRepositoryURL); err != nil {
		return fmt.Errorf("Failed to add Helm repository: %w", err)
	}

	if err := helmRepoUpdate(); err != nil {
		return fmt.Errorf("Failed to update Helm repository: %w", err)
	}

	if err := helmDeploy(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		repositoryName,
		commitHash,
		dryRun,
	); err != nil {
		return cli.Exit(err, 1)
	}

	if err := kubectlRolloutStatus(
		applicationName,
		systemName,
		workloadType,
	); err != nil {
		return cli.Exit(err, 1)
	}

	if err := kubectlGetEvents(applicationName, systemName); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func resolveCommitHash(possibleCommitHash string) (string, error) {
	if possibleCommitHash != "" {
		return possibleCommitHash, nil
	}

	hash, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return "",
			fmt.Errorf(
				"Failed to resolve commit hash: %w. Please verify you are currently in a Git repository, or manually specify the commit hash with --commit-hash.",
				err,
			)
	}

	return strings.TrimSpace(string(hash)), nil
}

func resolveRepositoryName(possibleRepositoryName string) (string, error) {
	if possibleRepositoryName != "" {
		return possibleRepositoryName, nil
	}

	gitTopLevel, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "",
			fmt.Errorf(
				"Failed to resolve repository name: %w. Please verify you are currently in a Git repository, or manually specify the repository with --repository-name.",
				err,
			)
	}

	return path.Base(strings.TrimSpace(string(gitTopLevel))), nil
}
