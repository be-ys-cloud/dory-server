package totp

import (
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/database"
	"github.com/sirupsen/logrus"
)

func DeleteTOTP(userDN string) error {
	if configuration.Configuration.Features.DisableTOTP {
		logrus.Warnf("User %s tried to revoke TOTP, but this function is disabled.", userDN)
		return errors.New("totp is disabled on this server")
	}

	// Delete token from database
	err := database.DeleteToken(encodeUser(userDN))
	if err != nil {
		logrus.Warnf("Failed to revoke TOTP for user %s. Error was: %s.", userDN, err.Error())
		return err
	}

	return nil

}
