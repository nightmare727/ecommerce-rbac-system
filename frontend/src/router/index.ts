import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页', icon: 'House' }
      },
      {
        path: 'system',
        name: 'System',
        meta: { title: '系统管理', icon: 'Setting' },
        children: [
          {
            path: 'user',
            name: 'User',
            component: () => import('@/views/system/User.vue'),
            meta: { title: '用户管理', permission: 'system:user:list' }
          },
          {
            path: 'role',
            name: 'Role',
            component: () => import('@/views/system/Role.vue'),
            meta: { title: '角色管理', permission: 'system:role:list' }
          },
          {
            path: 'permission',
            name: 'Permission',
            component: () => import('@/views/system/Permission.vue'),
            meta: { title: '权限管理', permission: 'system:permission:list' }
          },
          {
            path: 'dept',
            name: 'Dept',
            component: () => import('@/views/system/Dept.vue'),
            meta: { title: '部门管理', permission: 'system:dept:list' }
          },
          {
            path: 'log',
            name: 'Log',
            meta: { title: '日志管理', icon: 'Document' },
            children: [
              {
                path: 'login',
                name: 'LoginLog',
                component: () => import('@/views/system/LoginLog.vue'),
                meta: { title: '登录日志', permission: 'system:log:login' }
              },
              {
                path: 'operation',
                name: 'OperationLog',
                component: () => import('@/views/system/OperationLog.vue'),
                meta: { title: '操作日志', permission: 'system:log:operation' }
              }
            ]
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.path === '/login') {
    if (userStore.isLoggedIn) {
      next('/')
    } else {
      next()
    }
    return
  }

  if (!userStore.isLoggedIn) {
    next('/login')
    return
  }

  // 权限检查
  if (to.meta.permission) {
    if (!userStore.hasPermission(to.meta.permission as string)) {
      ElMessage.error('权限不足')
      next(from.path)
      return
    }
  }

  next()
})

export default router
