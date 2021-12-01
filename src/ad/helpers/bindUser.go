package helpers

import (
	"errors"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"structures"
)

func BindUser(l *ldap.Conn, username string, password string) structures.Error {
	err := l.Bind(username, password)

	if err != nil {
		logrus.Warnln("Bad login for user "+ username+" ! Error was : " + err.Error())
		return structures.Error{Error :errors.New("invalid credential given to ldap server"), HttpCode: 500}
	}

	return structures.Error{}
}