package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	token := os.Getenv("HASS_TOKEN")
	server := os.Getenv("HASS_SERVER")

	if token == "" || server == "" {
		log.Fatal("Fehler: HASS_TOKEN und HASS_SERVER müssen gesetzt sein.")
	}

	u, err := url.Parse(server)
	if err != nil {
		log.Fatal(err)
	}

	wsScheme := "ws"
	if u.Scheme == "https" {
		wsScheme = "wss"
	}
	wsURL := url.URL{Scheme: wsScheme, Host: u.Host, Path: "/api/websocket"}

	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Fatal("Dial-Fehler:", err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 1. Auth
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Fatal("Read-Fehler:", err)
	}

	authReq := map[string]string{
		"type":         "auth",
		"access_token": token,
	}
	if err := c.WriteJSON(authReq); err != nil {
		log.Fatal("Auth-Write-Fehler:", err)
	}

	_, message, err = c.ReadMessage()
	var authResp map[string]interface{}
	json.Unmarshal(message, &authResp)
	if authResp["type"] != "auth_ok" {
		log.Fatal("Auth-Fehler:", string(message))
	}

	// 2. Subscribe
	subReq := map[string]interface{}{
		"id":         1,
		"type":       "subscribe_events",
		"event_type": "mobile_app_notification_action",
	}
	if err := c.WriteJSON(subReq); err != nil {
		log.Fatal("Sub-Write-Fehler:", err)
	}

	fmt.Println("Lausche auf Action-Events...")

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal("Read-Fehler:", err)
		}

		var event map[string]interface{}
		json.Unmarshal(message, &event)

		if event["type"] == "event" {
			eventData := event["event"].(map[string]interface{})["data"].(map[string]interface{})
			action := eventData["action"].(string)
			fmt.Printf("ANTWORT ERHALTEN: %s\n", action)
			os.Exit(0)
		}
	}
}
