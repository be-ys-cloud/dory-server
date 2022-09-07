package service

import (
	"github.com/be-ys-cloud/dory-server/internal/ad"
	"github.com/be-ys-cloud/dory-server/internal/authentication"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/mailer"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
	"net/url"
)

func AskPasswordReinitialization(username string) error {
	// Get user email
	email, err := ad.GetUserEmail(username)

	if err != nil {
		logrus.Warnln("Could not get email for user " + username + " ! Unable to send reinitialization password link, aborting request.")
		return &structures.CustomError{Text: "could not find your email on active directory server", HttpCode: 500}
	}

	//Create new unique key
	key := authentication.CreateKey(username)

	//Send mail to user
	err = mailer.SendMail("mail_reinit", email, struct {
		Name string
		URL  string
	}{Name: username, URL: configuration.Configuration.FrontAddress + "reinitialize?user=" + url.QueryEscape(username) + "&token=" + url.QueryEscape(key)})

	if err != nil {
		logrus.Warnln("Could not send to user " + username + " ! Unable to send reinitialization password link, aborting request.")
		authentication.DeleteKey(username)
		return &structures.CustomError{Text: "email gateway is not reachable, try again later", HttpCode: 503}
	}

	return nil
}

func ReinitializePassword(username string, token string, newPassword string) error {
	//Check that token is valid
	if !authentication.CheckKey(username, token) {
		return &structures.CustomError{Text: "bad token supplied, or expired link", HttpCode: 400}
	}

	//Token is valid, now modifying password !
	err := ad.ReinitializePassword(username, newPassword)
	if err != nil {
		logrus.Warnln("Error while reinitializing password for user " + username + ". Error was : " + err.Error())
		return &structures.CustomError{Text: "an error occurred while reinitializing password", HttpCode: 500}
	}

	//Removing key
	authentication.DeleteKey(username)

	//Modification done, sending mail
	email, err := ad.GetUserEmail(username)

	if err != nil {
		logrus.Warnln("Could not send password changed mail to user " + username + " because there is no email associated to it on Active Directory.")
	} else {
		mailer.SendMail("mail_info_changed", email, struct{ Name string }{Name: username})
	}

	return nil
}
