import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from '~/entities/auth/model/authStore'
import type { UserProfileDto } from './types'

export const useUserProfileStore = defineStore('userProfile', () => {
  const item = ref<UserProfileDto | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const username = computed(() => item.value?.username ?? '')

  async function loadProfile(idOrUsername: string) {
    const authStore = useAuthStore()

    if (!idOrUsername) {
      item.value = null
      error.value = null
      return
    }

    isLoading.value = true
    error.value = null

    try {
      const currentUser = authStore.user
      const currentUserId = currentUser?.id ? String(currentUser.id) : ''
      const currentUsername = currentUser?.username ?? ''
      const matchesCurrentUser = idOrUsername === currentUsername || idOrUsername === currentUserId

      if (!currentUser || !matchesCurrentUser) {
        item.value = null
        error.value = 'User not found'
        return
      }

      item.value = {
        id: String(currentUser.id),
        username: currentUser.username,
        email: currentUser.email,
        avatar_url: currentUser.avatar_url,
        role: currentUser.role as UserProfileDto['role'],
      }
    } catch (loadError) {
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
      item.value = null
    } finally {
      isLoading.value = false
    }
  }

  function reset() {
    item.value = null
    isLoading.value = false
    error.value = null
  }

  return {
    item,
    isLoading,
    error,
    username,
    loadProfile,
    reset,
  }
})
