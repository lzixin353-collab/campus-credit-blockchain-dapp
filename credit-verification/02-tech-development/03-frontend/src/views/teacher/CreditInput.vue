<template>
  <div class="credit-input-container">
    <el-page-header @back="() => {}" content="学分录入" />
    <el-card style="margin-top: 20px;">
      <el-form
        :model="form"
        :rules="rules"
        ref="formRef"
        label-width="120px"
        size="large"
        style="max-width: 560px;"
      >
        <el-form-item label="学生地址/学号" prop="student_address">
          <el-input
            v-model="form.student_address"
            placeholder="学生钱包地址(0x...)或学号，与合约 studentId 一致"
            clearable
          />
        </el-form-item>
        <el-form-item label="课程名称" prop="course_name">
          <el-input
            v-model="form.course_name"
            placeholder="如：区块链原理、计算机网络"
            clearable
          />
        </el-form-item>
        <el-form-item label="成绩(0-100)" prop="score">
          <el-input-number
            v-model="form.score"
            :min="0"
            :max="100"
            :step="1"
            placeholder="0-100 整数"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitCredit" :loading="submitLoading">
            提交上链
          </el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { recordCredit } from '@/api/credit'
import { ElMessage } from 'element-plus'

const formRef = ref(null)
const submitLoading = ref(false)

const form = ref({
  student_address: '',
  course_name: '',
  score: 80
})

const rules = {
  student_address: [
    { required: true, message: '请输入学生地址或学号', trigger: 'blur' }
  ],
  course_name: [
    { required: true, message: '请输入课程名称', trigger: 'blur' }
  ],
  score: [
    { required: true, message: '请输入成绩', trigger: 'blur' },
    { type: 'number', min: 0, max: 100, message: '成绩需在 0-100 之间', trigger: 'blur' }
  ]
}

const submitCredit = async () => {
  try {
    await formRef.value.validate()
  } catch (e) {
    return
  }
  submitLoading.value = true
  try {
    const res = await recordCredit({
      student_address: form.value.student_address.trim(),
      course_name: form.value.course_name.trim(),
      score: Number(form.value.score)
    })
    if (res && res.code === 200) {
      ElMessage.success('学分已提交上链，可在录入列表查看审核状态')
      resetForm()
    } else {
      ElMessage.error(res?.msg || '提交失败')
    }
  } catch (error) {
    const msg = error?.response?.data?.msg || error?.message || '提交失败'
    console.error('提交失败：', error?.response?.data || error)
    ElMessage.error(typeof msg === 'string' ? msg : (msg?.msg || '提交失败'))
  } finally {
    submitLoading.value = false
  }
}

const resetForm = () => {
  formRef.value?.resetFields()
  form.value = { student_address: '', course_name: '', score: 80 }
}
</script>

<style scoped>
.credit-input-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  height: 100%;
}
</style>
