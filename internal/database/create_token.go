package database

import (
	"github.com/be-ys-cloud/dory-server/internal/configuration"
)

func CreateToken(userId string, key string) error {
	_, err := configuration.Database.Exec("INSERT INTO totp(id, key) VALUES(?, ?);", userId, key)
	return err
}
