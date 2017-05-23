package timer

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/context"
)

type (
	Timer struct {
		PrioritySetting int `mapstructure:"priority"`
	}
)

func New() *Timer {
	return &Timer{}
}

func (*Timer) Name() string {
	return "test"
}

func (t *Timer) Priority() int {
	return t.PrioritySetting
}

func (t *Timer) Initialize(config config.Configuration) error {
	return config.UnmarshalKey("timer", t)
}

func (*Timer) Evaluate(ctx context.Context, _ interface{}) error {
	defer func(begin time.Time) {
		log.Infof(
			"%s::%s took %dms",
			ctx.Service().Name,
			ctx.Operation().Name,
			time.Since(begin).Nanoseconds()/int64(time.Millisecond))
	}(time.Now())
	return ctx.Next()
}
