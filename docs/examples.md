## 📖 Examples

The CLI assumes that you are in the git repository of the project you are working on.

### Init

<details>

<summary>Expand</summary>

#### Create a new 3lv configuration file

```bash
# you will be prompted for input
3lv init
```

This will create a `3lv.yml` file in the root of your repository.
The advantage of having a configuration file is that you can avoid typing the same options every time you run a command.

### Build

<details>

<summary>Expand</summary>

#### Build a container for a .NET project

```bash
3lv build --project-file src/MyProject.csproj --system-name core my-cool-application
# or use shorthand
3lv build -f src/MyProject.csproj -s core my-cool-application
# or with 3lv.yml
3lv build my-cool-application
```

#### Build a container for a .NET project and push it to Elvia's registry

```bash
3lv build --project-file src/MyProject.csproj --system-name core --push my-cool-application
# or use shorthand
3lv build -f src/MyProject.csproj -s core -p my-cool-application
# or with 3lv.yml
3lv build -p my-cool-application
```

#### Build a container for a Go project and push it to GitHub Container Registry

```bash
3lv build --project-file src/MyProject.csproj --system-name core --push --registry ghcr my-cool-application
# or use shorthand
3lv build -f src/MyProject.csproj -s core -p -r ghcr my-cool-application
# or with 3lv.yml
3lv build -p -r ghcr my-cool-application
```

#### Generate a Dockerfile for a .NET project

```bash
3lv build --project-file src/MyProject.csproj --system-name core --generate-only my-cool-application
# or use shorthand
3lv build -f src/MyProject.csproj -s core -G my-cool-application
# or with 3lv.yml
3lv build -G my-cool-application
```

</details>

### Scan

<details>

<summary>Expand</summary>

#### Scan a container for vulnerabilities

```bash
3lv scan my-cool-image
```

#### Scan a container for critical vulnerabilities only

```bash
3lv scan --severity CRITICAL my-cool-image
# or use shorthand
3lv scan -S CRITICAL my-cool-image
```

#### Scan a container for vulnerabilities and output the results to JSON and Markdown

```bash
3lv scan --formats json,markdown my-cool-image
# or use shorthand
3lv scan -F json,markdown my-cool-image
```

</details>

### GitHub Actions

<details>

<summary>Expand</summary>

### Generate a GitHub Actions workflow for deploying to Kubernetes

```bash
3lv github-actions --system-name core --application-name my-cool-application --project-file src/MyProject.csproj
# or use shorthand
3lv gha -s core -a my-cool-application -f src/MyProject.csproj
```

Remember to also add your repository to [github-repositories-terraform](https://github.com/3lvia/github-repositories-terraform)
to enable access from GitHub Actions to Kubernetes.

### Generate a GitHub Actions workflow for deploying to ISS, using an existing Helm values file

```bash
3lv github-actions --system-name core --application-name my-cool-application --runtime-cloud-provider iss --helm-values-file CI/values.yml --project-file src/MyProject.csproj
# or use shorthand
3lv gha -s core -a my-cool-application -r iss -F CI/values.yml -f src/MyProject.csproj
```

Remember to also add your repository to the list of repositories in the `github-actions-deploy`-module
in [iss-terraform](https://github.com/3lvia/iss-terraform) to enable access from GitHub Actions to ISS.

</details>

### Create

<details>

<summary>Expand</summary>

#### Create a new .NET 8 API in the current directory

```bash
3lv create --system-name core --application-name my-cool-application --template dotnet8-webapi .
# or use shorthand
3lv create -s core -a my-cool-application -t dotnet8-webapi .
```

#### Create a new Python API in the current directory

```bash
3lv create --system-name core --application-name my-cool-application --template python-webapi .
# or use shorthand
3lv create -s core -a my-cool-application -t python-webapi .
```

#### Create a new Go API in the current directory

```bash
3lv create --system-name core --application-name my-cool-application --template go-webapi .
# or use shorthand
3lv create -s core -a my-cool-application -t go-webapi .
```

#### Create a new .NET 8 API in the `applications` directory of a monorepo, putting the GitHub Actions workflows in the the root of the repository

```bash
3lv create --system-name core --application-name my-cool-application --template dotnet8-webapi --github-actions-directory . applications
# or use shorthand
3lv create -s core -a my-cool-application -t dotnet8-webapi -G . applications
```

</details>
