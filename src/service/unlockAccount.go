package service

import (
	"ad"
	"authentication"
	"errors"
	"github.com/sirupsen/logrus"
	"mailer"
	"net/url"
	"structures"
)

func AskAccountUnlock(username string) structures.Error {
	// Get user email
	email, err := ad.GetUserEmail(username)

	if err.Error != nil {
		logrus.Warnln("Could not get email for user " + username + " ! Unable to send unlock link, aborting request.")
		return structures.Error{Error: errors.New("could not find your email on active directory server"), HttpCode: 500}
	}

	//Create new unique key
	key := authentication.CreateKey(username)

	//Send mail to user
	err = mailer.SendMail("mail_unlock", email, struct{Name string; URL string}{Name: username, URL: Conf.FrontAddress+"unlock?user="+url.QueryEscape(username)+"&token="+url.QueryEscape(key)})

	if err.Error != nil {
		logrus.Warnln("Could not send to user " + username + " ! Unable to send unlock link, aborting request.")
		authentication.DeleteKey(username)
		return structures.Error{Error: errors.New("email gateway is not reachable, try again later"), HttpCode: 503}
	}

	return structures.Error{}
}

func UnlockAccount(username string, token string) structures.Error {
	//Check that token is valid
	if !authentication.CheckKey(username, token){
		return structures.Error{Error: errors.New("bad token supplied, or expired link"), HttpCode: 400}
	}

	//Token is valid, now modifying password !
	err := ad.UnlockAccount(username)
	if err.Error != nil {
		logrus.Warnln("Error while unlocking user " + username + ". Error was : " + err.Error.Error())
		return structures.Error{Error: errors.New("an error occurred while unlocking account"), HttpCode: 500}
	}

	//Removing key
	authentication.DeleteKey(username)

	//Modification done, sending mail
	email, err := ad.GetUserEmail(username)

	if err.Error != nil {
		logrus.Warnln("Could not send unlocked account mail to user " + username + " because there is no email associated to it on Active Directory.")
	} else {
		mailer.SendMail("mail_info_unlocked", email, struct{Name string}{Name: username})
	}

	return structures.Error{}
}
