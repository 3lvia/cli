package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/3lvia/cli/pkg/style"
	"github.com/samber/lo"
)

type GrafanaAnnotation struct {
	What string   `json:"what"` // required
	Data string   `json:"data"` // required
	Tags []string `json:"tags"` // required
}

type FormatDeploymentMessageOptions struct {
	RunID string
}

func formatDeploymentMessage(
	repositoryName string,
	commitMessage string,
	options *FormatDeploymentMessageOptions,
) string {
	if options == nil {
		options = &FormatDeploymentMessageOptions{} //nolint:exhaustruct_v5
	}

	deployedFrom := func() string {
		if options.RunID == "" {
			return "Manually deployed with CLI"
		}

		return "Deployed from GitHub Actions run " + options.RunID
	}()

	deployLink := func() string {
		const GitHubOwner = "3lvia"

		if options.RunID == "" {
			return ""
		}

		return fmt.Sprintf(
			"<a href=\"https://github.com/%s/%s/actions/runs/%s\">Link</a>",
			GitHubOwner,
			repositoryName,
			options.RunID,
		)
	}()

	if deployLink == "" {
		return strings.Join(
			[]string{
				deployedFrom,
				commitMessage,
				repositoryName,
			},
			" - ",
		)
	}

	return strings.Join(
		[]string{
			deployedFrom,
			commitMessage,
			repositoryName,
			deployLink,
		},
		" - ",
	)
}

type PostGrafanaAnnotationOptions struct {
	RunID string
}

func resolveEnvironment(
	environment string,
	runtimeCloudProvider string,
) string {
	if runtimeCloudProvider == "gke" {
		return environment + "_gke"
	}

	if runtimeCloudProvider == "iss" {
		return environment + "_iss"
	}

	return environment
}

func addGrafanaDeploymentAnnotation(
	ctx context.Context,
	wasSuccessful bool,
	applicationName string,
	systemName string,
	environment string,
	runtimeCloudProvider string,
	repositoryName string,
	commitMessage string,
	grafanaURL string,
	grafanaSecret string,
	options *PostGrafanaAnnotationOptions,
) error {
	what := func() string {
		if wasSuccessful {
			return "Deploy successful."
		}

		return "Deploy failed."
	}()

	grafanaAnnotation := GrafanaAnnotation{
		What: what,
		Data: formatDeploymentMessage(
			repositoryName,
			commitMessage,
			&FormatDeploymentMessageOptions{
				RunID: options.RunID,
			},
		),
		Tags: []string{
			"app:" + applicationName,
			"system:" + systemName,
			"env:" + resolveEnvironment(environment, runtimeCloudProvider),
			"event:deploy",
		},
	}

	style.PrintInfo(
		fmt.Sprintf("Sending deploy annotation to Grafana: %v\n", grafanaAnnotation),
	)

	body, err := json.Marshal(grafanaAnnotation)
	if err != nil {
		return err
	}

	// TODO: actually find out why Grafana is returning 429 instead of just retrying
	const (
		RetryAttempts = 5
		RetryDelay    = 5 * time.Second
	)

	_, _, err = lo.AttemptWithDelay(
		RetryAttempts,
		RetryDelay,
		func(i int, _ time.Duration) error {
			style.PrintInfo(
				fmt.Sprintf("Sending deploy annotation to Grafana, attempt %d\n\n", i),
			)

			statusCode, err := sendRequest(
				ctx,
				grafanaURL+"annotations/graphite",
				grafanaSecret,
				body,
			)
			if err != nil {
				return err
			}

			if statusCode != 200 {
				return fmt.Errorf("Grafana returned status code %d", statusCode)
			}

			return nil
		},
	)
	if err != nil {
		style.PrintError(
			fmt.Sprintf("Failed to send deploy annotation to Grafana after %d attempts\n", RetryAttempts),
		)

		return err
	}

	style.PrintSuccess("Deploy annotation sent to Grafana!\n")

	return nil
}

func sendRequest(
	ctx context.Context,
	url string,
	secret string,
	body []byte,
) (int, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
