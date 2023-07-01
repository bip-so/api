/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/phonepost/bip-be-platform/bipdev/seeders"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
)

// migrateCmd represents the migrate command
var seedersCmd = &cobra.Command{
	Use:   "seed",
	Short: "seeds the database",
	Long:  `seeds the database from old database`,
	Run: func(cmd *cobra.Command, args []string) {
		core.InitCore(".env", ".")
		seeders.InitDB()
		seeders.StartUsersSeeding()
	},
}

func init() {
	rootCmd.AddCommand(seedersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
