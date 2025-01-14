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
	helmChartRepositoryURL string,
	runOptions *command.RunOptions,
) command.Output {
	url := chartsRepositoryURL
	if helmChartRepositoryURL != "" {
		url = helmChartRepositoryURL
	}

	return command.Run(
		*exec.Command(
			"helm",
			"repo",
			"add",
			chartsNamespace,
			url,
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
	imageDigest string,
	repositoryName string,
	commitHash string,
	dryRun bool,
	useISSChart bool,
	runOptions *command.RunOptions,
) command.Output {
	if workloadType != "deployment" && workloadType != "statefulset" {
		return command.Error(
			fmt.Errorf("Workload type must be either deployment or statefulset, got '%s'.", workloadType),
		)
	}

	chartName, err := func() (string, error) {
		if useISSChart {
			if workloadType == "deployment" {
				return "iss-" + workloadType, nil
			}

			return "", fmt.Errorf("Workload type '%s' is not supported with ISS chart.", workloadType)
		}

		return "elvia-" + workloadType, nil
	}()
	if err != nil {
		return command.Error(err)
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
		chartsNamespace+"/"+chartName,
		"--set-string",
		"environment="+environment,
		"--set-string",
		"labels.repositoryName="+repositoryName,
		"--set-string",
		"labels.commitHash=\""+commitHash+"\"",
	)

	if imageTag != "" {
		cmd.Args = append(cmd.Args, "--set-string", "image.tag="+imageTag)
	}

	if imageDigest != "" {
		cmd.Args = append(cmd.Args, "--set-string", "image.digest="+imageDigest)
	}

	if dryRun {
		cmd.Args = append(cmd.Args, "--dry-run")
	}

	return command.Run(*cmd, runOptions)
}
