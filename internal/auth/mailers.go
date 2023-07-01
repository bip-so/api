package auth

import (
	"bytes"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GetVerificationEmailHtml(cta string) string {
	var path string
	wd, _ := os.Getwd()
	if configs.GetConfigString("APP_MODE") == "local" {
		path = wd + "/templates/auth/auth-confirm-account.html"
	} else {
		path = "/bip/templates/auth/auth-confirm-account.html"
	}
	var templateBuffer bytes.Buffer
	type EmailData struct {
		CTALink string
	}
	data := EmailData{
		CTALink: cta,
	}
	htmlData, err := ioutil.ReadFile(path)
	htmlTemplate := template.Must(template.New("auth-confirm-account.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "auth-confirm-account.html", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return templateBuffer.String()
}

func (s *authService) EmailVerificationMailer(email string, uuid string, baseUrl string) {
	subject := "[BIP] Please confirm your email"
	VerificationUrl := baseUrl + "/api/v1/auth/verify-email/" + uuid
	body := GetVerificationEmailHtml(VerificationUrl)
	var mailer pkg.BipMailer
	toList := []string{email}
	emptyList := []string{}
	err2 := pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, body)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func GetOTPEmailHtml(otp string, client string, magicLink string) string {
	var path string
	wd, _ := os.Getwd()
	if configs.GetConfigString("APP_MODE") == "local" {
		path = wd + "/templates/auth/otp.html"
	} else {
		path = "/bip/templates/auth/otp.html"
	}
	var templateBuffer bytes.Buffer
	type EmailData struct {
		OTP       string
		MagicLink string
	}
	data := EmailData{
		OTP:       otp,
		MagicLink: magicLink,
	}
	htmlData, err := ioutil.ReadFile(path)
	htmlTemplate := template.Must(template.New("otp.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "otp.html", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return templateBuffer.String()
}

func (s *authService) OTPMailer(email string, otp string) {
	subject := "[BIP] Please use " + otp + " to login to your account."
	feUrl := ""
	ENVMAILER := configs.GetConfigString("ENV")
	magicLink := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + "auth/login/magic?token=" + otp + "&email=" + strings.ToLower(email)
	body := GetOTPEmailHtml(otp, feUrl, magicLink)
	var mailer pkg.BipMailer
	toList := []string{email}
	emptyList := []string{}
	err2 := pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, body)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func GetForgotPasswordEmailHtml(otp string, ctaLink string) string {
	var path string
	wd, _ := os.Getwd()
	if configs.GetConfigString("APP_MODE") == "local" {
		path = wd + "/templates/auth/forgotpass.html"
	} else {
		path = "/bip/templates/auth/forgotpass.html"
	}
	var templateBuffer bytes.Buffer
	type EmailData struct {
		OTP     string
		CTALINK string
	}
	data := EmailData{
		OTP:     otp,
		CTALINK: ctaLink,
	}
	htmlData, err := ioutil.ReadFile(path)
	htmlTemplate := template.Must(template.New("forgotpass.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "forgotpass.html", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return templateBuffer.String()
}

func (s *authService) ForgotPasswordMailer(email string, otp string) {
	subject := "[BIP] Forgot Password Email "
	ENVMAILER := configs.GetConfigString("ENV")
	bipAuth := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + "auth/reset/?token=" + otp + "&email=" + strings.ToLower(email)
	body := GetForgotPasswordEmailHtml(otp, bipAuth)
	var mailer pkg.BipMailer
	toList := []string{email}
	emptyList := []string{}
	err2 := pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, body)
	if err2 != nil {
		fmt.Println(err2)
	}
}
