package main

import (
	"database/sql"
	"net/http"
	"reflect"
	"strings"

	"github.com/darkartx/go-project-278/handlers"
	"github.com/go-playground/validator/v10"

	db "github.com/darkartx/go-project-278/db/generated"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/rollbar"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Debug       bool
	DatabaseUrl string
	Bind        string
}

func NewConfig(debug bool, databaseUrl string, bind string) *Config {
	return &Config{debug, databaseUrl, bind}
}

func Api(config *Config) error {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	setupValidator()

	database, err := setupDB(config)
	if err != nil {
		return err
	}

	defer func() {
		_ = database.Close()
	}()

	queries := db.New(database)
	router := setupRouter(queries, config)
	router.TrustedPlatform = gin.PlatformCloudflare

	if setupRollbar() {
		router.Use(rollbar.Recovery(true))
	}

	return router.Run(config.Bind)
}

func setupRouter(queries *db.Queries, config *Config) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:5173"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	router.Use(cors.New(corsConfig))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := router.Group("api")
	links := api.Group("links")
	linksHandler := handlers.NewLinkHandler(queries)
	linksHandler.Register(links)

	linkVisits := api.Group("link_visits")
	linkVisitsHandler := handlers.NewLinkVisitHandler(queries)
	linkVisitsHandler.Register(linkVisits)

	redirectHandler := handlers.NewRedirectHandler(queries)
	redirectHandler.Register(router)

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

func setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Регестрируем функцию для получения имени поля из тега JSON
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}
