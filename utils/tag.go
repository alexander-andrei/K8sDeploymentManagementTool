package utils

import "strings"

func ExtractTagVersion(imageName string) string {
	imageParts := strings.Split(imageName, ":")
	if len(imageParts) > 1 {
		return imageParts[1]
	}

	return ""
}
