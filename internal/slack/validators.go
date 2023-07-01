package slack2

type SlackAppMentionActionPayload struct {
	Payload string `json:"payload"`
}

type SlackAppMentionPayload struct {
	Type     string `json:"type"`
	Token    string `json:"token"`
	ActionTs string `json:"action_ts"`
	Team     struct {
		Id     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	User struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		TeamId   string `json:"team_id"`
		Name     string `json:"name"`
	} `json:"user"`
	Channel struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	IsEnterpriseInstall bool        `json:"is_enterprise_install"`
	Enterprise          interface{} `json:"enterprise"`
	CallbackId          string      `json:"callback_id"`
	TriggerId           string      `json:"trigger_id"`
	ResponseUrl         string      `json:"response_url"`
	MessageTs           string      `json:"message_ts"`
	Message             struct {
		ClientMsgId string `json:"client_msg_id"`
		Type        string `json:"type"`
		Text        string `json:"text"`
		User        string `json:"user"`
		Ts          string `json:"ts"`
		Team        string `json:"team"`
		Files       []struct {
			Id                 string `json:"id"`
			Created            int    `json:"created"`
			Timestamp          int    `json:"timestamp"`
			Name               string `json:"name"`
			Title              string `json:"title"`
			Mimetype           string `json:"mimetype"`
			Filetype           string `json:"filetype"`
			PrettyType         string `json:"pretty_type"`
			User               string `json:"user"`
			UserTeam           string `json:"user_team"`
			Editable           bool   `json:"editable"`
			Size               int    `json:"size"`
			Mode               string `json:"mode"`
			IsExternal         bool   `json:"is_external"`
			ExternalType       string `json:"external_type"`
			IsPublic           bool   `json:"is_public"`
			PublicUrlShared    bool   `json:"public_url_shared"`
			DisplayAsBot       bool   `json:"display_as_bot"`
			Username           string `json:"username"`
			UrlPrivate         string `json:"url_private"`
			UrlPrivateDownload string `json:"url_private_download"`
			MediaDisplayType   string `json:"media_display_type"`
			Thumb64            string `json:"thumb_64"`
			Thumb80            string `json:"thumb_80"`
			Thumb360           string `json:"thumb_360"`
			Thumb360W          int    `json:"thumb_360_w"`
			Thumb360H          int    `json:"thumb_360_h"`
			Thumb480           string `json:"thumb_480"`
			Thumb480W          int    `json:"thumb_480_w"`
			Thumb480H          int    `json:"thumb_480_h"`
			Thumb160           string `json:"thumb_160"`
			OriginalW          int    `json:"original_w"`
			OriginalH          int    `json:"original_h"`
			ThumbTiny          string `json:"thumb_tiny"`
			Permalink          string `json:"permalink"`
			PermalinkPublic    string `json:"permalink_public"`
			IsStarred          bool   `json:"is_starred"`
			HasRichPreview     bool   `json:"has_rich_preview"`
			FileAccess         string `json:"file_access"`
		} `json:"files"`
		Blocks []struct {
			Type     string `json:"type"`
			BlockId  string `json:"block_id"`
			Elements []struct {
				Type     string `json:"type"`
				Elements []struct {
					Type string `json:"type"`
					Text string `json:"text"`
					Url  string `json:"url"`
				} `json:"elements"`
			} `json:"elements"`
		} `json:"blocks"`
	} `json:"message"`
	Actions []struct {
		Type           string `json:"type"`
		ActionId       string `json:"action_id"`
		BlockId        string `json:"block_id"`
		SelectedOption struct {
			Text struct {
				Type  string `json:"type"`
				Text  string `json:"text"`
				Emoji bool   `json:"emoji"`
			} `json:"text"`
			Value string `json:"value"`
		} `json:"selected_option"`
		Placeholder struct {
			Type  string `json:"type"`
			Text  string `json:"text"`
			Emoji bool   `json:"emoji"`
		} `json:"placeholder"`
		ActionTs string `json:"action_ts"`
	} `json:"actions"`
}

type SlackEventTypePayload struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge,omitempty"`
	Token     string `json:"token,omitempty"`
	TeamId    string `json:"team_id,omitempty"`
	ApiAppId  string `json:"api_app_id,omitempty"`
	Event     struct {
		Type string `json:"type"`
	} `json:"event,omitempty"`
}

