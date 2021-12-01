package mailer

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"strconv"
	"structures"
	"text/template"
)

func SendMail(templateName string, destEmail string, args interface{}) structures.Error {

	// Receiver email address.
	to := []string{	destEmail }


	// Authentication.
	var auth smtp.Auth = nil
	if Conf.MailServer.Password != "" {
		auth = smtp.PlainAuth("", Conf.MailServer.SenderAddress, Conf.MailServer.Password, Conf.MailServer.Address)
	}

	//Templating
	t, _ := template.ParseFiles("templates/"+templateName + ".html")

	var body bytes.Buffer

	headers := make(map[string]string)
	headers["Subject"] = Conf.MailServer.Subject
	headers["From"] = Conf.MailServer.SenderName
	headers["To"] = destEmail
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html"
	headers["Charset"] = "\"UTF-8\""

	for k, v := range headers {
		body.WriteString(k + ": " + v + "\r\n")
	}


	err := t.Execute(&body, args)
	if err != nil {
		logrus.Warnln("Unable to parse template ! " + err.Error())
		return structures.Error{Error: errors.New("unable to parse template"), HttpCode: 500}
	}

	// Sending email.
	err = smtp.SendMail(Conf.MailServer.Address+":"+strconv.Itoa(Conf.MailServer.Port), auth, Conf.MailServer.SenderAddress, to, body.Bytes())
	if err != nil {
		logrus.Warnln("Failed to send mail to user ! error was : " + err.Error())
		return structures.Error{Error: errors.New("failed to send mail"), HttpCode: 500}
	}

	return structures.Error{}
}
