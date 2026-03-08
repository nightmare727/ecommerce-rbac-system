<template>
  <el-card>
    <template #header>
      <span>操作日志</span>
    </template>
    <el-table :data="logs" stripe v-loading="loading">
      <el-table-column prop="username" label="操作人" width="120" />
      <el-table-column prop="module" label="模块" width="100" />
      <el-table-column prop="operation" label="操作" width="100" />
      <el-table-column prop="method" label="方法" width="80" />
      <el-table-column prop="url" label="请求URL" />
      <el-table-column prop="ip" label="IP" width="140" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="costTime" label="耗时(ms)" width="100" />
      <el-table-column prop="createdAt" label="操作时间" width="180" />
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
import { operationLogApi } from '@/api'

const logs = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const fetchData = async () => {
  loading.value = true
  try {
    const res: any = await operationLogApi.list({ page: page.value, pageSize: pageSize.value })
    logs.value = res.data || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>
