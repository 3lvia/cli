package deploy

import (
	"fmt"
	"log"
	"os"
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

	helmOptions := HelmDeployOptions{
		ApplicationName: applicationName,
		SystemName:      systemName,
		HelmValuesFile:  helmValuesFile,
		Environment:     environment,
		WorkloadType:    workloadType,
		ImageTag:        imageTag,
		RepositoryName:  repositoryName,
		CommitHash:      commitHash,
	}
	if err := helmDeploy(helmOptions); err != nil {
		return cli.Exit(err, 1)
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

type AuthenticateGKEOptions struct {
	GKEProjectID       string
	GKEClusterName     string
	GKEClusterLocation string
}

func authenticateGKE(environment string, options AuthenticateGKEOptions) error {
	gkeProjectID := func() string {
		if options.GKEProjectID == "" {
			return "elvia-runtimeservice-" + environment
		}

		return options.GKEProjectID
	}()

	gkeClusterName := func() string {
		if options.GKEClusterName == "" {
			return "runtimeservice-gke-" + environment
		}

		return options.GKEClusterName
	}()

	gkeClusterLocation := func() string {
		if options.GKEClusterLocation == "" {
			return "europe-west1"
		}

		return options.GKEClusterLocation
	}()

	gcloudGetCredentialsCmd := exec.Command(
		"gcloud",
		"container",
		"clusters",
		"get-credentials",
		gkeClusterName,
		"--region",
		gkeClusterLocation,
		"--project",
		gkeProjectID,
	)
	gcloudGetCredentialsCmd.Stdout = os.Stdout
	gcloudGetCredentialsCmd.Stderr = os.Stderr

	if err := gcloudGetCredentialsCmd.Run(); err != nil {
		return fmt.Errorf("Failed to authenticate to GKE: %w", err)
	}

	return nil
}

type HelmDeployOptions struct {
	ApplicationName string // required
	SystemName      string // required
	HelmValuesFile  string // required
	Environment     string // required
	WorkloadType    string // required
	ImageTag        string // required
	RepositoryName  string // required
	CommitHash      string // required
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
		"elvia-charts/elvia-"+options.WorkloadType,
		"--set",
		"environment="+options.Environment,
		"--set",
		"image.tag="+options.ImageTag,
		"--set",
		"labels.repositoryName="+options.RepositoryName,
		"--set",
		"labels.commitHash="+options.CommitHash,
	)
	fmt.Println(helmDeployCmd.String())
	helmDeployCmd.Stdout = os.Stdout
	helmDeployCmd.Stderr = os.Stderr

	if err := helmDeployCmd.Run(); err != nil {
		return fmt.Errorf("Failed to deploy Helm chart: %w", err)
	}

	rolloutStatusCmd := exec.Command(
		"kubectl",
		"rollout",
		"status",
		"-n",
		options.SystemName,
		options.WorkloadType+"/"+options.ApplicationName,
	)
	rolloutStatusCmd.Stdout = os.Stdout
	rolloutStatusCmd.Stderr = os.Stderr

	if err := rolloutStatusCmd.Run(); err != nil {
		return fmt.Errorf("Failed to check rollout status: %w", err)
	}

	pipe := "kubectl get events -n " + options.SystemName + " --sort-by .lastTimestamp | grep " + options.ApplicationName
	eventsCmd := exec.Command(
		"bash",
		"-c",
		pipe,
	)
	eventsCmd.Stdout = os.Stdout
	eventsCmd.Stderr = os.Stderr

	if err := eventsCmd.Run(); err != nil {
		return fmt.Errorf("Failed to get events: %w", err)
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
