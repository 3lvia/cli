## 📖 Examples

The CLI assumes that you are in the git repository of the project you are working on.

### Build

<details>

<summary>Expand</summary>

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

#### Generate a Dockerfile for a .NET project

```bash
3lv build --project-file src/MyProject.csproj --system-name core --generate-only my-cool-application
# or use shorthand
3lv build -f src/MyProject.csproj -s core -G my-cool-application
```

</details>

### Scan

<details>

<summary>Expand</summary>

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

</details>

### GitHub Actions

<details>

<summary>Expand</summary>

### Generate GitHub Actions workflow for Kubernetes deploy

```bash
3lv github-actions --system-name core --application-name my-cool-application --runtime-cloud-provider aks --helm-values-file CI/values.yml
# or use shorthand
3lv gha -s core -a my-cool-application -r aks -f CI/values.yml
```

Remember to also add your repository to [github-repositories-terraform](https://github.com/3lvia/github-repositories-terraform)
to enable access from GitHub Actions to Kubernetes.

</details>

### Create

<details>

<summary>Expand</summary>

#### Create a new Elvia application in the current directory

```bash
3lv create --system-name core --application-name my-cool-application .
# or use shorthand
3lv create -s core -a my-cool-application .
```

#### Create a new Elvia application in the applications directory of a monorepo, putting the GitHub Actions workflows in the the root of the repository

```bash
3lv create --system-name core --application-name my-cool-application --github-actions-directory . applications
# or use shorthand
3lv create -s core -a my-cool-application -G . applications
```

</details>
