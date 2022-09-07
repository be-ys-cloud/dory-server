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

func AskAccountUnlock(username string) error {
	// Get user email
	email, err := ad.GetUserEmail(username)

	if err != nil {
		logrus.Warnln("Could not get email for user " + username + " ! Unable to send unlock link, aborting request.")
		return &structures.CustomError{Text: "could not find your email on active directory server", HttpCode: 500}
	}

	//Create new unique key
	key := authentication.CreateKey(username)

	//Send mail to user
	err = mailer.SendMail("mail_unlock", email, struct {
		Name string
		URL  string
	}{Name: username, URL: configuration.Configuration.FrontAddress + "unlock?user=" + url.QueryEscape(username) + "&token=" + url.QueryEscape(key)})

	if err != nil {
		logrus.Warnln("Could not send to user " + username + " ! Unable to send unlock link, aborting request.")
		authentication.DeleteKey(username)
		return &structures.CustomError{Text: "email gateway is not reachable, try again later", HttpCode: 503}
	}

	return nil
}

func UnlockAccount(username string, token string) error {
	//Check that token is valid
	if !authentication.CheckKey(username, token) {
		return &structures.CustomError{Text: "bad token supplied, or expired link", HttpCode: 400}
	}

	//Token is valid, now modifying password !
	err := ad.UnlockAccount(username)
	if err != nil {
		logrus.Warnln("Error while unlocking user " + username + ". Error was : " + err.Error())
		return &structures.CustomError{Text: "an error occurred while unlocking account", HttpCode: 500}
	}

	//Removing key
	authentication.DeleteKey(username)

	//Modification done, sending mail
	email, err := ad.GetUserEmail(username)

	if err != nil {
		logrus.Warnln("Could not send unlocked account mail to user " + username + " because there is no email associated to it on Active Directory.")
	} else {
		mailer.SendMail("mail_info_unlocked", email, struct{ Name string }{Name: username})
	}

	return nil
}
