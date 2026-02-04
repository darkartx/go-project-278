package code

import (
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

	setupRollbar()
	router := setupRouter(config)

	return router.Run(":" + config.Port)
}

func setupRouter(config *Config) *gin.Engine {
	router := gin.Default()
	router.Use(rollbar.Recovery(true))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
		panic("a problem")
	})

	return router
}
