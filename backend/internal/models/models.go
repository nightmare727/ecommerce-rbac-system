package models

import (
	"time"
)

type User struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password    string    `json:"-" gorm:"size:255;not null"`
	RealName    string    `json:"realName" gorm:"size:50;not null"`
	Email       string    `json:"email" gorm:"uniqueIndex;size:100"`
	Phone       string    `json:"phone" gorm:"size:20"`
	Avatar      string    `json:"avatar" gorm:"size:500"`
	DepartmentID *int64   `json:"departmentId" gorm:"index"`
	Status      int       `json:"status" gorm:"default:1;comment:1:启用 0:禁用 2:锁定"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	LastLoginIP string    `json:"lastLoginIp" gorm:"size:50"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// 关联
	Department *Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	Roles      []Role      `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

type Role struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:50;not null"`
	Code        string    `json:"uniqueIndex" gorm:"uniqueIndex;size:50;not null"`
	Description string    `json:"description" gorm:"type:text"`
	DataScope   int       `json:"dataScope" gorm:"default:1;comment:1:全部 2:自定义 3:本部门及以下 4:本部门 5:仅本人"`
	Status      int       `json:"status" gorm:"default:1;comment:1:启用 0:禁用"`
	SortOrder   int       `json:"sortOrder" gorm:"default:0"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// 关联
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
}

type Permission struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID  *int64     `json:"parentId" gorm:"index"`
	Name      string     `json:"name" gorm:"size:50;not null"`
	Code      string     `json:"code" gorm:"uniqueIndex;size:100;not null"`
	Type      int        `json:"type" gorm:"default:1;comment:1:目录 2:菜单 3:按钮"`
	Path      string     `json:"path" gorm:"size:200"`
	Component string     `json:"component" gorm:"size:200"`
	Icon      string     `json:"icon" gorm:"size:50"`
	SortOrder int        `json:"sortOrder" gorm:"default:0"`
	Visible   int        `json:"visible" gorm:"default:1;comment:1:显示 0:隐藏"`
	Status    int        `json:"status" gorm:"default:1;comment:1:启用 0:禁用"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// 关联
	Children    []Permission `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Roles       []Role       `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

type Department struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"size:100;not null"`
	Code      string     `json:"code" gorm:"uniqueIndex;size:50;not null"`
	ParentID  *int64     `json:"parentId" gorm:"index"`
	SortOrder int        `json:"sortOrder" gorm:"default:0"`
	Status    int        `json:"status" gorm:"default:1;comment:1:启用 0:禁用"`
	LeaderID  *int64     `json:"leaderId"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// 关联
	Children    []Department `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Users       []User       `json:"users,omitempty" gorm:"foreignKey:DepartmentID"`
}

type LoginLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    *int64    `json:"userId" gorm:"index"`
	Username  string    `json:"username" gorm:"size:50"`
	IP        string    `json:"ip" gorm:"size:50"`
	Location  string    `json:"location" gorm:"size:100"`
	Browser   string    `json:"browser" gorm:"size:100"`
	OS        string    `json:"os" gorm:"size:50"`
	Status    int       `json:"status" gorm:"comment:1:成功 0:失败"`
	Message   string    `json:"message" gorm:"size:200"`
	LoginTime time.Time `json:"loginTime" gorm:"autoCreateTime"`
}

type OperationLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    *int64    `json:"userId" gorm:"index"`
	Username  string    `json:"username" gorm:"size:50"`
	Module    string    `json:"module" gorm:"size:50"`
	Operation string    `json:"operation" gorm:"size:50"`
	Method    string    `json:"method" gorm:"size:10"`
	URL       string    `json:"url" gorm:"size:200"`
	IP        string    `json:"ip" gorm:"size:50"`
	Params    string    `json:"params" gorm:"type:text"`
	Result    string    `json:"result" gorm:"type:text"`
	Status    int       `json:"status" gorm:"comment:1:成功 0:失败"`
	ErrorMsg  string    `json:"errorMsg" gorm:"type:text"`
	CostTime  int       `json:"costTime" gorm:"comment:耗时(ms)"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

type PermissionChangeLog struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         *int64    `json:"userId" gorm:"index"`
	OperatorID     *int64    `json:"operatorId" gorm:"index"`
	TargetType     string    `json:"targetType" gorm:"size:20;comment:user|role"`
	TargetID       int64     `json:"targetId"`
	ChangeType     string    `json:"changeType" gorm:"size:20;comment:add|remove"`
	PermissionCode string    `json:"permissionCode" gorm:"size:100"`
	Reason         string    `json:"reason" gorm:"type:text"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