type SlackEventPayload struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	TeamId    string `json:"team_id"`
	ApiAppId  string `json:"api_app_id"`
	Event     struct {
		ClientMsgId string `json:"client_msg_id"`
		Type        string `json:"type"`
		Text        string `json:"text"`
		User        string `json:"user"`
		Ts          string `json:"ts"`
		Team        string `json:"team"`

		SubteamId          string   `json:"subteam_id"`
		TeamId             string   `json:"team_id"`
		DatePreviousUpdate int      `json:"date_previous_update"`
		DateUpdate         int      `json:"date_update"`
		AddedUsers         []string `json:"added_users"`
		AddedUsersCount    int      `json:"added_users_count"`
		RemovedUsers       []string `json:"removed_users"`
		RemovedUsersCount  int      `json:"removed_users_count"`

		Subteam struct {
			Id                  string      `json:"id"`
			TeamId              string      `json:"team_id"`
			IsUsergroup         bool        `json:"is_usergroup"`
			IsSubteam           bool        `json:"is_subteam"`
			Name                string      `json:"name"`
			Description         string      `json:"description"`
			Handle              string      `json:"handle"`
			IsExternal          bool        `json:"is_external"`
			DateCreate          int         `json:"date_create"`
			DateUpdate          int         `json:"date_update"`
			DateDelete          int         `json:"date_delete"`
			AutoType            interface{} `json:"auto_type"`
			AutoProvision       bool        `json:"auto_provision"`
			EnterpriseSubteamId string      `json:"enterprise_subteam_id"`
			CreatedBy           string      `json:"created_by"`
			UpdatedBy           string      `json:"updated_by"`
			DeletedBy           string      `json:"deleted_by"`
			Prefs               struct {
				Channels []interface{} `json:"channels"`
				Groups   []interface{} `json:"groups"`
			} `json:"prefs"`
			Users        []string `json:"users"`
			UserCount    int      `json:"user_count"`
			ChannelCount int      `json:"channel_count"`
		} `json:"subteam"`
		Blocks []struct {
			Type     string `json:"type"`
			BlockId  string `json:"block_id"`
			Elements []struct {
				Type     string `json:"type"`
				Elements []struct {
					Type   string `json:"type"`
					UserId string `json:"user_id"`
				} `json:"elements"`
			} `json:"elements"`
		} `json:"blocks"`
		ThreadTs     string `json:"thread_ts"`
		ParentUserId string `json:"parent_user_id"`
		Channel      string `json:"channel"`
		EventTs      string `json:"event_ts"`
	} `json:"event"`
	Type           string `json:"type"`
	EventId        string `json:"event_id"`
	EventTime      int    `json:"event_time"`
	Authorizations []struct {
		EnterpriseId        interface{} `json:"enterprise_id"`
		TeamId              string      `json:"team_id"`
		UserId              string      `json:"user_id"`
		IsBot               bool        `json:"is_bot"`
		IsEnterpriseInstall bool        `json:"is_enterprise_install"`
	} `json:"authorizations"`
	IsExtSharedChannel bool   `json:"is_ext_shared_channel"`
	EventContext       string `json:"event_context"`
	User               struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		TeamId   string `json:"team_id"`
	} `json:"user"`
	Channel struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	Team struct {
		Id     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
}

type SlackUserTeamJoinEventPayload struct {
	Token    string `json:"token"`
	TeamId   string `json:"team_id"`
	ApiAppId string `json:"api_app_id"`
	Event    struct {
		Type string `json:"type"`
		User struct {
			Id       string `json:"id"`
			TeamId   string `json:"team_id"`
			Name     string `json:"name"`
			Deleted  bool   `json:"deleted"`
			Color    string `json:"color"`
			RealName string `json:"real_name"`
			Tz       string `json:"tz"`
			TzLabel  string `json:"tz_label"`
			TzOffset int    `json:"tz_offset"`
			Profile  struct {
				Title                 string `json:"title"`
				Phone                 string `json:"phone"`
				Skype                 string `json:"skype"`
				RealName              string `json:"real_name"`
				RealNameNormalized    string `json:"real_name_normalized"`
				DisplayName           string `json:"display_name"`
				DisplayNameNormalized string `json:"display_name_normalized"`
				Fields                struct {
				} `json:"fields"`
				StatusText             string        `json:"status_text"`
				StatusEmoji            string        `json:"status_emoji"`
				StatusEmojiDisplayInfo []interface{} `json:"status_emoji_display_info"`
				StatusExpiration       int           `json:"status_expiration"`
				AvatarHash             string        `json:"avatar_hash"`
				Email                  string        `json:"email"`
				FirstName              string        `json:"first_name"`
				LastName               string        `json:"last_name"`
				Image24                string        `json:"image_24"`
				Image32                string        `json:"image_32"`
				Image48                string        `json:"image_48"`
				Image72                string        `json:"image_72"`
				Image192               string        `json:"image_192"`
				Image512               string        `json:"image_512"`
				StatusTextCanonical    string        `json:"status_text_canonical"`
				Team                   string        `json:"team"`
			} `json:"profile"`
			IsAdmin                bool   `json:"is_admin"`
			IsOwner                bool   `json:"is_owner"`
			IsPrimaryOwner         bool   `json:"is_primary_owner"`
			IsRestricted           bool   `json:"is_restricted"`
			IsUltraRestricted      bool   `json:"is_ultra_restricted"`
			IsBot                  bool   `json:"is_bot"`
			IsAppUser              bool   `json:"is_app_user"`
			Updated                int    `json:"updated"`
			IsEmailConfirmed       bool   `json:"is_email_confirmed"`
			WhoCanShareContactCard string `json:"who_can_share_contact_card"`
			Presence               string `json:"presence"`
		} `json:"user"`
		CacheTs int    `json:"cache_ts"`
		EventTs string `json:"event_ts"`
	} `json:"event"`
	Type           string `json:"type"`
	EventId        string `json:"event_id"`
	EventTime      int    `json:"event_time"`
	Authorizations []struct {
		EnterpriseId        interface{} `json:"enterprise_id"`
		TeamId              string      `json:"team_id"`
		UserId              string      `json:"user_id"`
		IsBot               bool        `json:"is_bot"`
		IsEnterpriseInstall bool        `json:"is_enterprise_install"`
	} `json:"authorizations"`
	IsExtSharedChannel bool `json:"is_ext_shared_channel"`
}

