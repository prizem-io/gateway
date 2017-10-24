package jwt

import (
	"strings"

	"github.com/prizem-io/gateway/context"
	ef "github.com/prizem-io/gateway/errorfactory"
	"github.com/prizem-io/gateway/oauth2"
)

var bearerPrefix = "Bearer "

func parseBearer(auth string) (string, bool) {
	if strings.HasPrefix(auth, bearerPrefix) {
		return string(auth[len(bearerPrefix):]), true
	}

	return "", false
}

func getOAuthCredential(ctx context.Context, id string) (*oauth2.OAuth2Credential, error) {
	_credential, err := ctx.GetCredential(id)
	if err != nil {
		return nil, err
	}

	credential, ok := _credential.(*oauth2.OAuth2Credential)
	if !ok {
		return nil, ef.New(ctx, "invalidCredential")
	}

	return credential, nil
}
