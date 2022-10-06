package service

import (
	"github.com/be-ys-cloud/dory-server/internal/authentication/totp"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/mailer"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
	"strings"
)

func CreateTOTP(user structures.UserCreateTOTP) (structures.TOTPToken, error) {

	validPassword, err := ldap.IsPasswordValid(user.Username, user.Password)
	if err != nil {
		return structures.TOTPToken{}, &structures.CustomError{HttpCode: 401, Text: err.Error()}
	}
	if !validPassword {
		return structures.TOTPToken{}, &structures.CustomError{HttpCode: 401, Text: err.Error()}
	}

	userDN, err := ldap.GetUserDN(user.Username)
	if err != nil {
		return structures.TOTPToken{}, err
	}

	result, err := totp.CreateTOTP(userDN, user.Username)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return structures.TOTPToken{}, &structures.CustomError{Text: "TOTP already active for this account; please revoke before.", HttpCode: 400}
		}
		return structures.TOTPToken{}, err
	}

	// Send email
	email, err := ldap.GetUserMail(user.Username)

	if err != nil {
		logrus.Warnf("Could not send totp created mail to user %s because there is no email associated to it on Active Directory.", user.Username)
	} else {
		_ = mailer.SendMail("mail_totp_created", email, struct {
			Name string
			URL  string
			LDAP string
		}{Name: user.Username, URL: configuration.Configuration.FrontAddress, LDAP: configuration.Configuration.LDAPServer.Address})
	}

	return structures.TOTPToken{TOTP: result}, nil
}
