package service

import (
	"context"
	"errors"
	"fmt"
	"ecommerce-rbac-system/internal/config"
	"ecommerce-rbac-system/internal/models"
	"ecommerce-rbac-system/internal/repository"
	"ecommerce-rbac-system/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	User                UserService
	Auth                AuthService
	Role                RoleService
	Permission          PermissionService
	Department           DepartmentService
	LoginLog            LoginLogService
	OperationLog        OperationLogService
	PermissionChangeLog PermissionChangeLogService
}

type UserService interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int64) error
	GetByID(id int64) (*models.User, error)
	List(page, pageSize int, keyword string) ([]*models.User, int64, error)
	AssignRoles(userID int64, roleIDs []int64) error
	GetPermissions(userID int64) ([]*models.Permission, error)
}

type RoleService interface {
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id int64) error
	GetByID(id int64) (*models.Role, error)
	List(page, pageSize int, keyword string) ([]*models.Role, int64, error)
	AssignPermissions(roleID int64, permissionIDs []int64) error
}

type PermissionService interface {
	Create(permission *models.Permission) error
	Update(permission *models.Permission) error
	Delete(id int64) error
	GetByID(id int64) (*models.Permission, error)
	List(page, pageSize int, keyword string) ([]*models.Permission, int64, error)
	Tree() ([]*models.Permission, error)
}

type DepartmentService interface {
	Create(department *models.Department) error
	Update(department *models.Department) error
	Delete(id int64) error
	GetByID(id int64) (*models.Department, error)
	List(page, pageSize int) ([]*models.Department, int64, error)
	Tree() ([]*models.Department, error)
}

type LoginLogService interface {
	List(page, pageSize int) ([]*models.LoginLog, int64, error)
}

type OperationLogService interface {
	List(page, pageSize int) ([]*models.OperationLog, int64, error)
}

type AuthService interface {
	Login(username, password string) (string, *models.User, error)
	Logout(token string) error
	GetUserFromToken(token string) (*models.User, error)
	HasPermission(userID int64, permissionCode string) (bool, error)
	GetPermissionCodes(userID int64) ([]string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Create(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.repo.Create(user)
}

func (s *userService) Update(user *models.User) error {
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	return s.repo.Update(user)
}

func (s *userService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *userService) GetByID(id int64) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) List(page, pageSize int, keyword string) ([]*models.User, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *userService) AssignRoles(userID int64, roleIDs []int64) error {
	return s.repo.AssignRoles(userID, roleIDs)
}

func (s *userService) GetPermissions(userID int64) ([]*models.Permission, error) {
	return s.repo.GetPermissions(userID)
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo}
}

func (s *roleService) Create(role *models.Role) error {
	return s.repo.Create(role)
}

func (s *roleService) Update(role *models.Role) error {
	return s.repo.Update(role)
}

func (s *roleService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *roleService) GetByID(id int64) (*models.Role, error) {
	return s.repo.GetByID(id)
}

func (s *roleService) List(page, pageSize int, keyword string) ([]*models.Role, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *roleService) AssignPermissions(roleID int64, permissionIDs []int64) error {
	return s.repo.AssignPermissions(roleID, permissionIDs)
}

type permissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(repo repository.PermissionRepository) PermissionService {
	return &permissionService{repo}
}

func (s *permissionService) Create(permission *models.Permission) error {
	return s.repo.Create(permission)
}

func (s *permissionService) Update(permission *models.Permission) error {
	return s.repo.Update(permission)
}

func (s *permissionService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *permissionService) GetByID(id int64) (*models.Permission, error) {
	return s.repo.GetByID(id)
}

func (s *permissionService) List(page, pageSize int, keyword string) ([]*models.Permission, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *permissionService) Tree() ([]*models.Permission, error) {
	return s.repo.Tree()
}

type departmentService struct {
	repo repository.DepartmentRepository
}

func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {
	return &departmentService{repo}
}

func (s *departmentService) Create(department *models.Department) error {
	return s.repo.Create(department)
}

func (s *departmentService) Update(department *models.Department) error {
	return s.repo.Update(department)
}

