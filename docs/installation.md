## 💾 Installation

Supported platforms:

- **Linux/WSL** (recommended)
- **macOS** (both Intel and M-series)
- **Windows**

We stronlgy recommend using [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/install) if you are on Windows.
This is because we depend on software tools (see [requirements](requirements.md)) that are much easier to install on a Linux distribution than on Windows.

Note however that the CLI works perfectly fine on Windows, as long as you manage to install the aforementioned required software yourself.

### Linux/WSL and macOS

Download the tarball file for your platform (and optionally the MD5 cheksum) from the [releases page](https://github.com/3lvia/cli/releases),
extract it and move the binary to a directory in your PATH.

If you are using WSL, you can install the Linux binary.

#### Example installation

```bash
# OPTIONAL: verify checksum first
md5sum -c 3lv-linux-amd64.tar.gz.md5

tar -xzf 3lv-linux-amd64.tar.gz
sudo install -Dm755 -t /usr/local/bin 3lv
```

For macOS, you can use the same commands as above, but replace `linux` with `macos`.
**If you have an M1 or newer mac, you can use the `macos-arm64` binary.**

### Windows

Download the MSI file (and optionally the MD5 checksum) from the [releases page](https://github.com/3lvia/cli/releases) and run it.

#### Verify checksum (OPTIONAL)

```pwsh
certutil -hashfile 3lv-windows-amd64.msi MD5
```

### Upgrades

To upgrade the CLI, simply run this command:

```bash
3lv upgrade
```

This will download the latest version and replace the existing binary.

You can also download the new binary manually and replace the old one using the same steps as in the installation section.

#### Upgrade on Windows

If you are on Windows (not WSL), the `3lv upgrade` command will not work. You have to download a new MSI installer manually and run it.
