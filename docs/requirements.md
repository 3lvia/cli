## 📋 Requirements

To use every part of the 3lv CLI, you need to have these dependencies installed:

- [Docker](https://docs.docker.com/engine/install): used for building
- [Helm](https://helm.sh/docs/intro/install): used for deploying
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl): used for deploying
- [Trivy](https://aquasecurity.github.io/trivy): used for scanning Docker images
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli): used for pushing to Azure Container Registry and deploying to Azure Kubernetes Service
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install): used for deploying to Google Kubernetes Engine
- [GitHub CLI](https://cli.github.com): used for pushing to GitHub Container Registry
- [cookiecutter](https://cookiecutter.readthedocs.io): used for creating new projects. If [pipx](https://pipxproject.github.io/pipx/) is installed, the CLI will prompt you to install cookiecutter for you.

**Any of these dependencies can be skipped if you dont't use the subcommands that require them.**

### Pushing to registries

If you want to push to a registry, you need to be authenticated to that registry.

#### Azure Container Registry

This is Elvia's default registry.
The CLI will automatically log you in if you have the Azure CLI installed.

#### GitHub Container Registry

Use the following command (with your GitHub username) to login:

```bash
gh auth token | docker login ghcr.io --username your-github-username --password-stdin
```
