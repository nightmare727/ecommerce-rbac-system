package main

import (
	"log"
	"ecommerce-rbac-system/internal/config"
	"ecommerce-rbac-system/internal/database"
	"ecommerce-rbac-system/internal/handler"
	"ecommerce-rbac-system/internal/middleware"
	"ecommerce-rbac-system/internal/repository"
	"ecommerce-rbac-system/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化Redis
	redisClient, err := database.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 初始化仓库层
	repos := repository.NewRepositories(db)

	// 初始化服务层
	services := service.NewServices(repos, redisClient, cfg)

	// 初始化处理器
	handlers := handler.NewHandlers(services)

	// 创建路由
	r := gin.Default()

	// 跨域配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 公开路由（无需认证）
	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", handlers.Auth.Login)
		public.POST("/auth/logout", handlers.Auth.Logout)
	}

	// 需要认证的路由
	auth := r.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware(services.Auth, cfg))
	{
		// 获取当前用户信息
		auth.GET("/user/info", handlers.User.GetInfo)
		auth.PUT("/user/info", handlers.User.UpdateInfo)

		// 用户管理
		users := auth.Group("/users")
		users.Use(middleware.PermissionMiddleware("system:user:list"))
		{
			users.GET("", handlers.User.List)
			users.GET("/:id", handlers.User.GetByID)
		}
		users.POST("", middleware.PermissionMiddleware("system:user:add"), handlers.User.Create)
		users.PUT("/:id", middleware.PermissionMiddleware("system:user:edit"), handlers.User.Update)
		users.DELETE("/:id", middleware.PermissionMiddleware("system:user:delete"), handlers.User.Delete)
		users.POST("/:id/roles", middleware.PermissionMiddleware("system:user:assign"), handlers.User.AssignRoles)

		// 角色管理
		roles := auth.Group("/roles")
		roles.Use(middleware.PermissionMiddleware("system:role:list"))
		{
			roles.GET("", handlers.Role.List)
			roles.GET("/:id", handlers.Role.GetByID)
		}
		roles.POST("", middleware.PermissionMiddleware("system:role:add"), handlers.Role.Create)
		roles.PUT("/:id", middleware.PermissionMiddleware("system:role:edit"), handlers.Role.Update)
		roles.DELETE("/:id", middleware.PermissionMiddleware("system:role:delete"), handlers.Role.Delete)
		roles.POST("/:id/permissions", middleware.PermissionMiddleware("system:role:assign"), handlers.Role.AssignPermissions)

		// 权限管理
		permissions := auth.Group("/permissions")
		permissions.Use(middleware.PermissionMiddleware("system:permission:list"))
		{
			permissions.GET("", handlers.Permission.List)
			permissions.GET("/tree", handlers.Permission.Tree)
			permissions.GET("/:id", handlers.Permission.GetByID)
		}
		permissions.POST("", middleware.PermissionMiddleware("system:permission:add"), handlers.Permission.Create)
		permissions.PUT("/:id", middleware.PermissionMiddleware("system:permission:edit"), handlers.Permission.Update)
		permissions.DELETE("/:id", middleware.PermissionMiddleware("system:permission:delete"), handlers.Permission.Delete)

		// 部门管理
		depts := auth.Group("/departments")
		depts.Use(middleware.PermissionMiddleware("system:dept:list"))
		{
			depts.GET("", handlers.Department.List)
			depts.GET("/tree", handlers.Department.Tree)
			depts.GET("/:id", handlers.Department.GetByID)
		}
		depts.POST("", middleware.PermissionMiddleware("system:dept:add"), handlers.Department.Create)
		depts.PUT("/:id", middleware.PermissionMiddleware("system:dept:edit"), handlers.Department.Update)
		depts.DELETE("/:id", middleware.PermissionMiddleware("system:dept:delete"), handlers.Department.Delete)

		// 日志管理
		loginLogs := auth.Group("/login-logs")
		loginLogs.Use(middleware.PermissionMiddleware("system:log:login"))
		{
			loginLogs.GET("", handlers.LoginLog.List)
		}
		operationLogs := auth.Group("/operation-logs")
		operationLogs.Use(middleware.PermissionMiddleware("system:log:operation"))
		{
			operationLogs.GET("", handlers.OperationLog.List)
		}
	}

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
