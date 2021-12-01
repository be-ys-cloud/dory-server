package ad

import (
	"ad/helpers"
	"errors"
	"github.com/sirupsen/logrus"
	"structures"
)

func GetUserEmail(username string) (string, structures.Error) {

	l, err := helpers.GetSession(Conf.ActiveDirectory.Address, Conf.ActiveDirectory.Port, Conf.ActiveDirectory.SkipTLSVerify)
	if err.Error != nil {
		logrus.Warnln("ChangePassword service : Could not connect to server")
		return "", err
	}

	defer l.Close()


	//Connect to Active Directory as user
	err = helpers.BindUser(l, Conf.ActiveDirectory.Admin.Username, Conf.ActiveDirectory.Admin.Password)
	if err.Error != nil {
		logrus.Warnln("ChangePassword service : Could not login to Active Directory : Bad AD Password supplied")
		return "", err
	}

	user, err := helpers.GetUser(l, Conf.ActiveDirectory.BaseDN, Conf.ActiveDirectory.FilterOn, username)
	if err.Error != nil {
		logrus.Warnln("ChangePassword service : Could not find user")
		return "", err
	}

	for i := range user.Attributes {
		if user.Attributes[i].Name == Conf.ActiveDirectory.EmailField {
			if user.Attributes[i].Values[0] == "" {
				return "", structures.Error{Error: errors.New("email not provided in active directory server"), HttpCode: 500}
			}

			return user.Attributes[i].Values[0], structures.Error{}
		}
	}

	return "", structures.Error{Error: errors.New(Conf.ActiveDirectory.EmailField+ " not found on this user"), HttpCode: 500}
}