package temp_emails

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

//var mailer pkg.BipMailer
//toList := []string{"chirax@gmail.com"}
//emptyList := []string{}
//subject := "Testing Dorothy 3"
//body := temp_emails.GetHTMLTestTemplate()
//_ = pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, body)

func GetHTMLTestTemplate() string {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wd)

	var templateBuffer bytes.Buffer
	type EmailData struct {
		FirstName string
		LastName  string
	}
	data := EmailData{
		FirstName: "John",
		LastName:  "Doe",
	}
	path := wd + "/templates/temp_emails/test.html"

	htmlData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	htmlTemplate := template.Must(template.New("test.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "test.html", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return templateBuffer.String()
}
