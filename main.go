package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

const PAGE_ACCESS_TOKEN = "EAARQbeXy1gsBO6mc42B2gGTiupf01XUAzRobv7rGVPbNMZCT33yXULtrtxyZBxyIZBmZAlKUKRu2tGuOWNIFkNdedU2VbnKDPZCZCKptUtEhWN2oYezwOtbpXJZBms3XX0ZCuI1iUkPyM6YPGrdfAWUewZAbCCjidvJZCWaNz4NWAO14XxFZC2on9FI1JZB6rvWLaeIbiwZDZD"
const VERIFY_TOKEN = "mito2003"

func main() {
	router := gin.Default()

	router.GET("/webhook", func(c *gin.Context) {
		mode := c.Query("hub.mode")
		token := c.Query("hub.verify_token")
		challenge := c.Query("hub.challenge")

		if mode == "subscribe" && token == VERIFY_TOKEN {
			c.String(http.StatusOK, challenge)
		} else {
			c.String(http.StatusForbidden, "Forbidden")
		}
	})

	router.POST("/webhook", func(c *gin.Context) {
		// Handle incoming messages here
		c.String(http.StatusOK, "EVENT_RECEIVED")
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "mito hi")
	})

	//port := os.Getenv("HTTP_PLATFORM_PORT")
	port := os.Getenv("PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
}
