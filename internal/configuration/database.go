package configuration

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func generateDatabase() {
	db := getDatabaseConnection()
	//Regenerating database
	_, _ = db.Exec("CREATE TABLE totp(`id` VARCHAR(512) PRIMARY KEY, `key` VARCHAR(512) NOT NULL);")
}

func getDatabaseConnection() *sql.DB {
	db, err := sql.Open("sqlite3", Configuration.Server.DatabasePath)
	if err != nil {
		logrus.Error("Unable to connect to database !")
		logrus.Fatal(err.Error())
	}
	return db
}
