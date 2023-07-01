package studio_integration

func (c studioIntegrationsController) GetSettings(studioID uint64) (map[string]interface{}, error) {

	discordFound := false
	discordIntegration, errDiscord := App.Repo.GetDiscordStudioIntegration(studioID)
	if errDiscord == nil {
		discordFound = true
	}

	slackFound := false
	slackIntegration, errSlack := App.Repo.GetSlackStudioIntegration(studioID)
	if errSlack == nil {
		slackFound = true
	}

	studio, _ := App.Repo.GetStudio(studioID)
	return map[string]interface{}{
		"discord":                  discordFound,
		"discordDm":                studio.DiscordNotificationsEnabled,
		"discordIntegrationStatus": discordIntegration.IntegrationStatus,
		"slack":                    slackFound,
		"slackDm":                  studio.SlackNotificationsEnabled,
		"slackIntegrationStatus":   slackIntegration.IntegrationStatus,
	}, nil
}
