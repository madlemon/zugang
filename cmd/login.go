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

// loginCmd represents the connect login
var loginCmd = &cobra.Command{
	Short: "Connect to a host via ssh",
	Long: `The login command enables you to connect to a remote host using credentials from your vault. 
If no Bitwarden session is active you may be prompted for the master password to initiate a session.
Credentials within your vault need the following URL pattern "ssh://{host}" to be eligible as credentials for the given host.
If multiple credentials exists for a host you need to specify a specific user using the --user or -u flag.
If successful, it attempts to open the SSH connection.
`,
	Use:  "login <host> [flags]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]

		specificUser, _ := cmd.Flags().GetString("user")
		hostKeyCheckEnabled, _ := cmd.Flags().GetBool("hostKeyCheck")

		username, password, err := bitwarden.FindSSHCredentials(host, specificUser)
		if err != nil {
			fmt.Println("Failed finding SSH credentials:", err)
			return
		}

		fmt.Println("Connecting to", host, "as", username)

		sshConf := &ssh.Conf{
			Host:                host,
			User:                username,
			Password:            password,
			HostKeyCheckEnabled: hostKeyCheckEnabled,
		}
		ssh.Connect(sshConf)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.
	loginCmd.Flags().StringP("user", "u", "", "connect with specified username")
	loginCmd.Flags().Bool("hostKeyCheck", true, "enable/disable host key check")
}