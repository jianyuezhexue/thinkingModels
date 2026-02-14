# ThinkingModels Backend Makefile

.PHONY: help build run test clean swagger

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
help: ## 显示帮助信息
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 构建相关
build: ## 构建应用
	go build -o bin/app main.go

run: ## 运行应用
	go run main.go

run-dev: ## 开发模式运行（使用 air 热重载）
	air

# 测试相关
test: ## 运行测试
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码质量
lint: ## 运行 linter
	golangci-lint run

fmt: ## 格式化代码
	go fmt ./...

vet: ## 运行 go vet
	go vet ./...

# Swagger 文档生成
swagger: ## 生成 Swagger API 文档
	swag init --parseDependency --parseInternal

swagger-serve: ## 生成并运行 Swagger UI（需要先启动服务）
	swag init --parseDependency --parseInternal
	@echo "Swagger UI available at: http://localhost:2500/swagger/index.html"

swagger-clean: ## 清理 Swagger 文档
	rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml

# 依赖管理
deps: ## 下载依赖
	go mod download

deps-tidy: ## 整理依赖
	go mod tidy

deps-vendor: ## 生成 vendor 目录
	go mod vendor

# 数据库相关（可选）
migrate-up: ## 执行数据库迁移
	@echo "TODO: 实现数据库迁移"

migrate-down: ## 回滚数据库迁移
	@echo "TODO: 实现数据库回滚"

# Docker 相关（可选）
docker-build: ## 构建 Docker 镜像
	docker build -t thinkingmodels-backend .

docker-run: ## 运行 Docker 容器
	docker run -p 2500:2500 thinkingmodels-backend

# 清理
clean: ## 清理构建产物
	rm -rf bin/ coverage.out coverage.html
	go clean

# 完整开发流程
dev-setup: ## 开发环境初始化
	go mod download
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
