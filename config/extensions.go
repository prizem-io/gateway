package config

import (
	"time"
)

type (
	Configuration interface {
		Get(key string) interface{}
		GetBool(key string) bool
		GetFloat64(key string) float64
		GetInt(key string) int
		GetString(key string) (string, error)
		GetStringMap(key string) map[string]interface{}
		GetStringMapString(key string) map[string]string
		GetStringSlice(key string) []string
		GetTime(key string) time.Time
		GetDuration(key string) time.Duration
		IsSet(key string) bool
		Unmarshal(rawVal interface{}) error
		UnmarshalKey(key string, dest interface{}) error
	}

	Initializable interface {
		Initialize(config Configuration) error
	}

	Configurable interface {
		DecodeConfig(config map[string]interface{}) (interface{}, error)
	}

	ConfigurationCombiner interface {
		Combine(configurations ...interface{}) (interface{}, error)
	}
)

func (c *PluginConfig) HandleConfig(f func(*PluginConfig) (interface{}, error)) error {
	conf, err := f(c)
	if err != nil {
		return err
	}

	c.Config = conf
	return nil
}
