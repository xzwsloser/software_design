import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '../components/AuthView.vue'
import SiteList from '../components/SiteList.vue'
import SiteDetail from '../components/SiteDetail.vue'
import LikedSites from '../components/LikedSites.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: AuthView
  },
  {
    path: '/sites',
    name: 'SiteList',
    component: SiteList,
    meta: { requiresAuth: true }
  },
  {
    path: '/sites/:siteIndex',
    name: 'SiteDetail',
    component: SiteDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/liked-sites',
    name: 'LikedSites',
    component: LikedSites,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router