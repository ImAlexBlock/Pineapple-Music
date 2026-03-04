import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
  },
  {
    path: '/bootstrap',
    name: 'Bootstrap',
    component: () => import('../views/Bootstrap.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
  },
  {
    path: '/tracks',
    name: 'Tracks',
    component: () => import('../views/TrackList.vue'),
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('../views/Search.vue'),
  },
  {
    path: '/playlists',
    name: 'Playlists',
    component: () => import('../views/Playlists.vue'),
  },
  {
    path: '/playlists/:id',
    name: 'PlaylistDetail',
    component: () => import('../views/PlaylistDetail.vue'),
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/AdminDashboard.vue'),
    meta: { requiresAdmin: true },
  },
  {
    path: '/admin/settings',
    name: 'AdminSettings',
    component: () => import('../views/AdminSettings.vue'),
    meta: { requiresAdmin: true },
  },
  {
    path: '/admin/upload',
    name: 'AdminUpload',
    component: () => import('../views/AdminUpload.vue'),
    meta: { requiresAdmin: true },
  },
  {
    path: '/admin/scan',
    name: 'AdminScan',
    component: () => import('../views/AdminScan.vue'),
    meta: { requiresAdmin: true },
  },
  {
    path: '/admin/audit',
    name: 'AdminAudit',
    component: () => import('../views/AdminAudit.vue'),
    meta: { requiresAdmin: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const auth = useAuthStore()

  // Check bootstrap status on first navigation
  if (!auth.checked) {
    await auth.checkStatus()
  }

  // Redirect to bootstrap if not bootstrapped
  if (!auth.bootstrapped && to.name !== 'Bootstrap') {
    return next({ name: 'Bootstrap' })
  }

  // Don't go to bootstrap if already bootstrapped
  if (auth.bootstrapped && to.name === 'Bootstrap') {
    return next({ name: 'Home' })
  }

  // Admin routes require admin role
  if (to.meta.requiresAdmin && auth.role !== 'admin') {
    return next({ name: 'Login' })
  }

  next()
})

export default router
