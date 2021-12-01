package ad

import (
	"ad/helpers"
	"errors"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"structures"
)

func UnlockAccount(username string) structures.Error {
	l, err := helpers.GetSession(Conf.ActiveDirectory.Address, Conf.ActiveDirectory.Port, Conf.ActiveDirectory.SkipTLSVerify)
	if err.Error != nil {
		logrus.Warnln("UnlockAccount service : Could not connect to server")
		return err
	}

	defer l.Close()

	// Search user in database
	err = helpers.BindUser(l, Conf.ActiveDirectory.Admin.Username, Conf.ActiveDirectory.Admin.Password)
	if err.Error != nil {
		logrus.Warnln("UnlockAccount service : Could not login to Active Directory : Bad AD Password supplied")
		return err
	}

	user, err := helpers.GetUser(l, Conf.ActiveDirectory.BaseDN, Conf.ActiveDirectory.FilterOn, username)
	if err.Error != nil {
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
		logrus.Warnln("Could not unlock account for user " + username +" : " + erro.Error())
		return structures.Error{Error: errors.New("could not unlock account"), HttpCode: 500}
	}

	l.Close()
	return structures.Error{}
}
