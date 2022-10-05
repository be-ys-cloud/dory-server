package service

import (
	"github.com/be-ys-cloud/dory-server/internal/authentication/token"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/mailer"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
)

func UnlockAccount(user structures.UserUnlock) error {
	//Check that token is valid
	if valid, err := checkAuth(user.Username, user.Authentication); !valid || err != nil {
		return &structures.CustomError{Text: "bad authentication provided", HttpCode: 401}
	}

	//Token is valid, now modifying password !
	err := ldap.UnlockAccount(user.Username)
	if err != nil {
		logrus.Warnf("Error while unlocking user %s. Error was: %s", user.Username, err.Error())
		return &structures.CustomError{Text: "an error occurred while unlocking account", HttpCode: 500}
	}

	//Removing key
	token.DeleteKey(user.Username)

	//Modification done, sending mail
	email, err := ldap.GetUserMail(user.Username)

	if err != nil {
		logrus.Warnf("Could not send unlocked account mail to user %s because there is no email associated to it on Active Directory.", user.Username)
	} else {
		_ = mailer.SendMail("mail_info_unlocked", email, struct {
			Name string
			URL  string
			LDAP string
		}{Name: user.Username, URL: configuration.Configuration.FrontAddress, LDAP: configuration.Configuration.LDAPServer.Address})
	}

	return nil
}
