/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"zugang/internal/bitwarden"

	"github.com/spf13/cobra"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Discards any stored session and locks your vault",
	Long: `The lock command manually locks your vault in Bitwarden and discards any stored session key.
This effectively revokes access to the vault until it is unlocked again.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := bitwarden.Lock()
		if err != nil {
			fmt.Println("Error locking bitwarden vault:", err)
		} else {
			fmt.Println("Your vault is locked.")
		}
		err = bitwarden.DiscardSession()
		if err != nil {
			fmt.Println("Error discarding session:", err)
		} else {
			fmt.Println("The session discarded.")
		}
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)
}
