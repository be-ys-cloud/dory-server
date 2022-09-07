package ad

import (
	"github.com/be-ys-cloud/dory-server/internal/ad/helpers"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
)

func GetUserEmail(username string) (string, error) {

	l, err := helpers.GetSession(configuration.Configuration.ActiveDirectory.Address, configuration.Configuration.ActiveDirectory.Port, configuration.Configuration.ActiveDirectory.SkipTLSVerify)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not connect to server")
		return "", err
	}

	defer l.Close()

	//Connect to Active Directory as user
	err = helpers.BindUser(l, configuration.Configuration.ActiveDirectory.Admin.Username, configuration.Configuration.ActiveDirectory.Admin.Password)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not login to Active Directory : Bad AD Password supplied")
		return "", err
	}

	user, err := helpers.GetUser(l, configuration.Configuration.ActiveDirectory.BaseDN, configuration.Configuration.ActiveDirectory.FilterOn, username)
	if err != nil {
		logrus.Warnln("ChangePassword service : Could not find user")
		return "", err
	}

	for i := range user.Attributes {
		if user.Attributes[i].Name == configuration.Configuration.ActiveDirectory.EmailField {
			if user.Attributes[i].Values[0] == "" {
				return "", &structures.CustomError{Text: "email not provided in active directory server", HttpCode: 500}
			}

			return user.Attributes[i].Values[0], nil
		}
	}

	return "", &structures.CustomError{Text: configuration.Configuration.ActiveDirectory.EmailField + " not found on this user", HttpCode: 500}
}
