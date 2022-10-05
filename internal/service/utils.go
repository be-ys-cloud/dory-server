package service

import (
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/authentication/token"
	"github.com/be-ys-cloud/dory-server/internal/authentication/totp"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/structures"
)

func checkAuth(username string, authentication structures.Authentication) (bool, error) {
	if authentication.Token != "" {
		if !token.VerifyKey(username, authentication.Token) {
			return false, &structures.CustomError{Text: "bad token supplied, or expired link", HttpCode: 400}
		} else {
			return true, nil
		}
	}
	if authentication.TOTP != "" {
		userDn, err := ldap.GetUserDN(username)
		if err != nil {
			return false, err
		}
		return totp.VerifyTOTP(userDn, authentication.TOTP)
	}

	return false, errors.New("not implemented")
}
