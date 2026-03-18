import { createRouter, createWebHistory } from 'vue-router'

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
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
