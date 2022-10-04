package ldap

import (
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap/helpers"
	"github.com/sirupsen/logrus"
)

func IsPasswordValid(username string, password string) (bool, error) {
	l, err := helpers.GetSession(configuration.Configuration.LDAPServer.Address, configuration.Configuration.LDAPServer.Port, configuration.Configuration.LDAPServer.SkipTLSVerify)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not connect to server")
		return false, err
	}

	defer l.Close()

	//Bind as admin user
	err = helpers.BindUser(l, configuration.Configuration.LDAPServer.Admin.Username, configuration.Configuration.LDAPServer.Admin.Password)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not login to LDAP : Bad password supplied")
		return false, err
	}

	user, err := helpers.GetUser(l, configuration.Configuration.LDAPServer.BaseDN, configuration.Configuration.LDAPServer.FilterOn, username)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not find user")
		return false, err
	}

	//Check user have provided correct password
	err = helpers.BindUser(l, user.DN, password)
	if err != nil {
		logrus.Warnln("ChangePassword service : Invalid old password")
		return false, err
	}

	return true, nil
}
