import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

// 路由规则（版本兼容：vue-router 4.2.4）
const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  // 学生端
  {
    path: '/student',
    name: 'StudentLayout',
    component: () => import('@/views/student/Layout.vue'),
    meta: { requireAuth: true, role: 'student' },
    children: [
      { path: 'credit', name: 'StudentCredit', component: () => import('@/views/student/Credit.vue') },
      { path: 'profile', name: 'StudentProfile', component: () => import('@/views/student/Profile.vue') }
    ]
  },
  // 教师端
  {
    path: '/teacher',
    name: 'TeacherLayout',
    component: () => import('@/views/teacher/Layout.vue'),
    meta: { requireAuth: true, role: 'teacher' },
    children: [
      { path: 'credit-input', name: 'TeacherCreditInput', component: () => import('@/views/teacher/CreditInput.vue') },
      { path: 'credit-list', name: 'TeacherCreditList', component: () => import('@/views/teacher/CreditList.vue') }
    ]
  },
  // 管理员端
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('@/views/admin/Layout.vue'),
    meta: { requireAuth: true, role: 'admin' },
    children: [
      { path: 'credit-audit', name: 'AdminCreditAudit', component: () => import('@/views/admin/CreditAudit.vue') },
      { path: 'role-manage', name: 'AdminRoleManage', component: () => import('@/views/admin/RoleManage.vue') }
    ]
  },
  // 404
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫（权限校验）
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  // 恢复本地存储的用户信息
  userStore.restoreUserInfo()

  if (to.meta.requireAuth) {
    if (!userStore.isLogin) {
      next('/login')
    } else {
      if (userStore.role === to.meta.role) {
        next()
      } else {
        next(`/${userStore.role}`)
      }
    }
  } else {
    next()
  }
})

export default router