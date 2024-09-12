package scan

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"
)

//go:embed *.tmpl
var trivyMarkdownTemplate embed.FS

const templateFile = "trivyMarkdown.tmpl"

type TrivyResult struct {
	SchemaVersion int    `json:"SchemaVersion"`
	CreatedAt     string `json:"CreatedAt"`
	ArtifactName  string `json:"ArtifactName"`
	ArtifactType  string `json:"ArtifactType"`
	Metadata      struct {
		OS struct {
			Family string `json:"Family"`
			Name   string `json:"Name"`
			EOSL   bool   `json:"EOSL"`
		} `json:"OS"`
		ImageID     string   `json:"ImageID"`
		DiffIDs     []string `json:"DiffIDs"`
		RepoTags    []string `json:"RepoTags"`
		RepoDigests []string `json:"RepoDigests"`
		ImageConfig struct {
			Architecture  string `json:"architecture"`
			Container     string `json:"container"`
			Created       string `json:"created"`
			DockerVersion string `json:"docker_version"`
			History       []struct {
				Created    string `json:"created"`
				CreatedBy  string `json:"created_by"`
				EmptyLayer bool   `json:"empty_layer"`
			} `json:"history"`
			OS     string `json:"os"`
			Rootfs struct {
				Type    string   `json:"type"`
				DiffIDs []string `json:"diff_ids"`
			} `json:"rootfs"`
			Config struct {
				Cmd   []string `json:"Cmd"`
				Env   []string `json:"Env"`
				Image string   `json:"Image"`
			} `json:"config"`
		} `json:"ImageConfig"`
	} `json:"Metadata"`
	Results []struct {
		Target          string `json:"Target"`
		Class           string `json:"Class"`
		Type            string `json:"Type"`
		Vulnerabilities []struct {
			VulnerabilityID string `json:"VulnerabilityID"`
			PkgID           string `json:"PkgID"`
			PkgName         string `json:"PkgName"`
			PkgIdentifier   struct {
				PURL string `json:"PURL"`
				UID  string `json:"UID"`
			} `json:"PkgIdentifier"`
			InstalledVersion string `json:"InstalledVersion"`
			Status           string `json:"Status"`
			Layer            struct {
				Digest string `json:"Digest"`
				DiffID string `json:"DiffID"`
			} `json:"Layer"`
			SeveritySource string `json:"SeveritySource"`
			PrimaryURL     string `json:"PrimaryURL"`
			DataSource     struct {
				ID   string `json:"ID"`
				Name string `json:"Name"`
				URL  string `json:"URL"`
			} `json:"DataSource"`
			Title          string   `json:"Title"`
			Description    string   `json:"Description"`
			Severity       string   `json:"Severity"`
			CweIDs         []string `json:"CweIDs"`
			VendorSeverity struct {
				Azure      int `json:"azure"`
				Nvd        int `json:"nvd"`
				OracleOval int `json:"oracle-oval"`
				Photon     int `json:"photon"`
				Redhat     int `json:"redhat"`
				Ubuntu     int `json:"ubuntu"`
			} `json:"VendorSeverity"`
			CVSS struct {
				Nvd struct {
					V2Vector string  `json:"V2Vector"`
					V3Vector string  `json:"V3Vector"`
					V2Score  float64 `json:"V2Score"`
					V3Score  float64 `json:"V3Score"`
				} `json:"nvd"`
				Redhat struct {
					V3Vector string  `json:"V3Vector"`
					V3Score  float64 `json:"V3Score"`
				} `json:"redhat"`
			} `json:"CVSS"`
			References       []string `json:"References"`
			PublishedDate    string   `json:"PublishedDate"`
			LastModifiedDate string   `json:"LastModifiedDate"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

func parseJSONOutput() (TrivyResult, error) {
	jsonFile, err := os.Open("trivy.json")
	if err != nil {
		return TrivyResult{}, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return TrivyResult{}, err
	}

	var trivyResult TrivyResult
	err = json.Unmarshal(byteValue, &trivyResult)
	if err != nil {
		return TrivyResult{}, err
	}

	return trivyResult, nil
}

func toMarkdown(result TrivyResult) ([]byte, error) {
	markdownTemplate, err := template.New(templateFile).ParseFS(trivyMarkdownTemplate, templateFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse markdown template: %v", err)
	}

	var markdownBuffer bytes.Buffer
	err = markdownTemplate.Execute(&markdownBuffer, result)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute markdown template: %v", err)
	}

	return markdownBuffer.Bytes(), nil
}
