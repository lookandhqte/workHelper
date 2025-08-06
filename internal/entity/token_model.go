package entity

type Token struct {
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	ServerTime    int    `json:"server_time"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	IntegrationID int    `json:"integration_id"`
}
