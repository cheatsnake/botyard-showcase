package main

import (
	"bytes"
	"math/rand"

	"github.com/skip2/go-qrcode"
)

func generateRandomQrcode(size int) (*bytes.Buffer, error) {
	randomContent := generateRandomString(64)

	qr, err := qrcode.New(randomContent, qrcode.Low)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	content, err := qr.PNG(size)
	if err != nil {
		return nil, err
	}

	body.Write(content)

	return &body, nil

}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
