<template>
  <div class="credit-list-container">
    <el-page-header @back="() => {}" content="录入列表" />
    <el-card style="margin-top: 20px;">
      <el-table
        :data="creditList"
        border
        stripe
        style="width: 100%"
        v-loading="tableLoading"
        empty-text="暂无录入记录，请先在「学分录入」提交"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="学生地址/学号" min-width="160">
          <template #default="scope">
            {{ formatAddress(scope.row.student_address) }}
          </template>
        </el-table-column>
        <el-table-column prop="course_name" label="课程名称" min-width="140" />
        <el-table-column prop="score" label="成绩" width="90" />
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

const loadList = async () => {
  tableLoading.value = true
  try {
    const res = await getCreditList()
    creditList.value = Array.isArray(res?.data) ? res.data : []
  } catch (error) {
    console.error('获取录入列表失败：', error)
    creditList.value = []
  } finally {
    tableLoading.value = false
  }
}

onMounted(() => {
  loadList()
})
</script>

<style scoped>
.credit-list-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  height: 100%;
}
</style>
