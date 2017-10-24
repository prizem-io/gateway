package logger

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
)

type (
	Logger struct {
		PrioritySetting int `mapstructure:"priority"`
	}
)

func New() *Logger {
	return &Logger{}
}

func (*Logger) Name() string {
	return "logger"
}

func (t *Logger) Priority() int {
	return t.PrioritySetting
}

func (t *Logger) Initialize(config config.Configuration) error {
	return config.UnmarshalKey("logger", t)
}

func (*Logger) Evaluate(ctx context.Context, _ interface{}) (err error) {
	defer func(begin time.Time) {
		if err == nil {
			log.Infof(
				"%s::%s took %dms",
				ctx.Service().Name,
				ctx.Operation().Name,
				time.Since(begin).Nanoseconds()/int64(time.Millisecond))
		} else {
			log.Errorf(
				"Error occurred inside %s::%s - %v",
				ctx.Service().Name,
				ctx.Operation().Name,
				err)
		}
	}(time.Now())
	err = ctx.Next()
	return
}
