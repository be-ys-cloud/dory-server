package authentication


func DeleteKey(username string) {
	Mutex.Lock()
	delete(Keys, username)
	Mutex.Unlock()
}