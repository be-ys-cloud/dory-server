package ldap

import (
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap/helpers"
	"github.com/sirupsen/logrus"
)

func ReinitializePassword(username string, newPassword string) error {
	l, err := helpers.GetSession(configuration.Configuration.LDAPServer.Address, configuration.Configuration.LDAPServer.Port, configuration.Configuration.LDAPServer.SkipTLSVerify)
	if err != nil {
		logrus.Warnln("reinitializePassword service : Could not connect to server")
		return err
	}

	defer l.Close()

	// Search user in database
	err = helpers.BindUser(l, configuration.Configuration.LDAPServer.Admin.Username, configuration.Configuration.LDAPServer.Admin.Password)
	if err != nil {
		logrus.Warnln("reinitializePassword service : Could not login to LDAP : Bad LDAP password supplied")
		return err
	}

	user, err := helpers.GetUser(l, configuration.Configuration.LDAPServer.BaseDN, configuration.Configuration.LDAPServer.FilterOn, username)
	if err != nil {
		logrus.Warnln("reinitializePassword service : Could not find user")
		return err
	}

	err = helpers.ChangePassword(l, user.DN, newPassword)
	if err != nil {
		logrus.Warnf("reinitializePassword service: Could not change password ! Error was %s", err.Error())
		return err
	}

	l.Close()
	return nil
}
