<template>
  <el-container class="layout">
    <el-aside width="240px">
      <div class="logo">权限管理</div>
      <el-menu :default-active="activeMenu" router>
        <template v-for="route in menuRoutes" :key="route.path">
          <el-sub-menu v-if="route.children" :index="route.path">
            <template #title>
              <el-icon v-if="route.meta?.icon"><component :is="route.meta.icon" /></el-icon>
              <span>{{ route.meta?.title }}</span>
            </template>
            <template v-for="child in route.children" :key="child.path">
              <el-sub-menu v-if="child.children" :index="child.path">
                <template #title>
                  <el-icon v-if="child.meta?.icon"><component :is="child.meta.icon" /></el-icon>
                  <span>{{ child.meta?.title }}</span>
                </template>
                <el-menu-item
                  v-for="sub in filterChildren(child.children)"
                  :key="sub.path"
                  :index="sub.path"
                >
                  {{ sub.meta?.title }}
                </el-menu-item>
              </el-sub-menu>
              <el-menu-item v-else-if="hasPermission(child.meta?.permission)" :index="child.path">
                {{ child.meta?.title }}
              </el-menu-item>
            </template>
          </el-sub-menu>
          <el-menu-item v-else-if="hasPermission(route.meta?.permission)" :index="route.path">
            <el-icon v-if="route.meta?.icon"><component :is="route.meta.icon" /></el-icon>
            <span>{{ route.meta?.title }}</span>
          </el-menu-item>
        </template>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header>
        <div class="header-content">
          <el-breadcrumb>
            <el-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">{{ item.meta?.title }}</el-breadcrumb-item>
          </el-breadcrumb>
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="userStore.user?.avatar">{{ userStore.user?.realName?.[0] }}</el-avatar>
              <span>{{ userStore.user?.realName }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const menuRoutes = computed(() => router.options.routes.find(r => r.path === '/')?.children || [])
const activeMenu = computed(() => route.path)

const breadcrumbs = computed(() => {
  const matched = route.matched.filter(r => r.meta?.title)
  return matched.length > 0 ? matched : [{ meta: { title: '首页' }, path: '/' }]
})

const hasPermission = (permission?: string) => {
  if (!permission) return true
  return userStore.hasPermission(permission)
}

const filterChildren = (children: any[]) => {
  return children.filter((child: any) => hasPermission(child.meta?.permission))
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
  }
}
</script>

<style scoped>
.layout {
  height: 100%;
}

.el-aside {
  background: #304156;
  overflow-y: auto;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
}

.el-menu {
  border: none;
  background: #304156;
}

:deep(.el-menu-item), :deep(.el-sub-menu__title) {
  color: #bfcbd9;
}

:deep(.el-menu-item:hover), :deep(.el-sub-menu__title:hover) {
  background: #263445;
}

:deep(.el-menu-item.is-active) {
  background: #409eff !important;
  color: #fff;
}

.el-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.el-main {
  background: #f0f2f5;
}
</style>
