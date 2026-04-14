## 🧑‍💻 Development

### Installation from source

Requires [Go](https://golang.org) and [Make](https://www.gnu.org/software/make).
You can then move the binary to a directory in your PATH.

#### Linux and macOS

Ensure you have Go and Make installed:

```bash
sudo apt install golang make
```

Clone the repository and install `3lv`:

```bash
git clone git@github.com:3lvia/cli
cd cli
make install # requires sudo
```

**macOS**: If `GOOS` and `GOARCH` are not properly set, you can use this command:

```bash
# for Intel macs
make install-macos-amd64
# for M1 and newer macs
make install-macos-arm64
```

### Running tests

Unit tests are written in Go and can be run with the following command:

```
make test
```

End-to-end tests are written in Bash. They are located in the `tests` directory, where each command has its own folder with a `run_tests.sh` script.
For example, to run the tests for the `build` command:

```bash
./tests/build/run_tests.sh
```

Note that the end-to-end tests run the `3lv` command directly.
This means that if you have changed some code and want to test it locally, you will need to install from source again (e.g. `make install`).

### Linter

We use the linter [golangci-lint](https://golangci-lint.run), which can be run with the following command:

```bash
make lint
```

Linter configuration can be found in `.golangci.yml`.

### Releasing a new version

Increase the number in the `VERSION` file, adhering to [semver](https://semver.org).
Before `3lv` is stable at `v1.0.0`, breaking changes will happen in minor versions (and possibly also patch versions).
Therefore, we will not increase the major version until `v1.0.0` is released.

When a commit increasing the version is merged or pushed to `trunk`, a GitHub Action will automatically create a new release with the new version number as the tag.
