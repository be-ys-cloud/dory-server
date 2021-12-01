package ad

import (
	"ad/helpers"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
	"structures"
)

func ChangePassword(username string, oldPassword string, newPassword string) structures.Error {

	l, err := helpers.GetSession(Conf.ActiveDirectory.Address, Conf.ActiveDirectory.Port, Conf.ActiveDirectory.SkipTLSVerify)
	if err.Error != nil {
		logrus.Warnln("ChangePassword service : Could not connect to server")
		return err
	}

	defer l.Close()


	//Bind as admin user
	err = helpers.BindUser(l, Conf.ActiveDirectory.Admin.Username, Conf.ActiveDirectory.Admin.Password)
	if err.Error != nil {
		logrus.Warnln("ChangePassword service : Could not login to Active Directory : Bad AD Password supplied")
		return err
	}

	user, err := helpers.GetUser(l, Conf.ActiveDirectory.BaseDN, Conf.ActiveDirectory.FilterOn, username)
	if err.Error != nil {
		logrus.Warnln("ChangePassword service : Could not find user")
		return err
	}

	//Get object name to bind with in order to check password validity
	userCn := ""
	for _,v := range user.Attributes {
		if v.Name == "cn" {
			userCn = v.Values[0]
		}
	}
	if userCn == "" {
		return structures.Error{Error: errors.New("could not get user cn"), HttpCode: 400}
	}

	//Check user have provided correct password
	err = helpers.BindUser(l, userCn, oldPassword)
	if err.Error != nil {
		fmt.Println(err)
		logrus.Warnln("ChangePassword service : Invalid old password")
		return err
	}

	//Re-rebind as admin
	err = helpers.BindUser(l, Conf.ActiveDirectory.Admin.Username, Conf.ActiveDirectory.Admin.Password)
	if err.Error != nil {
		logrus.Warnln("ChangePassword service : Could not login to Active Directory : Bad AD Password supplied")
		return err
	}

	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	encoded, erro := utf16.NewEncoder().String("\"" + newPassword + "\"")

	if erro != nil {
		logrus.Warnln("ChangePassword service : could not parse new password (wtf?!)")
		return structures.Error{Error: errors.New("could not parse password to utf16"), HttpCode: 500}
	}

	attrs := ldap.PartialAttribute{Type: "unicodePwd", Vals: []string{encoded}}

	passReq := &ldap.ModifyRequest{
		DN:      user.DN,
		Changes: []ldap.Change{{2, attrs}},
	}

	erro = l.Modify(passReq)
	if erro != nil {
		logrus.Warnln("Could not change password for user " + username + " : " + erro.Error())
		return structures.Error{Error: errors.New("could not change password"), HttpCode: 500}
	}

	l.Close()
	return structures.Error{}
}