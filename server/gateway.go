package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"

	"github.com/prizem-io/gateway/backend"
	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/filter"
)

type CredentialDecoder interface {
	Type() string
	DecodeCredential(map[string]interface{}) (interface{}, string, error)
}

type ConfigDecoder func(name string, config map[string]interface{}) (interface{}, error)

type GatewayConfig struct {
	Consumers   []config.Consumer        `yaml:"consumers"`
	Credentials []map[string]interface{} `yaml:"credentials"`
	Permissions []config.Permission      `yaml:"permissions"`
	Plans       []config.Plan            `yaml:"plans"`
	Plugins     []config.Plugin          `yaml:"plugin"`
	Services    []config.Service         `yaml:"services"`
}

type Gateway struct {
	Services            []config.Service
	Consumers           map[string]*config.Consumer
	Credentials         map[string]interface{}
	CredentialsByClient map[string]interface{}
	Permissions         map[string]*config.Permission
	Plans               map[string]*config.Plan
	Plugins             map[string]*config.Plugin
}

var (
	_credentialDecoders = map[string]CredentialDecoder{}
	_configDecoders     = []ConfigDecoder{}
)

func AddCredentialDecoders(decoders ...CredentialDecoder) {
	for _, decoder := range decoders {
		_credentialDecoders[decoder.Type()] = decoder
	}
}

func AddConfigDecoders(decoders ...ConfigDecoder) {
	_configDecoders = append(_configDecoders, decoders...)
}

func (g *Gateway) GetPlugin(name string) (*config.Plugin, error) {
	plugin, ok := g.Plugins[name]
	if !ok {
		return nil, fmt.Errorf("Could not find plugin: %s", name)
	}
	return plugin, nil
}

func (g *Gateway) GetConsumer(id string) (*config.Consumer, error) {
	consumer, ok := g.Consumers[id]
	if !ok {
		return nil, fmt.Errorf("Could not find consumer: %s", id)
	}
	return consumer, nil
}

func (g *Gateway) GetCredential(id string) (interface{}, error) {
	credential, ok := g.Credentials[id]
	if !ok {
		return nil, fmt.Errorf("Could not find credential: %s", id)
	}
	return credential, nil
}

func (g *Gateway) FindCredential(credentialType, clientId string) (interface{}, error) {
	key := credentialType + "|" + clientId
	credential, ok := g.CredentialsByClient[key]
	if !ok {
		return nil, fmt.Errorf("Could not find credential: %s: %s", credentialType, clientId)
	}
	return credential, nil
}

func (g *Gateway) GetPlan(id string) (*config.Plan, error) {
	plan, ok := g.Plans[id]
	if !ok {
		return nil, fmt.Errorf("Could not find consumer: %s", id)
	}
	return plan, nil
}

func (g *Gateway) GetPermission(id string) (*config.Permission, error) {
	permission, ok := g.Permissions[id]
	if !ok {
		return nil, fmt.Errorf("Could not find consumer: %s", id)
	}
	return permission, nil
}

func LoadGateway(configLocation string) (*Gateway, error) {
	gatewayConfig := LoadGatewayConfig(configLocation)
	return ProcessGatewayConfig(gatewayConfig)
}

func LoadGatewayConfig(configLocation string) *GatewayConfig {
	var data []byte
	var err error
	gatewayConfig := GatewayConfig{}

	if strings.HasPrefix(configLocation, "http://") ||
		strings.HasPrefix(configLocation, "https://") {
		resp, err := http.Get(configLocation)
		if err == nil {
			defer resp.Body.Close()
			data, err = ioutil.ReadAll(resp.Body)
		}
	} else {
		data, err = ioutil.ReadFile(configLocation)
	}
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}

	if filepath.Ext(configLocation) == ".json" {
		err = json.Unmarshal(data, &gatewayConfig)
	} else {
		err = yaml.Unmarshal(data, &gatewayConfig)
	}
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}

	return &gatewayConfig
}

