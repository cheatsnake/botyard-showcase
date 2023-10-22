package main

import (
	"botyard"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type StageStore struct {
	stages map[string]Stage
	mu     sync.Mutex
}

type Stage struct {
	Body   string
	Action uint8
}

const storeSize = 64 // How many user stages can store at once
const (
	sha256HashStage = iota
	sha256VerifyStage1
	sha256VerifyStage2
	sha512HashStage
	sha512VerifyStage1
	sha512VerifyStage2
	sha1HashStage
	sha1VerifyStage1
	sha1VerifyStage2
	md5HashStage
	md5VerifyStage1
	md5VerifyStage2
)

const (
	helloMsg       = "Hello, I am a crypto bot! 🔐\n\nI am able to hash any text messages and also perform verification of them."
	startHashMsg   = "💬 Type in some phrase and I will hash it for you:"
	enterHashMsg   = "🔒 Enter your hash:"
	enterPhraseMsg = "🔑 Enter your phrase:"
	unknownCmdMsg  = "🙃 Sorry, but I don't understand you."
	validHashMsg   = "✅ Your phrase is valid!"
	invalidHashMsg = "❌ Your phrase is invalid..."
)

var stageStore = &StageStore{
	stages: make(map[string]Stage, storeSize),
}

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println("Bot is running...")
	log.Fatal(http.ListenAndServe(":"+"80", nil))
}

func messageHandler(msg botyard.Message) {
	stageStore.mu.Lock()
	checkStoreSize(msg.SenderId)

	reply := ""

	switch msg.Body {
	case "/start":
		reply = helloMsg
	case "/sha256_hash":
		stageStore.stages[msg.SenderId] = Stage{Action: sha256HashStage}
		reply = startHashMsg
	case "/sha256_verify":
		stageStore.stages[msg.SenderId] = Stage{Action: sha256VerifyStage1}
		reply = enterHashMsg
	case "/sha512_hash":
		stageStore.stages[msg.SenderId] = Stage{Action: sha512HashStage}
		reply = startHashMsg
	case "/sha512_verify":
		stageStore.stages[msg.SenderId] = Stage{Action: sha512VerifyStage1}
		reply = enterHashMsg
	case "/sha1_hash":
		stageStore.stages[msg.SenderId] = Stage{Action: sha1HashStage}
		reply = startHashMsg
	case "/sha1_verify":
		stageStore.stages[msg.SenderId] = Stage{Action: sha1VerifyStage1}
		reply = enterHashMsg
	case "/md5_hash":
		stageStore.stages[msg.SenderId] = Stage{Action: md5HashStage}
		reply = startHashMsg
	case "/md5_verify":
		stageStore.stages[msg.SenderId] = Stage{Action: md5VerifyStage1}
		reply = enterHashMsg
	default:
		stg, ok := stageStore.stages[msg.SenderId]
		if !ok {
			reply = unknownCmdMsg
			break
		}

		sha256result := sha256Handler(msg, stg)
		if sha256result != "" {
			reply = sha256result
			break
		}

		sha512result := sha512Handler(msg, stg)
		if sha512result != "" {
			reply = sha512result
			break
		}

		sha1result := sha1Handler(msg, stg)
		if sha1result != "" {
			reply = sha1result
			break
		}

		md5result := md5Handler(msg, stg)
		if md5result != "" {
			reply = md5result
			break
		}

		reply = unknownCmdMsg
	}

	stageStore.mu.Unlock()
	botyard.SendMessage(msg.ChatId, reply, nil, os.Getenv("CRYPTO_BOT_KEY"))
}

func checkStoreSize(currentUserId string) {
	if len(stageStore.stages) < storeSize {
		return
	}

	for key := range stageStore.stages {
		if key != currentUserId {
			delete(stageStore.stages, key)
		}
	}
}

func sha256Handler(msg botyard.Message, stg Stage) string {
	if stg.Action == sha256HashStage {
		delete(stageStore.stages, msg.SenderId)
		return fmt.Sprintf("%x", sha256.Sum256([]byte(msg.Body)))
	}

	if stg.Action == sha256VerifyStage1 {
		stageStore.stages[msg.SenderId] = Stage{
			Body:   msg.Body,
			Action: sha256VerifyStage2,
		}

		return enterPhraseMsg
	}

	if stg.Action == sha256VerifyStage2 {
		newHash := fmt.Sprintf("%x", sha256.Sum256([]byte(msg.Body)))
		delete(stageStore.stages, msg.SenderId)

		if stg.Body == newHash {
			return validHashMsg
		} else {
			return invalidHashMsg
		}
	}

	return ""
}

func sha512Handler(msg botyard.Message, stg Stage) string {
	if stg.Action == sha512HashStage {
		delete(stageStore.stages, msg.SenderId)
		return fmt.Sprintf("%x", sha512.Sum512([]byte(msg.Body)))
	}

	if stg.Action == sha512VerifyStage1 {
		stageStore.stages[msg.SenderId] = Stage{
			Body:   msg.Body,
			Action: sha512VerifyStage2,
		}

		return enterPhraseMsg
	}

	if stg.Action == sha512VerifyStage2 {
		newHash := fmt.Sprintf("%x", sha512.Sum512([]byte(msg.Body)))
		delete(stageStore.stages, msg.SenderId)

		if stg.Body == newHash {
			return validHashMsg
		} else {
			return invalidHashMsg
		}
	}

	return ""
}

func sha1Handler(msg botyard.Message, stg Stage) string {
	if stg.Action == sha1HashStage {
		delete(stageStore.stages, msg.SenderId)
		return fmt.Sprintf("%x", sha1.Sum([]byte(msg.Body)))
	}

	if stg.Action == sha1VerifyStage1 {
		stageStore.stages[msg.SenderId] = Stage{
			Body:   msg.Body,
			Action: sha1VerifyStage2,
		}

		return enterPhraseMsg
	}

	if stg.Action == sha1VerifyStage2 {
		newHash := fmt.Sprintf("%x", sha1.Sum([]byte(msg.Body)))
		delete(stageStore.stages, msg.SenderId)

		if stg.Body == newHash {
			return validHashMsg
		} else {
			return invalidHashMsg
		}
	}

	return ""
}

func md5Handler(msg botyard.Message, stg Stage) string {
	if stg.Action == md5HashStage {
		delete(stageStore.stages, msg.SenderId)
		return fmt.Sprintf("%x", md5.Sum([]byte(msg.Body)))
	}

	if stg.Action == md5VerifyStage1 {
		stageStore.stages[msg.SenderId] = Stage{
			Body:   msg.Body,
			Action: md5VerifyStage2,
		}

		return enterPhraseMsg
	}

	if stg.Action == md5VerifyStage2 {
		newHash := fmt.Sprintf("%x", md5.Sum([]byte(msg.Body)))
		delete(stageStore.stages, msg.SenderId)

		if stg.Body == newHash {
			return validHashMsg
		} else {
			return invalidHashMsg
		}
	}

	return ""
}
