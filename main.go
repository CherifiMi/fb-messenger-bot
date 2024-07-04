/*package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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
		log.Println("EVENT_RECEIVED")
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mito": "nyx"})
		log.Println("MITOOOOOOOOOOOOOOOOOOOOOOOOOOO LOGS")
	})

	//port := os.Getenv("HTTP_PLATFORM_PORT")
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
}
*/
///
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (

	// VERIFY_TOKEN use to verify the incoming request
	VERIFY_TOKEN = "mito2003"
	// ACCESS_TOKEN use to access Messenger API
	ACCESS_TOKEN = "EAARQbeXy1gsBO6mc42B2gGTiupf01XUAzRobv7rGVPbNMZCT33yXULtrtxyZBxyIZBmZAlKUKRu2tGuOWNIFkNdedU2VbnKDPZCZCKptUtEhWN2oYezwOtbpXJZBms3XX0ZCuI1iUkPyM6YPGrdfAWUewZAbCCjidvJZCWaNz4NWAO14XxFZC2on9FI1JZB6rvWLaeIbiwZDZD"
	// GRAPHQL_URL is a base URL v12.0 for Messenger API
	GRAPHQL_URL = "https://graph.facebook.com/v12.0"
)

// Message data structure for message event
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
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

// SendMessage data structure for send message
type SendMessage struct {
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

// webhook is a handler for Webhook server
func webhook(w http.ResponseWriter, r *http.Request) {
	// return all with status code 200
	w.WriteHeader(http.StatusOK)

	// method that allowed are GET & POST
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		log.Printf("invalid method: not get or post")
		return
	}

	// if the method of request is GET
	if r.Method == http.MethodGet {
		// read token from query parameter
		verifyToken := r.URL.Query().Get("hub.verify_token")

		// verify the token included in the incoming request
		if verifyToken != VERIFY_TOKEN {
			log.Printf("invalid verification token: %s", verifyToken)
			return
		}

		// write string from challenge query parameter
		if _, err := w.Write([]byte(r.URL.Query().Get("hub.challenge"))); err != nil {
			log.Printf("failed to write response body: %v", err)
		}

		return
	}

	// ready body in the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		return
	}

	// initiate Message data structure to message variable
	// unmarshal []byte data into message
	var message Message
	if err := json.Unmarshal(body, &message); err != nil {
		log.Printf("failed to unmarshal body: %v", err)
		return
	}

	// send message to end-user
	err = sendMessage(message.Entry[0].Messaging[0].Sender.ID, "Automatically Reply 🙌🏻")
	if err != nil {
		log.Printf("failed to send message: %v", err)
	}

	return
}

// sendMessage sends a message to end-user
func sendMessage(senderId, message string) error {
	// configure the sender ID and message
	var request SendMessage
	request.Recipient.ID = senderId
	request.Message.Text = message

	// validate empty message
	if len(message) == 0 {
		return errors.New("message can't be empty")
	}

	// marshal request data
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshall request: %w", err)
	}

	// setup http request
	url := fmt.Sprintf("%s/%s?access_token=%s", GRAPHQL_URL, "me/messages", ACCESS_TOKEN)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed wrap request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// send http request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed send request: %w", err)
	}
	defer res.Body.Close()

	// print response
	log.Printf("message sent successfully?\n%#v", res)

	return nil
}

func main() {
	// create the handler
	handler := http.NewServeMux()
	handler.HandleFunc("/", webhook)

	// configure http server
	srv := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("localhost:%d", 3000),
	}

	// start http server
	log.Printf("http server listening at %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