func (s *departmentService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *departmentService) GetByID(id int64) (*models.Department, error) {
	return s.repo.GetByID(id)
}

func (s *departmentService) List(page, pageSize int) ([]*models.Department, int64, error) {
	return s.repo.List(page, pageSize)
}

func (s *departmentService) Tree() ([]*models.Department, error) {
	return s.repo.Tree()
}

type loginLogService struct {
	repo repository.LoginLogRepository
}

func NewLoginLogService(repo repository.LoginLogRepository) LoginLogService {
	return &loginLogService{repo}
}

func (s *loginLogService) List(page, pageSize int) ([]*models.LoginLog, int64, error) {
	return s.repo.List(page, pageSize)
}

type operationLogService struct {
	repo repository.OperationLogRepository
}

func NewOperationLogService(repo repository.OperationLogRepository) OperationLogService {
	return &operationLogService{repo}
}

func (s *operationLogService) List(page, pageSize int) ([]*models.OperationLog, int64, error) {
	return s.repo.List(page, pageSize)
}

type authService struct {
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	redis     *redis.Client
	cfg       *config.Config
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, redis *redis.Client, cfg *config.Config) AuthService {
	return &authService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		redis:     redis,
		cfg:       cfg,
	}
}

func (s *authService) Login(username, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("用户不存在")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", nil, errors.New("密码错误")
	}

	if user.Status != 1 {
		return "", nil, errors.New("账户已被禁用")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, s.cfg.JWT.Secret, s.cfg.JWT.ExpireTime)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) Logout(token string) error {
	ctx := context.Background()
	// 将token加入黑名单，设置过期时间为剩余时间
	return s.redis.Set(ctx, "blacklist:"+token, "1", s.cfg.JWT.ExpireTime*3600).Err()
}

func (s *authService) GetUserFromToken(token string) (*models.User, error) {
	claims, err := utils.ParseToken(token, s.cfg.JWT.Secret)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) HasPermission(userID int64, permissionCode string) (bool, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d:permissions", userID)

	// 先从Redis缓存获取
	cached, err := s.redis.Exists(ctx, cacheKey).Result()
	if err != nil {
		return false, err
	}

	var permissions []string
	if cached > 0 {
		// 从缓存读取
		results, err := s.redis.SMembers(ctx, cacheKey).Result()
		if err != nil {
			return false, err
		}
		permissions = results
	} else {
		// 从数据库查询
		perms, err := s.userRepo.GetPermissions(userID)
		if err != nil {
			return false, err
		}

		for _, perm := range perms {
			permissions = append(permissions, perm.Code)
		}

		// 写入缓存，过期时间1小时
		if len(permissions) > 0 {
			if err := s.redis.SAdd(ctx, cacheKey, permissions).Err(); err != nil {
				return false, err
			}
			s.redis.Expire(ctx, cacheKey, 3600)
		}
	}

	// 检查权限
	for _, perm := range permissions {
		if perm == permissionCode {
			return true, nil
		}
	}

	return false, nil
}

func (s *authService) GetPermissionCodes(userID int64) ([]string, error) {
	perms, err := s.userRepo.GetPermissions(userID)
	if err != nil {
		return nil, err
	}

	var codes []string
	for _, perm := range perms {
		codes = append(codes, perm.Code)
	}

	return codes, nil
}

type permissionChangeLogService struct {
	repo repository.PermissionChangeLogRepository
}

func NewPermissionChangeLogService(repo repository.PermissionChangeLogRepository) PermissionChangeLogService {
	return &permissionChangeLogService{repo}
}

func NewServices(repos *repository.Repositories, redis *redis.Client, cfg *config.Config) *Services {
	return &Services{
		User:                NewUserService(repos.User),
		Auth:                NewAuthService(repos.User, repos.Role, redis, cfg),
		Role:                NewRoleService(repos.Role),
		Permission:          NewPermissionService(repos.Permission),
		Department:           NewDepartmentService(repos.Department),
		LoginLog:            NewLoginLogService(repos.LoginLog),
		OperationLog:        NewOperationLogService(repos.OperationLog),
		PermissionChangeLog: NewPermissionChangeLogService(repos.PermissionChangeLog),
	}
}
