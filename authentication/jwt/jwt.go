package jwt

import (
	"io/ioutil"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
	ef "github.com/prizem-io/gateway/errorfactory"
	"github.com/prizem-io/gateway/identity"
	"github.com/prizem-io/gateway/oauth2"
)

type (
	Identifier func(subject string) (identity.Identity, error)

	JWTAuthenticator struct {
		hmacKey []byte
	}

	jwtConfig struct {
		Message string
	}
)

var (
	_identifier Identifier
)

func Initialize(identifier Identifier) {
	_identifier = identifier
}

func New() *JWTAuthenticator {
	return &JWTAuthenticator{}
}

func (a *JWTAuthenticator) Name() string {
	return "jwt"
}

func (a *JWTAuthenticator) Initialize(configuration config.Configuration) error {
	hmacKeyFile, err := configuration.GetString("filter.jwt.secretFile")
	if err != nil {
		return err
	}

	a.hmacKey, err = ioutil.ReadFile(hmacKeyFile)
	if err != nil {
		return err
	}

	return nil
}

func (a *JWTAuthenticator) DecodeConfig(input map[string]interface{}) (interface{}, error) {
	var config jwtConfig
	err := mapstructure.Decode(input, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (a *JWTAuthenticator) Authenticate(ctx context.Context, confuration interface{}) (*config.Credential, identity.Identity, error) {
	rq := ctx.Rq()
	auth := rq.Header("Authorization")
	bearer, ok := parseBearer(auth)
	// Bearer token not passed in Authorization header
	if !ok {
		return nil, nil, nil
	}
	// Not a JWT
	if strings.IndexRune(bearer, '.') == -1 {
		return nil, nil, nil
	}

	var credential *oauth2.OAuth2Credential
	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(bearer, claims, func(token *jwt.Token) (interface{}, error) {
		credentialID, ok := getClaimString(token.Header, "cid")
		if !ok {
			return a.hmacKey, nil
		}

		c, err := getOAuthCredential(ctx, credentialID)
		if err != nil {
			return nil, err
		}
		credential = c

		// TODO return credential signing info
		return a.hmacKey, nil
	})
	if token.Valid {
		var identity identity.Identity
		subject, ok := getClaimString(claims, "sub")
		if ok {
			identity, err = _identifier(subject)
			if err != nil {
				return nil, nil, err
			}
		}

		return &credential.Credential, identity, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, nil, ef.New(ctx, "tokenMalformed")
		} else if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			return nil, nil, ef.New(ctx, "tokenExpired")
		} else if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			return nil, nil, ef.New(ctx, "tokenNotYetActive")
		}
		return nil, nil, ef.New(ctx, "invalidCredential")
	} else {
		return nil, nil, err
	}
}

func getClaimString(claims map[string]interface{}, key string) (string, bool) {
	_value, ok := claims["cid"]
	if !ok {
		return "", false
	}
	return _value.(string), true
}
