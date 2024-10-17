package utils

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func RemoveZeroValues(slice []string) []string {
	var result []string
	for _, value := range slice {
		if value != "" {
			result = append(result, value)
		}
	}

	return result
}

func ResolveCommitHash(possibleCommitHash string) (string, error) {
	if possibleCommitHash != "" {
		return possibleCommitHash, nil
	}

	hash, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return "",
			fmt.Errorf(
				"Failed to resolve commit hash: %w. Please verify you are currently in a Git repository, or manually specify the commit hash with --commit-hash.",
				err,
			)
	}

	return strings.TrimSpace(string(hash)), nil
}

func ResolveRepositoryName(possibleRepositoryName string) (string, error) {
	if possibleRepositoryName != "" {
		return possibleRepositoryName, nil
	}

	gitTopLevel, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "",
			fmt.Errorf(
				"Failed to resolve repository name: %w. Please verify you are currently in a Git repository, or manually specify the repository with --repository-name.",
				err,
			)
	}

	return path.Base(strings.TrimSpace(string(gitTopLevel))), nil
}

func ResolveCommitMessage(possibleCommitMessage string) (string, error) {
	if possibleCommitMessage != "" {
		return possibleCommitMessage, nil
	}

	message, err := exec.Command("git", "log", "-1", "--no-merges", "--pretty=%B").Output()
	if err != nil {
		return "",
			fmt.Errorf(
				"Failed to resolve commit message: %w. Please verify you are currently in a Git repository, or manually specify the commit message with --commit-message.",
				err,
			)
	}

	return strings.TrimSpace(string(message)), nil
}

func StringWithDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}
