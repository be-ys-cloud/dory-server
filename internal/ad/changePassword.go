package ad

import (
	"github.com/be-ys-cloud/dory-server/internal/ad/helpers"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
)

func ChangePassword(username string, oldPassword string, newPassword string) error {

	l, err := helpers.GetSession(configuration.Configuration.ActiveDirectory.Address, configuration.Configuration.ActiveDirectory.Port, configuration.Configuration.ActiveDirectory.SkipTLSVerify)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not connect to server")
		return err
	}

	defer l.Close()

	//Bind as admin user
	err = helpers.BindUser(l, configuration.Configuration.ActiveDirectory.Admin.Username, configuration.Configuration.ActiveDirectory.Admin.Password)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not login to Active Directory : Bad AD Password supplied")
		return err
	}

	user, err := helpers.GetUser(l, configuration.Configuration.ActiveDirectory.BaseDN, configuration.Configuration.ActiveDirectory.FilterOn, username)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not find user")
		return err
	}

	//Get object name to bind with in order to check password validity
	userDn := ""
	for _, v := range user.Attributes {
		if v.Name == "distinguishedName" {
			userDn = v.Values[0]
		}
	}
	if userDn == "" {
		return &structures.CustomError{Text: "could not get user cn", HttpCode: 400}
	}

	//Check user have provided correct password
	err = helpers.BindUser(l, userDn, oldPassword)
	if err != nil {
		logrus.Warnln("ChangePassword service : Invalid old password")
		return err
	}

	//Re-rebind as admin
	err = helpers.BindUser(l, configuration.Configuration.ActiveDirectory.Admin.Username, configuration.Configuration.ActiveDirectory.Admin.Password)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not login to Active Directory : Bad AD Password supplied")
		return err
	}

	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	encoded, erro := utf16.NewEncoder().String("\"" + newPassword + "\"")

	if erro != nil {
		logrus.Warnln("ChangePassword service : could not parse new password (wtf?!)")
		return &structures.CustomError{Text: "could not parse password to utf16", HttpCode: 500}
	}

	attrs := ldap.PartialAttribute{Type: "unicodePwd", Vals: []string{encoded}}

	passReq := &ldap.ModifyRequest{
		DN:      user.DN,
		Changes: []ldap.Change{{2, attrs}},
	}

	erro = l.Modify(passReq)
	if erro != nil {
		logrus.Warnln("Could not change password for user " + username + " : " + erro.Error())
		return &structures.CustomError{Text: "could not change password", HttpCode: 500}
	}

	l.Close()
	return nil
}