package cmd

import (
	"fmt"
	"k8s/tool/checker"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var k8sCommand = &cobra.Command{
	Use:   "k8s-change [imageName] [version]",
	Short: "Revert latest deployment if kibana errors are too many",
	Long: `Change checks the kibana error rate for a particular 
					deployment and reverts the deployment to the previous tag if the error 
					rate in the last 30 mins is greater than expected`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Command started. For the next 30 mins there will be a check running on:")
		fmt.Println("Deployment:", args[0])
		fmt.Println("With image:", args[1])

		latestTag, previousTag, err := checker.LatestAndPreviousImageTags()
		if err != nil {
			log.Error().Str("Application", args[0]).Str("LatestTag:", latestTag).Str("PreviousTag:", previousTag).Err(err).Msg("An error occured while getting latest and previous tags")
			panic(err)
		}

		errorRate, err := checker.KibanaErrorRate()
		if err != nil {
			log.Error().Str("Application", args[0]).Str("LatestTag:", latestTag).Str("PreviousTag:", previousTag).Err(err).Msg("An error occured while checking kibana error rate")
			panic(err)
		}

		err = checker.VerifyAndChangeImage(errorRate, err, latestTag, previousTag, args[0], args[1])

		if err != nil {
			log.Error().Str("Application", args[0]).Str("LatestTag:", latestTag).Str("PreviousTag:", previousTag).Err(err).Msg("An error occured while verifying and changing the image")
		}
	},
}

func GetK8sCommand() *cobra.Command {
	return k8sCommand
}
