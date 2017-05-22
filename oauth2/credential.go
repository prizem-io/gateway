package oauth2

import (
	"github.com/mitchellh/mapstructure"

	"github.com/prizem-io/gateway/config"
)

type GrantSettings struct {
	Enabled             bool                   `enabled`
	AccessTokenTimeout  *uint64                `accessTokenTimeout`
	RefreshTokenTimeout *uint64                `refreshTokenTimeout`
	Lifespan            config.Lifespan        `lifespan`
	TokenType           string                 `tokenType`
	PermissionIds       []string               `permissionIds`
	Extended            map[string]interface{} `extended`
}

type OAuth2Credential struct {
	config.Credential `mapstructure:",squash"`
	ClientId          string                   `clientId`
	ClientSecret      string                   `clientSecret`
	GrantSettings     map[string]GrantSettings `grantSettings`
	RedirectUri       *string                  `redirectUri`
	PermissionIds     []string                 `permissionIds`
	Extended          map[string]interface{}   `extended`
}

type OAuthToken struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	ExpiresIn    uint64  `json:"expires_in"`
}

type OAuth2CredentialDecoder struct{}

func NewOAuth2CredentialDecoder() *OAuth2CredentialDecoder {
	return &OAuth2CredentialDecoder{}
}

func (d *OAuth2CredentialDecoder) Type() string {
	return "oauth2"
}

func (d *OAuth2CredentialDecoder) DecodeCredential(input map[string]interface{}) (interface{}, string, error) {
	credential := OAuth2Credential{}
	err := mapstructure.Decode(input, &credential)
	if err != nil {
		return nil, "", err
	}

	return &credential, credential.ClientId, nil
}
