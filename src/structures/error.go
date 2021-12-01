package structures


type Error struct {
	Error error `json:"error"`
	HttpCode int `json:"http_code"`
}
