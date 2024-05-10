/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"zugang/internal/bitwarden"

	"github.com/spf13/cobra"
)

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock your vault",
	Long: `The unlock command unlocks your vault in Bitwarden.
This command prompts the user to provide authentication details if necessary and then retrieves the session key required to access the vault.
The session key is subsequently stored for future use.
`,
	Run: func(cmd *cobra.Command, args []string) {
		sessionKey, err := bitwarden.Unlock()
		if err != nil {
			fmt.Println("Error unlocking your vault", err)
			return
		}
		err = bitwarden.SaveSession(sessionKey)
		if err != nil {
			fmt.Println("Error saving the session", err)
			return
		}

		fmt.Println("Your vault is unlocked and the session is stored")
		isVerbose, _ := cmd.Flags().GetBool("verbose")
		if isVerbose {
			fmt.Println("Session key:", sessionKey)
		}
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)

	unlockCmd.Flags().BoolP("verbose", "v", false, "Print session key")
}
