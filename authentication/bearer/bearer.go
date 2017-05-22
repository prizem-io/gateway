package bearer

import (
	"strings"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
	ef "github.com/prizem-io/gateway/errorfactory"
	"github.com/prizem-io/gateway/identity"
)

type (
	Identifier func(subject string) (identity.Identity, error)

	Tokener interface {
		Get(id string) (*config.Token, error)
		Touch(token *config.Token) error
	}

	BearerAuthenticator struct{}
)

var (
	_identifier Identifier
	_tokener    Tokener
)

func Initialize(identifier Identifier, tokener Tokener) {
	_identifier = identifier
	_tokener = tokener
}

func New() *BearerAuthenticator {
	return &BearerAuthenticator{}
}

func (a *BearerAuthenticator) Name() string {
	return "bearer"
}

func (a *BearerAuthenticator) Initialize(config config.Configuration) error {
	return nil
}

func (a *BearerAuthenticator) Authenticate(ctx context.Context, confuration interface{}) (*config.Credential, identity.Identity, error) {
	rq := ctx.Rq()
	auth := rq.Header("Authorization")
	bearer, ok := parseBearer(auth)
	// Bearer token not passed in Authorization header
	if !ok {
		return nil, nil, nil
	}
	// Not a bearer token
	if strings.IndexRune(bearer, '.') != -1 {
		return nil, nil, nil
	}

	token, err := _tokener.Get(bearer)
	if err != nil {
		return nil, nil, ef.NewError(ctx, "invalidToken")
	}

	identity, err := _identifier(token.Subject)
	if err != nil {
		return nil, nil, err
	}

	credential, err := getOAuthCredential(ctx, token.CredentialID)
	if err != nil {
		return nil, nil, err
	}

	if token.Lifespan == config.LifespanSession {
		_tokener.Touch(token)
	}

	return &credential.Credential, identity, err
}
