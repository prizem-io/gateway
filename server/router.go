package server

import (
	"github.com/prizem-io/gateway/context"
)

type (
	Handle func(context.Context)

	Router interface {
		GET(path string, handle Handle)
		HEAD(path string, handle Handle)
		OPTIONS(path string, handle Handle)
		POST(path string, handle Handle)
		PUT(path string, handle Handle)
		PATCH(path string, handle Handle)
		DELETE(path string, handle Handle)
		Handle(method, path string, handle Handle)
	}

	BuildRouterCallback func(router Router)
)

var (
	GatewayConfigLocation string
	BuildRouterCallbacks  = []BuildRouterCallback{}
)

func AddBuildRouterCallbacks(callbacks ...BuildRouterCallback) {
	BuildRouterCallbacks = append(BuildRouterCallbacks, callbacks...)
}
