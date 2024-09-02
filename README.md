# cli

Command Line Interface tool for developing, building and deploying Elvia applications.

# Installation

See the [releases page](https://github.com/3lvia/cli/releases) and download your platform's binary.

Supported platforms:

- **Linux**
- **macOS** (Intel and M-series)
- **Windows**

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

Optionally, you can build a Windows binary using the following command:

```bash
make build-windows-amd64
```

You can then move the binary to a directory in your PATH.

# Usage

```bash
3lv --help
```
