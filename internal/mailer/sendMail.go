package mailer

import (
	"bytes"
	"crypto/tls"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"strconv"
	"text/template"
)

func SendMail(templateName string, destEmail string, args interface{}) error {

	// Receiver email address.
	to := []string{destEmail}

	// Authentication.
	var auth smtp.Auth = nil
	if configuration.Configuration.MailServer.Password != "" {
		auth = smtp.PlainAuth("", configuration.Configuration.MailServer.SenderAddress, configuration.Configuration.MailServer.Password, configuration.Configuration.MailServer.Address)
	}

	//Templating
	t, _ := template.ParseFiles("templates/" + templateName + ".html")

	var body bytes.Buffer

	headers := make(map[string]string)
	headers["Subject"] = configuration.Configuration.MailServer.Subject
	headers["From"] = configuration.Configuration.MailServer.SenderName
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
		return &structures.CustomError{Text: "unable to parse template", HttpCode: 500}
	}

	// Sending email.
	err = sendMail(auth, configuration.Configuration.MailServer.SenderAddress, to, body.Bytes())
	if err != nil {
		logrus.Warnln("Failed to send mail to user ! error was : " + err.Error())
		return &structures.CustomError{Text: "failed to send mail", HttpCode: 500}
	}

	return nil
}

func sendMail(a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(configuration.Configuration.MailServer.Address + ":" + strconv.Itoa(configuration.Configuration.MailServer.Port))

	if err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{
			InsecureSkipVerify: configuration.Configuration.MailServer.SkipTLSVerify,
		}

		if err = c.StartTLS(config); err != nil {
			return err
		}

	}
	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
