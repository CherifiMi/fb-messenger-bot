package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

const VERIFY_TOKEN = "mito2003"
const KEY = "EAAFbboQd3EEBOZCG3kI2aut3T4kHt4GZC3BlgpvSwm39rQTZCP3AtlZAh7SerWCBcWNZBzVOurBQudIZBHxBeQvZA6AmKnMbogClbQerkZBqjArW4zdxH7xgJfOAlu2KImxsWAtbsmZCQ0sW7sqXEmUTFB5G4lbSScuJvoThZAk24RxoDKjZAIA4Hkb0zKUe7Dz6cGrpAZDZD"

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		mode := c.Query("hub.mode")
		token := c.Query("hub.verify_token")
		challenge := c.Query("hub.challenge")

		if mode == "subscribe" && token == VERIFY_TOKEN {
			c.String(http.StatusOK, challenge)
		} else {
			c.String(http.StatusForbidden, "Forbidden")
		}
	})

	router.POST("/", func(c *gin.Context) {
		// Handle incoming messages here
		c.String(http.StatusOK, "EVENT_RECEIVED")
		log.Println("EVENT_RECEIVED")
	})

	router.GET("/mito", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mito": "nyx"})
		log.Println("MITOOOOOOOOOOOOOOOOOOOOOOOOOOO LOGS")
	})

	port := os.Getenv("HTTP_PLATFORM_PORT")
	//port := "8080"

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
}
