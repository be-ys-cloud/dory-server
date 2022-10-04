package ldap

import (
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap/helpers"
	"github.com/sirupsen/logrus"
)

func ChangePassword(username string, oldPassword string, newPassword string) error {

	l, err := helpers.GetSession(configuration.Configuration.LDAPServer.Address, configuration.Configuration.LDAPServer.Port, configuration.Configuration.LDAPServer.SkipTLSVerify)
	if err != nil {
		logrus.Warnln("ChangePassword service: Could not connect to server")
		return err
	}

	defer l.Close()

	//Bind as admin user
	err = helpers.BindUser(l, configuration.Configuration.LDAPServer.Admin.Username, configuration.Configuration.LDAPServer.Admin.Password)
	if err != nil {
		logrus.Warnln("ChangePassword service: Could not login to LDAP : Bad AD Password supplied")
		return err
	}

	// Find user
	user, err := helpers.GetUser(l, configuration.Configuration.LDAPServer.BaseDN, configuration.Configuration.LDAPServer.FilterOn, username)
	if err != nil {
		logrus.Warnln("ChangePassword service: Could not find user")
		return err
	}

	//Check user have provided correct password
	err = helpers.BindUser(l, user.DN, oldPassword)
	if err != nil {
		logrus.Warnln("ChangePassword service: Invalid old password")
		return err
	}

	//Re-rebind as admin
	err = helpers.BindUser(l, configuration.Configuration.LDAPServer.Admin.Username, configuration.Configuration.LDAPServer.Admin.Password)
	if err != nil {
		logrus.Warnln("ChangePassword service: Could not login to LDAP: Bad password supplied")
		return err
	}

	err = helpers.ChangePassword(l, user.DN, newPassword)
	if err != nil {
		logrus.Warnf("ChangePassword service: Could not change password ! Error was %s", err.Error())
		return err
	}

	l.Close()
	return nil
}
