package middleware

import (
	"net/http"
	"strings"
	"time"
	"ecommerce-rbac-system/internal/config"
	"ecommerce-rbac-system/internal/models"
	"ecommerce-rbac-system/internal/service"
	"github.com/gin-gonic/gin"
)

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

		user, err := authService.GetUserFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("token", token)

		c.Next()
	}
}

// PermissionMiddleware 通过闭包接收 authService，从 context 取 *models.User
func PermissionMiddleware(authService service.AuthService, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		user, ok := userVal.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息异常"})
			c.Abort()
			return
		}

		hasPermission, err := authService.HasPermission(user.ID, requiredPermission)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func DataScopeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		user, ok := userVal.(*models.User)
		if !ok {
			c.Next()
			return
		}

		// 将用户数据权限范围存入 context 供后续使用
		if len(user.Roles) > 0 {
			minScope := user.Roles[0].DataScope
			for _, role := range user.Roles {
				if role.DataScope < minScope {
					minScope = role.DataScope
				}
			}
			c.Set("dataScope", minScope)
		}

		c.Next()
	}
}

func OperationLogMiddleware(logService service.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		if userVal, exists := c.Get("user"); exists {
			if user, ok := userVal.(*models.User); ok {
				costTime := time.Since(startTime).Milliseconds()
				status := 1
				if c.Writer.Status() >= 400 {
					status = 0
				}
				_ = user
				_ = costTime
				_ = status
				// TODO: 将操作日志写入数据库
			}
		}
	}
}
