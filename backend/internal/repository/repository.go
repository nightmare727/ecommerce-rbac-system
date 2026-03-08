package repository

import (
	"ecommerce-rbac-system/internal/models"
	"gorm.io/gorm"
)

type Repositories struct {
	User                UserRepository
	Role                RoleRepository
	Permission          PermissionRepository
	Department          DepartmentRepository
	LoginLog            LoginLogRepository
	OperationLog        OperationLogRepository
	PermissionChangeLog PermissionChangeLogRepository
}

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int64) error
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	List(page, pageSize int, keyword string) ([]*models.User, int64, error)
	AssignRoles(userID int64, roleIDs []int64) error
	GetPermissions(userID int64) ([]*models.Permission, error)
}

type RoleRepository interface {
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id int64) error
	GetByID(id int64) (*models.Role, error)
	GetByCode(code string) (*models.Role, error)
	List(page, pageSize int, keyword string) ([]*models.Role, int64, error)
	AssignPermissions(roleID int64, permissionIDs []int64) error
}

type PermissionRepository interface {
	Create(permission *models.Permission) error
	Update(permission *models.Permission) error
	Delete(id int64) error
	GetByID(id int64) (*models.Permission, error)
	GetByCode(code string) (*models.Permission, error)
	List(page, pageSize int, keyword string) ([]*models.Permission, int64, error)
	Tree() ([]*models.Permission, error)
}

type DepartmentRepository interface {
	Create(department *models.Department) error
	Update(department *models.Department) error
	Delete(id int64) error
	GetByID(id int64) (*models.Department, error)
	List(page, pageSize int) ([]*models.Department, int64, error)
	Tree() ([]*models.Department, error)
}

type LoginLogRepository interface {
	Create(log *models.LoginLog) error
	List(page, pageSize int) ([]*models.LoginLog, int64, error)
}

type OperationLogRepository interface {
	Create(log *models.OperationLog) error
	List(page, pageSize int) ([]*models.OperationLog, int64, error)
}

type PermissionChangeLogRepository interface {
	Create(log *models.PermissionChangeLog) error
}

// --- User ---

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int64) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Department").Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(page, pageSize int, keyword string) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.Model(&models.User{}).Preload("Department").Preload("Roles")
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *userRepository) AssignRoles(userID int64, roleIDs []int64) error {
	tx := r.db.Begin()
	if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, roleID := range roleIDs {
		if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *userRepository) GetPermissions(userID int64) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.Raw(`
		SELECT DISTINCT p.* FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND p.deleted_at IS NULL
		ORDER BY p.sort_order
	`, userID).Scan(&permissions).Error
	return permissions, err
}

// --- Role ---

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id int64) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *roleRepository) GetByID(id int64) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByCode(code string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) List(page, pageSize int, keyword string) ([]*models.Role, int64, error) {
	var roles []*models.Role
	var total int64

	query := r.db.Model(&models.Role{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) AssignPermissions(roleID int64, permissionIDs []int64) error {
	tx := r.db.Begin()
	if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, permissionID := range permissionIDs {
		if err := tx.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permissionID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// --- Permission ---

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db}
}

func (r *permissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id int64) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

func (r *permissionRepository) GetByID(id int64) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByCode(code string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) List(page, pageSize int, keyword string) ([]*models.Permission, int64, error) {
	var permissions []*models.Permission
	var total int64

	query := r.db.Model(&models.Permission{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) Tree() ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.Order("sort_order").Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	permissionMap := make(map[int64]*models.Permission)
	var roots []*models.Permission

	for _, perm := range permissions {
		perm.Children = nil
		permissionMap[perm.ID] = perm
	}

	for _, perm := range permissions {
		if perm.ParentID != nil && *perm.ParentID != 0 {
			if parent, exists := permissionMap[*perm.ParentID]; exists {
				parent.Children = append(parent.Children, perm)
			}
		} else {
			roots = append(roots, perm)
		}
	}

	return roots, nil
}

// --- Department ---

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db}
}

func (r *departmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

func (r *departmentRepository) Update(department *models.Department) error {
	return r.db.Save(department).Error
}

func (r *departmentRepository) Delete(id int64) error {
	return r.db.Delete(&models.Department{}, id).Error
}

func (r *departmentRepository) GetByID(id int64) (*models.Department, error) {
	var department models.Department
	err := r.db.First(&department, id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) List(page, pageSize int) ([]*models.Department, int64, error) {
	var departments []*models.Department
	var total int64

	query := r.db.Model(&models.Department{})
	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&departments).Error
	return departments, total, err
}

func (r *departmentRepository) Tree() ([]*models.Department, error) {
	var departments []*models.Department
	err := r.db.Order("sort_order").Find(&departments).Error
	if err != nil {
		return nil, err
	}

	deptMap := make(map[int64]*models.Department)
	var roots []*models.Department

	for _, dept := range departments {
		dept.Children = nil
		deptMap[dept.ID] = dept
	}

	for _, dept := range departments {
		if dept.ParentID != nil && *dept.ParentID != 0 {
			if parent, exists := deptMap[*dept.ParentID]; exists {
				parent.Children = append(parent.Children, dept)
			}
		} else {
			roots = append(roots, dept)
		}
	}

	return roots, nil
}

// --- Logs ---

type loginLogRepository struct {
	db *gorm.DB
}

func NewLoginLogRepository(db *gorm.DB) LoginLogRepository {
	return &loginLogRepository{db}
}

func (r *loginLogRepository) Create(log *models.LoginLog) error {
	return r.db.Create(log).Error
}

func (r *loginLogRepository) List(page, pageSize int) ([]*models.LoginLog, int64, error) {
	var logs []*models.LoginLog
	var total int64

	r.db.Model(&models.LoginLog{}).Count(&total)
	err := r.db.Order("login_time DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

type operationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) OperationLogRepository {
	return &operationLogRepository{db}
}

func (r *operationLogRepository) Create(log *models.OperationLog) error {
	return r.db.Create(log).Error
}

func (r *operationLogRepository) List(page, pageSize int) ([]*models.OperationLog, int64, error) {
	var logs []*models.OperationLog
	var total int64

	r.db.Model(&models.OperationLog{}).Count(&total)
	err := r.db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

type permissionChangeLogRepository struct {
	db *gorm.DB
}

func NewPermissionChangeLogRepository(db *gorm.DB) PermissionChangeLogRepository {
	return &permissionChangeLogRepository{db}
}

func (r *permissionChangeLogRepository) Create(log *models.PermissionChangeLog) error {
	return r.db.Create(log).Error
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:                NewUserRepository(db),
		Role:                NewRoleRepository(db),
		Permission:          NewPermissionRepository(db),
		Department:          NewDepartmentRepository(db),
		LoginLog:            NewLoginLogRepository(db),
		OperationLog:        NewOperationLogRepository(db),
		PermissionChangeLog: NewPermissionChangeLogRepository(db),
	}
}
