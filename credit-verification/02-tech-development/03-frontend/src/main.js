import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
// 全局样式（可选）
import './assets/style/global.css'

const app = createApp(App)
// 挂载Pinia和路由
app.use(createPinia())
app.use(router)
app.mount('#app')