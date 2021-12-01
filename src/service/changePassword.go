package service

import (
	"ad"
	"github.com/sirupsen/logrus"
	"mailer"
	"structures"
)

func ChangePassword(username string, oldPassword string, newPassword string) structures.Error {

	// Send request to AD
	err := ad.ChangePassword(username, oldPassword, newPassword)

	if err.Error != nil {
		return err
	}

	// Send email
	email, err := ad.GetUserEmail(username)

	if err.Error != nil {
		logrus.Warnln("Could not change password changed mail to user " + username + " because there is no email associated to it on Active Directory.")
	} else {
		mailer.SendMail("mail_info_changed", email, struct{Name string}{Name: username})
	}

	return structures.Error{}

}
