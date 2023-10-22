package main

import (
	"botyard"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Bot is running...")
	log.Fatal(http.ListenAndServe(":"+"80", nil))
}

func messageHandler(msg botyard.Message) {
	reply := ""

	switch msg.Body {
	case "/start":
		reply = "👋 Hello! I like ping pong.\n\nLet's play with me using /ping command."
	case "/ping":
		reply = "🏓 P" + strings.Repeat("O", rand.Intn(10-1)+1) + "NG"
	default:
		reply = "Sorry, but I don't understand you."
	}

	botyard.SendMessage(msg.ChatId, reply, nil, os.Getenv("PING_PONG_BOT_KEY"))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var msg botyard.Message

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go messageHandler(msg)

	w.WriteHeader(http.StatusOK)
}
