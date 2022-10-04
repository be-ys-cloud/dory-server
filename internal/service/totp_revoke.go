package service

import (
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/authentication/totp"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/structures"
)

func RevokeTOTP(user structures.UserCreateTOTP) error {

	validPassword, err := ldap.IsPasswordValid(user.Username, user.Password)
	if err != nil {
		return err
	}
	if !validPassword {
		return errors.New("invalid password")
	}

	userDN, err := ldap.GetUserDN(user.Username)
	if err != nil {
		return err
	}

	return totp.DeleteTOTP(userDN)
}
