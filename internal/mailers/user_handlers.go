package mailers

import (
	"bytes"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"path"
	"path/filepath"
	"text/template"
)

func (s mailersService) VerifyEmailHandler(mailData SendEmail) {
	mailData.Subject = "Verify Your Email"
	dir, err := filepath.Abs("./")
	fp := path.Join(dir, "templates/user/verify_email.html")
	parsedTemplate, err := template.ParseFiles(fp)
	data := map[string]interface{}{
		"verifyEmail": "https://bip.so",
	}
	buf := new(bytes.Buffer)
	if err = parsedTemplate.Execute(buf, data); err != nil {
		fmt.Println("error while compiling template with data:", err.Error())
	}
	mailData.BodyHtml = buf.String()
	err = s.Mailer.SendEmail(mailData.ToEmails, mailData.CcEmails, mailData.BccEmails, mailData.Subject, mailData.BodyHtml, mailData.BodyText)
	if err != nil {
		logger.Error(err.Error())
	}
}
