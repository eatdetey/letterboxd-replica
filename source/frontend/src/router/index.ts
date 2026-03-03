import { createRouter, createWebHistory } from 'vue-router'
import Authorization from "~/features/Authorization/Authorization.vue";

const routes = [
    {
        path: '/authorization',
        name: 'authorization',
        component: Authorization,
    },
    {
        path: '/',
        name: 'home',
        component:() => import('~/features/Home/Home.vue')
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router