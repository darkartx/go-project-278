package code

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router
}

func Api() error {
	if err := godotenv.Load(); err != nil {
		fmt.Printf(".env loading error: %s", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := setupRouter()

	return router.Run(":" + port)
}
