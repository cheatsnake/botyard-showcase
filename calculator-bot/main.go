package main

import (
	"botyard"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/maja42/goval"
)

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Bot is running...")
	log.Fatal(http.ListenAndServe(":"+"80", nil))
}

var zeroDiv = regexp.MustCompile(`\d+\s*/\s*0`)

func messageHandler(msg botyard.Message) {
	reply := ""

	switch msg.Body {
	case "/start":
		reply = "👋 Hi, I am a calculator bot and can calculate various maths expressions.\n\n"
		reply += "Try sending me something and I'll try to solve it. For example:\n\n"
		reply += "<b>2 + 2</b>\n\n"
		reply += "<b>25 * 90</b>\n\n"
		reply += "<b>(12 + 30) / 6 - 11</b>\n\n"
		reply += "<b>22.0 / 7.01</b>\n\n"
		reply += "<b>1 << 32</b>\n\n"
	default:
		if zeroDiv.MatchString(msg.Body) {
			reply = "Divison by zero is not allowed."
			break
		}

		eval := goval.NewEvaluator()
		result, err := eval.Evaluate(msg.Body, nil, nil)
		if err != nil {
			reply = "Sorry, but I don't understand you."
			break
		}

		reply = fmt.Sprintf("%v", result)
	}

	botyard.SendMessage(msg.ChatId, reply, nil, os.Getenv("CALCULATOR_BOT_KEY"))
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
