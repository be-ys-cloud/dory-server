package service

import (
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/authentication/totp"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"strings"
)

func CreateTOTP(user structures.UserCreateTOTP) (structures.TOTPToken, error) {

	validPassword, err := ldap.IsPasswordValid(user.Username, user.Password)
	if err != nil {
		return structures.TOTPToken{}, err
	}
	if !validPassword {
		return structures.TOTPToken{}, errors.New("invalid password")
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

	return structures.TOTPToken{TOTP: result}, nil
}
