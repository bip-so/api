package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// createAppCmd represents the createApp command
var SyncEnvCmd = &cobra.Command{
	Use:   "env-sync",
	Short: "Will Sync Env files.",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Module is created.")
	},
}

func init() {
	rootCmd.AddCommand(SyncEnvCmd)
}
