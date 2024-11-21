# cli

Command Line Interface tool for developing, building and securing Elvia applications ⚡

## 💾 Installation

Supported platforms:

- **Windows**
- **macOS** (Intel and M-series)
- **Linux** (any distribution)

### Windows

Download the MSI file (and optionally the MD5 checksum) from the [releases page](https://github.com/3lvia/cli/releases) and run it.

#### Verify checksum

```pwsh
certutil -hashfile 3lv-windows-amd64.msi MD5
```

### Linux/WSL and macOS

Download the tarball file for your platform (and optionally the MD5 cheksum) from the [releases page](https://github.com/3lvia/cli/releases),
extract it and move the binary to a directory in your PATH.

#### Example installation

```bash
# Optional: verify checksum first
md5sum -c 3lv-linux-amd64.tar.gz.md5

tar -xzf 3lv-linux-amd64.tar.gz
sudo install -Dm755 -t /usr/local/bin 3lv
```

For macOS, you can use the same commands as above, but replace `linux` with `macos`.
If you have an M1 or newer mac, you can use the `macos-arm64` binary.

## 📋 Requirements

To use every part of the 3lv CLI, you need to have these dependencies installed:

- [Docker](https://docs.docker.com/engine/install): used for building
- [Helm](https://helm.sh/docs/intro/install): used for deploying
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl): used for deploying
- [Trivy](https://aquasecurity.github.io/trivy): used for scanning Docker images
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli): used for pushing to Azure Container Registry and deploying to Azure Kubernetes Service
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install): used for deploying to Google Kubernetes Engine
- [GitHub CLI](https://cli.github.com): used for pushing to GitHub Container Registry

**Any of these dependencies can be skipped if you dont't use the subcommands that require them.**

## ❓ Usage

```bash
3lv --help
```

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

## 🚀 Breaking changes

Before version `v1.0.0` is released, breaking changes will happen in minor versions (and possibly also patch versions).

## 📖 Examples

### Build

#### Build a Docker image for a .NET project

```bash
3lv build --project-file src/MyProject.csproj --system-name core my-cool-application
# or use shorthand
3lv build -f src/MyProject.csproj -s core my-cool-application
```

#### Build a Docker image for a .NET project and push it to Elvias registry

```bash
3lv build --project-file src/MyProject.csproj --system-name core --push my-cool-application
# or use shorthand
3lv build -f src/MyProject.csproj -s core -p my-cool-application
```

#### Build a Docker image for a Go project and push it to GitHub Container Registry

```bash
3lv build --project-file src/MyProject.csproj --system-name core --push --registry ghcr my-cool-application
# or use shorthand
3lv build -f src/MyProject.csproj -s core -p -r ghcr my-cool-application
```

### Scan

#### Scan a Docker image for vulnerabilities

```bash
3lv scan my-cool-image
```

#### Scan a Docker image for critical vulnerabilities only

```bash
3lv scan --severity CRITICAL my-cool-image
# or use shorthand
3lv scan -S CRITICAL my-cool-image
```

#### Scan a Docker image for vulnerabilities and output the results to JSON and Markdown

```bash
3lv scan --formats json,markdown my-cool-image
# or use shorthand
3lv scan -F json,markdown my-cool-image
```

### Generate GitHub Actions workflow for Kubernetes deploy

```bash
3lv github-actions --system-name core --application-name my-cool-application --runtime-cloud-provider aks --helm-values-file CI/values.yml
# or use shorthand
3lv gha -s core -a my-cool-application -r aks -f CI/values.yml
```

Remember to also add your repository to [github-repositories-terraform](https://github.com/3lvia/github-repositories-terraform)
to enable access from GitHub Actions to Kubernetes.

## 🧑‍💻 Development

### Installation from source

Requires [Go](https://golang.org) and [Make](https://www.gnu.org/software/make).

#### Windows

Install [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) and follow the Linux instructions.

Optionally, you can build a Windows binary using the following command:

```bash
make build-windows-amd64
```

You can then move the binary to a directory in your PATH.

#### Linux and macOS

Ensure you have Go and Make installed:

```bash
sudo apt install golang make
```

Clone the repository and install the CLI:

```bash
git clone git@github.com:3lvia/cli
cd cli
sudo make install
```

**macOS**: If `GOOS` and `GOARCH` are not properly set, you can use this command:

```bash
# for Intel macs
sudo make install-macos-amd64
# for M1 and newer macs
sudo make install-macos-arm64
```

### Running tests

Unit tests are written in Go and can be run with the following command:

```bash
make test
```

### Linter

We use the linter [golangci-lint](https://golangci-lint.run) and can be run with the following command:

```bash
make lint
```

### Releasing a new version

Bump the number in the `VERSION` file and make a pull request.
When merged, the new version will be released automatically by GitHub Actions.
