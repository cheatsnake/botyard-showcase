package main

import (
	"botyard"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	errorMsg      = "Sorry, something went wrong, try again later..."
	unknownCmdMsg = "Sorry, but I don't understand you."
)

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Bot is running...")
	log.Fatal(http.ListenAndServe(":"+"80", nil))
}

func messageHandler(msg botyard.Message) {
	reply := ""
	attachmentIds := make([]string, 0, 1)

	switch msg.Body {
	case "/start":
		reply = "👋 Welcome!\n\nI can generate random images for you."
	case "/generate_pixelart":
		reply = "🎨 Random pixel image specially for you:"
		buf, err := newPixelart(512, 512, 32)
		if err != nil {
			fmt.Println(err)
			reply = errorMsg
			break
		}

		attach, err := botyard.UploadFile(buf, "pixelart.png", os.Getenv("IMAGE_BOT_KEY"))
		if err != nil {
			fmt.Println(err)
			reply = errorMsg
			break
		}

		attachmentIds = append(attachmentIds, attach[0].Id)
	case "/generate_gradient":
		reply = "🌈 Random gradient image specially for you:"
		buf, err := generateRandomGradient(800, 450)
		if err != nil {
			fmt.Println(err)
			reply = errorMsg
			break
		}

		attach, err := botyard.UploadFile(buf, "gradient.png", os.Getenv("IMAGE_BOT_KEY"))
		if err != nil {
			fmt.Println(err)
			reply = errorMsg
			break
		}

		attachmentIds = append(attachmentIds, attach[0].Id)
	case "/generate_qrcode":
		reply = "📷 Random QRcode image specially for you:"
		buf, err := generateRandomQrcode(512)
		if err != nil {
			fmt.Println(err)
			reply = errorMsg
			break
		}

		attach, err := botyard.UploadFile(buf, "qrcode.png", os.Getenv("IMAGE_BOT_KEY"))
		if err != nil {
			fmt.Println(err)
			reply = errorMsg
			break
		}

		attachmentIds = append(attachmentIds, attach[0].Id)
	default:
		reply = unknownCmdMsg
	}

	botyard.SendMessage(msg.ChatId, reply, attachmentIds, os.Getenv("IMAGE_BOT_KEY"))
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
