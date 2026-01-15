import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/catalog',
    name: 'Catalog',
    component: () => import('@/views/CatalogView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/authors',
    name: 'AuthorsCatalog',
    component: () => import('@/views/AuthorsCatalogView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/poem/:id',
    name: 'PoemDetail',
    component: () => import('@/views/PoemDetailView.vue'),
    props: true,
    meta: { requiresAuth: true }
  },
  {
    path: '/author/:name',
    name: 'Author',
    component: () => import('@/views/AuthorView.vue'),
    props: true,
    meta: { requiresAuth: true }
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/views/SearchView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresGuest: true }  // 只有未登录用户可以访问
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/ProfileView.vue'),
    meta: { requiresAuth: true }  // 需要登录才能访问
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const isAuthenticated = userStore.isAuthenticated

  // 需要登录的页面
  if (to.meta.requiresAuth && !isAuthenticated) {
    // 未登录，跳转到登录页，并保存原目标路径
    next({
      name: 'Login',
      query: { redirect: to.fullPath }
    })
    return
  }

  // 只有未登录用户可以访问的页面（如登录页）
  if (to.meta.requiresGuest && isAuthenticated) {
    // 已登录用户访问登录页，跳转到首页
    next({ name: 'Home' })
    return
  }

  next()
})

export default router
