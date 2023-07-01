package discord

/*
	if body.Type == InteractionApplicationCommandAutocomplete {
		users, _ := FindUsersByDiscordIDs([]string{body.Member.User.ID})
		user := users[0]
		guildId := body.GuildID
		studioIntegration, err := studiointegration.GetStudioIntegrationByDiscordTeamId(guildId)
		if err != nil || len(studioIntegration) == 0 {
			fmt.Println("error while retrieving studio from guild id")
			return
		}
		studioID := studioIntegration[0].ID
		var searchString string
		if body.Data.Name == "bip-new" {
			searchString = body.Data.Options[1].Value
		} else {
			searchString = body.Data.Options[0].Value
		}
		_, err = permissions.App.Repo.GetMember(map[string]interface{}{"studio_id": studioID, "user_id": user.ID})

		var canvasRepos *[]models.CanvasRepository
		if err == gorm.ErrRecordNotFound {
			canvasRepos, err = AnonymousGetAllCanvasController(parentCollectionID, parentCanvasRepoID)
		} else {
			canvasRepos, err = AuthUserGetAllCanvasController(parentCollectionID, parentCanvasRepoID, user.User, studioID)
		}

		var suggestions []map[string]interface{}
		var untitledsuggestions []map[string]interface{}
		for _, val := range *canvasRepos {
			if val.Key == "" {
				continue
			}
			var name string
			if body.Data.Name == "bip-search" && val.Type == "COLLECTION" {
				continue
			}
			if val.Name == "" {
				name = "Untitled"
				x := map[string]interface{}{
					"name":  name,
					"value": val.ID,
				}
				untitledsuggestions = append(untitledsuggestions, x)
			} else {
				if searchString == "" || strings.Contains(strings.ToLower(val.Name), strings.ToLower(searchString)) {
					name = val.Name
					x := map[string]interface{}{
						"name":  name,
						"value": val.ID,
					}
					if len(suggestions) >= 25 {
						break
					}
					suggestions = append(suggestions, x)
				}
			}
		}
		for _, val := range untitledsuggestions {
			if len(suggestions) >= 25 {
				break
			}
			suggestions = append(suggestions, val)
		}
		response.RenderCustomResponse(c, map[string]interface{}{
			"type": 8,
			"data": map[string]interface{}{
				"tts":     false,
				"content": "suggestions",
				"embeds":  []string{},
				"allowed_mentions": map[string]interface{}{
					"parse": []string{},
				},
				"flags":   1 << 6,
				"choices": suggestions,
			},
		})
		return
	}
*/

//if body.Type == InteractionApplicationCommand {

/*
	if body.Data.Type == 1 && body.Data.Name == "bip-new" {
		users, _ := FindUsersByDiscordIDs([]string{body.Member.User.ID})
		user := users[0]
		guildId := body.GuildID
		studioIntegration, err := studiointegration.GetStudioIntegrationByDiscordTeamId(guildId)
		if err != nil || len(studioIntegration) == 0 {
			fmt.Println("error while retrieving studio from guild id")
			return
		}
		studioID := studioIntegration[0].ID

		title := body.Data.Options[0].Value
		var collectionID uint64
		var parentCanvasRepositoryID uint64
		if len(body.Data.Options) > 1 {
			parentCanvasRepositoryID = body.Data.Options[1].Value
		} else {
			parentCanvasRepositoryID = 0
		}

		position := 0

		var parentPage *models.Page
		if parentPageID != "" {
			pPage, err := models.GetPageByID(c.Request.Context(), parentPageID)
			if err != nil {
				fmt.Println(err)
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Please choose parent canvas title from the options",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}

			parentPage = &pPage
		}

		studio, err := studio.App.StudioRepo.GetStudioByID(studioID)
		if err != nil {
			return
		}

		repo, err := CreateCanvasRepo(title, "", user.UserID, studioID, collectionID, uint(position), parentCanvasRepositoryID)

		if err != nil {
			println("\n\n---create page......\n\n", err)
			return
		}

		response.RenderCustomResponse(c, map[string]interface{}{
			"type": 4,
			"data": map[string]interface{}{
				"tts":     false,
				"content": configs.GetConfigString("SITEROOT") + "@" + url.PathEscape(studio.Handle) + "/" + slugify.Slugify(repo.Name) + "-" + repo.Key,
				"embeds":  []string{},
				"allowed_mentions": map[string]interface{}{
					"parse": []string{},
				},
				"flags": 1 << 6,
			},
		})
		return
	}

	if body.Data.Type == 1 && body.Data.Name == "bip-search" {
		guildId := body.GuildID
		studioIntegration, err := studiointegration.GetStudioIntegrationByDiscordTeamId(guildId)
		if err != nil || len(studioIntegration) == 0 {
			fmt.Println("erro while retrieving studio from guild id")
			return
		}
		studioID := studioIntegration[0].ID
		studio, err := studio.App.StudioRepo.GetStudioByID(studioID)
		if err != nil {
			return
		}
		pageID := body.Data.Options[0].Value

		pageTemp, err := models.GetPageByID(c.Request.Context(), pageID)
		page := &pageTemp
		var message string
		if err != nil {
			users, _ := FindUsersByDiscordIDs([]string{body.Member.User.ID})
			user := users[0]
			_, err = permissions.App.Repo.GetMember(map[string]interface{}{"studio_id": studioID, "user_id": user.ID})

			var canvasRepos *[]models.CanvasRepository
			if err == gorm.ErrRecordNotFound {
				canvasRepos, err = AnonymousGetAllCanvasController(parentCollectionID, parentCanvasRepoID)
			} else {
				canvasRepos, err = AuthUserGetAllCanvasController(parentCollectionID, parentCanvasRepoID, user.User, studioID)
			}

			if len(*canvasRepos) == 0 {
				message = "The Canvas does not exist. Kindly check"
			} else {
				for _, repo := range *canvasRepos {
					if repo.Name != "" && strings.Contains(strings.ToLower(page.Title), strings.ToLower(pageID)) {
						message += configs.GetConfigString("SITEROOT") + "@" + url.PathEscape(studio.Handle) + "/" + slugify.Slugify(repo.Name) + "-" + repo.Key + "\n"
					}
				}
				if message == "" {
					message = "The Canvas does not exist. Kindly check"
				}
			}
		} else {
			message = configs.GetConfigString("SITEROOT") + "@" + url.PathEscape(studio.Handle) + "/" + slugify.Slugify(page.Title) + "-" + page.Key
		}
		response.RenderCustomResponse(c, map[string]interface{}{
			"type": 4,
			"data": map[string]interface{}{
				"tts":     false,
				"content": message,
				"embeds":  []string{},
				"allowed_mentions": map[string]interface{}{
					"parse": []string{},
				},
				"flags": 1 << 6,
			},
		})
		return
	}
*/
