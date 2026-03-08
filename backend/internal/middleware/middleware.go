package middleware

import (
	"net/http"
	"strings"
	"ecommerce-rbac-system/internal/config"
	"ecommerce-rbac-system/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(authService service.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证token"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			c.Abort()
			return
		}

		// 验证token并获取用户
		user, err := authService.GetUserFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败"})
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("user", user)
		c.Set("token", token)

		c.Next()
	}
}

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从context获取用户和authService
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		authService, exists := c.Get("authService")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务错误"})
			c.Abort()
			return
		}

		// 检查权限
		hasPermission, err := authService.(service.AuthService).HasPermission(user.(*service.User).ID, requiredPermission)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// DataScopeMiddleware 数据权限过滤中间件
func DataScopeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		// 从用户的角色中获取数据权限范围
		// 这里可以添加数据过滤逻辑
		_ = user
		c.Next()
	}
}

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := c.GetTime("startTime")

		// 记录操作日志
		if user, exists := c.Get("user"); exists {
			_ = user
			_ = startTime
			// TODO: 实现日志记录逻辑
		}

		c.Next()
	}
}
