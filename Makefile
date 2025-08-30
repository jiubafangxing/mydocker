# Makefile for mydocker project

# 变量定义
BINARY_NAME=mydocker
BUILD_DIR=build
VERSION=v1.0.0
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

# Go相关变量
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# LDFLAGS
LDFLAGS=-ldflags "-s -w -X 'main.Version=$(VERSION)' -X 'main.GitCommit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'"

# 默认目标
.PHONY: all
all: clean build

# 构建
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build completed: $(BINARY_NAME)"

# 交叉编译
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux .

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows.exe .

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin .

# 构建所有平台
.PHONY: build-all
build-all: build-linux build-windows build-darwin
	@echo "All builds completed"

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# 清理
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -rf $(BUILD_DIR)

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# 代码格式化
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# 开发模式（监听文件变化并重新构建）
.PHONY: dev
dev:
	@echo "Starting development mode..."
	@which air > /dev/null || (echo "Installing air..." && $(GOGET) -u github.com/cosmtrek/air)
	air

# 安装到系统
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation completed"

# 显示版本信息
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

# 帮助
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  build-all   - Build for all platforms (Linux, Windows, macOS)"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  fmt         - Format code"
	@echo "  vet         - Run go vet"
	@echo "  install     - Install binary to system"
	@echo "  version     - Show version information"
	@echo "  help        - Show this help message"
