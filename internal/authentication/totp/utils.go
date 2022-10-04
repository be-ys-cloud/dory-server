package totp

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
)

func encodeUser(userDN string) string {
	hmac512 := hmac.New(sha512.New, []byte(configuration.Configuration.TOTP.Secret))
	hmac512.Write([]byte(userDN))
	secret := base64.StdEncoding.EncodeToString(hmac512.Sum(nil))

	return secret
}
