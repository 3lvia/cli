package utils

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/3lvia/cli/pkg/style"
	"github.com/google/go-github/v89/github"
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
				"Failed to resolve commit hash: %w."+
					" Please verify you are currently in a Git repository, or manually specify the commit hash with --commit-hash",
				err,
			)
	}

	return strings.TrimSpace(string(hash)), nil
}

func ResolveGitRepositoryTopLevelPath() (string, error) {
	gitTopLevel, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(gitTopLevel)), nil
}

func ResolveRepositoryName(possibleRepositoryName string) (string, error) {
	if possibleRepositoryName != "" {
		return possibleRepositoryName, nil
	}

	gitTopLevel, err := ResolveGitRepositoryTopLevelPath()
	if err != nil {
		return "",
			fmt.Errorf(
				"Failed to resolve repository name: %w."+
					" Please verify you are currently in a Git repository, or manually specify the repository with --repository-name",
				err,
			)
	}

	return path.Base(gitTopLevel), nil
}

func ResolveCommitMessage(possibleCommitMessage string) (string, error) {
	if possibleCommitMessage != "" {
		return possibleCommitMessage, nil
	}

	message, err := exec.Command("git", "log", "-1", "--no-merges", "--pretty=%B").Output()
	if err != nil {
		return "",
			fmt.Errorf(
				"Failed to resolve commit message: %w."+
					" Please verify you are currently in a Git repository,"+
					" or manually specify the commit message with --commit-message",
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

func WriteFileWithTemplate(
	directory string,
	fileName string,
	templateFile string,
	templates embed.FS,
	variables any,
) (string, error) {
	filePath := path.Join(directory, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to create file: %w", err)
	}

	defer file.Close()

	tmpl, err := template.New(templateFile).ParseFS(templates, templateFile)
	if err != nil {
		return "", fmt.Errorf("Failed to parse template: %w", err)
	}

	var fileBuffer bytes.Buffer

	err = tmpl.ExecuteTemplate(&fileBuffer, path.Base(templateFile), variables)
	if err != nil {
		return "", fmt.Errorf("Failed to execute template: %w", err)
	}

	if _, err := file.Write(fileBuffer.Bytes()); err != nil {
		return "", fmt.Errorf("Failed to write file: %w", err)
	}

	return filePath, nil
}

func PromptYesNo(question string, nonInteractive bool) (bool, error) {
	style.PrintInfo(
		question + " (y/n): ",
	)

	if nonInteractive {
		return true, nil
	}

	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, fmt.Errorf("Failed to read response: %w", err)
	}

	return strings.ToLower(response) == "y", nil
}

func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

func GetAuthenticatedGithubClient() *github.Client {
	if os.Getenv("CI") == "1" {
		return github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	}

	return github.NewClient(nil)
}
