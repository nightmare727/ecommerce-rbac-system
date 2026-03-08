<template>
  <el-card>
    <template #header>
      <span>登录日志</span>
    </template>
    <el-table :data="logs" stripe v-loading="loading">
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="ip" label="IP地址" width="140" />
      <el-table-column prop="location" label="登录位置" />
      <el-table-column prop="browser" label="浏览器" />
      <el-table-column prop="os" label="操作系统" width="120" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="message" label="消息" />
      <el-table-column prop="loginTime" label="登录时间" width="180" />
    </el-table>

    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      style="margin-top: 16px; justify-content: flex-end;"
      @size-change="fetchData"
      @current-change="fetchData"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { loginLogApi } from '@/api'

const logs = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const fetchData = async () => {
  loading.value = true
  try {
    const res: any = await loginLogApi.list({ page: page.value, pageSize: pageSize.value })
    logs.value = res.data || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>
