# cli

Command Line Interface tool for developing, building and securing Elvia applications ⚡

## 🚀 Features

- **Build** Docker images for .NET and Go projects without needing a Dockerfile.
- **Scan** Docker images for vulnerabilities using Trivy.
- **Deploy** to Azure Kubernetes Service, Google Kubernetes Engine and ISS.
- **Create** new projects with all batteries included.
- **Generate** GitHub Actions workflows for building and deploying to Elvia's clusters on Azure, Google Cloud and ISS.

The GitHub composite actions at [core-github-actions-templates](https://github.com/3lvia/core-github-actions-templates) are wrappers around many of the CLI commands.
Therefore it's useful to use the CLI when debugging or testing Elvias actions, since you can very easily reproduce the same commands locally.

## 📚 Documentation

- **[💾 Installation](docs/installation.md)**
- **[📋 Requirements](docs/requirements.md)**
- **[⚙️ Configuration](docs/configuration.md)**
- **[📖 Examples](docs/examples.md)**
- **[🧑‍💻 Development](docs/development.md)**

## 💥 Breaking changes

Before version `v1.0.0` is released, breaking changes will happen in minor versions (and possibly also patch versions).
