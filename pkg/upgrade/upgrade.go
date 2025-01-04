package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
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

const (
	commandName            = "upgrade"
	defaultInstallLocation = "/usr/local/bin"
)

func Command(version string) *cli.Command {
	return &cli.Command{
		Name:      commandName,
		Aliases:   []string{"u"},
		Usage:     "Upgrade 3lv to the latest version.",
		UsageText: "3lv upgrade [options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "install-location",
				Usage:   "Where to install the new 3lv binary. If set, the use-existing-binary-location flag will be ignored.",
				Aliases: []string{"l"},
				Value:   defaultInstallLocation,
			},
			&cli.BoolFlag{
				Name:    "use-existing-binary-location",
				Usage:   "Install the new 3lv binary to where the existing 3lv binary is located.",
				Aliases: []string{"e"},
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "force-reinstall",
				Usage:   "Force the reinstallation of the 3lv binary, even if the version is the same.",
				Aliases: []string{"f"},
				Value:   false,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			return Upgrade(ctx, c, version)
		},
	}
}

func Upgrade(ctx context.Context, c *cli.Command, version string) error {
	installLocation := func() string {
		if c.IsSet("install-location") {
			return c.String("install-location")
		}

		if c.Bool("use-existing-binary-location") {
			currentBinary, err := os.Executable()
			if err != nil {
				style.PrintError(
					"Could not determine the current binary location, will install to " + defaultInstallLocation,
				)

				return defaultInstallLocation
			}

			return path.Dir(currentBinary)
		}

		return c.String("install-location")
	}()

	latestVersion, err := GetLatestCLIVersion(ctx)
	if err != nil {
		return cli.Exit(err, 1)
	}

	if semver.Compare("v"+version, "v"+latestVersion) != -1 {
		style.PrintSuccess(
			"You are already using the latest version of 3lv: " + version,
		)

		if !c.Bool("force-reinstall") {
			style.PrintInfo(
				"Use the --force-reinstall flag to force the reinstallation of the 3lv binary.",
			)

			return cli.Exit("", 0)
		}
	}

	binaryURL, err := getLatestBinaryURL(ctx)
	if err != nil {
		return cli.Exit(err, 1)
	}

	style.PrintInfo(
		fmt.Sprintf("Upgrading 3lv from %s to %s...", version, latestVersion),
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

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		binaryURL,
		nil,
	)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
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

	if installCommandOutput := installCommand(
		nil,
		path.Join(tempDir, "3lv"),
		installLocation,
	); installCommandOutput.Error != nil {
		return cli.Exit(installCommandOutput.Error, 1)
	}

	style.PrintSuccess(
		fmt.Sprintf("Successfully upgraded 3lv to %s!", latestVersion),
	)

	return nil
}

func getLatestBinaryURL(ctx context.Context) (string, error) {
	if runtime.GOOS == "windows" {
		return "",
			errors.New(
				"Auto-upgrade is not supported on Windows. Please download the MSI installer from the GitHub releases page",
			)
	}

	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return "", fmt.Errorf("Auto-upgrade is not supported on %s", runtime.GOOS)
	}

	if runtime.GOARCH != "amd64" && (runtime.GOOS != "darwin" || runtime.GOARCH != "arm64") {
		return "", fmt.Errorf("Auto-upgrade is not supported on %s", runtime.GOARCH)
	}

	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, "3lvia", "cli")
	if err != nil {
		return "", err
	}

	for _, asset := range release.Assets {
		goos := func() string {
			if runtime.GOOS == "darwin" {
				return "macos"
			}

			return runtime.GOOS
		}()

		if strings.HasPrefix(*asset.Name, "3lv-") &&
			strings.HasSuffix(*asset.Name, fmt.Sprintf("-%s-%s.tar.gz", goos, runtime.GOARCH)) {
			return *asset.BrowserDownloadURL, nil
		}
	}

	return "", errors.New("No binary found for the latest release")
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
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)

	for {
		header, err := tarReader.Next()
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
			err := os.MkdirAll(target, 0o755)
			if err != nil {
				return err
			}
		} else if header.Typeflag == tar.TypeReg {
			file, err := os.Create(target)
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, tarReader); err != nil {
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
	installLocation string,
) command.Output {
	return command.Run(
		*exec.Command(
			"sudo",
			"install",
			"-Dm755",
			"-t",
			installLocation,
			binaryFile,
		),
		options,
	)
}
