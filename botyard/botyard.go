package botyard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

type Message struct {
	Id            string   `json:"id,omitempty"`
	ChatId        string   `json:"chatId"`
	SenderId      string   `json:"senderId,omitempty"`
	Body          string   `json:"body"`
	AttachmentIds []string `json:"attachmentIds,omitempty"`
	Timestamp     int64    `json:"timestamp,omitempty"`
}

type Attachment struct {
	Id       string `json:"id"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	MimeType string `json:"mimeType"`
}

const BotAPI = "/v1/bot-api"

func SendMessage(chatId, body string, attachmentIds []string, botKey string) {
	jsonBody, err := json.Marshal(&Message{
		ChatId:        chatId,
		Body:          body,
		AttachmentIds: attachmentIds,
	})
	if err != nil {
		fmt.Printf("can't marshal json %s\n", err.Error())
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		os.Getenv("API_HOST")+BotAPI+"/chat/message",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		fmt.Printf("can't make a new request %s\n", err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", botKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("can't send message to user %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("got response with error code: %d\n", resp.StatusCode)
	}
}

func UploadFile(content *bytes.Buffer, filename, botKey string) ([]Attachment, error) {
	if content == nil {
		return []Attachment{}, nil
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	headers := textproto.MIMEHeader{}
	headers.Set("Content-Disposition", `form-data; name="file"; `+`filename="`+filename+`"`)
	headers.Set("Content-Type", defineContentType(filename))

	part, err := writer.CreatePart(headers)
	if err != nil {
		return nil, err
	}

	_, err = part.Write(content.Bytes())
	if err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", os.Getenv("API_HOST")+BotAPI+"/files", body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", botKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("files upload failed with status code: %d", resp.StatusCode)
	}

	var attachments []Attachment
	err = json.NewDecoder(resp.Body).Decode(&attachments)
	if err != nil {
		return nil, err
	}

	return attachments, nil
}

func defineContentType(filename string) string {
	parts := strings.Split(filename, ".")
	extension := strings.ToLower(parts[len(parts)-1])
	contentType := mime.TypeByExtension("." + extension)

	if contentType == "" {
		contentType = "text/plain"
	}

	return contentType
}
