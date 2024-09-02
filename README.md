# cli

Command Line Interface tool for developing, building and deploying Elvia applications.

## Installation

See the releases page and download your platform's binary.

Supported platforms:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Installation from source

### Linux and macOS

Requires [Go](https://golang.org) and [Make](https://www.gnu.org/software/make).

```bash
# debian/ubuntu/WSL
sudo apt install golang make

git clone git@github.com:3lvia/cli
cd cli
sudo make install
```

**macOS**: If `GOOS` and `GOARCH` is not properly set, you can use this command:

```bash
make install-macos-amd64
# for M1 and newer macs
make install-macos-arm64
```

### Windows

Install [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) and follow the Linux instructions.

## Usage

```bash
3lv --help
```