type SlackMessagesPayload struct {
	Ok       bool `json:"ok"`
	Messages []struct {
		ClientMsgId string `json:"client_msg_id"`
		Type        string `json:"type"`
		Text        string `json:"text"`
		User        string `json:"user"`
		Ts          string `json:"ts"`
		Team        string `json:"team"`
		Files       []struct {
			Id                 string `json:"id"`
			Created            int    `json:"created"`
			Timestamp          int    `json:"timestamp"`
			Name               string `json:"name"`
			Title              string `json:"title"`
			Mimetype           string `json:"mimetype"`
			Filetype           string `json:"filetype"`
			PrettyType         string `json:"pretty_type"`
			User               string `json:"user"`
			UserTeam           string `json:"user_team"`
			Editable           bool   `json:"editable"`
			Size               int    `json:"size"`
			Mode               string `json:"mode"`
			IsExternal         bool   `json:"is_external"`
			ExternalType       string `json:"external_type"`
			IsPublic           bool   `json:"is_public"`
			PublicUrlShared    bool   `json:"public_url_shared"`
			DisplayAsBot       bool   `json:"display_as_bot"`
			Username           string `json:"username"`
			UrlPrivate         string `json:"url_private"`
			UrlPrivateDownload string `json:"url_private_download"`
			MediaDisplayType   string `json:"media_display_type"`
			Thumb64            string `json:"thumb_64"`
			Thumb80            string `json:"thumb_80"`
			Thumb360           string `json:"thumb_360"`
			Thumb360W          int    `json:"thumb_360_w"`
			Thumb360H          int    `json:"thumb_360_h"`
			Thumb480           string `json:"thumb_480"`
			Thumb480W          int    `json:"thumb_480_w"`
			Thumb480H          int    `json:"thumb_480_h"`
			Thumb160           string `json:"thumb_160"`
			OriginalW          int    `json:"original_w"`
			OriginalH          int    `json:"original_h"`
			ThumbTiny          string `json:"thumb_tiny"`
			Permalink          string `json:"permalink"`
			PermalinkPublic    string `json:"permalink_public"`
			IsStarred          bool   `json:"is_starred"`
			HasRichPreview     bool   `json:"has_rich_preview"`
			FileAccess         string `json:"file_access"`
		} `json:"files"`
		Blocks []struct {
			Type     string `json:"type"`
			BlockId  string `json:"block_id"`
			Elements []struct {
				Type     string `json:"type"`
				Elements []struct {
					Type   string `json:"type"`
					Text   string `json:"text,omitempty"`
					UserId string `json:"user_id,omitempty"`
				} `json:"elements"`
			} `json:"elements"`
		} `json:"blocks"`
		ThreadTs        string   `json:"thread_ts"`
		ReplyCount      int      `json:"reply_count,omitempty"`
		ReplyUsersCount int      `json:"reply_users_count,omitempty"`
		LatestReply     string   `json:"latest_reply,omitempty"`
		ReplyUsers      []string `json:"reply_users,omitempty"`
		IsLocked        bool     `json:"is_locked,omitempty"`
		Subscribed      bool     `json:"subscribed,omitempty"`
		LastRead        string   `json:"last_read,omitempty"`
		ParentUserId    string   `json:"parent_user_id,omitempty"`
	} `json:"messages"`
	HasMore bool `json:"has_more"`
}
