package authentication

import "time"

func CheckKey(username string, sk string) bool {
	Mutex.Lock()
	val, ok := Keys[username]
	Mutex.Unlock()

	if !ok {
		return false
	}

	if val.IssuedAd.Add(10 * time.Minute).Before(time.Now()) {
		DeleteKey(username)
		return false
	}

	if sk != val.SecretKey {
		return false
	}

	return true


}
