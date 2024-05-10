package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"zugang/internal/bitwarden"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zugang",
	Short: "Securely connect to remote hosts via SSH using credentials stored in your Bitwarden vault.",
	Long: `zugang is a CLI tool for securely connecting to remote hosts via SSH using credentialsstored in your Bitwarden vault.

- To connect to a remote host named "example.com":
    zugang login example.com

- To sync the latest vault data from the server:
    zugang sync`,
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
