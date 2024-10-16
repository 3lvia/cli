package deploy

import (
	"fmt"
	"os/exec"

	"github.com/3lvia/cli/pkg/command"
)

const (
	chartsNamespace     = "elvia-charts"
	chartsRepositoryURL = "https://raw.githubusercontent.com/3lvia/kubernetes-charts/master"
)

func checkHelmInstalledCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command("helm", "version"),
		runOptions,
	)
}

func helmRepoAddCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"helm",
			"repo",
			"add",
			chartsNamespace,
			chartsRepositoryURL,
		),
		runOptions,
	)
}

func helmRepoUpdateCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"helm",
			"repo",
			"update",
		),
		runOptions,
	)
}

func helmDeployCommand(
	applicationName string,
	systemName string,
	helmValuesFile string,
	environment string,
	workloadType string,
	imageTag string,
	repositoryName string,
	commitHash string,
	dryRun bool,
	runOptions *command.RunOptions,
) command.Output {
	if workloadType != "deployment" && workloadType != "statefulset" {
		return command.Error(
			fmt.Errorf("workloadType must be either deployment or statefulset, got %s", workloadType),
		)
	}

	cmd := exec.Command(
		"helm",
		"upgrade",
		"--debug",
		"--install",
		"-n",
		systemName,
		"-f",
		helmValuesFile,
		applicationName,
		chartsNamespace+"/elvia-"+workloadType,
		"--set-string",
		"environment="+environment,
		"--set-string",
		"image.tag="+imageTag,
		"--set-string",
		"labels.repositoryName="+repositoryName,
		"--set-string",
		"labels.commitHash=\""+commitHash+"\"",
	)

	if dryRun {
		cmd.Args = append(cmd.Args, "--dry-run")
	}

	return command.Run(*cmd, runOptions)
}
