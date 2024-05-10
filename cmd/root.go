/*
Copyright © 2024 NAME HERE klocke@volavis.de
*/
package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"zugang/internal/bitwarden"
)

var BWSessionKey string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zugang",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	bitwarden.CheckForExecutable()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zugang.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	//rootCmd.PersistentFlags().StringVar(&BWSessionKey, "bw_session", "", "Bitwarden session key")
	//_ = viper.BindPFlag("bw_session", rootCmd.PersistentFlags().Lookup("bw_session"))

}
