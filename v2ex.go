package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func runV2EX(cookie string, notifier Notifier) {
	if cookie == "" { return }
	
	client := &http.Client{}
	taskURL := "https://www.v2ex.com/mission/daily"
	
	// 1. 获取签到链接
	req, _ := http.NewRequest("GET", taskURL, nil)
	setHeaders(req, cookie)
	resp, err := client.Do(req)
	if err != nil { return }
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`/mission/daily/redeem\?once=\d+`)
	path := re.FindString(string(body))

	var msg string
	if path == "" {
		msg = "未找到签到链接（可能已签到）"
	} else {
		// 2. 执行签到
		req2, _ := http.NewRequest("GET", "https://www.v2ex.com"+path, nil)
		setHeaders(req2, cookie)
		req2.Header.Set("Referer", taskURL)
		client.Do(req2)
		msg = "签到执行成功"
	}

	fmt.Printf("[V2EX] %s\n", msg)
	if notifier != nil {
		_ = notifier.Send("V2EX 签到", msg)
	}
}