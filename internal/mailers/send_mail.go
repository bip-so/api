package mailers

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"reflect"
)

func (s mailersService) SendEmailEvent(event string, mailData SendEmail) {
	sendEmailString, _ := json.Marshal(mailData)
	s.kafka.Publish(configs.KAFKA_TOPICS_EMAILS, event, sendEmailString)
}

func (s mailersService) ReceiveEmailEvent(event string, mailData SendEmail) {
	eventHandler := event + "Handler"

	// populating userIDs to emails
	if len(mailData.userIDs) > 0 {
		users, _ := App.Repo.GetUsersByIDs(mailData.userIDs)
		toEmails := []string{}
		for _, user := range users {
			if user.Email.Valid {
				toEmails = append(toEmails, user.Email.String)
			}
		}
		mailData.ToEmails = append(mailData.ToEmails, toEmails...)
	}

	t := mailersService{}
	method := reflect.ValueOf(t).MethodByName(eventHandler)
	params := []reflect.Value{reflect.ValueOf(mailData)}
	method.Call(params)
}
