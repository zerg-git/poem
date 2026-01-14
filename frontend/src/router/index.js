import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue')
  },
  {
    path: '/catalog',
    name: 'Catalog',
    component: () => import('@/views/CatalogView.vue')
  },
  {
    path: '/authors',
    name: 'AuthorsCatalog',
    component: () => import('@/views/AuthorsCatalogView.vue')
  },
  {
    path: '/poem/:id',
    name: 'PoemDetail',
    component: () => import('@/views/PoemDetailView.vue'),
    props: true
  },
  {
    path: '/author/:name',
    name: 'Author',
    component: () => import('@/views/AuthorView.vue'),
    props: true
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/views/SearchView.vue')
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

export default router
