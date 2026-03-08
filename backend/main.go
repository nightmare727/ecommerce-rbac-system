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
	cfg := config.Load()

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient, err := database.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, redisClient, cfg)
	handlers := handler.NewHandlers(services)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// perm 是带 authService 闭包的权限中间件快捷方法
	perm := func(code string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(services.Auth, code)
	}

	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", handlers.Auth.Login)
	}

	auth := r.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware(services.Auth, cfg))
	{
		// logout 需要认证才能拿到 token
		auth.POST("/auth/logout", handlers.Auth.Logout)

		auth.GET("/user/info", handlers.User.GetInfo)
		auth.PUT("/user/info", handlers.User.UpdateInfo)

		users := auth.Group("/users")
		users.Use(perm("system:user:list"))
		{
			users.GET("", handlers.User.List)
			users.GET("/:id", handlers.User.GetByID)
		}
		users.POST("", perm("system:user:add"), handlers.User.Create)
		users.PUT("/:id", perm("system:user:edit"), handlers.User.Update)
		users.DELETE("/:id", perm("system:user:delete"), handlers.User.Delete)
		users.POST("/:id/roles", perm("system:user:assign"), handlers.User.AssignRoles)

		roles := auth.Group("/roles")
		roles.Use(perm("system:role:list"))
		{
			roles.GET("", handlers.Role.List)
			roles.GET("/:id", handlers.Role.GetByID)
		}
		roles.POST("", perm("system:role:add"), handlers.Role.Create)
		roles.PUT("/:id", perm("system:role:edit"), handlers.Role.Update)
		roles.DELETE("/:id", perm("system:role:delete"), handlers.Role.Delete)
		roles.POST("/:id/permissions", perm("system:role:assign"), handlers.Role.AssignPermissions)

		permissions := auth.Group("/permissions")
		permissions.Use(perm("system:permission:list"))
		{
			permissions.GET("", handlers.Permission.List)
			permissions.GET("/tree", handlers.Permission.Tree)
			permissions.GET("/:id", handlers.Permission.GetByID)
		}
		permissions.POST("", perm("system:permission:add"), handlers.Permission.Create)
		permissions.PUT("/:id", perm("system:permission:edit"), handlers.Permission.Update)
		permissions.DELETE("/:id", perm("system:permission:delete"), handlers.Permission.Delete)

		depts := auth.Group("/departments")
		depts.Use(perm("system:dept:list"))
		{
			depts.GET("", handlers.Department.List)
			depts.GET("/tree", handlers.Department.Tree)
			depts.GET("/:id", handlers.Department.GetByID)
		}
		depts.POST("", perm("system:dept:add"), handlers.Department.Create)
		depts.PUT("/:id", perm("system:dept:edit"), handlers.Department.Update)
		depts.DELETE("/:id", perm("system:dept:delete"), handlers.Department.Delete)

		loginLogs := auth.Group("/login-logs")
		loginLogs.Use(perm("system:log:login"))
		{
			loginLogs.GET("", handlers.LoginLog.List)
		}
		operationLogs := auth.Group("/operation-logs")
		operationLogs.Use(perm("system:log:operation"))
		{
			operationLogs.GET("", handlers.OperationLog.List)
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
