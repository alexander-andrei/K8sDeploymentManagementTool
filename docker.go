package main

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/heroku/docker-registry-client/registry"
)

func LatestAndPreviousImageTags() (string, string) {
	client, err := registry.New(GlobalConfig.Docker.RegistryURL, GlobalConfig.Docker.Username, GlobalConfig.Docker.Password)
	if err != nil {
		panic(err)
	}

	tags, err := client.Tags(GlobalConfig.Docker.RepoName)
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

	fmt.Printf("Latest valid image tag for repository '%s': %s\n", GlobalConfig.Docker.RepoName, latestTag)
	if previousTag != "" {
		fmt.Printf("Previous valid image tag for repository '%s': %s\n", GlobalConfig.Docker.RepoName, previousTag)
	} else {
		fmt.Println("No previous valid image tag found in the repository.")
	}

	return latestTag, previousTag
}
