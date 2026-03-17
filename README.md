# Multi-Platform Checkin Bot

这是一个基于 Go 语言编写的轻量化自动签到工具。采用接口化设计，支持多平台签到任务并行执行，并支持 Bark、Telegram 等多渠道通知。

## 🌟 核心特性

- **模块化设计**：V2EX 与 GLaDOS 逻辑完全解耦，易于扩展新平台。
- **多渠道通知**：内置 **Telegram Bot** 与 **Bark** 支持，采用组合模式，支持多渠道同步推送。
- **极简部署**：静态编译，无外部依赖，单个二进制文件即可在 Linux 服务器运行。
- **开发者友好**：提供 `Makefile` 管理构建，编译产物已通过 `strip` 优化体积。

---

## 🛠️ 项目架构

- `common.go`: 核心接口 `Notifier` 定义及通用网络请求封装。
- `v2ex.go`: 负责 V2EX 任务页面爬取及签到逻辑。
- `glados.go`: 负责 GLaDOS API 交互及结果解析。
- `main.go`: 程序入口，负责参数解析与任务调度。

---

## 📦 编译与安装

确保本地已安装 Go 1.18+。

1. 克隆仓库：
   ```shell
   git clone https://github.com/YOUR_USERNAME/checkin-bot.git
   cd checkin-bot
   ```

2. 使用 Makefile 编译：
   ```shell
   make        # 编译当前平台版本
   make linux  # 跨平台编译 Linux (amd64) 版本
   ```
---

## 📋 使用指南

### 1. 准备 Cookie
在根目录下创建文本文件（已被 .gitignore 忽略）：
- `v2ex.txt`: 放入 V2EX 的 Cookie 字符串。
- `glados.txt`: 放入 GLaDOS 的 Cookie 字符串。

### 2. 参数说明

| 参数 | 说明 | 示例 |
| :--- | :--- | :--- |
| -v2ex | V2EX Cookie 文件路径 | ./v2ex.txt |
| -glados | GLaDOS Cookie 文件路径 | ./glados.txt |
| -bark | Bark 推送 Key (选填) | b8Xmj... |
| -tg-token | Telegram Bot Token (选填) | 7777:AA... |
| -tg-chatid | Telegram Chat ID (选填) | 123456 |

### 3. 运行示例
```shell
./v2ex_bot -v2ex v2ex.txt -glados glados.txt -tg-token "TOKEN" -tg-chatid "ID"
```

---

## ⏰ 自动化部署 (Crontab)

建议在 Linux 服务器设置定时任务：
```shell
0 9 * * * /path/to/v2ex_bot -v2ex /path/to/v2ex.txt -tg-token "xxx" -tg-chatid "xxx" >> /path/to/log 2>&1
```

---

## ⚠️ 免责声明
本项目仅供学习交流，请勿用于非法用途。使用本项目造成的任何后果由使用者自行承担。

## 📄 开源协议
MIT License