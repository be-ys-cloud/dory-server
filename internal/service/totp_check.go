package service

import (
	"github.com/be-ys-cloud/dory-server/internal/authentication/totp"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
)

func CheckTOTP(user structures.UserVerifyTOTP) (bool, error) {

	ldapDN, err := ldap.GetUserDN(user.Username)
	if err != nil {
		logrus.Warnf("Unable to find user %s. Server responded: %s", user.Username, err.Error())
		return false, err
	}

	return totp.VerifyTOTP(ldapDN, user.TOTP)
}
