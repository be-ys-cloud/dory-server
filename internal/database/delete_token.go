package database

import "github.com/be-ys-cloud/dory-server/internal/configuration"

func DeleteToken(userId string) error {
	_, err := configuration.Database.Exec("DELETE FROM totp WHERE id=?;", userId)
	return err
}
