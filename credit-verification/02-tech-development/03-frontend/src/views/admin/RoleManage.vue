<template>
  <div class="role-manage-container">
    <el-page-header @back="() => {}" content="角色管理" />
    <el-card style="margin-top: 20px; max-width: 640px;">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" size="large">
        <el-form-item label="钱包地址" prop="user_address">
          <el-input
            v-model="form.user_address"
            placeholder="0x 开头的以太坊地址"
            clearable
          />
        </el-form-item>
        <el-form-item label="分配角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色" style="width: 100%;">
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="assignRoleHandle" :loading="assignLoading">
            分配角色（链上 + 后端）
          </el-button>
        </el-form-item>
      </el-form>
      <div style="color: #909399; font-size: 13px; margin-top: 12px;">
        分配后该地址在链上拥有对应权限；若需用该钱包登录并保持角色，请让用户先绑定钱包或由管理员在库中设置 address。
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { assignRole as apiAssignRole } from '@/api/credit'
import { ElMessage } from 'element-plus'

const formRef = ref(null)
const assignLoading = ref(false)
const form = ref({
  user_address: '',
  role: ''
})

const rules = {
  user_address: [
    { required: true, message: '请输入钱包地址', trigger: 'blur' },
    { pattern: /^0x[a-fA-F0-9]{40}$/i, message: '请输入有效的 0x+40 位十六进制地址', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

const assignRoleHandle = async () => {
  try {
    await formRef.value.validate()
  } catch (e) {
    return
  }
  assignLoading.value = true
  try {
    await apiAssignRole(form.value.user_address.trim(), form.value.role)
    ElMessage.success('角色分配成功')
    formRef.value.resetFields()
  } catch (error) {
    ElMessage.error(error?.response?.data?.msg || error?.message || '分配失败')
  } finally {
    assignLoading.value = false
  }
}
</script>

<style scoped>
.role-manage-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  height: 100%;
}
</style>
