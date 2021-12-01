package helpers

import (
	"crypto/tls"
	"errors"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"strconv"
	"structures"
)

func GetSession(addr string, port int, untrustedCert bool) (*ldap.Conn, structures.Error) {
	l, err := ldap.DialTLS("tcp", addr + ":" + strconv.Itoa(port), &tls.Config{InsecureSkipVerify: untrustedCert})

	if err != nil {
		logrus.Warnln("Unable to reach Active Directory server ! Error was : " + err.Error())
		return nil, structures.Error{Error :errors.New("could not connect to ldap server"), HttpCode: 503}
	}

	return l, structures.Error{}
}