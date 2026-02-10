package code

import (
	db "code/db/generated"
	"code/handlers"
	"database/sql"
	"net/http"

	"github.com/gin-contrib/rollbar"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	database, err := setupDB(config)
	if err != nil {
		return err
	}

	defer func() {
		_ = database.Close()
	}()

	queries := db.New(database)
	router := setupRouter(queries, config)

	if setupRollbar() {
		router.Use(rollbar.Recovery(true))
	}

	return router.Run(":" + config.Port)
}

func setupRouter(queries *db.Queries, config *Config) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := router.Group("api")
	links := api.Group("links")
	linksHandler := handlers.NewLinkHandler(queries)
	linksHandler.Register(links)

	return router
}

func setupDB(config *Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DatabaseUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
