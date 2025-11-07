# Go 电商项目 Makefile
.PHONY: help build run clean test lint deps dev prod docker-build docker-run swagger

# 变量定义
BINARY_NAME=ecommerce
MAIN_FILE=main.go
BUILD_DIR=build
AIR_TMP_DIR=.air
DOCKER_IMAGE=ecommerce-app
DOCKER_CONTAINER=ecommerce-container
GO_VERSION=$(shell go version | awk '{print $$3}')

# 默认目标
help: ## 显示帮助信息
	@echo "Go 电商项目 Makefile"
	@echo "====================="
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 构建相关
deps: ## 下载依赖
	go mod download
	go mod tidy

build: ## 构建二进制文件
	@echo "正在构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

build-linux: ## 构建 Linux 二进制文件
	@echo "正在构建 Linux 版 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@echo "Linux 构建完成: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

# 开发相关
dev: ## 开发模式运行（使用 Air 热重载）
	@echo "启动开发模式..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air 未安装，正在安装..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

run: build ## 构建并运行应用
	@echo "运行应用..."
	$(BUILD_DIR)/$(BINARY_NAME)

# 测试相关
test: ## 运行测试
	@echo "运行测试..."
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 代码质量
lint: ## 运行代码检查
	@echo "运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，建议安装: https://golangci-lint.run/usage/install/"; \
		go vet ./...; \
	fi

fmt: ## 格式化代码
	@echo "格式化代码..."
	go fmt ./...

# 清理
clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(AIR_TMP_DIR)
	@rm -f coverage.out coverage.html
	@echo "清理完成"

# Docker 相关
docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	docker build -t $(DOCKER_IMAGE):latest .

docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	docker run -d --name $(DOCKER_CONTAINER) -p 8080:8080 $(DOCKER_IMAGE):latest

docker-stop: ## 停止 Docker 容器
	@echo "停止 Docker 容器..."
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

# 数据库相关
db-migrate: ## 运行数据库迁移（如果有的话）
	@echo "运行数据库迁移..."
	@echo "注意：请根据你的数据库迁移工具实现此功能"

# 工具安装
install-tools: ## 安装开发工具
	@echo "安装开发工具..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "开发工具安装完成"

# 文档
swagger: ## 生成 Swagger 文档（如果有的话）
	@echo "生成 Swagger 文档..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init; \
	else \
		echo "swag 未安装，请先运行: make install-tools"; \
	fi

# 生产环境
prod: build ## 生产环境构建
	@echo "生产环境构建..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w -X main.version=$(shell git describe --tags --always 2>/dev/null || echo 'dev')" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-prod $(MAIN_FILE)
	@echo "生产环境构建完成: $(BUILD_DIR)/$(BINARY_NAME)-prod"

# 系统信息
info: ## 显示系统信息
	@echo "系统信息:"
	@echo "Go 版本: $(GO_VERSION)"
	@echo "项目模块: $(shell go list -m)"
	@echo "构建目录: $(BUILD_DIR)"
	@echo "二进制名称: $(BINARY_NAME)"

# 快速开始
setup: deps install-tools ## 项目初始化设置
	@echo "项目初始化完成！"
	@echo "可以运行: make dev 开始开发"

# 默认目标
.DEFAULT_GOAL := help