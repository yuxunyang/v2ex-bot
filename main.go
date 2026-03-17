package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 1. 定义命令行参数
	vPath := flag.String("v2ex", "", "V2EX Cookie文件路径")
	gPath := flag.String("glados", "", "GLaDOS Cookie文件路径")

	// 通知参数
	barkKey := flag.String("bark", "", "Bark Key (选填)")
	tgToken := flag.String("tg-token", "", "Telegram Bot Token (选填)")
	tgChatID := flag.String("tg-chatid", "", "Telegram Chat ID (选填)")

	// 随机睡眠参数 (单位：分钟)
	maxSleepMin := flag.Int("sleep", 0, "最大随机睡眠时长 (分钟)，设置为 0 则不睡眠")

	flag.Parse()

	// 2. 执行随机睡眠逻辑 (分钟转秒)
	if *maxSleepMin > 0 {
		// 初始化随机种子
		rand.Seed(time.Now().UnixNano())

		// 将分钟转换为秒进行更细粒度的计算
		maxSeconds := *maxSleepMin * 60
		sleepSeconds := rand.Intn(maxSeconds + 1)

		if sleepSeconds > 0 {
			// 格式化打印：如果超过60秒则显示分秒，否则只显示秒
			if sleepSeconds >= 60 {
				fmt.Printf("⏰ 随机等待中: %d 分 %d 秒...\n", sleepSeconds/60, sleepSeconds%60)
			} else {
				fmt.Printf("⏰ 随机等待中: %d 秒...\n", sleepSeconds)
			}
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
		}
	}

	// 3. 初始化通知渠道 (多渠道组合)
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

	// 4. 执行任务逻辑
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
