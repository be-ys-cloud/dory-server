package ad

import (
	"github.com/be-ys-cloud/dory-server/internal/ad/helpers"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
)

func UnlockAccount(username string) error {
	l, err := helpers.GetSession(configuration.Configuration.ActiveDirectory.Address, configuration.Configuration.ActiveDirectory.Port, configuration.Configuration.ActiveDirectory.SkipTLSVerify)
	if err != nil {
		logrus.Warnln("UnlockAccount service : Could not connect to server")
		return err
	}

	defer l.Close()

	// Search user in database
	err = helpers.BindUser(l, configuration.Configuration.ActiveDirectory.Admin.Username, configuration.Configuration.ActiveDirectory.Admin.Password)
	if err != nil {
		logrus.Warnln("UnlockAccount service : Could not login to Active Directory : Bad AD Password supplied")
		return err
	}

	user, err := helpers.GetUser(l, configuration.Configuration.ActiveDirectory.BaseDN, configuration.Configuration.ActiveDirectory.FilterOn, username)
	if err != nil {
		logrus.Warnln("UnlockAccount service : Could not find user")
		return err
	}

	attrs := ldap.PartialAttribute{Type: "lockoutTime", Vals: []string{"0"}}

	passReq := &ldap.ModifyRequest{
		DN:      user.DN,
		Changes: []ldap.Change{{2, attrs}},
	}

	erro := l.Modify(passReq)
	if erro != nil {
		logrus.Warnln("Could not unlock account for user " + username + " : " + erro.Error())
		return &structures.CustomError{Text: "could not unlock account", HttpCode: 500}
	}

	l.Close()
	return nil
}
