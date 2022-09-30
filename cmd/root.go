package cmd

import (
	"github.com/rs/zerolog"

	"os"

	"github.com/spf13/cobra"
)

var (
	log     = zerolog.New(os.Stderr).With().Timestamp().Logger()
	rootCmd = &cobra.Command{
		Use:           "envoy-control-plane",
		Short:         "Bare bones Implementation of an Envoy xDS server",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
