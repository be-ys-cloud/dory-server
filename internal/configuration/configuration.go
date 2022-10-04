package configuration

import (
	"database/sql"
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Configuration structures.Configuration
var Database *sql.DB

func init() {

	logrus.Infoln("Acquiring configuration from configuration.json file...")

	file, err := os.Open("configuration.json")
	if err != nil {

		logrus.Fatal("Unable to load configuration.json file, now exiting !")
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		logrus.Fatal(err)
		logrus.Fatal("Unable to load configuration.json file because of invalid permissions, now exiting !")
	}

	err = json.Unmarshal(fileContent, &Configuration)
	if err != nil {
		logrus.Fatal("Malformed configuration.json file. Please check documentation. Program is now exiting.")
	}
	_ = file.Close()

	if Configuration.LDAPServer.Kind != "openldap" && Configuration.LDAPServer.Kind != "ad" {
		logrus.Fatal("Unsupported LDAP Server ! Please set ldap_server.kind to \"openldap\" or \"ad\".")
	}

	// If OpenLDAP, we must disable unlock.
	if Configuration.LDAPServer.Kind == "openldap" {
		Configuration.Features.DisableUnlock = true
	}

	if Configuration.Server.DatabasePath == "" {
		Configuration.Server.DatabasePath = "./database.sql"
	}

	// If TOTP is enabled, check that secret is not "" and have a decent length, and populate default name
	if !Configuration.Features.DisableTOTP {
		if Configuration.TOTP.CustomServiceName == "" {
			Configuration.TOTP.CustomServiceName = "DORY " + Configuration.LDAPServer.Address
		}
		if len(Configuration.TOTP.Secret) < 25 {
			logrus.Warnln("TOTP Secret key must be >= 25 characters! Disabling TOTP feature.")
			Configuration.Features.DisableTOTP = true
		} else {
			// If TOTP is active, we must have a database to store TOTP secrets
			generateDatabase()
			Database = getDatabaseConnection()
		}
	}
}
