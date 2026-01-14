import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { title: '仪表板' }
  },
  {
    path: '/smtp',
    name: 'SmtpConfig',
    component: () => import('../views/SmtpConfig.vue'),
    meta: { title: 'SMTP配置管理' }
  },
  {
    path: '/compose',
    name: 'ComposeEmail',
    component: () => import('../views/ComposeEmail.vue'),
    meta: { title: '邮件撰写' }
  },
  {
    path: '/templates',
    name: 'Templates',
    component: () => import('../views/Templates.vue'),
    meta: { title: '邮件模板' }
  },
  {
    path: '/history',
    name: 'History',
    component: () => import('../views/History.vue'),
    meta: { title: '发送历史' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫 - 设置页面标题
router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - SMTP邮件管理系统` : 'SMTP邮件管理系统'
  next()
})

export default router