package nats

import (
	"fmt"

	"github.com/nats-io/nats"

	"github.com/prizem-io/gateway"
)

type natsConfig struct {
	Url string `mapstructure:"url"`
}

func Connect(configuration prizem.Configuration) (*nats.Conn, error) {
	var config natsConfig
	configuration.UnmarshalKey("nats", &config)

	nc, err := nats.Connect(config.Url)
	if err != nil {
		return nil, fmt.Errorf("Can't connect: %v", err)
	}

	return nc, nil
}
