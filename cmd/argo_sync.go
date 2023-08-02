package cmd

import (
	"fmt"
	"k8s/tool/argo"
	"k8s/tool/checker"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var argoSyncCommand = &cobra.Command{
	Use:   "argo-revert-and-sync",
	Short: "Revert latest deployment if kibana errors are too many",
	Long: `Change checks the kibana error rate for a particular 
					deployment and reverts the deployment to the previous tag if the error 
					rate in the last 30 mins is greater than expected`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Command started. For the next 30 mins there will be a check running on:")
		fmt.Println("Deployment:", args[0])

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

		err = argo.CheckAndRevertTags(latestTag, previousTag, errorRate, args[0])
		if err != nil {
			log.Error().Str("Application", args[0]).Str("LatestTag:", latestTag).Str("PreviousTag:", previousTag).Err(err).Msg("An error occured while checking and reverting tags")
			panic(err)
		}

		err = argo.TriggerDeploymentSync(args[0])
		if err != nil {
			log.Error().Str("Application", args[0]).Err(err).Msg("An error occured while triggering deployment sync")
			panic(err)
		}
	},
}

func GetArgoSyncCommand() *cobra.Command {
	return argoSyncCommand
}
