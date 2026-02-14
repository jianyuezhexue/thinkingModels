package main

import (
	_ "thinkingModels/docs"
	"thinkingModels/engine"
)

// @title Thinking Models Platform API
// @version 1.0
// @description 思维模型平台 - 基于结构化思维模型和AI的复杂问题深度思考平台
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@thinkingmodels.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:2500
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT token认证，格式：Bearer {token}

// apifox token = afxp_731980BYiStMxas7r5wJnWXl5thuEaPchjfg
func main() {
	// 启动服务
	engine.Run()
}
