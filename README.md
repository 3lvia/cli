# cli

Command Line Interface tool for developing, building and deploying Elvia applications ⚡

## 💾 Installation

See the [releases page](https://github.com/3lvia/cli/releases) and download your platform's binary.

Supported platforms:

- **Linux**
- **macOS** (Intel and M-series)
- **Windows**

## 📋 Requirements

To use the 3lv CLI, you need to have these dependencies installed:

- [Docker](https://docs.docker.com/engine/install) (used for building)
- [Helm](https://helm.sh/docs/intro/install) (used for deploying)
- [Trivy](https://aquasecurity.github.io/trivy) (used for scanning)
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) (used for pushing to Azure Container Registry)
- [GitHub CLI](https://cli.github.com) (used for pushing to GitHub Container Registry)

### Pusing to registries

If you want to push to a registry, you need to be authenticated to that registry.

#### Azure Container Registry

This is Elvia's default registry.
Use the following command to login to Elvias registry:

```bash
az acr login -n containerregistryelvia
```

#### GitHub Container Registry

Use the following command (with your GitHub username) to login:

```bash
gh auth token | docker login ghcr.io --username your-github-username --password-stdin
```

## ❓ Usage

```bash
3lv --help
```

## ℹ️ Examples

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

## 🧑‍💻 Development

### Installation from source

#### Linux and macOS

Requires [Go](https://golang.org) and [Make](https://www.gnu.org/software/make).
These can be installed on Debian/Ubuntu/WSL with the following command:

```bash
sudo apt install golang make
```

Clone the repository and install the CLI:

```bash
git clone git@github.com:3lvia/cli
cd cli
sudo make install
```

**macOS**: If `GOOS` and `GOARCH` is not properly set, you can use this command:

```bash
# for Intel macs
sudo make install-macos-amd64
# for M1 and newer macs
sudo make install-macos-arm64
```

#### Windows

Install [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) and follow the Linux instructions.

Optionally, you can build a Windows binary using the following command:

```bash
sudo make build-windows-amd64
```

You can then move the binary to a directory in your PATH.

### Releasing a new version

Bump the number in the `VERSION` file and make a pull request.
When merged, the new version will be released automatically by GitHub Actions.
