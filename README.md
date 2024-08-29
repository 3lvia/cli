# cli

Command Line Interface tool for developing, building and deploying Elvia applications.

## Installation

Requires [Go](https://golang.org).

### Install from GitHub

```bash
go install github.com/3lvia/cli@latest
```

### Install from source

#### Script

```bash
./scripts/install.sh
```

#### Go

This will cause the CLI to be installed with the name `cli` (thanks https://github.com/golang/go/issues/44469).

```bash
git clone git@github.com:3lvia/cli.git
cd cli
go install .
```

## Usage

```bash
3lv --help
```
