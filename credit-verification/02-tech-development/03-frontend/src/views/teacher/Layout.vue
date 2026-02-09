<template>
  <el-container style="height: 100vh;">
    <el-aside width="200px" style="background-color: #2e3b4e;">
      <el-menu
        :default-active="$route.path"
        class="el-menu-vertical-demo"
        background-color="#2e3b4e"
        text-color="#fff"
        active-text-color="#ffd04b"
        @select="onMenuSelect"
      >
        <el-menu-item index="/teacher/credit-input">
          <span>学分录入</span>
        </el-menu-item>
        <el-menu-item index="/teacher/credit-list">
          <span>录入列表</span>
        </el-menu-item>
        <el-menu-item index="logout">
          <span>退出登录</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="text-align: right; font-size: 12px">
        <el-dropdown>
          <i class="el-icon-setting" style="margin-right: 15px"></i>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>{{ address }}</el-dropdown-item>
              <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <span>{{ shortAddress }}</span>
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
const userStore = useUserStore()
const router = useRouter()

const address = ref(userStore.address || '')
const shortAddress = computed(() => {
  return address.value ? `${String(address.value).slice(0, 6)}...${String(address.value).slice(-4)}` : ''
})

const onMenuSelect = (index) => {
  if (index === 'logout') {
    logout()
  } else {
    router.push(index)
  }
}

const logout = () => {
  userStore.logout()
  ElMessage.success('退出登录成功！')
  router.push('/login')
}
</script>

<style scoped>
.el-header {
  background-color: #fff;
  color: #333;
  line-height: 60px;
  border-bottom: 1px solid #e6e6e6;
}
.el-aside {
  color: #333;
}
.el-main {
  padding: 20px;
  background-color: #f5f5f5;
  height: calc(100vh - 60px);
  overflow-y: auto;
}
</style>