package helpers

import (
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
)

func ChangePassword(l *ldap.Conn, userDN string, newPassword string) error {
	switch configuration.Configuration.LDAPServer.Kind {
	case "openldap":
		req := ldap.PasswordModifyRequest{UserIdentity: userDN, NewPassword: newPassword}
		_, err := l.PasswordModify(&req)

		if err != nil {
			logrus.Warnf("Could not change password for user %s. Error was: %s", userDN, err.Error())
			return &structures.CustomError{Text: "could not change password", HttpCode: 500}
		}
		break
	case "ad":
		utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		encoded, err := utf16.NewEncoder().String("\"" + newPassword + "\"")

		if err != nil {
			logrus.Warnln("ChangePassword service : could not parse new password (wtf?!)")
			return &structures.CustomError{Text: "could not parse password to utf16", HttpCode: 500}
		}

		attrs := ldap.PartialAttribute{Type: "unicodePwd", Vals: []string{encoded}}

		passReq := &ldap.ModifyRequest{
			DN:      userDN,
			Changes: []ldap.Change{{2, attrs}},
		}

		err = l.Modify(passReq)
		if err != nil {
			logrus.Warnf("Could not change password for user %s. Error was: %s", userDN, err.Error())
			return &structures.CustomError{Text: "could not change password", HttpCode: 500}
		}
		break
	default:
		return &structures.CustomError{Text: "unknown provider", HttpCode: 500}
	}
	return nil
}
