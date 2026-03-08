<template>
  <el-row :gutter="20">
    <el-col :span="6" v-for="card in cards" :key="card.title">
      <el-card class="stat-card">
        <el-statistic :title="card.title" :value="card.value">
          <template #suffix>
            <el-icon :size="20" :color="card.color"><component :is="card.icon" /></el-icon>
          </template>
        </el-statistic>
      </el-card>
    </el-col>
  </el-row>

  <el-row :gutter="20" style="margin-top: 20px;">
    <el-col :span="12">
      <el-card>
        <template #header>
          <span>最近登录</span>
        </template>
        <el-table :data="loginLogs" stripe>
          <el-table-column prop="username" label="用户名" />
          <el-table-column prop="ip" label="IP地址" />
          <el-table-column prop="location" label="位置" />
          <el-table-column prop="loginTime" label="登录时间" />
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </el-col>

    <el-col :span="12">
      <el-card>
        <template #header>
          <span>系统信息</span>
        </template>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="系统版本">v1.0.0</el-descriptions-item>
          <el-descriptions-item label="框架">Go + Vue3</el-descriptions-item>
          <el-descriptions-item label="数据库">PostgreSQL</el-descriptions-item>
          <el-descriptions-item label="缓存">Redis</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { loginLogApi } from '@/api'

const cards = ref([
  { title: '用户总数', value: 100, icon: 'User', color: '#409eff' },
  { title: '角色总数', value: 10, icon: 'UserFilled', color: '#67c23a' },
  { title: '权限总数', value: 50, icon: 'Lock', color: '#e6a23c' },
  { title: '今日登录', value: 25, icon: 'Calendar', color: '#f56c6c' }
])

const loginLogs = ref([])

onMounted(async () => {
  try {
    const res: any = await loginLogApi.list({ page: 1, pageSize: 5 })
    loginLogs.value = res.data || []
  } catch (error) {
    // 加载失败
  }
})
</script>

<style scoped>
.stat-card {
  margin-bottom: 20px;
}
</:deep/>
