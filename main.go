package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	VERIFY_TOKEN = "mito2003"
	KEY          = "EAAFbboQd3EEBOZCG3kI2aut3T4kHt4GZC3BlgpvSwm39rQTZCP3AtlZAh7SerWCBcWNZBzVOurBQudIZBHxBeQvZA6AmKnMbogClbQerkZBqjArW4zdxH7xgJfOAlu2KImxsWAtbsmZCQ0sW7sqXEmUTFB5G4lbSScuJvoThZAk24RxoDKjZAIA4Hkb0zKUe7Dz6cGrpAZDZD"
	GRAPHQL_URL  = "https://graph.facebook.com/v12.0/me/messages"
)

type Message struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   *struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
			} `json:"message,omitempty"`
		} `json:"messaging"`
	} `json:"entry"`
}

func main() {
	router := gin.Default()

	router.GET("/", verifyWebhook)
	router.POST("/", receiveMessage)
	router.GET("/mito", mitoHandler)

	port := os.Getenv("HTTP_PLATFORM_PORT")
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
}

func receiveMessage(c *gin.Context) {
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		log.Printf("failed to bind JSON: %v", err)
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}
	for _, entry := range message.Entry {
		for _, messaging := range entry.Messaging {
			if messaging.Message != nil {
				switch messaging.Message.Text {
				case "date":
					sendMessage(messaging.Sender.ID, time.Now().Format("2006-01-02"))
				case "hi":
					sendMessage(messaging.Sender.ID, "Hello")
				default:
					sendMessage(messaging.Sender.ID, messaging.Message.Text)
				}
			}
		}
	}

	c.String(http.StatusOK, "EVENT_RECEIVED")
	log.Println("EVENT_RECEIVED")
}

func verifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == VERIFY_TOKEN {
		c.String(http.StatusOK, challenge)
	} else {
		c.String(http.StatusForbidden, "Forbidden")
	}
}

func mitoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"mito": "nyx"})
	log.Println("MITOOOOOOOOOOOOOOOOOOOOOOOOOOO LOGS")
}

func sendMessage(recipientID, text string) {

	suggestions := map[string]string{
		"What is the date?": "date",
		"Say Hi":            "hi",
		"MITO":              "mito",
	}

	buttons := make([]interface{}, len(suggestions))
	for title, payload := range suggestions {
		button := map[string]interface{}{
			"content_type": "text",
			"title":        title,
			"payload":      payload,
		}
		buttons = append(buttons, button)
	}

	//suggestions := []string{"What is the date?", "Hi", "MITO"}
	//buttons := make([]interface{}, len(suggestions))
	//for i, suggestion := range suggestions {
	//	buttons[i] = map[string]interface{}{
	//		"content_type": "text",
	//		"title":        suggestion,
	//		"payload":      suggestion,
	//	}
	//}

	message := struct {
		Recipient struct {
			ID string `json:"id"`
		} `json:"recipient"`
		Message struct {
			Text         string        `json:"text"`
			QuickReplies []interface{} `json:"quick_replies"`
		} `json:"message"`
	}{
		Recipient: struct {
			ID string `json:"id"`
		}{
			ID: recipientID,
		},
		Message: struct {
			Text         string        `json:"text"`
			QuickReplies []interface{} `json:"quick_replies"`
		}{
			Text:         text,
			QuickReplies: buttons,
		},
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, GRAPHQL_URL, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = "access_token=" + KEY

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("failed to send message, response code: %d, response body: %s", resp.StatusCode, body)
	} else {
		log.Println("message sent successfully")
	}
}
