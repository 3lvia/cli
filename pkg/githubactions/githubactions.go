package githubactions

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v3"
)

const (
	commandName            = "github-actions"
	exampleWorkflowBaseURL = "https://raw.githubusercontent.com/3lvia/.github/refs/heads/trunk/workflow-templates"
)

//go:embed values.yml.tmpl
var helmValuesFileTemplate embed.FS

var Command *cli.Command = &cli.Command{
	Name:      commandName,
	Aliases:   []string{"gha"},
	Usage:     "Add build and deploy with GitHub Actions to an exisiting project.",
	UsageText: "3lv build [options] <project-directory>",
	Flags: []cli.Flag{
		shared.SystemNameFlag(
			"The name of your system (Kubernetes namespace) you want to deploy to.",
			true,
		),
		shared.ApplicationNameFlag(
			"The name of the application you want to build and deploy.",
		),
		shared.RuntimeCloudProviderFlag(),
		shared.HelmValuesFileFlag(),
		shared.ProjectFileFlag(),
		&cli.StringFlag{
			Name:    "default-branch",
			Aliases: []string{"b"},
			Usage:   "The default branch of the repository",
			Value:   "trunk",
		},
	},
	Action: GitHubActions,
}

func GitHubActions(ctx context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	projectDirectory := func() string {
		first := c.Args().First()
		if first == "" {
			return "."
		}

		return first
	}()

	err := CreateDeployWorkflow(
		projectDirectory,
		c.String("project-file"),
		c.String("runtime-cloud-provider"),
		c.String("system-name"),
		c.String("application-name"),
		c.String("helm-values-file"),
		c.String("default-branch"),
		c.Bool("non-interactive"),
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	style.Print(
		"Successfully added GitHub Actions to the project!\n",
		&style.PrintOptions{Color: "green"},
	)

	return nil
}

func CreateDeployWorkflow(
	outputDirectory string,
	projectFile string,
	runtimeCloudProvider string,
	systemName string,
	applicationName string,
	helmValuesFile string,
	defaultBranch string,
	nonInteractive bool,
) error {
	const githubActionsDir = ".github/workflows"
	fullGithubActionsDir := path.Join(outputDirectory, githubActionsDir)

	if _, err := os.Stat(fullGithubActionsDir); os.IsNotExist(err) {
		style.Print(
			fmt.Sprintf("Creating directory '%s'.\n", githubActionsDir),
			nil,
		)
		if err := os.MkdirAll(fullGithubActionsDir, 0755); err != nil {
			return fmt.Errorf("Failed to create directory '%s'.", fullGithubActionsDir)
		}
	}

	language, err := getLanguageFromProjectFile(path.Base(projectFile))
	if err != nil {
		return err
	}

	helmValuesFile_, err := resolveHelmValuesFile(
		applicationName,
		systemName,
		&ResolveHelmValuesFileOptions{
			OutputDirectory: outputDirectory,
			HelmValuesFile:  helmValuesFile,
			NonInteractive:  nonInteractive,
		},
	)
	if err != nil {
		return err
	}

	exampleWorkflowFileURL, err := getExampleWorkflowFileURL(language, runtimeCloudProvider)
	if err != nil {
		return err
	}

	workflowFileName := fmt.Sprintf("build-deploy-%s.yml", applicationName)
	workflowFilePath := filepath.Join(fullGithubActionsDir, workflowFileName)

	if err := downloadFile(exampleWorkflowFileURL, workflowFilePath); err != nil {
		return err
	}

	style.Print(
		fmt.Sprintf(
			"Replacing placeholders in workflow file '%s'.\nYou may need to manually fill in some values yourself.\n",
			workflowFilePath,
		),
		&style.PrintOptions{Color: "yellow"},
	)
	replaceWorkflowPlaceholdersOptions := &ReplaceWorkflowPlaceholdersOptions{
		SystemName:      systemName,
		ApplicationName: applicationName,
		HelmValuesFile:  helmValuesFile_,
		DefaultBranch:   defaultBranch,
	}
	if err := replaceWorkflowPlaceholders(
		workflowFilePath,
		projectFile,
		replaceWorkflowPlaceholdersOptions,
	); err != nil {
		return err
	}

	terraformReminder := func() string {
		if runtimeCloudProvider == "iss" {
			return "NOTE: if you have not done so already, you will need to add your repository to the Terraform module 'github-actions-deploy' at https://github.com/3lvia/iss-terraform to enable deployments from GitHub Actions."
		}
		return "NOTE: if you have not done so already, you will need to add your system/repository to https://github.com/3lvia/github-repositories-terraform to enable deployments from GitHub Actions."
	}()
	style.Print(
		fmt.Sprintf("%s\n", terraformReminder),
		&style.PrintOptions{Color: "yellow"},
	)

	return nil
}

type ReplaceWorkflowPlaceholdersOptions struct {
	DefaultBranch   string
	SystemName      string
	ApplicationName string
	HelmValuesFile  string
}

func replaceWorkflowPlaceholders(
	workflowFilePath string,
	projectFile string,
	options *ReplaceWorkflowPlaceholdersOptions,
) error {
	if options == nil {
		options = &ReplaceWorkflowPlaceholdersOptions{}
	}

	file, err := os.Open(workflowFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	contentsString := string(contents)

	defaultBranch := utils.StringWithDefault(options.DefaultBranch, "trunk")
	contentsString = strings.ReplaceAll(contentsString, "$default-branch", defaultBranch)

	contentsString = strings.ReplaceAll(
		contentsString,
		"<your project file path here>",
		projectFile,
	)

	if options.ApplicationName != "" {
		contentsString = strings.ReplaceAll(
			contentsString,
			"<your application name here>",
			options.ApplicationName,
		)
	}

	if options.SystemName != "" {
		contentsString = strings.ReplaceAll(
			contentsString,
			"<your system name here>",
			options.SystemName,
		)
	}

	if options.HelmValuesFile != "" {
		contentsString = strings.ReplaceAll(
			contentsString,
			".github/deploy/values.yml",
			options.HelmValuesFile,
		)
	}

	// For analyze job
	contentsString = strings.ReplaceAll(
		contentsString,
		fmt.Sprintf(
			"# This can be set to a more specific path if you want to analyze only a part of the repository.\n%sworking-directory: '.'",
			strings.Repeat(" ", 10),
		),
		fmt.Sprintf(
			"# This can be set to a more specific path if you want to analyze only a part of the repository.\n%sworking-directory: '%s'",
			strings.Repeat(" ", 10),
			path.Dir(projectFile),
		),
	)

	// For integration-tests job
	contentsString = strings.ReplaceAll(
		contentsString,
		fmt.Sprintf(
			"# This can be set to a more specific path if you want to search for tests in only a part of the repository.\n%sworking-directory: '.'",
			strings.Repeat(" ", 10),
		),
		fmt.Sprintf(
			"# This can be set to a more specific path if you want to search for tests in only a part of the repository.\n%sworking-directory: '%s'",
			strings.Repeat(" ", 10),
			path.Dir(projectFile),
		),
	)

	contentsString = fmt.Sprintf(
		"# This file was generated by the 3lvia CLI: https://github.com/3lvia/cli\n\n%s",
		contentsString,
	)

	if err := os.WriteFile(workflowFilePath, []byte(contentsString), 0644); err != nil {
		return err
	}

	return nil
}

func downloadFile(url string, outputFilePath string) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getLanguageFromProjectFile(projectFile string) (string, error) {
	if strings.HasSuffix(projectFile, ".csproj") {
		return "dotnet", nil
	}

	if projectFile == "go.mod" {
		return "go", nil
	}

	if strings.Contains(projectFile, "Dockerfile") {
		return "dockerfile", nil
	}

	return "", fmt.Errorf("Unsupported project file '%s'", projectFile)
}

func getExampleWorkflowFileURL(language string, runtimeCloudProvider string) (string, error) {
	// .NET
	if language == "dotnet" && runtimeCloudProvider == "aks" {
		return fmt.Sprintf("%s/build-deploy-dotnet.yml", exampleWorkflowBaseURL), nil
	}
	if language == "dotnet" && runtimeCloudProvider == "gke" {
		return fmt.Sprintf("%s/build-deploy-dotnet-google.yml", exampleWorkflowBaseURL), nil
	}
	if language == "dotnet" && runtimeCloudProvider == "iss" {
		return fmt.Sprintf("%s/build-deploy-dotnet-iss.yml", exampleWorkflowBaseURL), nil
	}

	// Go
	if language == "go" && runtimeCloudProvider == "aks" {
		return fmt.Sprintf("%s/build-deploy-go.yml", exampleWorkflowBaseURL), nil
	}
	if language == "go" && runtimeCloudProvider == "gke" {
		return fmt.Sprintf("%s/build-deploy-go-google.yml", exampleWorkflowBaseURL), nil
	}
	if language == "go" && runtimeCloudProvider == "iss" {
		return "",
			fmt.Errorf("Example workflow is not implemented yet for language '%s' and runtime cloud provider '%s'",
				language,
				runtimeCloudProvider,
			)
	}

	// Dockerfile
	if language == "dockerfile" && runtimeCloudProvider == "aks" {
		return fmt.Sprintf("%s/build-deploy-dockerfile.yml", exampleWorkflowBaseURL), nil
	}
	if language == "dockerfile" && runtimeCloudProvider == "gke" {
		return fmt.Sprintf("%s/build-deploy-dockerfile-google.yml", exampleWorkflowBaseURL), nil
	}
	if language == "dockerfile" && runtimeCloudProvider == "iss" {
		return "",
			fmt.Errorf("Example workflow is not implemented yet for language '%s' and runtime cloud provider '%s'",
				language,
				runtimeCloudProvider,
			)
	}

	return "",
		fmt.Errorf(
			"No example workflow file found for language '%s' and runtime cloud provider '%s'",
			language,
			runtimeCloudProvider,
		)
}

type ResolveHelmValuesFileOptions struct {
	OutputDirectory string
	HelmValuesFile  string
	NonInteractive  bool
}

func resolveHelmValuesFile(
	applicationName string,
	systemName string,
	options *ResolveHelmValuesFileOptions,
) (string, error) {
	if options == nil {
		options = &ResolveHelmValuesFileOptions{}
	}
	if options.HelmValuesFile == "" {
		defaultHelmValuesFile := fmt.Sprintf(".github/deploy/values-%s.yml", applicationName)

		yes, err := utils.PromptYesNo(
			"You have not provided a Helm values file, which is required for the deployment. Do you want to create a default Helm values file?",
			options.NonInteractive,
		)
		if err != nil {
			return "", err
		}

		if !yes {
			return "", fmt.Errorf("Helm values file is required for the deployment.")
		}

		if err := os.MkdirAll(
			filepath.Dir(
				filepath.Join(options.OutputDirectory, defaultHelmValuesFile),
			),
			0755,
		); err != nil {
			return "", fmt.Errorf("Failed to create directory for Helm values file: %w", err)
		}

		const templateFile = "values.yml.tmpl"
		newHelmValuesFile, err := utils.WriteFileWithTemplate(
			options.OutputDirectory,
			defaultHelmValuesFile,
			templateFile,
			helmValuesFileTemplate,
			struct {
				ApplicationName string
				SystemName      string
			}{
				ApplicationName: applicationName,
				SystemName:      systemName,
			},
		)
		if err != nil {
			return "", err
		}

		style.Print(
			fmt.Sprintf("Created Helm values file at '%s'.\n", newHelmValuesFile),
			nil,
		)

		return defaultHelmValuesFile, nil
	}

	return options.HelmValuesFile, nil
}
