package deploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func checkHelmInstalled() error {
	if err := exec.Command("helm", "version").Run(); err != nil {
		return fmt.Errorf("Helm is not installed: %w", err)
	}

	return nil
}

func helmRepoAdd(chartsNamespace string, chartsRepositoryURL string) error {
	helmRepoAddCmd := exec.Command(
		"helm",
		"repo",
		"add",
		chartsNamespace,
		chartsRepositoryURL,
	)

	if err := helmRepoAddCmd.Run(); err != nil {
		return fmt.Errorf("Failed to add Helm repository: %w", err)
	}

	return nil
}

func helmRepoUpdate() error {
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

	return nil
}

func helmDeploy(
	applicationName string,
	systemName string,
	helmValuesFile string,
	environment string,
	workloadType string,
	imageTag string,
	repositoryName string,
	commitHash string,
	dryRun bool,
) error {
	helmDeployCmd := exec.Command(
		"helm",
		"upgrade",
		"--debug",
		"--install",
		"-n",
		systemName,
		"-f",
		helmValuesFile,
		applicationName,
		"elvia-charts/elvia-"+workloadType,
		"--set-string",
		"environment="+environment,
		"--set-string",
		"image.tag="+imageTag,
		"--set-string",
		"labels.repositoryName="+repositoryName,
		"--set-string",
		"labels.commitHash=\""+commitHash+"\"",
	)
	helmDeployCmd.Stdout = os.Stdout
	helmDeployCmd.Stderr = os.Stderr

	if dryRun {
		helmDeployCmd.Args = append(helmDeployCmd.Args, "--dry-run")
	}

	log.Print(helmDeployCmd.String())

	if err := helmDeployCmd.Run(); err != nil {
		return fmt.Errorf("Failed to deploy Helm chart: %w", err)
	}

	return nil
}
