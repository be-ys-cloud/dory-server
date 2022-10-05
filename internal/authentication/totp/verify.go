package totp

import (
	"encoding/base32"
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/database"
	"github.com/sirupsen/logrus"
)
import "github.com/pquerna/otp/totp"

func VerifyTOTP(userId string, token string) (bool, error) {
	if configuration.Configuration.Features.DisableTOTP {
		logrus.Warnf("User %s tried to verify TOTP, but this function is disabled.", userId)
		return false, errors.New("totp is disabled on this server")
	}

	// Get token from database
	tokenUser, err := database.GetToken(encodeUser(userId))
	if err != nil {
		logrus.Warnf("Failed to generate TOTP for user %s. Error was: %s.", userId, err.Error())
		return false, err
	}

	// Return validation state
	return totp.Validate(token, base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(tokenUser))), nil
}
