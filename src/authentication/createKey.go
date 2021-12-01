package authentication

import (
	"github.com/thanhpk/randstr"
	"structures"
	"time"
)

func CreateKey(username string) string {
	sk := randstr.String(64)

	Mutex.Lock()
	Keys[username] = structures.Key{
		IssuedAd: time.Now(),
		Username: username,
		SecretKey: sk,
	}
	Mutex.Unlock()

	return sk
}