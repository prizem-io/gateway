package backend

import (
	"fmt"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
	ef "github.com/prizem-io/gateway/errorfactory"
)

type (
	Handler interface {
		Name() string
		Handle(context.Context) error
	}
)

var (
	backends        = map[string]Handler{}
	defaultUpstream = "http"
)

func Register(handlers ...Handler) {
	lookup := make(map[string]Handler, len(handlers))
	for _, router := range handlers {
		lookup[router.Name()] = router
	}
	backends = lookup
}

func GetConfig(name string, properties map[string]interface{}) (interface{}, error) {
	backend, ok := backends[name]
	if !ok {
		return nil, fmt.Errorf("Could not find upstream: %s", name)
	}

	configurable, ok := backend.(config.Configurable)
	if !ok {
		return nil, nil
	}

	return configurable.DecodeConfig(properties)
}

func GetHandler(ctx context.Context, name *string) (Handler, error) {
	var nameStr string
	if name != nil {
		nameStr = *name
	} else {
		nameStr = defaultUpstream
	}

	backend, ok := backends[nameStr]
	if !ok {
		return nil, ef.NewError(ctx, "routerUnrecognized", ef.Params{
			"name": nameStr,
		})
	}

	return backend, nil
}
