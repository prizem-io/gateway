package fasthttp

import (
	"bytes"
	"sync/atomic"
	"unsafe"

	log "github.com/Sirupsen/logrus"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/prizem-io/gateway/config"
	ef "github.com/prizem-io/gateway/errorfactory"
	"github.com/prizem-io/gateway/server"
)

type (
	fastHttpRouter struct {
		router  *fasthttprouter.Router
		gateway *server.Gateway
	}

	operationRoute struct {
		Gateway   *server.Gateway
		Service   *config.Service
		Operation *config.Operation
		Path      string
	}
)

var (
	// Pointer to a func(ctx *fasthttp.RequestCtx)
	fastHttpRouterHandler unsafe.Pointer
)

func LoadGatewayRouter() error {
	gateway, err := server.LoadGateway(server.GatewayConfigLocation)
	if err != nil {
		return err
	}

	LoadRouter(gateway)

	return nil
}

func LoadRouter(gateway *server.Gateway) {
	router := fasthttprouter.New()
	BuildFastHttpRouter(router, gateway)

	pr := &fastHttpRouter{router: router, gateway: gateway}
	for _, callback := range server.BuildRouterCallbacks {
		callback(pr)
	}
	router.NotFound = notFound
	router.HandleMethodNotAllowed = true
	router.MethodNotAllowed = methodNotAllowed
	router.PanicHandler = internalError

	f := router.Handler
	atomic.StorePointer(&fastHttpRouterHandler, unsafe.Pointer(&f))
}

func Serve(ctx *fasthttp.RequestCtx) {
	ptr := atomic.LoadPointer(&fastHttpRouterHandler)
	handler := *(*func(ctx *fasthttp.RequestCtx))(ptr)
	handler(ctx)
}

func BuildFastHttpRouter(router *fasthttprouter.Router, gateway *server.Gateway) {
	for j := range gateway.Services {
		service := &gateway.Services[j]

		for i := range service.Operations {
			operation := &service.Operations[i]
			sourceSize := len(operation.URIPattern)
			targetSize := len(operation.URIPattern)
			if service.ContextRoot != nil {
				targetSize += len(*service.ContextRoot)
			}
			if service.URIPrefix != nil {
				sourceSize += len(*service.URIPrefix)
				targetSize += len(*service.URIPrefix)
			}

			sourceBuffer := bytes.NewBuffer(make([]byte, 0, sourceSize))
			targetBuffer := bytes.NewBuffer(make([]byte, 0, targetSize))
			if service.ContextRoot != nil {
				targetBuffer.WriteString(*service.ContextRoot)
			}
			if service.URIPrefix != nil {
				sourceBuffer.WriteString(*service.URIPrefix)
				targetBuffer.WriteString(*service.URIPrefix)
			}
			sourceBuffer.WriteString(operation.URIPattern)
			targetBuffer.WriteString(operation.URIPattern)

			route := operationRoute{
				Gateway:   gateway,
				Service:   service,
				Operation: operation,
				Path:      targetBuffer.String(),
			}
			router.Handle(
				operation.Method.String(),
				sourceBuffer.String(),
				route.handleRouter)
		}
	}
}

func (o *operationRoute) handleRouter(frc *fasthttp.RequestCtx) {
	ctx := AcquireFastHttpContext(frc, "consumer")
	ctx.SetDataAccessor(o.Gateway)
	ctx.SetService(o.Service)
	ctx.SetOperation(o.Operation)
	server.Serve(ctx)
	ctx.Reset()
	ReleaseFastHttpContext(ctx)
}

func notFound(frc *fasthttp.RequestCtx) {
	ctx := AcquireFastHttpContext(frc, "consumer")
	err := ef.New(ctx, "notFound")
	ctx.SetStatusCode(err.Status)
	server.WriteEntity(ctx, err)
	ctx.Reset()
	ReleaseFastHttpContext(ctx)
}

func methodNotAllowed(frc *fasthttp.RequestCtx) {
	ctx := AcquireFastHttpContext(frc, "consumer")
	err := ef.New(ctx, "methodNotAllowed")
	ctx.SetStatusCode(err.Status)
	server.WriteEntity(ctx, err)
	ctx.Reset()
	ReleaseFastHttpContext(ctx)
}

func internalError(frc *fasthttp.RequestCtx, rcv interface{}) {
	log.Error(rcv)
	ctx := AcquireFastHttpContext(frc, "consumer")
	err := ef.New(ctx, "internalError")
	ctx.SetStatusCode(err.Status)
	server.WriteEntity(ctx, err)
	ctx.Reset()
	ReleaseFastHttpContext(ctx)
}

func (r *fastHttpRouter) GET(path string, handle server.Handle) {
	r.Handle("GET", path, handle)
}

func (r *fastHttpRouter) HEAD(path string, handle server.Handle) {
	r.Handle("HEAD", path, handle)
}

func (r *fastHttpRouter) OPTIONS(path string, handle server.Handle) {
	r.Handle("OPTIONS", path, handle)
}

func (r *fastHttpRouter) POST(path string, handle server.Handle) {
	r.Handle("POST", path, handle)
}

func (r *fastHttpRouter) PUT(path string, handle server.Handle) {
	r.Handle("PUT", path, handle)
}

func (r *fastHttpRouter) PATCH(path string, handle server.Handle) {
	r.Handle("PATCH", path, handle)
}

func (r *fastHttpRouter) DELETE(path string, handle server.Handle) {
	r.Handle("DELETE", path, handle)
}

func (r *fastHttpRouter) Handle(method, path string, handle server.Handle) {
	r.router.Handle(method, path, func(frc *fasthttp.RequestCtx) {
		ctx := AcquireFastHttpContext(frc, "consumer")
		ctx.SetDataAccessor(r.gateway)
		handle(ctx)
		ctx.Reset()
		ReleaseFastHttpContext(ctx)
	})
}
