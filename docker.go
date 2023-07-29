package main

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/heroku/docker-registry-client/registry"
)

func LatestAndPreviousImageTags() (string, string) {
	// Replace these with your actual values
	registryURL := "https://registry.k8s.io"
	repoName := "e2e-test-images/agnhost"
	username := "" // Leave empty if the repository is public
	password := "" // Leave empty if the repository is public

	client, err := registry.New(registryURL, username, password)
	if err != nil {
		panic(err)
	}

	tags, err := client.Tags(repoName)
	if err != nil {
		panic(err)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found in the repository.")
		return "", ""
	}

	var validTags []string
	tagPattern := regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	for _, tag := range tags {
		if tagPattern.MatchString(tag) {
			validTags = append(validTags, tag)
		}
	}

	if len(validTags) == 0 {
		fmt.Println("No valid tags found in the repository.")
		return "", ""
	}

	sort.Strings(validTags)
	latestTag := validTags[len(validTags)-1]
	previousTag := ""
	if len(validTags) > 1 {
		previousTag = validTags[len(validTags)-2]
	}

	fmt.Printf("Latest valid image tag for repository '%s': %s\n", repoName, latestTag)
	if previousTag != "" {
		fmt.Printf("Previous valid image tag for repository '%s': %s\n", repoName, previousTag)
	} else {
		fmt.Println("No previous valid image tag found in the repository.")
	}

	return latestTag, previousTag
}
