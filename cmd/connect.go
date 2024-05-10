/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"zugang/internal/bitwarden"
	"zugang/internal/ssh"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Short: "Connect to a host",
	Long: `The connect command enables you to connect to a remote host using credentials from your vault. 
If no Bitwarden session is active you may be prompted for the master password to initiate a session.
Credentials within your vault need the following URL pattern "ssh://{host}" to be eligible as credentials for the given host.
If multiple credentials exists for a host you need to specify a specific user using the --user or -u flag.
If successful, it attempts to open the SSH connection.
`,
	Use:  "connect <host> [flags]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]

		specificUser, _ := cmd.Flags().GetString("user")

		username, password, err := bitwarden.FindSSHCredentials(host, specificUser)
		if err != nil {
			fmt.Println("Failed finding SSH credentials:", err)
			return
		}

		fmt.Println("Connecting to", host, "as", username)

		ssh.OpenWithGolang(host, username, password)
		//ssh.OpenWithPattern(host, username, password)
		//if err != nil {
		//	fmt.Println("Error opening SSH connection:", err)
		//	return
		//}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.
	connectCmd.Flags().StringP("user", "u", "", "connect with specified username")
}
