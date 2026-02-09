<template>
  <div class="profile-container">
    <el-page-header @back="() => {}" content="个人信息" />
    <el-card style="margin-top: 20px; max-width: 560px;">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="钱包地址">
          <span>{{ displayAddress }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="角色">
          <el-tag type="info">{{ displayRole }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
      <div style="margin-top: 20px; color: #909399; font-size: 13px;">
        学籍信息由教师录入学分后与链上数据同步，可在「我的学分」中查看已录入课程与审核状态。
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

const displayAddress = computed(() => {
  return userStore.address || '未绑定'
})

const displayRole = computed(() => {
  const r = userStore.role
  if (r === 'student') return '学生'
  if (r === 'teacher') return '教师'
  if (r === 'admin') return '管理员'
  return r || '未知'
})

onMounted(() => {
  userStore.restoreUserInfo()
})
</script>

<style scoped>
.profile-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  height: 100%;
}
</style>
