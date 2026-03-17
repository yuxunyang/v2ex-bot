package main

import (
	"flag"
	"fmt"
)

func main() {
	// 任务参数
	vPath := flag.String("v2ex", "", "V2EX Cookie文件路径")
	gPath := flag.String("glados", "", "GLaDOS Cookie文件路径")

	// 通知参数
	barkKey := flag.String("bark", "", "Bark Key (选填)")
	tgToken := flag.String("tg-token", "", "Telegram Bot Token (选填)")
	tgChatID := flag.String("tg-chatid", "", "Telegram Chat ID (选填)")

	flag.Parse()

	// 初始化通知渠道
	multi := &MultiNotifier{}

	if *barkKey != "" {
		multi.Notifiers = append(multi.Notifiers, &BarkNotifier{Key: *barkKey})
		fmt.Println("已启用 Bark 通知")
	}

	if *tgToken != "" && *tgChatID != "" {
		multi.Notifiers = append(multi.Notifiers, &TelegramNotifier{
			Token:  *tgToken,
			ChatID: *tgChatID,
		})
		fmt.Println("已启用 Telegram 通知")
	}

	// 执行任务
	if *vPath != "" {
		vCookie := readSecret(*vPath)
		if vCookie != "" {
			runV2EX(vCookie, multi)
		}
	}

	if *gPath != "" {
		gCookie := readSecret(*gPath)
		if gCookie != "" {
			runGLaDOS(gCookie, multi)
		}
	}
}
