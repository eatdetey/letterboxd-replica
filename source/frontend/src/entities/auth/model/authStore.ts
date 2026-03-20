import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { authApi, type AuthUserDto } from '../api/authApi'
import { authTokenStorage } from '@shared/api/authTokenStorage'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUserDto | null>(authTokenStorage.getUser<AuthUserDto>())
  const accessToken = ref<string | null>(authTokenStorage.getAccessToken())
  const refreshToken = ref<string | null>(authTokenStorage.getRefreshToken())
  const isReady = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value && !!accessToken.value)

  const username = computed(() => user.value?.username ?? '')
  const userId = computed(() => user.value?.id ?? '')
  const role = computed(() => user.value?.role ?? null)

  function setSession(nextUser: AuthUserDto, access: string, refresh: string) {
    user.value = nextUser
    accessToken.value = access
    refreshToken.value = refresh
    authTokenStorage.setSession(nextUser, access, refresh)
  }

  async function login(payload: { username: string; password: string }) {
    isLoading.value = true
    error.value = null

    try {
      const response = await authApi.login(payload)

      setSession(
        {
          id: response.id,
          username: response.username,
          role: response.role,
        },
        response.tokens.access,
        response.tokens.refresh,
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
        {
          id: response.id,
          username: response.username,
          role: response.role,
        },
        response.tokens.access,
        response.tokens.refresh,
      )
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Register failed'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function refresh() {
    if (!refreshToken.value) {
      throw new Error('No refresh token')
    }

    const response = await authApi.refresh({
      refresh_token: refreshToken.value,
    })

    accessToken.value = response.access
    refreshToken.value = response.refresh

    if (user.value) {
      authTokenStorage.setSession(user.value, response.access, response.refresh)
    }

    return response.access
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    authTokenStorage.clearSession()
  }

  async function init() {
    if (!refreshToken.value) {
      isReady.value = true
      return
    }

    try {
      if (!accessToken.value) {
        await refresh()
      } else if (!user.value) {
        const storedUser = authTokenStorage.getUser<AuthUserDto>()
        if (storedUser) {
          user.value = storedUser
        } else {
          await refresh()
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
    refreshToken,
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