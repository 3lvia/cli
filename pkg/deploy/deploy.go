package deploy

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strings"

	"github.com/3lvia/cli/pkg/auth"
	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v3"
)

var Command *cli.Command = &cli.Command{
	Name:      "deploy",
	Aliases:   []string{"d"},
	Usage:     "Deploy an application to a one of Atlas' Kubernetes clusters.",
	UsageText: "3lv deploy [options] <application-name>",
	Hidden:    true,
	Flags: []cli.Flag{
		shared.SystemNameFlag(
			"The name of the system (Kubernetes namespace) to deploy to.",
		),
		shared.RuntimeCloudProviderFlag(),
		shared.HelmValuesFileFlag(),
		&cli.StringFlag{
			Name:    "image-tag",
			Aliases: []string{"i"},
			Usage:   "The image tag to deploy.",
		},
		&cli.StringFlag{
			Name:    "environment",
			Aliases: []string{"e"},
			Usage:   "The environment to deploy to: sandbox, dev, test or prod",
			Value:   "dev",
			Action: func(_ context.Context, _ *cli.Command, environment string) error {
				allowedEnvironments := []string{"sandbox", "dev", "test", "prod"}
				if !slices.Contains(allowedEnvironments, environment) {
					return cli.Exit(fmt.Sprintf("Invalid environment provided: must be one of %v", allowedEnvironments), 1)
				}

				return nil
			},
		},
		&cli.StringFlag{
			Name:    "workload-type",
			Aliases: []string{"w"},
			Usage:   "The Kubernetes workload type to use: deployment or statefulset",
			Value:   "deployment",
			Action: func(_ context.Context, _ *cli.Command, workloadType string) error {
				allowedWorkloadTypes := []string{"deployment", "statefulset"}
				if !slices.Contains(allowedWorkloadTypes, workloadType) {
					return cli.Exit(fmt.Sprintf("Invalid workload type provided: must be one of %v", allowedWorkloadTypes), 1)
				}

				return nil
			},
		},
		&cli.StringFlag{
			Name:    "commit-hash",
			Aliases: []string{"c"},
			Usage: "The commit hash of the commit being deployed. Used for deployment annotations." +
				" If you are running this command from a git repository," +
				" the commit hash of the latest commit in the currently checked out branch will be used by default.",
		},
		&cli.StringFlag{
			Name:    "commit-message",
			Aliases: []string{"m"},
			Usage: "The commit message of the commit being deployed. Used for deployment annotations." +
				" If you are running this command from a git repository," +
				" the commit message of the latest commit in the currently checked out branch will be used by default.",
		},
		&cli.StringFlag{
			Name:    "repository-name",
			Aliases: []string{"n"},
			Usage: "Name of the repository the code of the application is stored in. Used for deployment annotations." +
				" If you are running this command from a git repository, that repository name will be used by default.",
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"D"},
			Usage:   "Simulate the deployment without actually deploying.",
		},
		&cli.StringFlag{
			Name:    "azure-tenant-id",
			Usage:   "The AKS tenant ID to use",
			Sources: cli.EnvVars("3LV_AZURE_TENANT_ID"),
		},
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
			Name:    "aks-subscription-id",
			Usage:   "Subscription ID of the AKS cluster to deploy to.",
			Sources: cli.EnvVars("3LV_AKS_SUBSCRIPTION_ID"),
		},
		&cli.StringFlag{
			Name:    "aks-cluster-name",
			Usage:   "Name of the AKS cluster to deploy to.",
			Sources: cli.EnvVars("3LV_AKS_CLUSTER_NAME"),
		},
		&cli.StringFlag{
			Name:    "aks-resource-group-name",
			Usage:   "Resource group name of the AKS cluster to deploy to.",
			Sources: cli.EnvVars("3LV_AKS_RESOURCE_GROUP_NAME"),
		},
		&cli.StringFlag{
			Name:    "gke-project-id",
			Usage:   "Project ID of the GKE cluster to deploy to.",
			Sources: cli.EnvVars("3LV_GKE_PROJECT_ID"),
		},
		&cli.StringFlag{
			Name:    "gke-cluster-name",
			Usage:   "Name of the GKE cluster to deploy to.",
			Sources: cli.EnvVars("3LV_GKE_CLUSTER_NAME"),
		},
		&cli.StringFlag{
			Name:    "gke-cluster-location",
			Usage:   "Location of the GKE cluster to deploy to.",
			Hidden:  true,
			Sources: cli.EnvVars("3LV_GKE_CLUSTER_LOCATION"),
		},
		&cli.BoolFlag{
			Name:  "add-deployment-annotation",
			Usage: "Add a deployment annotation to Grafana. Requires --grafana-url and --grafana-api-key to be set.",
		},
		&cli.StringFlag{
			Name:  "grafana-url",
			Usage: "The Grafana URL to use for deployment annotations.",
		},
		&cli.StringFlag{
			Name:  "grafana-api-key",
			Usage: "The Grafana API key to use for deployment annotations.",
		},
		&cli.StringFlag{
			Name:  "run-id",
			Usage: "The GitHub Actions run ID to use for deployment annotations.",
		},
		&cli.StringFlag{
			Name: "helm-chart-repository-url",
			Usage: "Override the helm chart repository where the elvia-charts are located." +
				" Useful for testing feature branches." +
				" For instance 'https://raw.githubusercontent.com/3lvia/kubernetes-charts/feature/cool-new-charts'.",
			Sources: cli.EnvVars("3LV_HELM_CHART_REPOSITORY_URL"),
		},
		&cli.BoolFlag{
			Name:    "allow-deploy",
			Hidden:  true,
			Sources: cli.EnvVars("CI"),
		},
	},
	Action: Deploy,
}

