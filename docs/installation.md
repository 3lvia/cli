## 💾 Installation

Supported platforms:

- **Windows**
- **macOS** (Intel and M-series)
- **Linux** (any distribution)

### Windows

Download the MSI file (and optionally the MD5 checksum) from the [releases page](https://github.com/3lvia/cli/releases) and run it.

#### Verify checksum (OPTIONAL)

```pwsh
certutil -hashfile 3lv-windows-amd64.msi MD5
```

### Linux/WSL and macOS

Download the tarball file for your platform (and optionally the MD5 cheksum) from the [releases page](https://github.com/3lvia/cli/releases),
extract it and move the binary to a directory in your PATH.

#### Example installation

```bash
# OPTIONAL: verify checksum first
md5sum -c 3lv-linux-amd64.tar.gz.md5

tar -xzf 3lv-linux-amd64.tar.gz
sudo install -Dm755 -t /usr/local/bin 3lv
```

For macOS, you can use the same commands as above, but replace `linux` with `macos`.
**If you have an M1 or newer mac, you can use the `macos-arm64` binary.**
