package main

import (
	"k8s/tool/cmd"
	"k8s/tool/config"

	"github.com/spf13/cobra"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	var rootCmd = &cobra.Command{Use: "tool"}
	rootCmd.AddCommand(cmd.GetK8sCommand())
	rootCmd.Execute()
}
