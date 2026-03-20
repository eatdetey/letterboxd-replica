import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { authApi, type AuthUserDto } from '../api/authApi'
import { authTokenStorage } from '@shared/api/authTokenStorage'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUserDto | null>(authTokenStorage.getUser<AuthUserDto>())
  const accessToken = ref<string | null>(authTokenStorage.getAccessToken())
  const isReady = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value && !!accessToken.value)

  const username = computed(() => user.value?.username ?? '')
  const userId = computed(() => (user.value?.id ? String(user.value.id) : ''))
  const role = computed(() => user.value?.role ?? null)

  function setSession(nextUser: AuthUserDto, access: string) {
    user.value = nextUser
    accessToken.value = access
    authTokenStorage.setSession(nextUser, access)
  }

  async function login(payload: { username: string; password: string }) {
    isLoading.value = true
    error.value = null

    try {
      const response = await authApi.login(payload)

      setSession(
        response.user,
        response.access_token,
      )
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Login failed'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function register(payload: { username: string; password: string; email: string }) {
    isLoading.value = true
    error.value = null

    try {
      const response = await authApi.register(payload)

      setSession(
        response.user,
        response.access_token,
      )
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Register failed'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function refresh() {
    const response = await authApi.refresh()

    accessToken.value = response.access_token

    if (user.value) {
      authTokenStorage.setSession(user.value, response.access_token)
    }

    return response.access_token
  }

  function logout() {
    user.value = null
    accessToken.value = null
    authTokenStorage.clearSession()
  }

  async function init() {
    if (!accessToken.value) {
      isReady.value = true
      return
    }

    try {
      if (!user.value) {
        const storedUser = authTokenStorage.getUser<AuthUserDto>()
        if (storedUser) {
          user.value = storedUser
        } else {
          logout()
        }
      }
    } catch {
      logout()
    } finally {
      isReady.value = true
    }
  }

  return {
    user,
    accessToken,
    isReady,
    isLoading,
    error,
    isAuthenticated,
    username,
    userId,
    role,
    login,
    register,
    refresh,
    logout,
    init,
  }
})
