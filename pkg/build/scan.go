package build

import (
	"fmt"
	"os"
	"os/exec"
)

type ScanImageOptions struct {
	ImageName string // required
	Severity  string // required
}

func scanImage(options ScanImageOptions) error {
	scanCmd := exec.Command(
		"trivy",
		"image",
		"--severity",
		options.Severity,
		options.ImageName,
	)
	scanCmd.Stdout = os.Stdout
	scanCmd.Stderr = os.Stderr

	if err := scanCmd.Run(); err != nil {
		return fmt.Errorf("Failed to scan Docker image: %w", err)
	}

	return nil
}
