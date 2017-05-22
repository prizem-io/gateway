package redis

import (
	"encoding/json"
	"strings"

	"github.com/go-redis/redis"

	"github.com/prizem-io/gateway/command"
)

func CommandSubscribe(redisClient *redis.Client) error {
	pubSub := redisClient.Subscribe("prizem")

	defer pubSub.Unsubscribe()

	for true {
		message, err := pubSub.ReceiveMessage()
		if err != nil {
			return err
		}

		index := strings.IndexAny(message.Payload, " \t\r\n")
		var commandName string
		var payload command.Params
		if index == -1 {
			commandName = message.Payload
		} else {
			commandName = message.Payload[0:index]
			jsonStr := strings.TrimSpace(message.Payload[index:len(message.Payload)])
			payload = command.Params{}
			json.Unmarshal([]byte(jsonStr), &payload)
		}

		command.Notify(commandName, payload)
	}

	return nil
}
