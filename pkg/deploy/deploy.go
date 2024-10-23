package deploy

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"slices"
	"strings"

	"github.com/3lvia/cli/pkg/command"
	"github.com/urfave/cli/v2"
)

const commandName = "deploy"

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
			Name:    "commit-message",
			Aliases: []string{"m"},
			Usage:   "The commit message to use",
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
			Name:    "aks-tenant-id",
			Usage:   "The AKS tenant ID to use",
			Hidden:  true,
			EnvVars: []string{"3LV_AKS_TENANT_ID"},
		},
		&cli.StringFlag{
			Name:    "aks-subscription-id",
			Usage:   "The AKS subscription ID to use",
			Hidden:  true,
			EnvVars: []string{"3LV_AKS_SUBSCRIPTION_ID"},
		},
		&cli.StringFlag{
			Name:    "aks-cluster-name",
			Usage:   "The AKS cluster name to use",
			Hidden:  true,
			EnvVars: []string{"3LV_AKS_CLUSTER_NAME"},
		},
		&cli.StringFlag{
			Name:    "aks-resource-group-name",
			Usage:   "The AKS resource group name to use",
			Hidden:  true,
			EnvVars: []string{"3LV_AKS_RESOURCE_GROUP_NAME"},
		},
		&cli.StringFlag{
			Name:    "gke-project-id",
			Usage:   "The GKE project ID to use",
			Hidden:  true,
			EnvVars: []string{"3LV_GKE_PROJECT_ID"},
		},
		&cli.StringFlag{
			Name:    "gke-cluster-name",
			Usage:   "The GKE cluster name to use",
			Hidden:  true,
			EnvVars: []string{"3LV_GKE_CLUSTER_NAME"},
		},
		&cli.StringFlag{
			Name:    "gke-cluster-location",
			Usage:   "The GKE cluster location to use",
			Hidden:  true,
			EnvVars: []string{"3LV_GKE_CLUSTER_LOCATION"},
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
	},
	Action: Deploy,
}

func Deploy(c *cli.Context) error {
	if c.NArg() <= 0 {
		return cli.ShowCommandHelp(c, commandName)
	}

	applicationName := c.Args().First()
	if applicationName == "" {
		log.Println("Application name not provided")
		return cli.ShowCommandHelp(c, commandName)
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

	addDeploymentAnnotation := c.Bool("add-deployment-annotation")
	commitMessage, err := resolveCommitMessage(c.String("commit-message"))
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
	skipAuthentication := c.Bool("skip-authentication")
	dryRun := c.Bool("dry-run")
	runID := c.String("run-id")

	checkKubectlInstalledOutput := checkKubectlInstalledCommand(nil)
	if command.IsError(checkKubectlInstalledOutput) {
		return cli.Exit(fmt.Errorf("kubectl is not installed: %w", checkKubectlInstalledOutput.Error), 1)
	}

	checkHelmInstalledOutput := checkHelmInstalledCommand(nil)
	if command.IsError(checkHelmInstalledOutput) {
		return cli.Exit(fmt.Errorf("helm is not installed: %w", checkHelmInstalledOutput.Error), 1)
	}

	if runtimeCloudProvider == "aks" {
		setupOptions := SetupAKSOptions{
			AKSTenantID:          c.String("aks-tenant-id"),
			AKSSubscriptionID:    c.String("aks-subscription-id"),
			AKSClusterName:       c.String("aks-cluster-name"),
			AKSResourceGroupName: c.String("aks-resource-group-name"),
		}
		if err := setupAKS(environment, skipAuthentication, setupOptions); err != nil {
			return cli.Exit(err, 1)
		}

	} else if runtimeCloudProvider == "gke" {
		authOptions := SetupGKEOptions{
			GKEProjectID:       c.String("gke-project-id"),
			GKEClusterName:     c.String("gke-cluster-name"),
			GKEClusterLocation: c.String("gke-cluster-location"),
		}
		if err := setupGKE(environment, skipAuthentication, authOptions); err != nil {
			return cli.Exit(err, 1)
		}
	}

	helmRepoAddOutput := helmRepoAddCommand(nil)
	if command.IsError(helmRepoAddOutput) {
		return cli.Exit(fmt.Errorf("Failed to add Helm repository: %w", helmRepoAddOutput.Error), 1)
	}

	helmRepoUpdateOutput := helmRepoUpdateCommand(nil)
	if command.IsError(helmRepoUpdateOutput) {
		return cli.Exit(fmt.Errorf("Failed to update Helm repository: %w", helmRepoUpdateOutput.Error), 1)
	}

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
		nil,
	)
	if command.IsError(helmDeployOutput) {
		// If the deployment failed, we still want to post the Grafana annotation, but we add a failure message to the annotation.
		if err := addGrafanaDeploymentAnnotation(
			false,
			applicationName,
			systemName,
			environment,
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

	if addDeploymentAnnotation {
		if err := addGrafanaDeploymentAnnotation(
			true,
			applicationName,
			systemName,
			environment,
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

func resolveCommitMessage(possibleCommitMessage string) (string, error) {
	if possibleCommitMessage != "" {
		return possibleCommitMessage, nil
	}

	message, err := exec.Command("git", "log", "-1", "--no-merges", "--pretty=%B").Output()
	if err != nil {
		return "",
			fmt.Errorf(
				"Failed to resolve commit message: %w. Please verify you are currently in a Git repository, or manually specify the commit message with --commit-message.",
				err,
			)
	}

	return strings.TrimSpace(string(message)), nil
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
