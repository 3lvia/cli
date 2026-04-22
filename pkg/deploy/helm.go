package deploy

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/3lvia/cli/pkg/command"
)

const (
	chartsRepositoryURL = "oci://ghcr.io/3lvia/charts"
)

func checkHelmInstalledCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command("helm", "version"),
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
	if workloadType != "deployment" && workloadType != "statefulset" && workloadType != "job" {
		return command.Error(
			fmt.Errorf("Workload type must be either deployment, statefulset or job, got '%s'.", workloadType),
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
		chartsRepositoryURL+"/"+chartName,
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

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		cmd.Args = append(cmd.Args, "--set-string", "labels.deployedBy=github-actions")
	}

	if dryRun {
		cmd.Args = append(cmd.Args, "--dry-run")
	}

	return command.Run(*cmd, runOptions)
}
