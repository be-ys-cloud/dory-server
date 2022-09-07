package authentication

import (
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/thanhpk/randstr"
	"time"
)

func CreateKey(username string) string {
	sk := randstr.String(64)

	Mutex.Lock()
	Keys[username] = structures.Key{
		IssuedAd:  time.Now(),
		Username:  username,
		SecretKey: sk,
	}
	Mutex.Unlock()

	return sk
}
