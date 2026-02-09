<template>
  <div class="credit-audit-container">
    <el-page-header @back="() => {}" content="学分审核" />
    <el-card style="margin-top: 20px;">
      <el-button type="primary" @click="loadPending" :loading="tableLoading" style="margin-bottom: 16px;">
        刷新待审核列表
      </el-button>
      <el-table
        :data="auditList"
        border
        stripe
        style="width: 100%"
        v-loading="tableLoading"
        empty-text="暂无待审核学分"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="学生地址/学号" min-width="160">
          <template #default="scope">
            {{ formatAddress(scope.row.student_address) }}
          </template>
        </el-table-column>
        <el-table-column prop="course_name" label="课程名称" min-width="140" />
        <el-table-column prop="score" label="成绩" width="90" />
        <el-table-column label="录入教师" min-width="140">
          <template #default="scope">
            {{ formatAddress(scope.row.teacher_address) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button
              type="success"
              size="small"
              @click="auditCredit(scope.row.id, true)"
              :loading="scope.row._loading"
            >
              通过
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="auditCredit(scope.row.id, false)"
              :loading="scope.row._loading"
            >
              驳回
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCreditPending, approveCredit, rejectCredit } from '@/api/credit'
import { ElMessage } from 'element-plus'

const tableLoading = ref(false)
const auditList = ref([])

function formatAddress(addr) {
  if (!addr) return '-'
  const s = String(addr)
  if (s.length <= 14) return s
  return s.slice(0, 6) + '...' + s.slice(-4)
}

const loadPending = async () => {
  tableLoading.value = true
  try {
    const res = await getCreditPending()
    const list = Array.isArray(res?.data) ? res.data : []
    auditList.value = list.map(item => ({ ...item, _loading: false }))
  } catch (error) {
    console.error('获取待审核列表失败：', error)
    auditList.value = []
  } finally {
    tableLoading.value = false
  }
}

const auditCredit = async (creditId, isApproved) => {
  const row = auditList.value.find(item => item.id === creditId)
  if (row) row._loading = true
  try {
    if (isApproved) {
      await approveCredit(creditId)
      ElMessage.success('已通过')
    } else {
      await rejectCredit(creditId)
      ElMessage.success('已驳回')
    }
    await loadPending()
  } catch (error) {
    ElMessage.error(error?.response?.data?.msg || error?.message || '操作失败')
  } finally {
    if (row) row._loading = false
  }
}

onMounted(() => {
  loadPending()
})
</script>

<style scoped>
.credit-audit-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  height: 100%;
}
</style>
