package helpers

import (
	"crypto/tls"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetSession(addr string, port int, untrustedCert bool) (*ldap.Conn, error) {
	l, err := ldap.DialTLS("tcp", addr+":"+strconv.Itoa(port), &tls.Config{InsecureSkipVerify: untrustedCert})

	if err != nil {
		logrus.Warnln("Unable to reach LDAP server ! Error was : " + err.Error())
		return nil, &structures.CustomError{Text: "could not connect to LDAP server", HttpCode: 503}
	}

	return l, nil
}
