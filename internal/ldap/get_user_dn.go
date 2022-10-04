package ldap

import (
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap/helpers"
	"github.com/sirupsen/logrus"
)

func GetUserDN(username string) (string, error) {

	l, err := helpers.GetSession(configuration.Configuration.LDAPServer.Address, configuration.Configuration.LDAPServer.Port, configuration.Configuration.LDAPServer.SkipTLSVerify)
	if err != nil {
		logrus.Warnln("GetUserMail service : Could not connect to server")
		return "", err
	}

	defer l.Close()

	//Connect to Active Directory as user
	err = helpers.BindUser(l, configuration.Configuration.LDAPServer.Admin.Username, configuration.Configuration.LDAPServer.Admin.Password)
	if err != nil {
		logrus.Warnln("GetUserMail service : Could not login to LDAP : Bad password supplied")
		return "", err
	}

	user, err := helpers.GetUser(l, configuration.Configuration.LDAPServer.BaseDN, configuration.Configuration.LDAPServer.FilterOn, username)
	if err != nil {
		logrus.Warnln("GetUserMail service : Could not find user")
		return "", err
	}

	return user.DN, nil
}
