package api

import (
	"github.com/axidex/Unknown/pkg/logger"
	"github.com/axidex/elliptic/config"
	_ "github.com/axidex/elliptic/docs"
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type App struct {
	config *config.Config
	logger logger.Logger
}

func CreateApp(config *config.Config, logger logger.Logger) *App {
	return &App{
		config: config,
		logger: logger,
	}
}

func (app *App) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithFormatter(app.loggerMiddleware))
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger/*any", app.swagger)

	api := router.Group("/api")
	{

		health := api.Group("/health")
		{
			health.GET("/ping", app.health)
		}

		cyphers := api.Group("/cypher")
		{
			elliptic := cyphers.Group("/elliptic")
			{
				elliptic.POST("/encrypt", app.encrypt)
				elliptic.GET("/keys", app.generateKey)
				elliptic.POST("/decrypt", app.decrypt)
			}
		}
	}

	router.Use(ginzerolog.Logger("gin"))

	return router
}

func (app *App) swagger(c *gin.Context) {
	switch c.Param("any") {
	case "/", "docs":
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	default:
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	}
}
