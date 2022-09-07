package helpers

import (
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
)

func BindUser(l *ldap.Conn, username string, password string) error {
	err := l.Bind(username, password)

	if err != nil {
		logrus.Warnln("Bad login for user " + username + " ! Error was : " + err.Error())
		return &structures.CustomError{Text: "invalid credential given to ldap server", HttpCode: 500}
	}

	return nil
}
