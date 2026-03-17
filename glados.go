package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func runGLaDOS(cookie string, notifier Notifier) {
	if cookie == "" {
		return
	}

	apiURL := "https://glados.cloud/api/user/checkin"
	payload := []byte(`{"token":"glados.cloud"}`)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	setHeaders(req, cookie)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://glados.cloud")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Message string `json:"message"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Printf("[GLaDOS] %s\n", result.Message)
	if notifier != nil {
		_ = notifier.Send("GLaDOS 签到", result.Message)
	}
}
