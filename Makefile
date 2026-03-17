# 项目名称
BINARY_NAME=v2ex_bot

# 编译参数：-s 移除符号表，-w 移除调试信息
LDFLAGS=-ldflags="-s -w"
# 输出目录
DIST_DIR=bin/dist

.PHONY: all build clean help build-all linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64

# 默认目标：编译
all: build

## build: 编译二进制文件 (优化体积)
build:
	@echo "正在编译..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .
	@echo "编译完成！可执行文件为: bin/$(BINARY_NAME)"
	@file bin/$(BINARY_NAME)

## build-all: 一键编译所有主流平台的架构 (Linux, macOS, Windows)
build-all: clean linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64
	@echo "所有架构编译完成，请查看 $(DIST_DIR) 目录"

## linux-amd64: 编译 Linux 64位 (常见服务器)
linux-amd64:
	@echo "编译 Linux amd64..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_linux_amd64 .

## linux-arm64: 编译 Linux ARM64 (树莓派、甲骨文 ARM 云服务器)
linux-arm64:
	@echo "编译 Linux arm64..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_linux_arm64 .

## darwin-amd64: 编译 macOS Intel 架构
darwin-amd64:
	@echo "编译 macOS amd64..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_darwin_amd64 .

## darwin-arm64: 编译 macOS M1/M2/M3 芯片架构
darwin-arm64:
	@echo "编译 macOS arm64..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_darwin_arm64 .

## windows-amd64: 编译 Windows 64位
windows-amd64:
	@echo "编译 Windows amd64..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_windows_amd64.exe .

## clean: 清理编译产生的文件
clean:
	@echo "清理中..."
	@rm -f bin/$(BINARY_NAME)
	@rm -rf $(DIST_DIR)
	@echo "清理完成。"

## run: 直接编译并运行 (示例运行命令)
# 使用方法: make run V2EX_PATH=v2ex.txt BARK_KEY=xxx
run: build
	./bin/$(BINARY_NAME) -v2ex $(V2EX_PATH) -glados $(GLADOS_PATH) -bark $(BARK_KEY) -tg-token $(TG_TOKEN) -tg-chatid $(TG_CHATID)

## linux: 跨平台编译 Linux 版本 (适用于在 macOS/Windows 上为 Linux 编译)
linux:
	@echo "正在为 Linux (x86_64) 编译..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)_linux .
	@echo "编译完成: bin/$(BINARY_NAME)_linux"

## help: 显示帮助信息
help:
	@echo "用法:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'