package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/style"
	"gopkg.in/yaml.v3"
)

const acrHostSuffix = ".azurecr.io"

// Why the Deployment image is pinned by digest:
//
// Admission webhooks such as sigstore policy-controller rewrite image tags to tag@digest on ReplicaSets.
// The Deployment controller only reuses a ReplicaSet whose pod template equals the Deployment's, so if
// the Deployment template still has the bare tag it never matches, and a new ReplicaSet is created
// (and rewritten) forever. The webhook also rewrites Deployments, but skips that on updates that lower
// spec.replicas (which a helm upgrade does whenever the HPA has scaled up) and on registry lookup errors.
// Resolving the digest here makes the Deployment template carry it from the start, for every repository
// using 3lv, without workflow changes.
//
// Resolution is best effort: on any failure we warn and deploy by tag only, exactly as before.

// deploymentImageDigest returns the digest to pin the Deployment image to: the one given explicitly, or else
// the one resolved from the registry. Only AKS deploys are logged in to Azure, which the registry lookup needs.
func deploymentImageDigest(
	imageDigest string,
	resolve bool,
	workloadType string,
	runtimeCloudProvider string,
	applicationName string,
	systemName string,
	helmValuesFile string,
	environment string,
	imageTag string,
	repositoryName string,
	commitHash string,
	useISSChart bool,
) string {
	if imageDigest != "" || !resolve || workloadType != "deployment" || runtimeCloudProvider != "aks" {
		return imageDigest
	}

	return resolveDeploymentImageDigest(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		imageTag,
		repositoryName,
		commitHash,
		useISSChart,
	)
}

// resolveDeploymentImageDigest renders the chart, finds the main container's image and asks Azure
// Container Registry which digest its tag points to. Returns "" when the digest can't or shouldn't be resolved.
func resolveDeploymentImageDigest(
	applicationName string,
	systemName string,
	helmValuesFile string,
	environment string,
	imageTag string,
	repositoryName string,
	commitHash string,
	useISSChart bool,
) string {
	templateOutput := helmTemplateCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		"deployment",
		imageTag,
		repositoryName,
		commitHash,
		useISSChart,
		&command.RunOptions{Quiet: true},
	)
	if command.IsError(templateOutput) {
		style.PrintWarning(fmt.Sprintf(
			"Could not render the chart to resolve the image digest, deploying by tag only: %v\n",
			templateOutput.Error,
		))

		return ""
	}

	image, err := mainContainerImage(templateOutput.Output)
	if err != nil {
		style.PrintWarning(fmt.Sprintf(
			"Could not find the image to resolve the digest for, deploying by tag only: %v\n",
			err,
		))

		return ""
	}

	registryName, repositoryWithTag, ok := splitACRImage(image)
	if !ok {
		style.PrintInfo(fmt.Sprintf("Not resolving digest for %s: not a tag reference in Azure Container Registry.\n", image))

		return ""
	}

	digestOutput := azAcrImageDigestCommand(registryName, repositoryWithTag, &command.RunOptions{Quiet: true})
	if command.IsError(digestOutput) {
		style.PrintWarning(fmt.Sprintf(
			"Could not resolve the digest of %s, deploying by tag only: %v %s\n",
			image,
			digestOutput.Error,
			strings.TrimSpace(digestOutput.Output),
		))

		return ""
	}

	digest := strings.TrimSpace(digestOutput.Output)
	if !strings.HasPrefix(digest, "sha256:") {
		style.PrintWarning(fmt.Sprintf("Unexpected digest %q for %s, deploying by tag only.\n", digest, image))

		return ""
	}

	style.PrintInfo(fmt.Sprintf("Pinning %s to %s.\n", image, digest))

	return digest
}

type renderedWorkload struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Template struct {
			Spec struct {
				Containers []struct {
					Name  string `yaml:"name"`
					Image string `yaml:"image"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

// mainContainerImage returns the image of the Deployment's main container in rendered chart output.
// The chart names the main container after the Deployment; sidecars have other names and are not pinned.
func mainContainerImage(renderedManifests string) (string, error) {
	decoder := yaml.NewDecoder(bytes.NewBufferString(renderedManifests))

	for {
		var workload renderedWorkload

		err := decoder.Decode(&workload)
		if errors.Is(err, io.EOF) {
			return "", errors.New("no Deployment found in rendered chart")
		}

		if err != nil {
			return "", fmt.Errorf("failed to parse rendered chart: %w", err)
		}

		if workload.Kind != "Deployment" {
			continue
		}

		for _, container := range workload.Spec.Template.Spec.Containers {
			if container.Name == workload.Metadata.Name {
				return container.Image, nil
			}
		}

		return "", fmt.Errorf("Deployment %q has no container named after it", workload.Metadata.Name)
	}
}

// splitACRImage splits "<registry>.azurecr.io/<repository>:<tag>" into registry name and "<repository>:<tag>".
// ok is false for other registries, for images without a tag and for images already pinned by digest.
func splitACRImage(image string) (string, string, bool) {
	host, repositoryWithTag, found := strings.Cut(image, "/")
	if !found || !strings.HasSuffix(host, acrHostSuffix) || strings.Contains(repositoryWithTag, "@") {
		return "", "", false
	}

	lastSlash := strings.LastIndex(repositoryWithTag, "/")

	colon := strings.LastIndex(repositoryWithTag, ":")
	if colon <= lastSlash || colon == len(repositoryWithTag)-1 {
		return "", "", false
	}

	return strings.TrimSuffix(host, acrHostSuffix), repositoryWithTag, true
}

func azAcrImageDigestCommand(
	registryName string,
	repositoryWithTag string,
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"az",
			"acr",
			"repository",
			"show",
			"--name",
			registryName,
			"--image",
			repositoryWithTag,
			"--query",
			"digest",
			"--output",
			"tsv",
		),
		runOptions,
	)
}
