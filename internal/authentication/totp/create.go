package totp

import (
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/database"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
)
import "github.com/pquerna/otp/totp"

func CreateTOTP(userDN string, userName string) (string, error) {
	if configuration.Configuration.Features.DisableTOTP {
		logrus.Warnf("User %s tried to generate TOTP, but this function is disabled.", userDN)
		return "", errors.New("totp is disabled on this server")
	}

	// Create random security key and store it into database
	secret := randstr.String(256)
	err := database.CreateToken(encodeUser(userDN), secret)
	if err != nil {
		logrus.Warnf("Failed to generate TOTP for user %s. Error was: %s.", userName, err.Error())
		return "", err
	}

	// Generate TOTP
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      configuration.Configuration.TOTP.CustomServiceName,
		AccountName: userName,
		Secret:      []byte(secret),
		Period:      30,
	})
	if err != nil {
		logrus.Warnf("Failed to generate TOTP for user %s. Error was: %s.", userName, err.Error())
		return "", err
	}

	return key.String(), nil

}
