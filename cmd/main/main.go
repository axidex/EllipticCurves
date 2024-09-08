package main

import (
	"fmt"
	"github.com/axidex/Unknown/pkg/logger"
	"github.com/axidex/elliptic/config"
	"github.com/axidex/elliptic/internal/api"
	"github.com/joho/godotenv"
)

func main() {
	appConfig, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("Got error when reading config from file - %s", err)
		return
	}
	fmt.Printf("Config: %+v\n", appConfig)

	appLogger, err := logger.CreateNewZapLogger(appConfig.Logger)
	if err != nil {
		fmt.Printf("Got error when initializing logger - %s", err)
		return
	}

	err = godotenv.Load()

	// App
	app := api.CreateApp(appConfig, appLogger)
	engine := app.InitRoutes()
	for _, item := range engine.Routes() {
		appLogger.Info("method:", item.Method, "\tpath:", item.Path)
	}
	err = engine.Run(fmt.Sprintf(":%d", appConfig.Server.Port))
	if err != nil {
		appLogger.Fatal("Failed to start server - %s", err)
		return
	}

}
