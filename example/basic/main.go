package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"

	"github.com/prizem-io/gateway/authentication"
	"github.com/prizem-io/gateway/authentication/bearer"
	"github.com/prizem-io/gateway/authentication/jwt"
	"github.com/prizem-io/gateway/authorization"
	"github.com/prizem-io/gateway/backend"
	"github.com/prizem-io/gateway/backend/http"
	"github.com/prizem-io/gateway/command"
	"github.com/prizem-io/gateway/connect/redis"
	ef "github.com/prizem-io/gateway/errorfactory"
	"github.com/prizem-io/gateway/filter"
	"github.com/prizem-io/gateway/filter/timer"
	"github.com/prizem-io/gateway/identity/simple"
	"github.com/prizem-io/gateway/oauth2"
	"github.com/prizem-io/gateway/server"
	fasthttpserver "github.com/prizem-io/gateway/server/fasthttp"
	"github.com/prizem-io/gateway/utils"
)

const (
	envEnvironment = "PRIZEM_ENV"
	envWorkingDir  = "PRIZEM_WD"
)

var (
	environment = os.Getenv(envEnvironment)
)

func main() {
	if wd := os.Getenv(envWorkingDir); wd != "" {
		os.Chdir(wd)
	}

	var configName = "config"
	if environment != "" {
		configName += "." + environment
	}

	log.Println("Reading configuration")
	viper.SetConfigName(configName)
	viper.AddConfigPath("./etc/")
	viper.SetEnvPrefix("core")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error reading config file: %s", err))
	}

	ef.Initialize("etc/errors")
	configuration := &utils.ViperConfiguration{}

	redisClient := redis.Connect(configuration)
	tokener := redis.NewTokener(redisClient)

	server.Initialize(configuration)
	filter.Initialize(configuration)
	authentication.Initialize(configuration)
	oauth2.Initialize(configuration, tokener)
	bearer.Initialize(simple.New, tokener)
	jwt.Initialize(simple.New)

	server.AddCredentialDecoders(
		oauth2.NewOAuth2CredentialDecoder(),
	)

	server.AddConfigDecoders(
		authentication.DecodeConfig,
	)

	authentication.SetAuthenticators(
		jwt.New(),
		bearer.New(),
	)

	filter.Register(
		timer.New(),
	)

	backend.Register(
		http.New(),
	)

	server.SetProcessingHandlers(
		authentication.Handler,
		authorization.Handler,
		filter.Handler,
	)

	server.GatewayConfigLocation = viper.GetString("gateway.config")

	server.AddBuildRouterCallbacks(func(router server.Router) {
		router.POST("/oauth2/token", oauth2.GrantHandler)
	})

	command.AddListener("reload", func(params command.Params) {
		err := fasthttpserver.LoadGatewayRouter()
		if err != nil {
			log.Warn("Could not reload gateway configuration: " + err.Error())
		}
	})

	err = fasthttpserver.LoadGatewayRouter()
	if err != nil {
		panic(fmt.Errorf("Error processing gateway config: %s", err))
	}
	go redis.CommandSubscribe(redisClient)

	log.Fatal(fasthttp.ListenAndServe(viper.GetString("gateway.listen"), fasthttpserver.Serve))
}
