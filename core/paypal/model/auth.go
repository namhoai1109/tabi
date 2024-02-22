package model

type ClientMetadata struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	LogoURI     string   `json:"logo_uri"`
	Scopes      []string `json:"scopes"`
	UIType      string   `json:"ui_type"`
}

type AccessTokenResponse struct {
	Scope                 string         `json:"scope"`
	AccessToken           string         `json:"access_token"`
	TokenType             string         `json:"token_type"`
	AppID                 string         `json:"app_id"`
	ExpiresIn             int            `json:"expires_in"`
	SupportedAuthnSchemes []string       `json:"supported_authn_schemes"`
	Nonce                 string         `json:"nonce"`
	ClientMetadata        ClientMetadata `json:"client_metadata"`
}
