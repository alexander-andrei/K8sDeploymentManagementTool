package main

func main() {
	if err := LoadConfig(); err != nil {
		panic(err)
	}

	latestTag, previousTag := LatestAndPreviousImageTags()
	errorRate, err := KibanaErrorRate()

	VerifyAndChangeImage(errorRate, err, latestTag, previousTag)
}
