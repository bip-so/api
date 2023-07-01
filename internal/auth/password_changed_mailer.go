package auth

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg"
)

func (s authService) PasswordChangedMailer(UserWhosePasswordChanged *models.User) {
	//ENVMAILER := configs.GetConfigString("ENV")
	body := "Hello " + UserWhosePasswordChanged.Username + ",<br> You have successfully changed password.</div>. "
	bodyPlainText := "Hello " + UserWhosePasswordChanged.Username + ", You have successfully changed password"
	subject := UserWhosePasswordChanged.Username + " successfully changed password."
	var mailer pkg.BipMailer
	toList := []string{UserWhosePasswordChanged.Email.String}
	emptyList := []string{}
	err := pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, bodyPlainText)
	if err != nil {
		fmt.Println(err)
	}
}
