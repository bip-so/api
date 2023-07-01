package studio

import (
	"bytes"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

/*

Invited new user
Subject:
<User> invited you to studio <Studio Name>
Body:
Hello,
You have been invited to studio <Studio Name> by <user>
Sign up to join the studio
[Sign Up] (button)
*/

// bipAuth,studioName, InvitedByUser.FullName
func InviteNewUserSendMailerHTML(cta string, studioName string, username string) string {
	var path string
	wd, _ := os.Getwd()
	if configs.GetConfigString("APP_MODE") == "local" {
		path = wd + "/templates/auth/new-user-invite.html"
	} else {
		path = "/bip/templates/auth/new-user-invite.html"
	}
	var templateBuffer bytes.Buffer
	type EmailData struct {
		CTALink    string
		StudioName string
		Username   string
	}
	data := EmailData{
		CTALink:    cta,
		StudioName: studioName,
		Username:   username,
	}
	htmlData, err := ioutil.ReadFile(path)
	htmlTemplate := template.Must(template.New("new-user-invite.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "new-user-invite.html", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return templateBuffer.String()
}

func (s studioService) InviteNewUserSendMailer(email string, InvitedByUser *models.User, studioID uint64) {
	// get studio instance
	var studio models.Studio
	ENVMAILER := configs.GetConfigString("ENV")
	err := postgres.GetDB().Model(models.Studio{}).Where("id = ?", studioID).First(&studio).Error
	//studioName := studio.Handle
	//studioURL := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + "@" + studioName
	bipAuth := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + studio.Handle
	//bipAuth := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + "auth/signin"

	subject := InvitedByUser.FullName + " invited you to studio " + studio.DisplayName
	// InviteNewUserSendMailerHTML
	body := InviteNewUserSendMailerHTML(bipAuth, studio.DisplayName, InvitedByUser.FullName)
	//body := "Hello, <br> You have been invite to Studio " + studioURL + " by " + InvitedByUser.Username + "</strong></div>"
	//body = body + "<br><br> Please follow the link " + bipAuth + ". Create your account to join the studio."
	//bodyPlainText := "Hello, <br> You have been invite to canvas in " + studioURL + " by " + InvitedByUser.Username + ".  Please follow the link " + bipAuth + " and create your account to join the studio. "

	var mailer pkg.BipMailer
	toList := []string{email}
	emptyList := []string{}
	err = pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, body)
	if err != nil {
		fmt.Println(err)
	}
}

func (s studioService) InformUserAddedToStudioSendMailer(email string, InvitedByUser *models.User, studioID uint64) {
	// get studio instance
	var studio models.Studio
	var existingInvitedUser models.User
	ENVMAILER := configs.GetConfigString("ENV")
	err := postgres.GetDB().Model(models.Studio{}).Where("id = ?", studioID).First(&studio).Error
	_ = postgres.GetDB().Model(models.User{}).Where("email = ?", email).First(&existingInvitedUser).Error

	studioName := studio.Handle
	studioURL := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + "@" + studioName
	body := "Hello " + existingInvitedUser.Username + ",<br> You have been added to Studio " + studioURL + " by " + InvitedByUser.Username + "</strong></div>"
	bodyPlainText := "Hello " + existingInvitedUser.Username + ", You have been added to Studio " + studioURL + " by " + InvitedByUser.Username
	subject := existingInvitedUser.Username + " added to studio" + studio.DisplayName
	var mailer pkg.BipMailer
	toList := []string{email}
	emptyList := []string{}
	err = pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, bodyPlainText)
	if err != nil {
		fmt.Println(err)
	}
}
