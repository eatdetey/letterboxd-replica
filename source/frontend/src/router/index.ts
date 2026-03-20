import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '~/entities/auth/model/authStore'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('~/pages/home/ui/HomePage.vue'),
  },
  {
    path: '/movies/:id',
    name: 'movie',
    component: () => import('~/pages/movie/ui/MoviePage.vue'),
  },
  {
    path: '/playlists',
    name: 'playlists',
    component: () => import('~/pages/playlists/ui/PlaylistsPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/playlists/:id',
    name: 'playlist',
    component: () => import('~/pages/playlist/ui/PlaylistPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/user/:username',
    name: 'user',
    component: () => import('~/pages/user/ui/UserPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('~/pages/login/ui/LoginPage.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (!authStore.isReady) {
    await authStore.init()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return {
      name: 'login',
      query: { redirect: to.fullPath },
    }
  }

  if (to.name === 'login' && authStore.isAuthenticated) {
    return { name: 'home' }
  }

  return true
})

export default router