func Deploy(ctx context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	if !c.Bool("allow-deploy") {
		return cli.Exit("Deploy command is disabled when not running 3lv from GitHub Actions.", 1)
	}

	applicationName := c.Args().First()
	if applicationName == "" {
		return cli.Exit("Application name not provided.", 1)
	}

	systemName := c.String("system-name")
	helmValuesFile := c.String("helm-values-file")
	imageTag := c.String("image-tag")

	commitHash, err := utils.ResolveCommitHash(c.String("commit-hash"))
	if err != nil {
		return cli.Exit(err, 1)
	}

	repositoryName, err := utils.ResolveRepositoryName(c.String("repository-name"))
	if err != nil {
		return cli.Exit(err, 1)
	}

	addDeploymentAnnotation := c.Bool("add-deployment-annotation")
	commitMessage, err := utils.ResolveCommitMessage(c.String("commit-message"))

	if err != nil && addDeploymentAnnotation {
		return cli.Exit(err, 1)
	}

	grafanaURL := c.String("grafana-url")
	grafanaAPIKey := c.String("grafana-api-key")

	if addDeploymentAnnotation && (grafanaURL == "" || grafanaAPIKey == "") {
		return cli.Exit("Grafana URL and API key must be set when adding a deployment annotation", 1)
	}

	environment := strings.ToLower(c.String("environment"))
	workloadType := strings.ToLower(c.String("workload-type"))
	runtimeCloudProvider := strings.ToLower(c.String("runtime-cloud-provider"))
	dryRun := c.Bool("dry-run")
	runID := c.String("run-id")
	helmChartRepositoryURL := c.String("helm-chart-repository-url")

	checkKubectlInstalledOutput := checkKubectlInstalledCommand(nil)
	if command.IsError(checkKubectlInstalledOutput) {
		return cli.Exit(fmt.Errorf("kubectl is not installed: %w", checkKubectlInstalledOutput.Error), 1)
	}

	checkHelmInstalledOutput := checkHelmInstalledCommand(nil)
	if command.IsError(checkHelmInstalledOutput) {
		return cli.Exit(fmt.Errorf("helm is not installed: %w", checkHelmInstalledOutput.Error), 1)
	}

	if runtimeCloudProvider == "aks" {
		azureTenantID := utils.StringWithDefault(
			c.String("azure-tenant-id"),
			auth.ElviaTenantID,
		)
		loginOptions := &auth.AzLoginCommandOptions{
			ClientID:       c.String("azure-client-id"),
			FederatedToken: c.String("azure-federated-token"),
		}

		setupOptions := &SetupAKSOptions{
			SubscriptionID:    c.String("aks-subscription-id"),
			ClusterName:       c.String("aks-cluster-name"),
			ResourceGroupName: c.String("aks-resource-group-name"),
			AzLoginOptions:    loginOptions,
		}
		if err := setupAKS(
			azureTenantID,
			environment,
			setupOptions,
		); err != nil {
			return cli.Exit(err, 1)
		}
	} else if runtimeCloudProvider == "gke" {
		authOptions := &SetupGKEOptions{
			ProjectID:       c.String("gke-project-id"),
			ClusterName:     c.String("gke-cluster-name"),
			ClusterLocation: c.String("gke-cluster-location"),
		}
		if err := setupGKE(
			environment,
			authOptions,
		); err != nil {
			return cli.Exit(err, 1)
		}
	} else if runtimeCloudProvider == "iss" { //nolint:revive
		// do nothing
	} else {
		// This should never happen, as the runtimeCloudProvider flag is validated in the cli.Action function.
		return cli.Exit(fmt.Errorf("Invalid runtime cloud provider: %s", runtimeCloudProvider), 1)
	}

	helmRepoAddOutput := helmRepoAddCommand(helmChartRepositoryURL, nil)
	if command.IsError(helmRepoAddOutput) {
		return cli.Exit(fmt.Errorf("Failed to add Helm repository: %w", helmRepoAddOutput.Error), 1)
	}

	helmRepoUpdateOutput := helmRepoUpdateCommand(nil)
	if command.IsError(helmRepoUpdateOutput) {
		return cli.Exit(fmt.Errorf("Failed to update Helm repository: %w", helmRepoUpdateOutput.Error), 1)
	}

	useISSChart := runtimeCloudProvider == "iss"

	helmDeployOutput := helmDeployCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		repositoryName,
		commitHash,
		dryRun,
		useISSChart,
		nil,
	)
	if command.IsError(helmDeployOutput) {
		if !dryRun {
			// If the deployment failed, we still want to post the Grafana annotation,
			// but we add a failure message to the annotation.
			if err := addGrafanaDeploymentAnnotation(
				ctx,
				false,
				applicationName,
				systemName,
				environment,
				runtimeCloudProvider,
				repositoryName,
				commitMessage,
				grafanaURL,
				grafanaAPIKey,
				&PostGrafanaAnnotationOptions{
					RunID: runID,
				},
			); err != nil {
				return cli.Exit(
					fmt.Errorf("Failed to deploy Helm chart %w and post Grafana annotation: %w", helmDeployOutput.Error, err),
					1,
				)
			}
		}

		return cli.Exit(fmt.Errorf("Failed to deploy Helm chart: %w", helmDeployOutput.Error), 1)
	}

	kubectlRolloutStatusOutput := kubectlRolloutStatusCommand(
		applicationName,
		systemName,
		workloadType,
		nil,
	)
	if command.IsError(kubectlRolloutStatusOutput) {
		return cli.Exit(kubectlRolloutStatusOutput.Error, 1)
	}

	kubectlGetEventsOutput := kubectlGetEventsCommand(
		systemName,
		nil,
	)
	if command.IsError(kubectlGetEventsOutput) {
		return cli.Exit(kubectlGetEventsOutput.Error, 1)
	}

	events := strings.Split(kubectlGetEventsOutput.Output, "\n")
	for _, event := range events {
		if strings.Contains(event, applicationName) {
			log.Print(event)
		}
	}

	if addDeploymentAnnotation && !dryRun {
		if err := addGrafanaDeploymentAnnotation(
			ctx,
			true,
			applicationName,
			systemName,
			environment,
			runtimeCloudProvider,
			repositoryName,
			commitMessage,
			grafanaURL,
			grafanaAPIKey,
			&PostGrafanaAnnotationOptions{
				RunID: runID,
			},
		); err != nil {
			return cli.Exit(fmt.Errorf("Failed to post Grafana annotation: %w", err), 1)
		}
	}

	return nil
}

func checkKubectlInstalledCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"kubectl",
			"version",
			"--client",
			"true",
		),
		runOptions,
	)
}

func kubectlRolloutStatusCommand(
	applicationName string,
	systemName string,
	workloadType string,
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"kubectl",
			"rollout",
			"status",
			"-n",
			systemName,
			workloadType+"/"+applicationName,
		),
		runOptions,
	)
}

func kubectlGetEventsCommand(
	systemName string,
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"kubectl",
			"get",
			"events",
			"-n",
			systemName,
			"--sort-by",
			".lastTimestamp",
		),
		runOptions,
	)
}
