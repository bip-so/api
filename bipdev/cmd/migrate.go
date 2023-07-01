/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"gitlab.com/phonepost/bip-be-platform/internal/models"

	"github.com/spf13/cobra"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates the models to the database.",
	Long:  `We need to add the models in the automigrate function`,
	Run: func(cmd *cobra.Command, args []string) {
		core.InitCore(".env", ".")
		err := postgres.GetDB().Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
		if err != nil {
			os.Exit(1)
		}

		// CreateTable(&models.CanvasBranch{})

		postgres.GetDB().AutoMigrate(
			// User
			&models.User{},
			&models.UserSocialAuth{},
			&models.UserSettings{},
			&models.UserContact{},
			&models.UserProfile{},

			// Studio
			&models.Studio{},
			&models.StudioPermission{},
			&models.Topic{},
			&models.UserAssociatedStudio{},

			// collection models
			&models.Collection{},

			// Roles and Members
			&models.Role{},
			&models.Member{},

			&models.StudioPermission{},
			&models.FollowUser{},
			&models.FollowStudio{},
			&models.CollectionPermission{},

			// Canvas Models
			&models.CanvasRepository{},
			&models.CanvasBranch{},
			//&models.CanvasBranchPublicRequest{},
			&models.CanvasBranchPermission{},

			&models.Block{},
			&models.BlockThread{},
			&models.BlockComment{},

			&models.StudioIntegration{},

			&models.Reel{},
			&models.ReelComment{},

			&models.Message{},
			// Commented this PW as this was breaking the Migarte: CC
			&models.ExternalReference{},

			// Reactions
			&models.BlockReaction{},
			&models.BlockThreadReaction{},
			&models.BlockCommentReaction{},
			&models.ReelReaction{},
			&models.ReelCommentReaction{},
			// MR
			&models.MergeRequest{},
			&models.PublishRequest{},

			// Notifications
			&models.Notification{},
			&models.NotificationCount{},

			// Casnvas Branch AccessRequest
			&models.AccessRequest{},

			//Attributions
			&models.Attribution{},
			&models.Short{},
			&models.BranchAccessToken{},

			// Integrations
			&models.IntegrationReference{},

			// Invite vis email
			&models.BranchInviteViaEmail{},
			&models.StudioInviteViaEmail{},

			&models.Post{},
			&models.PostComment{},
			&models.PostReaction{},
			&models.PostCommentReaction{},

			&models.StudioMembersRequest{},
			&models.StudioVendor{},
		)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
