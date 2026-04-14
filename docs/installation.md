## 💾 Installation

### TL;DR

```bash
curl -fsSL https://raw.githubusercontent.com/3lvia/cli/trunk/install.sh | bash
```

Supported platforms:

- **Linux/WSL**
- **macOS** (both Intel and M-series)

**Windows is not supported natively; use WSL**

[Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/install).

### Linux/WSL and macOS

Download the tarball file for your platform (and optionally the MD5 checksum) from the [releases page](https://github.com/3lvia/cli/releases),
extract it and move the binary to a directory in your PATH.

If you are using WSL, you can install the Linux binary.

#### Example Linux installation

```bash
# OPTIONAL: verify checksum first
md5sum -c 3lv-linux-amd64.tar.gz.md5

tar -xzf 3lv-linux-amd64.tar.gz
sudo install -Dm755 -t /usr/local/bin 3lv
```

#### Example Mac installation

The installation is done by first installing the app, secondly trying (and failing) to open it, and thirdly [allowing the opening of it](https://support.apple.com/en-gb/guide/mac-help/mh40616/mac).

```bash
tar -xzf 3lv-macos-arm64.tar.gz
sudo install -Dm755 3lv /usr/local/bin/
```

**If you have an M1 or newer mac, you can use the `macos-arm64` binary.**

### Upgrades

To upgrade `3lv`, simply run this command:

```bash
3lv upgrade
```

This will download the latest version and replace the existing binary.

You can also download the new binary manually and replace the old one using the same steps as in the installation section.
