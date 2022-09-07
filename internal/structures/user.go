package structures

type User struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Token       string `json:"token"`
}

// Structures for swaggo only

type UserAsk struct {
	Username string `json:"username"`
}

type UserUnlock struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserReinitialize struct {
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
	Token       string `json:"token"`
}

type UserChangePassword struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
