<template>
  <div class="credit-container">
    <el-page-header @back="() => {}" content="我的学分" />
    <el-card style="margin-top: 20px;">
      <el-button
        type="primary"
        @click="loadCredits"
        :loading="tableLoading"
        style="margin-bottom: 16px;"
      >
        刷新列表
      </el-button>
      <el-table
        :data="creditList"
        border
        stripe
        style="width: 100%"
        v-loading="tableLoading"
        empty-text="暂无学分记录，请联系教师录入学分后在列表中查看"
      >
        <el-table-column prop="course_name" label="课程名称" min-width="160" />
        <el-table-column prop="score" label="成绩" width="90" />
        <el-table-column label="授课教师" min-width="140">
          <template #default="scope">
            {{ formatAddress(scope.row.teacher_address) }}
          </template>
        </el-table-column>
        <el-table-column label="审核状态" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.status === 'pending'" type="warning">待审核</el-tag>
            <el-tag v-else-if="scope.row.status === 'approved'" type="success">已通过</el-tag>
            <el-tag v-else type="danger">已驳回</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="录入时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCreditList } from '@/api/credit'
import { ElMessage } from 'element-plus'

const tableLoading = ref(false)
const creditList = ref([])

function formatAddress(addr) {
  if (!addr) return '-'
  const s = String(addr)
  if (s.length <= 14) return s
  return s.slice(0, 6) + '...' + s.slice(-4)
}

function formatTime(t) {
  if (!t) return '-'
  try {
    return new Date(t).toLocaleString('zh-CN')
  } catch {
    return t
  }
}

const loadCredits = async () => {
  tableLoading.value = true
  try {
    const res = await getCreditList()
    creditList.value = Array.isArray(res?.data) ? res.data : []
  } catch (error) {
    console.error('获取学分失败：', error)
    creditList.value = []
    ElMessage.error('获取学分失败，请重试')
  } finally {
    tableLoading.value = false
  }
}

onMounted(() => {
  loadCredits()
})
</script>

<style scoped>
.credit-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  height: 100%;
}
</style>
