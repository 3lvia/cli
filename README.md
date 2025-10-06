# cli

**3lv** – the Command Line Interface tool for creating, building and securing Elvia applications ⚡

## 📦 Installation

### TL;DR

```bash
curl -fsSL https://raw.githubusercontent.com/3lvia/cli/trunk/install.sh | bash
```

### See [Installation](docs/installation.md) for more options.

## 🚀 Features

- **Build** a container for .NET, Go or Python projects without needing a Dockerfile.
- **Scan** a container for vulnerabilities using Trivy.
- **Deploy** to Azure Kubernetes Service, Google Kubernetes Engine and ISS.
- **Create** new projects from Elvia templates, with all batteries included.
- **Generate** a GitHub Actions workflow for building and deploying to Elvia's clusters on Azure, Google Cloud and ISS.

The GitHub Actions composite actions at [core-github-actions-templates](https://github.com/3lvia/core-github-actions-templates) are wrappers around many of the `3lv` commands.
Therefore it's useful to use `3lv` when debugging or testing Elvia's actions, since you can very easily reproduce the same commands locally.

## 📚 Documentation

- **[💾 Installation](docs/installation.md)**
- **[📋 Requirements](docs/requirements.md)**
- **[⚙️ Configuration](docs/configuration.md)**
- **[📖 Examples](docs/examples.md)**
- **[🧑‍💻 Development](docs/development.md)**

## 💥 Breaking changes

Before version `v1.0.0` is released, breaking changes will happen in minor versions (and possibly also patch versions).
