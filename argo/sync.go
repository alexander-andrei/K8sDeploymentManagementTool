package argo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"k8s/tool/config"
	"k8s/tool/gitops"
	"k8s/tool/utils"
)

type Application struct {
	Status Status `json:"status"`
}

type Status struct {
	Summary Summary `json:"summary"`
}

type Summary struct {
	Images []string `json:"images"`
}

func CheckAndRevertTags(latestTag string, previousTag string, errorRate float64, applicationName string) {
	image := getArgoCdDeploymentInfo(applicationName)

	if utils.ExtractTagVersion(image) != latestTag {
		os.Exit(3)
	}

	gitImage := gitops.GetServerConfigurationsFromGit(applicationName)

	if image != gitImage {
		os.Exit(3)
	}

	if errorRate > config.GlobalConfig.Kibana.ErrorRate {
		gitops.ChangeAndCommit(applicationName, previousTag, gitImage)
	}
}

func TriggerDeploymentSync(applicationName string) error {
	apiURL := fmt.Sprintf("%s/api/v1/applications/%s/sync", config.GlobalConfig.ArgoCD.Location, applicationName)

	payload := map[string]interface{}{
		"revision": "HEAD",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: customTransport,
	}

	req.Header.Set("Authorization", "Bearer "+config.GlobalConfig.ArgoCD.Token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to trigger deployment sync: %s", resp.Status)
	}

	fmt.Println("Deployment sync triggered successfully.")
	return nil
}

func getArgoCdDeploymentInfo(applicationName string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/applications/%s", config.GlobalConfig.ArgoCD.Location, applicationName), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+config.GlobalConfig.ArgoCD.Token)
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(err)
	}

	var application Application

	if err := json.NewDecoder(resp.Body).Decode(&application); err != nil {
		panic(err)
	}

	return application.Status.Summary.Images[0]
}
