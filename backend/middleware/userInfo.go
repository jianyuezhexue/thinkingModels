package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 密钥（应该从配置文件读取，这里暂时硬编码）
var jwtSecret = "your-secret-key-change-in-production"

// UserInfo JWT认证中间件
func UserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有 token，设置为默认系统用户（base库要求currUserId不能为空）
			c.Set("currUserId", "0")
			c.Set("currUserName", "系统")
			c.Set("currRoleIds", "")
			c.Next()
			return
		}

		// 提取 token（Bearer token 格式）
		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = authHeader
		}

		// 解析 JWT token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			// token 无效，设置为默认系统用户
			c.Set("currUserId", "0")
			c.Set("currUserName", "系统")
			c.Set("currRoleIds", "")
			c.Next()
			return
		}

		// 从 claims 中提取用户信息
		if claims, ok := token.Claims.(*UserClaims); ok {
			// 将用户ID转换为字符串（base库期望字符串类型）
			c.Set("currUserId", fmt.Sprintf("%.0f", claims.UserID))
			c.Set("currUserName", claims.Username)
			c.Set("currRoleIds", claims.RoleIds)
		}

		c.Next()
	}
}

// UserClaims JWT Claims 定义（与 domain/iam/user/model.go 中的定义一致）
type UserClaims struct {
	UserID       float64 `json:"sub"` // JWT 解析数字为 float64
	Username     string  `json:"username"`
	EnterpriseID float64 `json:"enterprise_id"`
	RoleIds      string  `json:"role_ids"`
	jwt.RegisteredClaims
}

// GetUserID 将 UserID 转换为 uint64
func (c *UserClaims) GetUserID() uint64 {
	return uint64(c.UserID)
}