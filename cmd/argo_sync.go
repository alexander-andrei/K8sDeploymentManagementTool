package cmd

import (
	"fmt"
	"k8s/tool/argo"
	"k8s/tool/checker"

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

		latestTag, previousTag := checker.LatestAndPreviousImageTags()
		errorRate, err := checker.KibanaErrorRate()

		if err != nil {
			panic(err)
		}

		argo.CheckAndRevertTags(latestTag, previousTag, errorRate)
		argo.TriggerDeploymentSync(args[0])
	},
}

func GetArgoSyncCommand() *cobra.Command {
	return argoSyncCommand
}
