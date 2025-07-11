package githubactions

import (
	"context"
	"embed"
	"errors"
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

//go:embed values.yaml.tmpl
var helmValuesFileTemplate embed.FS

func Command() *cli.Command {
	return &cli.Command{
		Name:      commandName,
		Aliases:   []string{"gha"},
		Usage:     "Add build and deploy with GitHub Actions to an exisiting project.",
		UsageText: "3lv github-actions [options] <project-directory>",
		Flags: []cli.Flag{
			shared.SystemNameFlag(
				"The name of your system (Kubernetes namespace) you want to deploy to.",
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

	config, err := shared.GetConfig()
	if err != nil {
		style.PrintWarning(err.Error() + "\n")
	}

	applicationName := c.String("application-name")

	configForApplication, err := config.GetConfigForApplication(applicationName)
	if !config.IsEmpty() && err != nil { // Ignore error if config is empty, will default to flags.
		style.PrintWarning(err.Error() + "\n")
	}

	projectFile := utils.FirstNonEmpty(c.String("project-file"), configForApplication.ProjectFile)
	if projectFile == "" {
		return cli.Exit("Project file not provided.", 1)
	}

	systemName := utils.FirstNonEmpty(c.String("system-name"), config.System)
	helmValuesFile := utils.FirstNonEmpty(c.String("helm-values-file"), configForApplication.HelmValuesFile)

	_, err = CreateDeployWorkflow(
		ctx,
		projectDirectory,
		projectFile,
		c.String("runtime-cloud-provider"),
		systemName,
		applicationName,
		helmValuesFile,
		c.String("default-branch"),
		c.Bool("non-interactive"),
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	style.PrintSuccess("Successfully added GitHub Actions to the project!\n")

	return nil
}

func CreateDeployWorkflow(
	ctx context.Context,
	outputDirectory string,
	projectFile string,
	runtimeCloudProvider string,
	systemName string,
	applicationName string,
	helmValuesFile string,
	defaultBranch string,
	nonInteractive bool,
) (string, error) {
	const githubActionsDir = ".github/workflows"

	fullGithubActionsDir := path.Join(outputDirectory, githubActionsDir)

	if _, err := os.Stat(fullGithubActionsDir); os.IsNotExist(err) {
		style.PrintInfo(fmt.Sprintf("Creating directory '%s'.\n", githubActionsDir))

		if err := os.MkdirAll(fullGithubActionsDir, 0o755); err != nil {
			return "", fmt.Errorf("Failed to create directory '%s'", fullGithubActionsDir)
		}
	}

	language, err := getLanguageFromProjectFile(path.Base(projectFile))
	if err != nil {
		return "", err
	}

	resolvedHelmValuesFile, err := resolveHelmValuesFile(
		applicationName,
		systemName,
		&ResolveHelmValuesFileOptions{
			OutputDirectory: outputDirectory,
			HelmValuesFile:  helmValuesFile,
			NonInteractive:  nonInteractive,
		},
	)
	if err != nil {
		return "", err
	}

	exampleWorkflowFileURL, err := getExampleWorkflowFileURL(language, runtimeCloudProvider)
	if err != nil {
		return "", err
	}

	workflowFileName := fmt.Sprintf("build-deploy-%s.yaml", applicationName)
	workflowFilePath := filepath.Join(fullGithubActionsDir, workflowFileName)

	if err := downloadFile(ctx, exampleWorkflowFileURL, workflowFilePath); err != nil {
		return "", err
	}

	style.PrintWarning(
		fmt.Sprintf(
			"Replacing placeholders in workflow file '%s'.\nYou may need to manually fill in some values yourself.\n",
			workflowFilePath,
		),
	)

	replaceWorkflowPlaceholdersOptions := &ReplaceWorkflowPlaceholdersOptions{
		SystemName:      systemName,
		ApplicationName: applicationName,
		HelmValuesFile:  resolvedHelmValuesFile,
		DefaultBranch:   defaultBranch,
	}
	if err := replaceWorkflowPlaceholders(
		workflowFilePath,
		projectFile,
		replaceWorkflowPlaceholdersOptions,
	); err != nil {
		return "", err
	}

	terraformReminder := func() string {
		if runtimeCloudProvider == "iss" {
			return "NOTE: if you have not done so already, you will need to add your repository to the Terraform module" +
				" 'github-actions-deploy' at https://github.com/3lvia/iss-terraform to enable deployments from GitHub Actions."
		}

		return "NOTE: if you have not done so already, you will need to add your system/repository to" +
			" https://github.com/3lvia/github-repositories-terraform to enable deployments from GitHub Actions."
	}()
	style.PrintWarning(terraformReminder + "\n")

	return helmValuesFile, nil
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
		options = &ReplaceWorkflowPlaceholdersOptions{} //nolint:exhaustruct
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

		contentsString = func() string {
			lines := strings.Split(contentsString, "\n")
			if len(lines) > 0 {
				lines[0] = "name: Build and deploy " + options.ApplicationName
			}

			return strings.Join(lines, "\n")
		}()
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
			".github/deploy/values.yaml",
			options.HelmValuesFile,
		)
	}

	// For analyze job
	contentsString = strings.ReplaceAll(
		contentsString,
		fmt.Sprintf(
			"# This can be set to a more specific path if you want to analyze only a part of the repository."+
				"\n%sworking-directory: '.'",
			strings.Repeat(" ", 10),
		),
		fmt.Sprintf(
			"# This can be set to a more specific path if you want to analyze only a part of the repository."+
				"\n%sworking-directory: '%s'",
			strings.Repeat(" ", 10),
			path.Dir(projectFile),
		),
	)

	// For integration-tests job
	contentsString = strings.ReplaceAll(
		contentsString,
		fmt.Sprintf(
			"# This can be set to a more specific path if you want to search for tests in only a part of the repository."+
				"\n%sworking-directory: '.'",
			strings.Repeat(" ", 10),
		),
		fmt.Sprintf(
			"# This can be set to a more specific path if you want to search for tests in only a part of the repository."+
				"\n%sworking-directory: '%s'",
			strings.Repeat(" ", 10),
			path.Dir(projectFile),
		),
	)

	contentsString = "# This file was generated by the 3lvia CLI: https://github.com/3lvia/cli\n\n" + contentsString

	if err := os.WriteFile(workflowFilePath, []byte(contentsString), 0o644); err != nil {
		return err
	}

	return nil
}

func downloadFile(ctx context.Context, url string, outputFilePath string) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(req)
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

	if projectFile == "pyproject.toml" {
		return "python", nil
	}

	if strings.Contains(projectFile, "Dockerfile") {
		return "dockerfile", nil
	}

	return "", fmt.Errorf("Unsupported project file '%s'", projectFile)
}

func getExampleWorkflowFileURL(language string, runtimeCloudProvider string) (string, error) {
	// .NET
	if language == "dotnet" && runtimeCloudProvider == "aks" {
		return exampleWorkflowBaseURL + "/build-deploy-dotnet.yaml", nil
	}

	if language == "dotnet" && runtimeCloudProvider == "gke" {
		return exampleWorkflowBaseURL + "/build-deploy-dotnet-google.yaml", nil
	}

	if language == "dotnet" && runtimeCloudProvider == "iss" {
		return exampleWorkflowBaseURL + "/build-deploy-dotnet-iss.yaml", nil
	}

	// Go
	if language == "go" && runtimeCloudProvider == "aks" {
		return exampleWorkflowBaseURL + "/build-deploy-go.yaml", nil
	}

	if language == "go" && runtimeCloudProvider == "gke" {
		return exampleWorkflowBaseURL + "/build-deploy-go-google.yaml", nil
	}

	if language == "go" && runtimeCloudProvider == "iss" {
		return "",
			fmt.Errorf("Example workflow is not implemented yet for language '%s' and runtime cloud provider '%s'",
				language,
				runtimeCloudProvider,
			)
	}

	// Python
	if language == "python" && runtimeCloudProvider == "aks" {
		return exampleWorkflowBaseURL + "/build-deploy-python.yaml", nil
	}

	if language == "python" && runtimeCloudProvider == "gke" {
		return exampleWorkflowBaseURL + "/build-deploy-python-google.yaml", nil
	}

	if language == "python" && runtimeCloudProvider == "iss" {
		return "",
			fmt.Errorf("Example workflow is not implemented yet for language '%s' and runtime cloud provider '%s'",
				language,
				runtimeCloudProvider,
			)
	}

	// Dockerfile
	if language == "dockerfile" && runtimeCloudProvider == "aks" {
		return exampleWorkflowBaseURL + "/build-deploy-dockerfile.yaml", nil
	}

	if language == "dockerfile" && runtimeCloudProvider == "gke" {
		return exampleWorkflowBaseURL + "/build-deploy-dockerfile-google.yaml", nil
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
		options = &ResolveHelmValuesFileOptions{} //nolint:exhaustruct
	}

	if options.HelmValuesFile == "" {
		defaultHelmValuesFile := fmt.Sprintf(".github/deploy/values-%s.yaml", applicationName)

		yes, err := utils.PromptYesNo(
			"You have not provided a Helm values file, which is required for the deployment."+
				" Do you want to create a default Helm values file?",
			options.NonInteractive,
		)
		if err != nil {
			return "", err
		}

		if !yes {
			return "", errors.New("Helm values file is required for the deployment")
		}

		if err := os.MkdirAll(
			filepath.Dir(
				filepath.Join(options.OutputDirectory, defaultHelmValuesFile),
			),
			0o755,
		); err != nil {
			return "", fmt.Errorf("Failed to create directory for Helm values file: %w", err)
		}

		const templateFile = "values.yaml.tmpl"

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

		style.PrintInfo(fmt.Sprintf("Created Helm values file at '%s'.\n", newHelmValuesFile))

		return defaultHelmValuesFile, nil
	}

	return options.HelmValuesFile, nil
}
