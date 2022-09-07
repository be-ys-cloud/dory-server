package service

import (
	"github.com/be-ys-cloud/dory-server/internal/ad"
	"github.com/be-ys-cloud/dory-server/internal/mailer"
	"github.com/sirupsen/logrus"
)

func ChangePassword(username string, oldPassword string, newPassword string) error {

	// Send request to AD
	err := ad.ChangePassword(username, oldPassword, newPassword)

	if err != nil {
		return err
	}

	// Send email
	email, err := ad.GetUserEmail(username)

	if err != nil {
		logrus.Warnln("Could not change password changed mail to user " + username + " because there is no email associated to it on Active Directory.")
	} else {
		mailer.SendMail("mail_info_changed", email, struct{ Name string }{Name: username})
	}

	return nil

}
