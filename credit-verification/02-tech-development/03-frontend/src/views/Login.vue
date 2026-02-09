<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2 class="title">校园区块链学分存证系统</h2>
      <el-button 
        type="primary" 
        size="large" 
        class="connect-btn"
        @click="connectWallet"
        :loading="loading"
      >
        {{ address ? `切换钱包 (${shortAddress})` : '连接Metamask钱包' }}
      </el-button>
      <el-button 
        v-if="address"
        type="success" 
        size="large" 
        class="login-btn"
        @click="handleLogin"
        :loading="loginLoading"
      >
        登录系统
      </el-button>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { initWeb3, getCurrentAddress } from '@/utils/web3'
import { login } from '@/api/credit'

const router = useRouter()
const userStore = useUserStore()

// 状态（版本兼容：Vue 3.3.4响应式）
const loading = ref(false)
const loginLoading = ref(false)
const address = ref('')

// 短地址显示
const shortAddress = computed(() => {
  if (!address.value) return ''
  return `${address.value.slice(0, 6)}...${address.value.slice(-4)}`
})

// 连接钱包
const connectWallet = async () => {
  loading.value = true
  try {
    const web3 = await initWeb3()
    if (web3) {
      address.value = await getCurrentAddress(web3)
    }
  } catch (error) {
    console.error('连接钱包失败：', error)
  } finally {
    loading.value = false
  }
}

// 登录系统（角色以后端返回为准，不再依赖前端合约 getRole，避免 ABI 解码/未部署报错）
const handleLogin = async () => {
  loginLoading.value = true
  try {
    const web3 = await initWeb3()
    if (!web3) return

    // 直接调后端登录，后端按地址查库或链上决定角色
    const res = await login(address.value)
    const token = res?.data?.token
    const user = res?.data?.user
    if (!token || !user) {
      alert('登录失败，未返回 token 或用户信息')
      return
    }

    const role = (user.role || 'student').toLowerCase()
    userStore.login({
      address: address.value,
      role,
      jwtToken: token
    })
    router.push(`/${role}`)
  } catch (error) {
    console.error('登录失败：', error)
  } finally {
    loginLoading.value = false
  }
}

// 页面加载时恢复钱包连接
onMounted(() => {
  userStore.restoreUserInfo()
  if (userStore.address) {
    address.value = userStore.address
  }
})
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}
.login-card {
  width: 400px;
  padding: 20px;
  text-align: center;
}
.title {
  margin-bottom: 30px;
  color: #1989fa;
}
.btn-wrap {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 16px;
}
.btn-item {
  width: 100%;
}
</style>