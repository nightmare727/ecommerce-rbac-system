-- 企业级电商权限管理系统 - 数据库脚本
-- PostgreSQL

-- 1. 部门表
CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    parent_id BIGINT REFERENCES departments(id) ON DELETE SET NULL,
    sort_order INT DEFAULT 0,
    status SMALLINT DEFAULT 1 COMMENT '1:启用 0:禁用',
    leader_id BIGINT COMMENT '部门负责人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 2. 用户表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    real_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    avatar VARCHAR(500),
    department_id BIGINT REFERENCES departments(id) ON DELETE SET NULL,
    status SMALLINT DEFAULT 1 COMMENT '1:启用 0:禁用 2:锁定',
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 3. 角色表
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    data_scope SMALLINT DEFAULT 1 COMMENT '1:全部 2:自定义 3:本部门及以下 4:本部门 5:仅本人',
    status SMALLINT DEFAULT 1 COMMENT '1:启用 0:禁用',
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 4. 用户-角色关联表（支持多角色）
CREATE TABLE user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- 5. 权限（菜单）表
CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT REFERENCES permissions(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(100) UNIQUE NOT NULL COMMENT '权限标识，如:user:list',
    type SMALLINT DEFAULT 1 COMMENT '1:目录 2:菜单 3:按钮',
    path VARCHAR(200) COMMENT '路由路径',
    component VARCHAR(200) COMMENT '组件路径',
    icon VARCHAR(50),
    sort_order INT DEFAULT 0,
    visible SMALLINT DEFAULT 1 COMMENT '1:显示 0:隐藏',
    status SMALLINT DEFAULT 1 COMMENT '1:启用 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 6. 角色-权限关联表
CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

-- 7. 数据权限规则表
CREATE TABLE data_scope_rules (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    department_id BIGINT REFERENCES departments(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. 登录日志表
CREATE TABLE login_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    username VARCHAR(50),
    ip VARCHAR(50),
    location VARCHAR(100),
    browser VARCHAR(100),
    os VARCHAR(50),
    status SMALLINT COMMENT '1:成功 0:失败',
    message VARCHAR(200),
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 9. 操作日志表
CREATE TABLE operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    username VARCHAR(50),
    module VARCHAR(50),
    operation VARCHAR(50),
    method VARCHAR(10),
    url VARCHAR(200),
    ip VARCHAR(50),
    params TEXT,
    result TEXT,
    status SMALLINT COMMENT '1:成功 0:失败',
    error_msg TEXT,
    cost_time INT COMMENT '耗时(ms)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. 权限变更日志表
CREATE TABLE permission_change_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    operator_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    target_type VARCHAR(20) COMMENT 'user|role',
    target_id BIGINT,
    change_type VARCHAR(20) COMMENT 'add|remove',
    permission_code VARCHAR(100),
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_department_id ON users(department_id);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_parent_id ON permissions(parent_id);
CREATE INDEX idx_login_logs_user_id ON login_logs(user_id);
CREATE INDEX idx_login_logs_login_time ON login_logs(login_time);
CREATE INDEX idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX idx_operation_logs_created_at ON operation_logs(created_at);

-- 初始化默认数据

-- 默认超级管理员密码：admin123 (BCrypt加密)
INSERT INTO users (username, password, real_name, email, status) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '超级管理员', 'admin@example.com', 1);

-- 默认部门
INSERT INTO departments (name, code, sort_order) VALUES
('总公司', 'ROOT', 0),
('技术部', 'TECH', 1),
('运营部', 'OPS', 2),
('财务部', 'FINANCE', 3);

-- 更新部门层级
UPDATE departments SET parent_id = 1 WHERE id IN (2, 3, 4);

-- 默认角色
INSERT INTO roles (name, code, description, data_scope) VALUES
('超级管理员', 'SUPER_ADMIN', '系统超级管理员，拥有所有权限', 1),
('管理员', 'ADMIN', '系统管理员', 3),
('普通用户', 'USER', '普通用户', 4);

-- 给超级管理员分配超级管理员角色
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);

-- 默认权限树
-- 系统管理
INSERT INTO permissions (parent_id, name, code, type, path, icon, sort_order) VALUES
(0, '系统管理', 'system', 1, '/system', 'Setting', 1);

-- 用户管理
INSERT INTO permissions (parent_id, name, code, type, path, component, sort_order) VALUES
(1, '用户管理', 'system:user', 2, '/system/user', 'system/User', 1);
INSERT INTO permissions (parent_id, name, code, type, sort_order) VALUES
(2, '查看用户', 'system:user:list', 3, 1),
(2, '新增用户', 'system:user:add', 3, 2),
(2, '编辑用户', 'system:user:edit', 3, 3),
(2, '删除用户', 'system:user:delete', 3, 4),
(2, '分配角色', 'system:user:assign', 3, 5);

-- 角色管理
INSERT INTO permissions (parent_id, name, code, type, path, component, sort_order) VALUES
(1, '角色管理', 'system:role', 2, '/system/role', 'system/Role', 2);
INSERT INTO permissions (parent_id, name, code, type, sort_order) VALUES
(7, '查看角色', 'system:role:list', 3, 1),
(7, '新增角色', 'system:role:add', 3, 2),
(7, '编辑角色', 'system:role:edit', 3, 3),
(7, '删除角色', 'system:role:delete', 3, 4),
(7, '分配权限', 'system:role:assign', 3, 5);

-- 权限管理
INSERT INTO permissions (parent_id, name, code, type, path, component, sort_order) VALUES
(1, '权限管理', 'system:permission', 2, '/system/permission', 'system/Permission', 3);
INSERT INTO permissions (parent_id, name, code, type, sort_order) VALUES
(12, '查看权限', 'system:permission:list', 3, 1),
(12, '新增权限', 'system:permission:add', 3, 2),
(12, '编辑权限', 'system:permission:edit', 3, 3),
(12, '删除权限', 'system:permission:delete', 3, 4);

-- 部门管理
INSERT INTO permissions (parent_id, name, code, type, path, component, sort_order) VALUES
(1, '部门管理', 'system:dept', 2, '/system/dept', 'system/Dept', 4);
INSERT INTO permissions (parent_id, name, code, type, sort_order) VALUES
(16, '查看部门', 'system:dept:list', 3, 1),
(16, '新增部门', 'system:dept:add', 3, 2),
(16, '编辑部门', 'system:dept:edit', 3, 3),
(16, '删除部门', 'system:dept:delete', 3, 4);

-- 日志管理
INSERT INTO permissions (parent_id, name, code, type, path, component, sort_order) VALUES
(1, '日志管理', 'system:log', 2, '/system/log', 'system/Log', 5);
INSERT INTO permissions (parent_id, name, code, type, sort_order) VALUES
(20, '登录日志', 'system:log:login', 3, 1),
(20, '操作日志', 'system:log:operation', 3, 2);

-- 给超级管理员分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE deleted_at IS NULL;

COMMENT ON TABLE departments IS '部门表';
COMMENT ON TABLE users IS '用户表';
COMMENT ON TABLE roles IS '角色表';
COMMENT ON TABLE user_roles IS '用户角色关联表';
COMMENT ON TABLE permissions IS '权限表';
COMMENT ON TABLE role_permissions IS '角色权限关联表';
COMMENT ON TABLE data_scope_rules IS '数据权限规则表';
COMMENT ON TABLE login_logs IS '登录日志表';
COMMENT ON TABLE operation_logs IS '操作日志表';
COMMENT ON TABLE permission_change_logs IS '权限变更日志表';
