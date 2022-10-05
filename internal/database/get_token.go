package database

import (
	"database/sql"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/structures"
)

func GetToken(userId string) (string, error) {
	result, err := configuration.Database.Query("SELECT key FROM totp WHERE id=? LIMIT 1;", userId)
	if err != nil {
		return "", &structures.CustomError{Text: "could not execute query", HttpCode: 500}
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(result)

	var token string
	for result.Next() {
		result.Scan(&token)
	}

	return token, nil

}
