package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
		options = &FormatDeploymentMessageOptions{}
	}

	deployedFrom := func() string {
		if options.RunID == "" {
			return "Manually deployed with CLI"
		}
		return fmt.Sprintf("Deployed from GitHub Actions run %s", options.RunID)
	}()

	deployLink := func() string {
		const GITHUB_OWNER = "3lvia"

		if options.RunID == "" {
			return ""
		}
		return fmt.Sprintf(
			"<a href=\"https://github.com/%s/%s/actions/runs/%s\">Link</a>",
			GITHUB_OWNER,
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

func addGrafanaDeploymentAnnotation(
	wasSuccessful bool,
	applicationName string,
	systemName string,
	environment string,
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
			"env:" + environment,
			"event:deploy",
		},
	}

	log.Printf("Sending deploy annotation to Grafana: %v\n", grafanaAnnotation)
	body, err := json.Marshal(grafanaAnnotation)
	if err != nil {
		return err
	}

	// TODO: actually find out why Grafana is returning 429 instead of just retrying
	const RETRY_ATTEMPTS = 5
	const RETRY_DELAY = 5 * time.Second

	_, _, err = lo.AttemptWithDelay(
		RETRY_ATTEMPTS,
		RETRY_DELAY,
		func(i int, duration time.Duration) error {
			log.Printf("Sending deploy annotation to Grafana, attempt %d\n", i)

			statusCode, err := sendRequest(
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
		log.Printf("Failed to send deploy annotation to Grafana after %d attempts\n", RETRY_ATTEMPTS)
		return err
	}

	log.Println("Deploy annotation sent to Grafana!")

	return nil
}

func sendRequest(
	url string,
	secret string,
	body []byte,
) (int, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+secret)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
