package structures

type Configuration struct {
	LDAPServer struct {
		Admin struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"admin"`
		BaseDN        string `json:"base_dn"`
		FilterOn      string `json:"filter_on"`
		Address       string `json:"address"`
		Port          int    `json:"port"`
		Kind          string `json:"kind"`
		SkipTLSVerify bool   `json:"skip_tls_verify"`
		EmailField    string `json:"email_field"`
	} `json:"ldap_server"`
	Server struct {
		Port         int    `json:"port"`
		BasePath     string `json:"base_path"`
		DatabasePath string `json:"database_path"`
	} `json:"server"`
	TOTP struct {
		CustomServiceName string `json:"custom_service_name"`
		Secret            string `json:"secret"`
	} `json:"totp"`
	MailServer struct {
		Address       string `json:"address"`
		Port          int    `json:"port"`
		Password      string `json:"password"`
		SenderAddress string `json:"sender_address"`
		SenderName    string `json:"sender_name"`
		Subject       string `json:"subject"`
		SkipTLSVerify bool   `json:"skip_tls_verify"`
	} `json:"mail_server"`
	FrontAddress string   `json:"front_address"`
	Features     Features `json:"features"`
}

type Features struct {
	DisableUnlock                   bool `json:"disable_unlock"`
	DisablePasswordUpdate           bool `json:"disable_password_update"`
	DisablePasswordReinitialization bool `json:"disable_password_reinitialization"`
	DisableTOTP                     bool `json:"disable_totp"`
}
