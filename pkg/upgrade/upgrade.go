package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/style"
	"github.com/google/go-github/v67/github"
	"github.com/urfave/cli/v3"
	"golang.org/x/mod/semver"
)

const commandName = "upgrade"

func Command(version string) *cli.Command {
	return &cli.Command{
		Name:    commandName,
		Aliases: []string{"u"},
		Usage:   "Upgrade the Elvia CLI",
		Action: func(ctx context.Context, c *cli.Command) error {
			return Upgrade(ctx, c, version)
		},
	}
}

func Upgrade(ctx context.Context, c *cli.Command, version string) error {
	latestVersion, err := GetLatestCLIVersion(ctx)
	if err != nil {
		return cli.Exit(err, 1)
	}

	if semver.Compare("v"+version, "v"+latestVersion) != -1 {
		style.Print(
			fmt.Sprintf("You are already using the latest version of 3lv: %s", version),
			&style.PrintOptions{Color: "green"},
		)
	}

	binaryURL, err := getLatestBinaryURL()
	if err != nil {
		return cli.Exit(err, 1)
	}

	style.Print(
		fmt.Sprintf("Upgrading 3lv from %s to %s...", version, latestVersion),
		&style.PrintOptions{Color: "yellow"},
	)

	tempDir, err := os.MkdirTemp("", "3lv-upgrade")
	if err != nil {
		return cli.Exit(err, 1)
	}
	defer os.RemoveAll(tempDir)

	out, err := os.Create(path.Join(tempDir, "3lv.tar.gz"))
	if err != nil {
		return cli.Exit(err, 1)
	}
	defer out.Close()

	response, err := http.Get(binaryURL)
	if err != nil {
		return cli.Exit(err, 1)
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = decompress(out.Name(), tempDir)
	if err != nil {
		return cli.Exit(err, 1)
	}

	installCommandOutput := installCommand(nil, path.Join(tempDir, "3lv"))
	if command.IsError(installCommandOutput) {
		return cli.Exit(installCommandOutput.Error, 1)
	}

	style.Print(
		fmt.Sprintf("Successfully upgraded 3lv to %s!", latestVersion),
		&style.PrintOptions{Color: "green"},
	)

	return nil
}

func getLatestBinaryURL() (string, error) {
	if runtime.GOOS == "windows" {
		return "", fmt.Errorf("Auto-upgrade is not supported on Windows. Please download the MSI installer from the GitHub releases page.")
	}

	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return "", fmt.Errorf("Auto-upgrade is not supported on %s.", runtime.GOOS)
	}

	if runtime.GOARCH != "amd64" && (runtime.GOOS != "darwin" || runtime.GOARCH != "arm64") {
		return "", fmt.Errorf("Auto-upgrade is not supported on %s.", runtime.GOARCH)
	}

	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "3lvia", "cli")
	if err != nil {
		return "", err
	}

	for _, asset := range release.Assets {
		os_ := func() string {
			if runtime.GOOS == "darwin" {
				return "macos"
			}
			return runtime.GOOS
		}()

		if strings.HasPrefix(*asset.Name, "3lv-") &&
			strings.HasSuffix(*asset.Name, fmt.Sprintf("-%s-%s.tar.gz", os_, runtime.GOARCH)) {
			return *asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("No binary found for the latest release")
}

func GetLatestCLIVersion(ctx context.Context) (string, error) {
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, "3lvia", "cli")
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(*release.TagName, "v"), nil
}

func decompress(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Clean the path and ensure it is within the destination directory
		cleanedPath := filepath.Clean(header.Name)
		target := filepath.Join(dest, cleanedPath)
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(target, 0755)
			if err != nil {
				return err
			}
		} else if header.Typeflag == tar.TypeReg {
			file, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return err
			}
			file.Close()
		}
	}
	return nil
}

func installCommand(
	options *command.RunOptions,
	binaryFile string,
) command.Output {
	return command.Run(
		*exec.Command(
			"sudo",
			"install",
			"-Dm755",
			"-t",
			"/usr/local/bin",
			binaryFile,
		),
		options,
	)
}
