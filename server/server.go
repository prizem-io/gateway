package server

import (
	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
	ef "github.com/prizem-io/gateway/errorfactory"
)

type (
	ProcessingHandler     func(ctx context.Context) error
	PostProcessingHandler func(ctx context.Context)
)

var (
	_config            config.Configuration
	processingHandlers = []ProcessingHandler{}
	successHandlers    = []PostProcessingHandler{}
	errorHandlers      = []PostProcessingHandler{}
)

func Initialize(config config.Configuration) {
	_config = config
}

func SetProcessingHandlers(_handlers ...ProcessingHandler) {
	processingHandlers = _handlers
}

func SetSuccessHandlers(_handlers ...PostProcessingHandler) {
	successHandlers = _handlers
}

func SetErrorHandlers(_handlers ...PostProcessingHandler) {
	errorHandlers = _handlers
}

func Serve(ctx context.Context) {
	err := invokeProcessingHandlers(ctx, processingHandlers)

	// Send error payload
	if err != nil {
		if apiErr, found := err.(*ef.APIError); found {
			ctx.Rs().SetStatusCode(apiErr.Status)
		} else {
			ctx.Rs().SetStatusCode(500)
		}
		ctx.SendEntity(err)
	}

	// Invoke post processing hanlders, if available
	if err != nil && len(errorHandlers) > 0 {
		invokePostProcessingHandlers(ctx, errorHandlers)
	} else if len(successHandlers) > 0 {
		invokePostProcessingHandlers(ctx, successHandlers)
	}
}

func invokeProcessingHandlers(ctx context.Context, handlers []ProcessingHandler) error {
	for _, handler := range handlers {
		err := handler(ctx)
		if err != nil {
			ctx.SetError(err)
			return err
		}
	}

	return nil
}

func invokePostProcessingHandlers(ctx context.Context, handlers []PostProcessingHandler) {
	for _, handler := range handlers {
		handler(ctx)
	}
}
