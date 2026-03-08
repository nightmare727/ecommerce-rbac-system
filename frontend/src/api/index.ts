import request from '@/utils/request'

export interface User {
  id: number
  username: string
  realName: string
  email?: string
  phone?: string
  avatar?: string
  departmentId?: number
  status: number
  department?: Department
  roles?: Role[]
}

export interface Role {
  id: number
  name: string
  code: string
  description?: string
  dataScope: number
  status: number
  sortOrder: number
  permissions?: Permission[]
}

export interface Permission {
  id: number
  parentId?: number
  name: string
  code: string
  type: number
  path?: string
  component?: string
  icon?: string
  sortOrder: number
  visible: number
  status: number
  children?: Permission[]
}

export interface Department {
  id: number
  name: string
  code: string
  parentId?: number
  sortOrder: number
  status: number
  leaderId?: number
  children?: Department[]
}

export const authApi = {
  login: (data: { username: string; password: string }) => 
    request.post('/auth/login', data),
  logout: () => 
    request.post('/auth/logout')
}

export const userApi = {
  getInfo: () => 
    request.get('/user/info'),
  updateInfo: (data: Partial<User>) => 
    request.put('/user/info', data),
  list: (params: { page: number; pageSize: number; keyword?: string }) => 
    request.get('/users', { params }),
  getById: (id: number) => 
    request.get(`/users/${id}`),
  create: (data: Partial<User> & { password: string }) => 
    request.post('/users', data),
  update: (id: number, data: Partial<User>) => 
    request.put(`/users/${id}`, data),
  delete: (id: number) => 
    request.delete(`/users/${id}`),
  assignRoles: (id: number, roleIds: number[]) => 
    request.post(`/users/${id}/roles`, { roleIds })
}

export const roleApi = {
  list: (params: { page: number; pageSize: number; keyword?: string }) => 
    request.get('/roles', { params }),
  getById: (id: number) => 
    request.get(`/roles/${id}`),
  create: (data: Partial<Role>) => 
    request.post('/roles', data),
  update: (id: number, data: Partial<Role>) => 
    request.put(`/roles/${id}`, data),
  delete: (id: number) => 
    request.delete(`/roles/${id}`),
  assignPermissions: (id: number, permissionIds: number[]) => 
    request.post(`/roles/${id}/permissions`, { permissionIds })
}

export const permissionApi = {
  list: (params: { page: number; pageSize: number; keyword?: string }) => 
    request.get('/permissions', { params }),
  tree: () => 
    request.get('/permissions/tree'),
  getById: (id: number) => 
    request.get(`/permissions/${id}`),
  create: (data: Partial<Permission>) => 
    request.post('/permissions', data),
  update: (id: number, data: Partial<Permission>) => 
    request.put(`/permissions/${id}`, data),
  delete: (id: number) => 
    request.delete(`/permissions/${id}`)
}

export const departmentApi = {
  list: (params: { page: number; pageSize: number }) => 
    request.get('/departments', { params }),
  tree: () => 
    request.get('/departments/tree'),
  getById: (id: number) => 
    request.get(`/departments/${id}`),
  create: (data: Partial<Department>) => 
    request.post('/departments', data),
  update: (id: number, data: Partial<Department>) => 
    request.put(`/departments/${id}`, data),
  delete: (id: number) => 
    request.delete(`/departments/${id}`)
}

export const loginLogApi = {
  list: (params: { page: number; pageSize: number }) => 
    request.get('/login-logs', { params })
}

export const operationLogApi = {
  list: (params: { page: number; pageSize: number }) => 
    request.get('/operation-logs', { params })
}
