package code

import (
	"code/handlers"
	"net/http"

	"github.com/gin-contrib/rollbar"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Debug       bool
	DatabaseUrl string
	Port        string
}

func NewConfig(debug bool, databaseUrl string, port string) *Config {
	return &Config{debug, databaseUrl, port}
}

func Api(config *Config) error {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := setupRouter(config)
	if setupRollbar() {
		router.Use(rollbar.Recovery(true))
	}

	return router.Run(":" + config.Port)
}

func setupRouter(config *Config) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := router.Group("api")
	links := api.Group("links")
	linksHandler := handlers.NewLinkHandler()
	linksHandler.Register(links)

	return router
}
