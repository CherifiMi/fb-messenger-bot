package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := gin.Default()

	router.GET("/webhook", func(c *gin.Context) {
		verifyToken := os.Getenv("VERIFY_TOKEN")
		mode := c.Query("hub.mode")
		token := c.Query("hub.verify_token")
		challenge := c.Query("hub.challenge")

		if mode == "subscribe" && token == verifyToken {
			c.String(http.StatusOK, challenge)
		} else {
			c.String(http.StatusForbidden, "Forbidden")
		}
	})

	router.POST("/webhook", func(c *gin.Context) {
		// Handle incoming messages here
		c.String(http.StatusOK, "EVENT_RECEIVED")
	})

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "mito hi")
	})

	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)

}
