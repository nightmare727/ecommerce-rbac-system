import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userApi, type User } from '@/api'
import { ElMessage } from 'element-plus'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(JSON.parse(localStorage.getItem('user') || 'null'))
  const permissions = ref<string[]>(JSON.parse(localStorage.getItem('permissions') || '[]'))

  const isLoggedIn = computed(() => !!token.value)

  const login = async (username: string, password: string) => {
    try {
      const res: any = await authApi.login({ username, password })
      token.value = res.token
      user.value = res.user
      permissions.value = res.permissions || []

      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(res.user))
      localStorage.setItem('permissions', JSON.stringify(permissions.value))

      ElMessage.success('登录成功')
      router.push('/')
    } catch (error) {
      ElMessage.error('登录失败')
      throw error
    }
  }

  const fetchUserInfo = async () => {
    try {
      const res: any = await userApi.getInfo()
      user.value = res.user
      permissions.value = res.permissions || []
      localStorage.setItem('user', JSON.stringify(res.user))
      localStorage.setItem('permissions', JSON.stringify(permissions.value))
    } catch {
      logout()
    }
  }

  const logout = async () => {
    try {
      await authApi.logout()
    } catch {
      // 忽略登出错误
    }
    token.value = ''
    user.value = null
    permissions.value = []
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('permissions')
    router.push('/login')
  }

  const hasPermission = (code: string): boolean => {
    return permissions.value.includes(code)
  }

  const hasAnyPermission = (codes: string[]): boolean => {
    return codes.some(code => hasPermission(code))
  }

  return {
    token,
    user,
    permissions,
    isLoggedIn,
    login,
    logout,
    fetchUserInfo,
    hasPermission,
    hasAnyPermission
  }
})
