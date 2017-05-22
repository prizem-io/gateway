package oauth2

import (
	"net/http"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
)

func handleClientCredentials(ctx context.Context, c *OAuth2Credential, scope []string) {
	grantSettings, ok := c.GrantSettings["client_credentials"]

	if !ok || !grantSettings.Enabled {
		errorResponse(ctx, http.StatusUnauthorized, "unauthorized_client")
		return
	}

	token := config.Token{
		CredentialID:  c.ID,
		TokenType:     grantSettings.TokenType,
		GrantType:     "client_credentials",
		Expiry:        int64(*grantSettings.AccessTokenTimeout),
		Lifespan:      grantSettings.Lifespan,
		PermissionIds: grantSettings.PermissionIds,
	}

	createTokenResponse(ctx, c, &grantSettings, &token, false)
}
