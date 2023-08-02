package checker

import (
	"fmt"
	"k8s/tool/config"
	"regexp"
	"sort"

	"github.com/heroku/docker-registry-client/registry"
)

func LatestAndPreviousImageTags() (string, string, error) {
	client, err := registry.New(config.GlobalConfig.Docker.RegistryURL, config.GlobalConfig.Docker.Username, config.GlobalConfig.Docker.Password)
	if err != nil {
		return "", "", err
	}

	tags, err := client.Tags(config.GlobalConfig.Docker.RepoName)
	if err != nil {
		return "", "", err
	}

	if len(tags) == 0 {
		fmt.Println("No tags found in the repository.")
		return "", "", nil
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
		return "", "", nil
	}

	sort.Strings(validTags)
	latestTag := validTags[len(validTags)-1]
	previousTag := ""
	if len(validTags) > 1 {
		previousTag = validTags[len(validTags)-2]
	}

	fmt.Printf("Latest valid image tag for repository '%s': %s\n", config.GlobalConfig.Docker.RepoName, latestTag)
	if previousTag != "" {
		fmt.Printf("Previous valid image tag for repository '%s': %s\n", config.GlobalConfig.Docker.RepoName, previousTag)
	} else {
		fmt.Println("No previous valid image tag found in the repository.")
	}

	return latestTag, previousTag, nil
}
