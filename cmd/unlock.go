/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"zugang/internal/bitwarden"
)

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock your vault and retrieve the session key",
	Long: `The unlock command unlocks your vault in Bitwarden.
This command prompts the user to provide authentication details and then retrieves the session key required to access the vault.
`,
	Run: func(cmd *cobra.Command, args []string) {

		sessionKey, err := bitwarden.Unlock()
		if err != nil {
			fmt.Println("Error unlocking your vault", err)
			return
		}

		isRaw, _ := cmd.Flags().GetBool("raw")
		if isRaw {
			fmt.Println(sessionKey)
		} else {
			fmt.Printf(`
To unlock your vault, set your session key to the 'BW_SESSION' environment variable. ex:
	$ export BW_SESSION=%q
	> $env:BW_SESSION=%q
`, sessionKey, sessionKey)
		}
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)

	unlockCmd.Flags().BoolP("raw", "", false, "Only print session key")
}
