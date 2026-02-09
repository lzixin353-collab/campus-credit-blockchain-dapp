import axios from 'axios'
import { useUserStore } from '@/store/user'

// 版本兼容：axios 1.4.0配置
const service = axios.create({
  baseURL: '/api',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})

// 请求拦截器：添加JWT令牌
service.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.jwtToken) {
      config.headers['Authorization'] = `Bearer ${userStore.jwtToken}`
    }
    return config
  },
  (error) => {
    console.error('请求错误：', error)
    return Promise.reject(error)
  }
)

// 响应拦截器：统一错误处理（后端统一返回 { code, msg, data }）
service.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code !== 200) {
      alert(res.msg || res.message || '请求失败')
      return Promise.reject(res)
    }
    return res
  },
  (error) => {
    const status = error.response?.status
    const data = error.response?.data
    // 500 时把完整响应打到控制台，便于排查
    if (status === 500) {
      console.error('500 响应 body：', typeof data === 'object' ? JSON.stringify(data) : data)
    }
    console.error('请求失败：', status, data, error.message)
    let msg = '服务器错误'
    if (data && typeof data === 'object') {
      msg = data.msg || data.message || msg
    } else if (data && typeof data === 'string' && data) {
      msg = data
    } else if (error.message) {
      msg = error.message
    }
    if (status === 500) {
      msg = msg + '（请确认后端已用最新代码重启，并查看控制台/后端日志）'
    }
    alert(msg)
    return Promise.reject(error)
  }
)

export default service