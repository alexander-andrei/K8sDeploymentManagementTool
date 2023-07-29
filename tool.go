package main

func main() {
	latestTag, previousTag := LatestAndPreviousImageTags()
	errorRate, err := KibanaErrorRate()

	VerifyAndChangeImage(errorRate, err, latestTag, previousTag)
}
