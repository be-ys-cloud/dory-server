package structures

import "time"

type Key struct {
	IssuedAd time.Time
	Username string
	SecretKey string
}
