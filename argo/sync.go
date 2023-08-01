package argo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"k8s/tool/config"
)

func CheckAndRevertTags(latestTag string, previousTag string, errorRate float64) {
	// check what tag is curently deployed
	// check what tag is defined in git
	// if error rate is bigger than it was supposed to be, change tag in git
	// if all is good exit command
}

func TriggerDeploymentSync(applicationName string) error {
	apiURL := fmt.Sprintf("%s/api/v1/applications/%s/sync", config.GlobalConfig.ArgoCD.Location, applicationName)

	// Create a JSON payload for the request body
	payload := map[string]interface{}{
		"revision": "HEAD",
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	// Set the content type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Create a custom HTTP transport with InsecureSkipVerify set to true
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Create a custom HTTP client with the custom transport
	client := &http.Client{
		Transport: customTransport,
	}

	req.Header.Set("Authorization", "Bearer "+config.GlobalConfig.ArgoCD.Token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to trigger deployment sync: %s", resp.Status)
	}

	fmt.Println("Deployment sync triggered successfully.")
	return nil
}
