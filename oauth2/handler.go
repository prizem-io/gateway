package oauth2

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
)

type (
	Tokener interface {
		Create(*config.Token) (*config.Token, error)
	}
)

var (
	_tokener Tokener
	hmacKey  []byte

	grantHandlers = map[string]func(context.Context, *OAuth2Credential, []string){
		"client_credentials": handleClientCredentials,
	}

	basicAuthPrefix = "Basic "
)

func Initialize(config config.Configuration, tokener Tokener) error {
	hmacKeyFile, err := config.GetString("filter.jwt.secretFile")
	if err != nil {
		return err
	}

	hmacKey, err = ioutil.ReadFile(hmacKeyFile)
	if err != nil {
		return err
	}

	_tokener = tokener

	return nil
}

func GrantHandler(ctx context.Context) {
	rq := ctx.Rq()
	grantType := rq.FormValue("grant_type")
	handler, ok := grantHandlers[grantType]

	if !ok {
		errorResponse(ctx, http.StatusBadRequest, "unsupported_grant_type")
		return
	}

	auth := rq.Header("Authorization")
	clientID, clientSecret, ok := parseBasicAuth(auth)

	if !ok {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized_client")
		return
	}

	cred, err := ctx.FindCredential("oauth2", clientID)
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized_client")
		return
	}

	c, ok := cred.(*OAuth2Credential)
	if ok && c.Enabled && c.SubjectType == "consumer" &&
		c.ClientSecret == clientSecret {
		var scopeArray = strings.Split(rq.FormValue("scope"), " ")
		handler(ctx, c, scopeArray)
	} else {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized_client")
	}
}

func createTokenResponse(ctx context.Context, c *OAuth2Credential, g *GrantSettings, t *config.Token, allowRefreshTokens bool) {
	if g.Lifespan == config.LifespanFinite &&
		g.AccessTokenTimeout != nil &&
		"jwt" == t.TokenType {

		claims := jwt.MapClaims{}
		for key, value := range t.Claims {
			claims[key] = value
		}

		now := time.Now()
		claims["jti"] = uuid.NewV4().String()
		claims["iat"] = now.Unix()
		claims["exp"] = now.Add(time.Second * time.Duration(*g.AccessTokenTimeout)).Unix()
		claims["cid"] = c.ID

		// Create the token
		method := jwt.SigningMethodHS256
		token := &jwt.Token{
			Header: map[string]interface{}{
				"typ": "JWT",
				"alg": method.Alg(),
				"cid": c.ID,
			},
			Claims: claims,
			Method: method,
		}

		// Sign and get the complete encoded token as a string
		tokenString, _ := token.SignedString(hmacKey)
		t.ID = tokenString

		tokenResponse(ctx, c, g, t, allowRefreshTokens)
	} else {
		token, err := _tokener.Create(t)
		if err != nil {
			log.Error("Could not create token: " + err.Error())
			errorResponse(ctx, http.StatusInternalServerError, "An internal error occured")
			return
		}

		tokenResponse(ctx, c, g, token, allowRefreshTokens)
	}
}

func tokenResponse(ctx context.Context, c *OAuth2Credential, g *GrantSettings, t *config.Token, allowRefreshTokens bool) {
	tokenResponse := OAuthToken{
		AccessToken: t.ID,
		TokenType:   t.TokenType,
		ExpiresIn:   uint64(t.Expiry),
	}

	// Create a seperate refresh token, if allowed
	if allowRefreshTokens &&
		g.Lifespan != config.LifespanSession &&
		g.RefreshTokenTimeout != nil {
		refreshToken := config.Token{
			CredentialID:  c.ID,
			TokenType:     "refresh",
			GrantType:     t.GrantType,
			Issuer:        t.Issuer,
			Subject:       t.Subject,
			Audience:      t.Audience,
			IssuedAt:      time.Now(),
			Expiry:        int64(*g.RefreshTokenTimeout),
			Lifespan:      config.LifespanFinite,
			FromToken:     *&t.ID,
			PermissionIds: t.PermissionIds,
			Claims:        t.Claims,
			State:         t.State,
			URI:           t.URI,
			Extended:      t.Extended,
			ExternalID:    t.ExternalID,
		}

		created, err := _tokener.Create(&refreshToken)
		if err != nil {
			log.Error("Could not create token: " + err.Error())
			errorResponse(ctx, http.StatusInternalServerError, "An internal error occured")
			return
		}
		tokenResponse.RefreshToken = &created.ID
	}

	rs := ctx.Rs()
	jsonValue, err := json.Marshal(tokenResponse)
	if err != nil {
		log.Error("Could not marshal token: " + err.Error())
		errorResponse(ctx, http.StatusInternalServerError, "An internal error occured")
		return
	}

	rs.SetHeader("Content-Type", "application/json")
	rs.SetBody(jsonValue)
}

func errorResponse(ctx context.Context, status int, error string) {
	response := map[string]interface{}{
		"error": error,
	}

	rs := ctx.Rs()
	jsonValue, err := json.Marshal(response)
	if err != nil {
		log.Error("Could not marshal error: " + err.Error())
	} else {
		rs.SetStatusCode(status)
		rs.SetHeader("Content-Type", "application/json")
		rs.SetBody(jsonValue)
	}
}

func parseBasicAuth(auth string) (string, string, bool) {
	if strings.HasPrefix(auth, basicAuthPrefix) {
		payload, err := base64.StdEncoding.DecodeString(string(auth[len(basicAuthPrefix):]))
		if err == nil {
			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) == 2 {
				return pair[0], pair[1], true
			}
		}
	}

	return "", "", false
}
