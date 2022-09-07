package structures

type Configuration struct {
	ActiveDirectory struct {
		Admin struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"admin"`
		BaseDN   string `json:"base_dn"`
		FilterOn string `json:"filter_on"`
		Address string `json:"address"`
		Port int `json:"port"`
		SkipTLSVerify bool `json:"skip_tls_verify"`
		EmailField string `json:"email_field"`
	} `json:"active_directory"`
	Server struct {
		Port int `json:"port"`
		BasePath string `json:"base_path"`
	}
	MailServer struct {
		Address string `json:"address"`
		Port int `json:"port"`
		Password string `json:"password"`
		SenderAddress string `json:"sender_address"`
		SenderName string `json:"sender_name"`
		Subject string `json:"subject"`
	} `json:"mail_server"`
	FrontAddress string `json:"front_address"`
}
