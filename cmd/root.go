package cmd

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	verbose bool
	rootCmd = &cobra.Command{
		Use:   "clerk-gopher",
		Short: "Clerk Gopher is a simple command line launcher for Toontown Rewritten",
		Long:  "A simple command line launcher written in Go to allow simple and fast login with login saving functionality",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "sets logging to verbose")
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	log.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
		DisableSorting:  true,
	})

	if verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
