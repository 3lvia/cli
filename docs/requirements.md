## 📋 Requirements

To use every part of the 3lv CLI, you need to have these dependencies installed:

- [git](https://git-scm.com/downloads): used for several commands.
- [Docker](https://docs.docker.com/engine/install): used for building containers.
- [Helm](https://helm.sh/docs/intro/install): used for deploying containers.
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl): used for deploying containers.
- [Trivy](https://aquasecurity.github.io/trivy): used for scanning containers.
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli): used for pushing containers to Azure Container Registry and deploying containers to Azure Kubernetes Service.
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install): used for deploying containers to Google Kubernetes Engine.
- [cookiecutter](https://cookiecutter.readthedocs.io/en/stable/installation.html): used for creating new projects; if [pipx](https://pipx.pypa.io/stable) is installed, the CLI will prompt you to install cookiecutter for you.

**Any of these dependencies can be skipped if you dont't use the subcommands that require them.**

### Pushing to registries

If you want to push to a registry, you need to be authenticated to that registry.

#### Azure Container Registry

This is Elvia's default registry.
The CLI will automatically log you in if you have the Azure CLI installed.

#### GitHub Container Registry

Install the [GitHub CLI](https://cli.github.com).

Use the following command (with your GitHub username) to login:

```bash
gh auth token | docker login ghcr.io --username your-github-username --password-stdin
```
