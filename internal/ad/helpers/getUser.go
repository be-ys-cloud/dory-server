package helpers

import (
	"fmt"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
)

func GetUser(l *ldap.Conn, baseDN string, filterOn string, username string) (*ldap.Entry, error) {

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(filterOn, ldap.EscapeFilter(username)),
		[]string{},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		logrus.Warnln("Unable to search into Active Directory. Detailed error : " + err.Error())
		return nil, &structures.CustomError{Text: "could not search into active directory", HttpCode: 503}
	}

	if len(sr.Entries) == 0 {
		logrus.Warnln("No user matched.")
		return nil, &structures.CustomError{Text: "user not found in active directory", HttpCode: 404}
	}

	if len(sr.Entries) > 1 {
		logrus.Warnln("Too many user matched.")
		return nil, &structures.CustomError{Text: "too many user matched active directory. could not process to avoid modifying one undesired account", HttpCode: 404}
	}

	return sr.Entries[0], nil
}
