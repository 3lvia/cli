package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strings"

	"github.com/3lvia/cli/pkg/auth"
	"github.com/3lvia/cli/pkg/build"
	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
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
				Usage:   "The environment to deploy to: sandbox, dev, kptest, test or prod",
				Value:   "dev",
				Action: func(_ context.Context, _ *cli.Command, environment string) error {
					allowedEnvironments := []string{"sandbox", "dev", "kptest", "test", "prod"}
					if !slices.Contains(allowedEnvironments, environment) {
						return cli.Exit(fmt.Sprintf("Invalid environment provided: must be one of %v", allowedEnvironments), 1)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:    "workload-type",
				Aliases: []string{"w"},
				Usage:   "The Kubernetes workload type to use: deployment, statefulset or job",
				Value:   "deployment",
				Action: func(_ context.Context, _ *cli.Command, workloadType string) error {
					allowedWorkloadTypes := []string{"deployment", "statefulset", "job"}
					if !slices.Contains(allowedWorkloadTypes, workloadType) {
						return cli.Exit(fmt.Sprintf("Invalid workload type provided: must be one of %v", allowedWorkloadTypes), 1)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:        "commit-hash",
				Aliases:     []string{"c"},
				Usage:       "The commit hash of the commit being deployed. Used for deployment annotations.",
				DefaultText: "latest commit hash of the repository the command is run in",
			},
			&cli.StringFlag{
				Name:        "commit-message",
				Aliases:     []string{"m"},
				Usage:       "The commit message of the commit being deployed. Used for deployment annotations.",
				DefaultText: "latest commit message of the repository the command is run in",
			},
			&cli.StringFlag{
				Name:        "repository-name",
				Aliases:     []string{"n"},
				Usage:       "Name of the repository the code of the application is stored in. Used for deployment annotations.",
				DefaultText: "name of the repository the command is run in",
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"D"},
				Usage:   "Simulate the deployment without actually deploying.",
				Sources: cli.EnvVars("3LV_DRY_RUN"),
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
				Name:    "add-deployment-annotation",
				Usage:   "Add a deployment annotation to Grafana. Requires --grafana-url and --grafana-api-key to be set.",
				Sources: cli.EnvVars("3LV_ADD_DEPLOYMENT_ANNOTATION"),
			},
			&cli.StringFlag{
				Name:    "grafana-url",
				Usage:   "The Grafana URL to use for deployment annotations.",
				Sources: cli.EnvVars("3LV_GRAFANA_URL"),
			},
			&cli.StringFlag{
				Name:    "grafana-api-key",
				Usage:   "The Grafana API key to use for deployment annotations.",
				Sources: cli.EnvVars("3LV_GRAFANA_API_KEY"),
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
}

//nolint:gocyclo
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

	config, err := shared.GetConfig()
	if err != nil {
		style.PrintWarning(err.Error() + "\n")
	}

	configForApplication, err := config.GetConfigForApplication(applicationName)
	if !config.IsEmpty() && err != nil { // Ignore error if config is empty, will default to flags.
		style.PrintWarning(err.Error() + "\n")
	}

	systemName := utils.FirstNonEmpty(c.String("system-name"), config.System)
	helmValuesFile := func() string {
		if c.IsSet("helm-values-file") {
			return c.String("helm-values-file")
		}

		return configForApplication.HelmValuesFile
	}()
	imageTag := c.String("image-tag")
	// NOTE: Not fully implemented yet.
	// imageDigest := getImageDigest(systemName, applicationName, imageTag)
	const imageDigest = ""

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
		return cli.Exit("Grafana URL and API key must be set when adding a deployment annotation.", 1)
	}

	environment := strings.ToLower(c.String("environment"))
	workloadType := strings.ToLower(c.String("workload-type"))
	runtimeCloudProvider := strings.ToLower(c.String("runtime-cloud-provider"))
	dryRun := c.Bool("dry-run")
	runID := c.String("run-id")
	helmChartRepositoryURL := c.String("helm-chart-repository-url")

	if checkKubectlInstalledOutput := checkKubectlInstalledCommand(nil); command.IsError(checkKubectlInstalledOutput) {
		return cli.Exit(fmt.Errorf("kubectl is not installed: %w", checkKubectlInstalledOutput.Error), 1)
	}

	if checkHelmInstalledOutput := checkHelmInstalledCommand(nil); command.IsError(checkHelmInstalledOutput) {
		return cli.Exit(fmt.Errorf("helm is not installed: %w", checkHelmInstalledOutput.Error), 1)
	}

	err = setupKubernetes(c, runtimeCloudProvider, environment)
	if err != nil {
		return cli.Exit(err, 1)
	}

	if helmRepoAddOutput := helmRepoAddCommand(helmChartRepositoryURL, nil); command.IsError(helmRepoAddOutput) {
		return cli.Exit(fmt.Errorf("Failed to add Helm repository: %w", helmRepoAddOutput.Error), 1)
	}

	if helmRepoUpdateOutput := helmRepoUpdateCommand(nil); command.IsError(helmRepoUpdateOutput) {
		return cli.Exit(fmt.Errorf("Failed to update Helm repository: %w", helmRepoUpdateOutput.Error), 1)
	}

	useISSChart := runtimeCloudProvider == "iss"

	if helmDeployOutput := helmDeployCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		imageDigest,
		repositoryName,
		commitHash,
		dryRun,
		useISSChart,
		nil,
	); command.IsError(helmDeployOutput) {
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

	// Jobs do not have a rollout status.
	if workloadType != "job" {
		if kubectlRolloutStatusOutput := kubectlRolloutStatusCommand(
			applicationName,
			systemName,
			workloadType,
			nil,
		); command.IsError(kubectlRolloutStatusOutput) {
			return cli.Exit(kubectlRolloutStatusOutput.Error, 1)
		}
	}

	kubectlGetEventsOutput := kubectlGetEventsCommand(
		systemName,
		nil,
	)
	if command.IsError(kubectlGetEventsOutput) {
		return cli.Exit(kubectlGetEventsOutput.Error, 1)
	}

	events := strings.SplitSeq(kubectlGetEventsOutput.Output, "\n")
	for event := range events {
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

func setupKubernetes(
	c *cli.Command,
	runtimeCloudProvider string,
	environment string,
) error {
	if runtimeCloudProvider == "aks" {
		return setupAKS(
			utils.StringWithDefault(
				c.String("azure-tenant-id"),
				auth.ElviaTenantID,
			),
			environment,
			&SetupAKSOptions{
				SubscriptionID:    c.String("aks-subscription-id"),
				ClusterName:       c.String("aks-cluster-name"),
				ResourceGroupName: c.String("aks-resource-group-name"),
				AzLoginOptions: &auth.AzLoginCommandOptions{
					ClientID:       c.String("azure-client-id"),
					FederatedToken: c.String("azure-federated-token"),
				},
			},
		)
	} else if runtimeCloudProvider == "gke" {
		return setupGKE(
			environment,
			&SetupGKEOptions{
				ProjectID:       c.String("gke-project-id"),
				ClusterName:     c.String("gke-cluster-name"),
				ClusterLocation: c.String("gke-cluster-location"),
			},
		)
	}

	return nil
}

func dockerManifestInspectCommand(
	imageNameWithTag string,
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"docker",
			"manifest",
			"inspect",
			imageNameWithTag,
			"-v",
		),
		runOptions,
	)
}

// Incomplete, since we only need digest.
type DockerManifest struct {
	Ref        string `json:"Ref"`
	Descriptor struct {
		Digest string `json:"digest"`
	} `json:"Descriptor"`
}

func getImageDigest( //nolint:unused
	systemName string,
	applicationName string,
	imageTag string,
) string {
	imageName, err := build.GetImageName(
		build.DefaultElviaContainerRegistry,
		systemName,
		applicationName,
	)
	if err != nil {
		return ""
	}

	dockerManifestInspectCommand := dockerManifestInspectCommand(imageName+":"+imageTag, nil)
	if command.IsError(dockerManifestInspectCommand) {
		return ""
	}

	var manifest DockerManifest

	err = json.Unmarshal([]byte(dockerManifestInspectCommand.Output), &manifest)
	if err != nil {
		return ""
	}

	return manifest.Descriptor.Digest
}