func ProcessGatewayConfig(gatewayConfig *GatewayConfig) (*Gateway, error) {
	var gateway Gateway
	operationCount := 0
	gateway.Services = gatewayConfig.Services

	for i := 0; i < len(gateway.Services); i++ {
		service := &gateway.Services[i]
		err := handleConfigurations(service.Filters)
		if err != nil {
			return nil, err
		}

		// Service-level upstream config
		backendConfig, err := backend.GetConfig(service.Backend.Name, service.Backend.Properties)
		if err != nil {
			return nil, err
		}
		service.Backend.Config = backendConfig

		for j := 0; j < len(service.Operations); j++ {
			operation := &service.Operations[j]
			err := handleConfigurations(operation.Filters)
			if err != nil {
				return nil, err
			}

			// Operation-level upstream config
			if operation.Backend != nil {
				backendConfig, err := backend.GetConfig(service.Backend.Name, service.Backend.Properties)
				if err != nil {
					return nil, err
				}
				operation.BackendConfig = backendConfig
			}
			operationCount++
		}
	}

	gateway.Consumers = make(map[string]*config.Consumer, len(gatewayConfig.Consumers))
	for _, consumer := range gatewayConfig.Consumers {
		err := handleConfigurations(consumer.Filters)
		if err != nil {
			return nil, err
		}
		gateway.Consumers[consumer.ID] = &consumer
	}

	gateway.Credentials = make(map[string]interface{}, len(gatewayConfig.Credentials))
	gateway.CredentialsByClient = make(map[string]interface{}, len(gatewayConfig.Credentials))
	for _, credMap := range gatewayConfig.Credentials {
		typeStr, ok := getCredentialField(credMap, "type")
		if !ok {
			log.Warn("Missing or invalid credential type")
			continue
		}
		decoder, ok := _credentialDecoders[typeStr]
		if !ok {
			log.WithFields(log.Fields{
				"type": typeStr,
			}).Warn("Unknown credential type")
			continue
		}
		idStr, ok := getCredentialField(credMap, "id")
		if !ok {
			log.Warn("Missing or invalid credential id")
			continue
		}
		credential, credKey, err := decoder.DecodeCredential(credMap)
		if err != nil {
			log.Warn("Error unmashalling credential", err)
		} else {
			key := typeStr + "|" + credKey
			gateway.Credentials[idStr] = credential
			gateway.CredentialsByClient[key] = credential
		}
	}

	gateway.Permissions = make(map[string]*config.Permission, len(gatewayConfig.Permissions))
	for _, permission := range gatewayConfig.Permissions {
		gateway.Permissions[permission.ID] = &permission
	}

	gateway.Plans = make(map[string]*config.Plan, len(gatewayConfig.Plans))
	for _, plan := range gatewayConfig.Plans {
		err := handleConfigurations(plan.Filters)
		if err != nil {
			return nil, err
		}
		gateway.Plans[plan.ID] = &plan
	}

	gateway.Plugins = make(map[string]*config.Plugin, len(gatewayConfig.Plugins))
	for _, plugin := range gatewayConfig.Plugins {
		for _, decoder := range _configDecoders {
			config, err := decoder(plugin.Name, plugin.Properties)
			if err != nil {
				return nil, err
			}
			if config != nil {
				plugin.Config = config
				continue
			}
		}
		if filter.HasFilter(plugin.Name) {
			err := plugin.HandleConfig(filter.GetConfig)
			if err != nil {
				return nil, err
			}
		}
		gateway.Plugins[plugin.Name] = &plugin
	}

	log.WithFields(log.Fields{
		"services":    len(gatewayConfig.Services),
		"operations":  operationCount,
		"consumers":   len(gateway.Consumers),
		"credentials": len(gateway.Credentials),
		"permissions": len(gateway.Permissions),
		"plans":       len(gateway.Plans),
		"plugins":     len(gateway.Plugins),
	}).Info("Processed gateway configuration succeeded")

	return &gateway, nil
}

func getCredentialField(cred map[string]interface{}, key string) (string, bool) {
	_value, ok := cred[key]
	if !ok {
		return "", false
	}
	value, ok := _value.(string)
	if !ok {
		return "", false
	}

	return value, true
}

func handleConfigurations(configs []config.PluginConfig) error {
	for i := 0; i < len(configs); i++ {
		err := configs[i].HandleConfig(filter.GetConfig)
		if err != nil {
			return err
		}
	}

	return nil
}
