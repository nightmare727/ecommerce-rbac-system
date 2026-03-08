# 电商权限管理系统

企业级电商权限管理系统，基于 RBAC 模型实现。

## 技术栈

### 后端
- Go 1.21+
- Gin 框架
- PostgreSQL (数据库)
- Redis (缓存)
- GORM (ORM)
- JWT (认证)

### 前端
- Vue 3
- TypeScript
- Vite
- Element Plus
- Pinia (状态管理)
- Vue Router
- Axios (HTTP 请求)

## 功能特性

### 核心功能
- ✅ 用户管理 (CRUD)
- ✅ 角色管理 (CRUD)
- ✅ 权限管理 (菜单权限、按钮权限)
- ✅ 部门管理 (组织架构树)
- ✅ RBAC 权限模型
- ✅ 登录日志
- ✅ 操作日志
- ✅ 权限变更日志

### 安全特性
- JWT 认证
- 权限拦截器
- 密码加密
- Token 黑名单
- 权限缓存

## 项目结构

```
ecommerce-rbac-system/
├── backend/              # 后端代码
│   ├── main.go           # 入口文件
│   ├── config.yaml       # 配置文件
│   ├── internal/         # 内部代码
│   │   ├── config/       # 配置加载
│   │   ├── database/     # 数据库连接
│   │   ├── models/       # 数据模型
│   │   ├── repository/   # 数据访问层
│   │   ├── service/      # 业务逻辑层
│   │   ├── handler/      # HTTP 处理器
│   │   ├── middleware/   # 中间件
│   │   └── utils/        # 工具函数
│   └── go.mod            # Go 模块文件
├── frontend/             # 前端代码
│   ├── src/
│   │   ├── api/          # API 接口
│   │   ├── components/   # 组件
│   │   ├── views/        # 页面
│   │   ├── router/       # 路由
│   │   ├── stores/       # 状态管理
│   │   └── utils/        # 工具函数
│   ├── package.json      # 依赖配置
│   └── vite.config.ts    # Vite 配置
└── docs/                 # 文档
    └── database.sql      # 数据库初始化脚本
```

## 快速开始

### 1. 环境准备

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Redis 6+

### 2. 数据库初始化

```bash
# 创建数据库
createdb ecommerce_rbac

# 导入表结构和初始数据
psql ecommerce_rbac < docs/database.sql
```

默认管理员账号：
- 用户名：`admin`
- 密码：`admin123`

### 3. 后端启动

```bash
cd backend

# 安装依赖
go mod download

# 修改配置文件 config.yaml（数据库、Redis等）

# 启动服务
go run main.go
```

后端服务运行在：`http://localhost:8080`

### 4. 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务运行在：`http://localhost:5173`

## API 接口文档

### 认证相关
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/logout` - 登出

### 用户管理
- `GET /api/v1/users` - 用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `POST /api/v1/users` - 创建用户
- `PUT /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户
- `POST /api/v1/users/:id/roles` - 分配角色

### 角色管理
- `GET /api/v1/roles` - 角色列表
- `GET /api/v1/roles/:id` - 获取角色详情
- `POST /api/v1/roles` - 创建角色
- `PUT /api/v1/roles/:id` - 更新角色
- `DELETE /api/v1/roles/:id` - 删除角色
- `POST /api/v1/roles/:id/permissions` - 分配权限

### 权限管理
- `GET /api/v1/permissions` - 权限列表
- `GET /api/v1/permissions/tree` - 权限树
- `GET /api/v1/permissions/:id` - 获取权限详情

### 部门管理
- `GET /api/v1/departments` - 部门列表
- `GET /api/v1/departments/tree` - 部门树
- `GET /api/v1/departments/:id` - 获取部门详情

### 日志管理
- `GET /api/v1/login-logs` - 登录日志
- `GET /api/v1/operation-logs` - 操作日志

## 数据权限说明

系统支持 5 种数据权限范围：

1. **全部数据权限**：可以查看所有数据
2. **自定义数据权限**：可以查看指定部门的数据
3. **本部门及以下数据权限**：可以查看本部门及子部门的数据
4. **本部门数据权限**：只能查看本部门的数据
5. **仅本人数据权限**：只能查看自己的数据

## 权限编码规范

权限编码采用 `模块:资源:操作` 格式，例如：
- `system:user:list` - 系统模块用户列表权限
- `system:user:add` - 系统模块用户新增权限
- `system:user:edit` - 系统模块用户编辑权限
- `system:user:delete` - 系统模块用户删除权限

## 部署建议

### 生产环境配置

1. **配置文件**
   - 修改 JWT secret 为强密钥
   - 启用 HTTPS
   - 配置正确的数据库连接池
Redis 持久化

2. **安全加固**
   - 限制数据库访问权限
   - 启用 Redis 密码
   - 配置防火墙规则
   - 定期备份数据库

3. **性能优化**
   - 启用 Redis 缓存
   - 配置数据库索引
   - 使用 CDN 加载静态资源

## 许可证

MIT License
