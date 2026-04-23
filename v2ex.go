package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

var vlog = log.New(os.Stdout, "[V2EX]   ", log.Ldate|log.Ltime|log.Lmsgprefix)

func fetchBalance(client *http.Client, cookie string) string {
	req, _ := http.NewRequest("GET", "https://www.v2ex.com/balance", nil)
	setHeaders(req, cookie)
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// 页面中余额格式：28 <img ... alt="S" ...> 81 <img ... alt="B" ...>
	reS := regexp.MustCompile(`(\d+) <img [^>]*alt="S"`)
	reB := regexp.MustCompile(`(\d+) <img [^>]*alt="B"`)
	silver, bronze := "", ""
	if m := reS.FindStringSubmatch(string(body)); len(m) > 1 {
		silver = m[1]
	}
	if m := reB.FindStringSubmatch(string(body)); len(m) > 1 {
		bronze = m[1]
	}
	if silver == "" && bronze == "" {
		return ""
	}
	return fmt.Sprintf("%s 银币 / %s 铜币", silver, bronze)
}

func runV2EX(cookie string, notifier Notifier) {
	if cookie == "" {
		vlog.Printf("cookie 未配置，跳过")
		return
	}

	client := &http.Client{}
	taskURL := "https://www.v2ex.com/mission/daily"

	// 1. 获取签到页
	vlog.Printf("请求签到页: %s", taskURL)
	req, _ := http.NewRequest("GET", taskURL, nil)
	setHeaders(req, cookie)
	resp, err := client.Do(req)
	if err != nil {
		vlog.Printf("请求签到页失败: %v", err)
		return
	}
	defer resp.Body.Close()
	vlog.Printf("签到页响应状态: %s", resp.Status)

	// 检查是否被重定向到登录页（cookie 无效）
	if resp.Request.URL.Path == "/signin" {
		msg := "cookie 无效或已过期，已被重定向到登录页"
		vlog.Printf(msg)
		if notifier != nil {
			_ = notifier.Send("V2EX 签到", msg)
		}
		return
	}

	// 提取新的 PB3_SESSION（与 once token 配对，必须用它发签到请求）
	updatedCookie := cookie
	rePB3 := regexp.MustCompile(`PB3_SESSION="[^"]*"`)
	for _, sc := range resp.Header["Set-Cookie"] {
		if m := rePB3.FindString(sc); m != "" {
			updatedCookie = rePB3.ReplaceAllString(cookie, m)
			vlog.Printf("已更新 PB3_SESSION")
			break
		}
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	reRedeem := regexp.MustCompile(`/mission/daily/redeem\?once=\d+`)
	reCoins := regexp.MustCompile(`领取 (\d+) 铜币`)
	path := reRedeem.FindString(bodyStr)

	var msg string
	if path == "" {
		msg = "未找到签到链接（可能已签到）"
		vlog.Printf(msg)
	} else {
		// 找到 redeem 链接，先记录签到前余额，再尝试签到
		coins := reCoins.FindStringSubmatch(bodyStr)
		balanceBefore := fetchBalance(client, cookie)
		if balanceBefore != "" {
			vlog.Printf("签到前余额: %s", balanceBefore)
		}

		redeemURL := "https://www.v2ex.com" + path
		vlog.Printf("执行签到: %s", redeemURL)
		req2, _ := http.NewRequest("GET", redeemURL, nil)
		setHeaders(req2, updatedCookie)
		req2.Header.Set("Referer", taskURL)
		noRedirectClient := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp2, err2 := noRedirectClient.Do(req2)
		if err2 != nil {
			vlog.Printf("签到请求失败: %v", err2)
			msg = "签到请求失败"
		} else {
			resp2.Body.Close()
			vlog.Printf("签到响应: %s", resp2.Status)
			if resp2.StatusCode == http.StatusFound {
				if len(coins) > 1 {
					msg = fmt.Sprintf("签到成功，领取 %s 铜币", coins[1])
				} else {
					msg = "签到成功"
				}
				balanceAfter := fetchBalance(client, cookie)
				if balanceBefore != "" && balanceAfter != "" {
					balanceChange := fmt.Sprintf("%s → %s", balanceBefore, balanceAfter)
					vlog.Printf("余额变动: %s", balanceChange)
					msg += "\n余额变动: " + balanceChange
				} else if balanceAfter != "" {
					vlog.Printf("当前余额: %s", balanceAfter)
					msg += "\n当前余额: " + balanceAfter
				}
			} else {
				msg = fmt.Sprintf("签到失败（预期 302，实际 %d）", resp2.StatusCode)
			}
			vlog.Printf(msg)
		}
	}

	if notifier != nil {
		_ = notifier.Send("V2EX 签到", msg)
	}
}