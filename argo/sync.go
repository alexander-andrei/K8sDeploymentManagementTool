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

	log "github.com/rs/zerolog/log"
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

func CheckAndRevertTags(latestTag string, previousTag string, errorRate float64, applicationName string) error {
	image, err := getArgoCdDeploymentInfo(applicationName)

	if err != nil {
		return err
	}

	imageTag := utils.ExtractTagVersion(image)

	if imageTag != latestTag {
		log.Warn().Str("ArgoCdTag", imageTag).Str("LatestTag", latestTag).Msg("Image tag from ArgoCD does not equal latest tag")
		fmt.Print("Image tag from ArgoCD does not equal latest tag")
		os.Exit(3)
	}

	gitImage := gitops.GetServerConfigurationsFromGit(applicationName)

	if image != gitImage {
		log.Warn().Str("ArgoCdImage", image).Str("GitOpsImage", gitImage).Msg("Image from GitOps repo does not equal image from ArgoCd")
		fmt.Print("Image from GitOps repo does not equal image from ArgoCd")
		os.Exit(3)
	}

	if errorRate > config.GlobalConfig.Kibana.ErrorRate {
		gitops.ChangeAndCommit(applicationName, previousTag, gitImage)
	}

	return nil
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
		return fmt.Errorf("Failed to trigger deployment sync: %s", resp.Status)
	}

	return nil
}

func getArgoCdDeploymentInfo(applicationName string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/applications/%s", config.GlobalConfig.ArgoCD.Location, applicationName), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+config.GlobalConfig.ArgoCD.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status code is not 200, status code: %d", resp.StatusCode)
	}

	var application Application

	if err := json.NewDecoder(resp.Body).Decode(&application); err != nil {
		return "", err
	}

	return application.Status.Summary.Images[0], nil
}
