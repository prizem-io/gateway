package command

import (
	log "github.com/Sirupsen/logrus"
)

type Params map[string]interface{}
type Listener func(Params)

var (
	listenersByCommandName = map[string][]Listener{}
)

func AddListener(command string, listener Listener) {
	listeners := listenersByCommandName[command]
	if listeners == nil {
		listeners = make([]Listener, 0, 10)
	}

	listenersByCommandName[command] = append(listeners, listener)
}

func Notify(command string, payload Params) {
	listeners := listenersByCommandName[command]
	if listeners != nil {
		for _, listener := range listeners {
			listener(payload)
		}
	} else {
		log.WithFields(log.Fields{
			"command": command,
		}).Warn("Received unknown command")
	}
}
