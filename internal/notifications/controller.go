package notifications

func (c notificationController) MarkNotificationsAsRead(userID uint64) error {
	updates := map[string]interface{}{"seen": true}
	err := App.Repo.UpdateNotifications(map[string]interface{}{"notifier_id": userID}, updates)
	if err != nil {
		return err
	}
	err = App.Service.MakeNotificationCountAllRead(userID)
	if err != nil {
		return err
	}
	return nil
}
