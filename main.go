package main

import (
	"k8s/tool/cmd"
	"k8s/tool/config"
	"os"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening log file")
	}

	log.Logger = zerolog.New(file).With().Timestamp().Logger()
}

func main() {
	var rootCmd = &cobra.Command{Use: "tool"}
	rootCmd.AddCommand(cmd.GetK8sCommand())
	rootCmd.AddCommand(cmd.GetArgoSyncCommand())
	rootCmd.Execute()
}
