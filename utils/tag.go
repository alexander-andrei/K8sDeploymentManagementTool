package utils

import (
	"strings"
)

func ExtractTagVersion(imageName string) string {
	imageParts := strings.Split(imageName, ":")
	if len(imageParts) > 1 {
		return imageParts[1]
	}

	return ""
}

func ReplaceImageTagVersion(imageName string, newTag string) string {
	parts := strings.Split(imageName, ":")
	if len(parts) != 2 {
		panic("Invalid image string format")
	}

	return parts[0] + ":" + newTag
}
