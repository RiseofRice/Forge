package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const Version = "0.2.0"

var (
	verbose    bool
	outputFmt  string
)

var rootCmd = &cobra.Command{
	Use:     "forge",
	Short:   "ForgeCLI — a modular data transformation and analysis framework",
	Long:    `ForgeCLI is a lightweight, modular, terminal-based data transformation and analysis tool inspired by CyberChef.`,
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		} else {
			zerolog.SetGlobalLevel(zerolog.Disabled)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "text", "Output format: text or json")
}
