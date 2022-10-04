package structures

type UserAsk struct {
	Username string `json:"username"`
}

type UserUnlock struct {
	Username       string         `json:"username"`
	Authentication Authentication `json:"authentication"`
}

type UserReinitialize struct {
	Username       string         `json:"username"`
	NewPassword    string         `json:"new_password"`
	Authentication Authentication `json:"authentication"`
}

type Authentication struct {
	Token string `json:"token"`
	TOTP  string `json:"totp"`
}

type UserChangePassword struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserCreateTOTP struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserVerifyTOTP struct {
	Username string `json:"username"`
	TOTP     string `json:"totp"`
}

type TOTPToken struct {
	TOTP string `json:"TOTP"`
}
