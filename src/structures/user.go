package structures

type User struct {
	Username string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Token string `json:"token"`
}
