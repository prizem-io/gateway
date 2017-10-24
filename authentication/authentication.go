package authentication

import (
	"fmt"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
	ef "github.com/prizem-io/gateway/errorfactory"
	"github.com/prizem-io/gateway/identity"
)

type (
	Authenticator interface {
		Name() string
		Authenticate(context.Context, interface{}) (*config.Credential, identity.Identity, error)
	}
)

var (
	_config          config.Configuration
	authenticatorMap = map[string]Authenticator{}
	authenticators   = []Authenticator{}
)

func Initialize(config config.Configuration) {
	_config = config
}

func HasAuthenticator(name string) bool {
	_, exists := authenticatorMap[name]
	return exists
}

func GetAuthenticatorConfig(conf *config.PluginConfig) (interface{}, error) {
	authenticator, ok := authenticatorMap[conf.Name]
	if !ok {
		return nil, fmt.Errorf("Could not find plugin: %s", conf.Name)
	}

	configurable, ok := authenticator.(config.Configurable)
	if !ok {
		return nil, nil
	}

	return configurable.DecodeConfig(conf.Properties)
}

func SetAuthenticators(_authenticators ...Authenticator) {
	authenticators = _authenticators
	authenticatorMap = map[string]Authenticator{}

	for _, authenticator := range authenticators {
		if initializable, ok := authenticator.(config.Initializable); ok {
			initializable.Initialize(_config)
		}
		authenticatorMap[authenticator.Name()] = authenticator
	}
}

func DecodeConfig(name string, conf map[string]interface{}) (interface{}, error) {
	authenticator, ok := authenticatorMap[name]
	if !ok {
		return nil, nil
	}

	configurable, ok := authenticator.(config.Configurable)
	if !ok {
		return nil, nil
	}

	return configurable.DecodeConfig(conf)
}

func Handler(ctx context.Context) error {
	for _, authenticator := range authenticators {
		name := authenticator.Name()

		var config interface{}
		configuration, err := ctx.GetPlugin(name)
		if err == nil {
			config = configuration.Config
		}

		credential, identity, err := authenticator.Authenticate(ctx, config)
		if err != nil {
			return err
		}

		// The authenticator did not find any valid credential,
		// continue on to the next
		if credential == nil {
			continue
		}

		//if ctx.SubjectType() != credential.SubjectType {
		if credential.SubjectType != "consumer" {
			return ef.New(ctx, "invalidCredential")
		}

		if !credential.Enabled {
			return ef.New(ctx, "credentialDisabled")
		}

		ctx.SetCredential(credential)
		ctx.SetIdentity(identity)

		consumer, err := ctx.GetConsumer(credential.SubjectID)
		if err != nil {
			return err
		}

		if consumer == nil {
			return ef.New(ctx, "invalidCredential")
		}

		ctx.SetConsumer(consumer)

		if consumer.PlanID != nil {
			plan, err := ctx.GetPlan(*consumer.PlanID)
			if err != nil {
				return err
			}

			ctx.SetPlan(plan)
		}

		break
	}

	authenticationType := ctx.Service().AuthenticationType

	if authenticationType != config.AuthenticationTypeNone && ctx.Consumer() == nil {
		return ef.New(ctx, "notAuthenticated")
	}

	if authenticationType == config.AuthenticationTypeThreeLegged && ctx.Identity() == nil {
		return ef.New(ctx, "notAuthenticated")
	}

	return nil
}
