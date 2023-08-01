package cmd

import (
	"fmt"
	"k8s/tool/checker"

	"github.com/spf13/cobra"
)

var k8sCommand = &cobra.Command{
	Use:   "k8s-change [imageName] [version]",
	Short: "A subcommand of the tool command",
	Long: `Change checks the kibana error rate for a particular 
					deployment and reverts the deployment to the previous tag if the error 
					rate in the last 30 mins is greater than expected`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Command started. For the next 30 mins there will be a check running on:")
		fmt.Println("Deployment:", args[0])
		fmt.Println("With image:", args[1])

		latestTag, previousTag := checker.LatestAndPreviousImageTags()
		errorRate, err := checker.KibanaErrorRate()

		checker.VerifyAndChangeImage(errorRate, err, latestTag, previousTag, args[0], args[1])
	},
}

func GetK8sCommand() *cobra.Command {
	return k8sCommand
}
