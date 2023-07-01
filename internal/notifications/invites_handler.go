package notifications

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s notificationService) StudioInviteByNameHandler(notification *PostNotification) {
	// Invited to Canvas (mention access to X sub-canvases too is possible and relevant)
	// NotifierID -> who is invited
	event := InvitesEntity.Events[notification.Event]
	notification.Entity = Invite
	notification.Activity = event.Activity
	user, _ := App.Repo.GetUser(notification.CreatedByID)
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, studio.DisplayName)
	notification.Priority = event.Priority
	notification.IsPersonal = InvitesEntity.IsPersonal
	notification.ExtraData.AppUrl = s.GenerateStudioUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		MessageBtnComponent{
			Type:  2,
			Label: "📄 View",
			Style: 5,
			Url:   notification.ExtraData.AppUrl,
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	if notification.ExtraData.Message != "" {
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, notification.ExtraData.Message)
	}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) StudioInviteByGroupHandler(notification *PostNotification) {
	// Invited to Canvas (mention access to X sub-canvases too is possible and relevant)
	// NotifierID -> who is invited & And all other members of the studioRole
	event := InvitesEntity.Events[notification.Event]
	notification.Entity = Invite
	notification.Activity = event.Activity
	user, _ := App.Repo.GetUser(notification.CreatedByID)
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, studio.DisplayName)
	notification.Priority = event.Priority
	notification.IsPersonal = InvitesEntity.IsPersonal
	notification.ExtraData.AppUrl = s.GenerateStudioUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		MessageBtnComponent{
			Type:  2,
			Label: "📄 View",
			Style: 5,
			Url:   notification.ExtraData.AppUrl,
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	if notification.ExtraData.Message != "" {
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, notification.ExtraData.Message)
	}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) CollectionInviteByNameHandler(notification *PostNotification) {
	// Invited to collection (mention access to X canvases)
	// NotifierID -> who is invited
	event := InvitesEntity.Events[notification.Event]
	notification.Entity = Invite
	notification.Activity = event.Activity
	user, _ := App.Repo.GetUser(notification.CreatedByID)
	collection, _ := App.Repo.GetCollectionByID(notification.ExtraData.CollectionID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, collection.Name)
	notification.Priority = event.Priority
	notification.IsPersonal = InvitesEntity.IsPersonal
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.ExtraData.AppUrl = s.GenerateCollectionUrl(studio.Handle, collection.ID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "📄 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	if notification.ExtraData.Message != "" {
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, notification.ExtraData.Message)
	}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) CollectionInviteByGroupHandler(notification *PostNotification) {
	// Invited to collection (mention access to X canvases)
	// NotifierID -> who is invited & And all other members of the collectionRole
	event := InvitesEntity.Events[notification.Event]
	notification.Entity = Invite
	notification.Activity = event.Activity
	user, _ := App.Repo.GetUser(notification.CreatedByID)
	collection, _ := App.Repo.GetCollectionByID(notification.ExtraData.CollectionID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, collection.Name)
	notification.Priority = event.Priority
	notification.IsPersonal = InvitesEntity.IsPersonal
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.ExtraData.AppUrl = s.GenerateCollectionUrl(studio.Handle, collection.ID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "📄 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	if notification.ExtraData.Message != "" {
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, notification.ExtraData.Message)
	}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) CanvasInviteByNameHandler(notification *PostNotification) {
	// Invited to canvas (mention access to X canvases)
	// NotifierID -> who is invited
	event := InvitesEntity.Events[notification.Event]
	notification.Entity = Invite
	notification.Activity = event.Activity
	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, canvasRepo.Name)
	notification.Priority = event.Priority
	notification.IsPersonal = InvitesEntity.IsPersonal
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "📄 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	if notification.ExtraData.Message != "" {
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, notification.ExtraData.Message)
	}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) CanvasInviteByGroupHandler(notification *PostNotification) {
	// Invited to canvas (mention access to X canvases)
	// NotifierID -> who is invited & And all other members
	event := InvitesEntity.Events[notification.Event]
	notification.Entity = Invite
	notification.Activity = event.Activity
	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, canvasRepo.Name)
	notification.Priority = event.Priority
	notification.IsPersonal = InvitesEntity.IsPersonal
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "📄 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	roles, _ := App.Repo.GetRolesByID(notification.RoleIDs)
	if len(roles) > 0 && roles[0].Name == models.SYSTEM_ROLE_MEMBER {
		notification.Text = fmt.Sprintf("%s invited studio members to %s", user.Username, canvasRepo.Name)
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	if notification.ExtraData.Message != "" {
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, notification.ExtraData.Message)
	}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}
