# 基础变量
BINARY_NAME=media-organizer
VERSION=1.0.0
BUILD_DIR=build
MAIN_PATH=cmd/media-organizer/main.go

# Go命令
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# 编译标记
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# 目标平台列表
PLATFORMS=linux-amd64 linux-arm64 linux-armv7 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64 synology-amd64 synology-arm64 synology-armv7

# 清理构建目录
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# 创建构建目录
.PHONY: init
init:
	mkdir -p $(BUILD_DIR)

# 构建所有平台
.PHONY: all
all: clean init $(PLATFORMS)

# Linux AMD64 (适用于大多数Linux系统和群晖)
.PHONY: linux-amd64
linux-amd64: init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

# Linux ARM64 (适用于树莓派和ARM架构的NAS)
.PHONY: linux-arm64
linux-arm64: init
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

# Linux ARMv7 (适用于32位ARM设备，如较老的树莓派)
.PHONY: linux-armv7
linux-armv7: init
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME)-linux-armv7 $(MAIN_PATH)

# macOS AMD64
.PHONY: darwin-amd64
darwin-amd64: init
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

# macOS ARM64 (Apple Silicon)
.PHONY: darwin-arm64
darwin-arm64: init
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

# Windows AMD64
.PHONY: windows-amd64
windows-amd64: init
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Windows ARM64
.PHONY: windows-arm64
windows-arm64: init
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)

# Synology AMD64 (适用于Intel/AMD架构的群晖NAS)
.PHONY: synology-amd64
synology-amd64: init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME)-synology-amd64 $(MAIN_PATH)
	cd $(BUILD_DIR) && tar czf $(BINARY_NAME)-synology-amd64.spk $(BINARY_NAME)-synology-amd64

# Synology ARM64 (适用于ARM64架构的群晖NAS)
.PHONY: synology-arm64
synology-arm64: init
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME)-synology-arm64 $(MAIN_PATH)
	cd $(BUILD_DIR) && tar czf $(BINARY_NAME)-synology-arm64.spk $(BINARY_NAME)-synology-arm64

# Synology ARMv7 (适用于32位ARM架构的群晖NAS)
.PHONY: synology-armv7
synology-armv7: init
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME)-synology-armv7 $(MAIN_PATH)
	cd $(BUILD_DIR) && tar czf $(BINARY_NAME)-synology-armv7.spk $(BINARY_NAME)-synology-armv7

# 运行测试
.PHONY: test
test:
	$(GOTEST) -v ./...

# 构建当前平台版本
.PHONY: build
build: init
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# 打包发布文件
.PHONY: release
release: all
	cd $(BUILD_DIR) && \
	for file in $(BINARY_NAME)-* ; do \
		zip $${file}.zip $$file ; \
	done

# 帮助信息
.PHONY: help
help:
	@echo "可用的构建目标:"
	@echo "  make build    - 构建当前平台的二进制文件"
	@echo "  make all      - 构建所有平台的二进制文件"
	@echo "  make release  - 构建并打包所有平台的二进制文件"
	@echo "  make test     - 运行测试"
	@echo "  make clean    - 清理构建目录"
	@echo ""
	@echo "单平台构建目标:"
	@echo "  make linux-amd64   - 构建 Linux AMD64 版本"
	@echo "  make linux-arm64   - 构建 Linux ARM64 版本"
	@echo "  make linux-armv7   - 构建 Linux ARMv7 版本"
	@echo "  make darwin-amd64  - 构建 macOS AMD64 版本"
	@echo "  make darwin-arm64  - 构建 macOS ARM64 版本"
	@echo "  make windows-amd64 - 构建 Windows AMD64 版本"
	@echo "  make windows-arm64 - 构建 Windows ARM64 版本"
	@echo ""
	@echo "群晖NAS构建目标:"
	@echo "  make synology-amd64  - 构建 群晖 Intel/AMD 版本"
	@echo "  make synology-arm64  - 构建 群晖 ARM64 版本"
	@echo "  make synology-armv7  - 构建 群晖 ARMv7 版本" 