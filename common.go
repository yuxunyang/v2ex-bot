package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36"

// Notifier 定义通知接口
type Notifier interface {
	Send(title, message string) error
}

// MultiNotifier 支持同时发送多个通知渠道
type MultiNotifier struct {
	Notifiers []Notifier
}

func (m *MultiNotifier) Send(title, message string) error {
	for _, n := range m.Notifiers {
		_ = n.Send(title, message)
	}
	return nil
}

// BarkNotifier 实现
type BarkNotifier struct {
	Key string
}

func (b *BarkNotifier) Send(title, message string) error {
	apiURL := fmt.Sprintf("https://api.day.app/%s/%s/%s", b.Key, url.PathEscape(title), url.PathEscape(message))
	_, err := http.Get(apiURL)
	return err
}

// TelegramNotifier 实现
type TelegramNotifier struct {
	Token  string
	ChatID string
}

func (t *TelegramNotifier) Send(title, message string) error {
	text := fmt.Sprintf("<b>%s</b>\n%s", title, message)
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=HTML",
		t.Token, t.ChatID, url.QueryEscape(text))
	_, err := http.Get(apiURL)
	return err
}

// 通用读取文件方法
func readSecret(path string) string {
	if path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("读取文件 [%s] 失败: %v", path, err)
		return ""
	}
	return strings.TrimSpace(string(data))
}

// 通用 Header 设置
func setHeaders(req *http.Request, cookie string) {
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Cookie", cookie)
}
