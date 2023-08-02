package main

import (
	"k8s/tool/argo"
	"k8s/tool/config"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	argo.CheckAndRevertTags("2.39", "2", 2.2, "agnhost")

	// var rootCmd = &cobra.Command{Use: "tool"}
	// rootCmd.AddCommand(cmd.GetK8sCommand())
	// rootCmd.AddCommand(cmd.GetArgoSyncCommand())
	// rootCmd.Execute()
}
