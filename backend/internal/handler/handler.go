package handler

import (
	"ecommerce-rbac-system/internal/models"
	"ecommerce-rbac-system/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	User         UserHandler
	Auth         AuthHandler
	Role         RoleHandler
	Permission   PermissionHandler
	Department   DepartmentHandler
	LoginLog     LoginLogHandler
	OperationLog OperationLogHandler
}

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	permissions, _ := h.authService.GetPermissionCodes(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"token":       token,
		"user":        user,
		"permissions": permissions,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// logout 现在在认证路由组内，可以从 context 取 token；
	// 同时兼容直接从 header 解析的方式
	token := c.GetString("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供token"})
		return
	}

	if err := h.authService.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "退出成功"})
}

type UserHandler struct {
	userService service.UserService
	authService service.AuthService
}

func NewUserHandler(userService service.UserService, authService service.AuthService) *UserHandler {
	return &UserHandler{userService, authService}
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	permissions, err := h.authService.GetPermissionCodes(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":        user,
		"permissions": permissions,
	})
}

func (h *UserHandler) UpdateInfo(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(*models.User)
	req.ID = user.ID

	if err := h.userService.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	users, total, err := h.userService.List(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, err := h.userService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		models.User
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.User.Password = req.Password
	if err := h.userService.Create(&req.User); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func (h *UserHandler) Update(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	req.ID = id

	if err := h.userService.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.userService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *UserHandler) AssignRoles(c *gin.Context) {
	var req struct {
		RoleIDs []int64 `json:"roleIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.userService.AssignRoles(id, req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色分配成功"})
}

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService}
}

func (h *RoleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	roles, total, err := h.roleService.List(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     roles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	role, err := h.roleService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req models.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roleService.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func (h *RoleHandler) Update(c *gin.Context) {
	var req models.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	req.ID = id

	if err := h.roleService.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.roleService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	var req struct {
		PermissionIDs []int64 `json:"permissionIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.roleService.AssignPermissions(id, req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限分配成功"})
}

type PermissionHandler struct {
	permissionService service.PermissionService
}

func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService}
}

func (h *PermissionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	permissions, total, err := h.permissionService.List(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     permissions,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *PermissionHandler) Tree(c *gin.Context) {
	tree, err := h.permissionService.Tree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tree})
}

func (h *PermissionHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	permission, err := h.permissionService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "权限不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permission})
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var req models.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.permissionService.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func (h *PermissionHandler) Update(c *gin.Context) {
	var req models.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	req.ID = id

	if err := h.permissionService.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.permissionService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

type DepartmentHandler struct {
	departmentService service.DepartmentService
}

func NewDepartmentHandler(departmentService service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{departmentService}
}

func (h *DepartmentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	departments, total, err := h.departmentService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     departments,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *DepartmentHandler) Tree(c *gin.Context) {
	tree, err := h.departmentService.Tree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tree})
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	department, err := h.departmentService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": department})
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	var req models.Department
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.departmentService.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	var req models.Department
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	req.ID = id

	if err := h.departmentService.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.departmentService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

type LoginLogHandler struct {
	loginLogService service.LoginLogService
}

func NewLoginLogHandler(loginLogService service.LoginLogService) *LoginLogHandler {
	return &LoginLogHandler{loginLogService}
}

func (h *LoginLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	logs, total, err := h.loginLogService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

type OperationLogHandler struct {
	operationLogService service.OperationLogService
}

func NewOperationLogHandler(operationLogService service.OperationLogService) *OperationLogHandler {
	return &OperationLogHandler{operationLogService}
}

func (h *OperationLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	logs, total, err := h.operationLogService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		User:         *NewUserHandler(services.User, services.Auth),
		Auth:         *NewAuthHandler(services.Auth),
		Role:         *NewRoleHandler(services.Role),
		Permission:   *NewPermissionHandler(services.Permission),
		Department:   *NewDepartmentHandler(services.Department),
		LoginLog:     *NewLoginLogHandler(services.LoginLog),
		OperationLog: *NewOperationLogHandler(services.OperationLog),
	}
}
