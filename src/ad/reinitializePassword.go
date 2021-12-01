package ad

import (
	"ad/helpers"
	"errors"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
	"structures"
)

func ReinitializePassword(username string, newPassword string) structures.Error {
	l, err := helpers.GetSession(Conf.ActiveDirectory.Address, Conf.ActiveDirectory.Port, Conf.ActiveDirectory.SkipTLSVerify)
	if err.Error != nil {
		logrus.Warnln("reinitializePassword service : Could not connect to server")
		return err
	}

	defer l.Close()

	// Search user in database
	err = helpers.BindUser(l, Conf.ActiveDirectory.Admin.Username, Conf.ActiveDirectory.Admin.Password)
	if err.Error != nil {
		logrus.Warnln("reinitializePassword service : Could not login to Active Directory : Bad AD Password supplied")
		return err
	}

	user, err := helpers.GetUser(l, Conf.ActiveDirectory.BaseDN, Conf.ActiveDirectory.FilterOn, username)
	if err.Error != nil {
		logrus.Warnln("reinitializePassword service : Could not find user")
		return err
	}


	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	encoded, erro := utf16.NewEncoder().String("\""+newPassword+"\"")

	if erro != nil {
		logrus.Warnln("reinitializePassword service : could not parse new password (wtf?!)")
		return structures.Error{ Error: errors.New("could not parse password to utf16"), HttpCode: 500}
	}

	attrs := ldap.PartialAttribute{Type: "unicodePwd", Vals: []string{encoded}}

	passReq := &ldap.ModifyRequest{
		DN:      user.DN,
		Changes: []ldap.Change{{2, attrs}},
	}

	erro = l.Modify(passReq)
	if erro != nil {
		logrus.Warnln("Could not change password for user " + username +" : " + erro.Error())
		return structures.Error{Error: errors.New("could not change password"), HttpCode: 500}
	}

	l.Close()
	return structures.Error{}
}
