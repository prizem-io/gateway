package oauth2

import (
	"github.com/mitchellh/mapstructure"

	"github.com/prizem-io/gateway/config"
)

type GrantSettings struct {
	Enabled             bool                   `mapstructure:"enabled"`
	AccessTokenTimeout  *uint64                `mapstructure:"accessTokenTimeout"`
	RefreshTokenTimeout *uint64                `mapstructure:"refreshTokenTimeout"`
	Lifespan            config.Lifespan        `mapstructure:"lifespan"`
	TokenType           string                 `mapstructure:"tokenType"`
	PermissionIds       []string               `mapstructure:"permissionIds"`
	Extended            map[string]interface{} `mapstructure:"extended"`
}

type OAuth2Credential struct {
	config.Credential `mapstructure:",squash"`
	ClientID          string                   `mapstructure:"clientId"`
	ClientSecret      string                   `mapstructure:"clientSecret"`
	GrantSettings     map[string]GrantSettings `mapstructure:"grantSettings"`
	RedirectURI       *string                  `mapstructure:"redirectUri"`
	PermissionIds     []string                 `mapstructure:"permissionIds"`
	Extended          map[string]interface{}   `mapstructure:"extended"`
}

type OAuthToken struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	ExpiresIn    uint64  `json:"expires_in"`
}

type CredentialDecoder struct{}

func NewCredentialDecoder() *CredentialDecoder {
	return &CredentialDecoder{}
}

func (d *CredentialDecoder) Type() string {
	return "oauth2"
}

func (d *CredentialDecoder) DecodeCredential(input map[string]interface{}) (interface{}, string, error) {
	var credential OAuth2Credential
	err := mapstructure.Decode(input, &credential)
	if err != nil {
		return nil, "", err
	}

	return &credential, credential.ClientID, nil
}
