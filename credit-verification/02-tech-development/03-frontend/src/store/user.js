import { defineStore } from 'pinia'

// 版本兼容：Pinia 2.1.7语法
export const useUserStore = defineStore('user', {
  state: () => ({
    address: '', // 钱包地址
    role: '', // student/teacher/admin
    isLogin: false,
    jwtToken: ''
  }),
  actions: {
    login(userInfo) {
      this.address = userInfo.address
      this.role = userInfo.role
      this.isLogin = true
      this.jwtToken = userInfo.jwtToken
      localStorage.setItem('userInfo', JSON.stringify(userInfo))
    },
    logout() {
      this.address = ''
      this.role = ''
      this.isLogin = false
      this.jwtToken = ''
      localStorage.removeItem('userInfo')
    },
    restoreUserInfo() {
      const userInfo = localStorage.getItem('userInfo')
      if (userInfo) {
        const info = JSON.parse(userInfo)
        this.address = info.address
        this.role = info.role
        this.isLogin = true
        this.jwtToken = info.jwtToken
      }
    }
  }
})