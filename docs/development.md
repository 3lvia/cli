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
