package main

import (
	"botyard"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const port = "4000"
const botKey = "PASTE_BOT_KEY_HERE"

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Bot is running...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func messageHandler(msg botyard.Message) {
	checkStoreSize(msg.SenderId)

	reply := ""

	switch msg.Body {
	case "/start":
		reply = "👋 Welcome! I am YouTube bot and I can search any videos. So what do you want to watch?"
	default:
		index, err := strconv.Atoi(msg.Body)
		if err == nil && index > 0 {
			vData, ok := videoStore.data[msg.SenderId]
			if !ok {
				reply = "Let's do a search first. What do you want to watch?"
				break
			}

			if index <= len(vData) {
				reply = "<b>" + vData[index-1].Title + "</b><br></br>"
				reply += embedVideoPlayer(vData[index-1].Url)

				break
			}

			reply = "I didn't understand you, I wanted the number on the list."
			break
		}

		videos, list, err := searchVideos(msg.Body)
		if err != nil {
			reply = "💥 Search failed! Try again latter of type other query."
			break
		}

		reply = "🔎 Search results:\n\n\n"
		reply += list
		reply += "\nSend me any number on the list..."

		videoStore.mu.Lock()
		videoStore.data[msg.SenderId] = *videos
		videoStore.mu.Unlock()
	}

	botyard.SendMessage(msg.ChatId, reply, nil, botKey)
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
