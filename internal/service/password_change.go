package service

import (
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/mailer"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
)

func ChangePassword(user structures.UserChangePassword) error {

	// Send request to AD
	err := ldap.ChangePassword(user.Username, user.OldPassword, user.NewPassword)

	if err != nil {
		return err
	}

	// Send email
	email, err := ldap.GetUserMail(user.Username)

	if err != nil {
		logrus.Warnf("Could not change password changed mail to user %s because there is no email associated to it on Active Directory.", user.Username)
	} else {
		_ = mailer.SendMail("mail_info_changed", email, struct {
			Name string
			URL  string
			LDAP string
		}{Name: user.Username, URL: configuration.Configuration.FrontAddress, LDAP: configuration.Configuration.LDAPServer.Address})
	}

	return nil

}
