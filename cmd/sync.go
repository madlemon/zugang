package cmd

import (
	"fmt"
	"zugang/internal/bitwarden"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pull the latest vault data from server.",

	Run: func(cmd *cobra.Command, args []string) {
		err := bitwarden.Sync()
		if err != nil {
			fmt.Println("Error pulling vault data:", err)
			return
		}
		fmt.Println("Syncing complete.")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
