import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { playlistsApi } from '../api/playlistsApi'
import { playlistsMockApi } from '../api/playlistsMockApi'
import type { PlaylistDto } from './types'

export const usePlaylistsStore = defineStore('playlists', () => {
  const items = ref<PlaylistDto[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const totalMovies = computed(() => {
    return items.value.reduce((sum, playlist) => sum + playlist.movies_count, 0)
  })

  async function loadPlaylists() {
    isLoading.value = true
    error.value = null

    try {
      const api = env.usePlaylistMocks ? playlistsMockApi : playlistsApi
      const response = await api.getPlaylists()
      items.value = response.items
    } catch (loadError) {
      items.value = []
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoading.value = false
    }
  }

  return {
    items,
    isLoading,
    error,
    totalMovies,
    loadPlaylists,
  }
})
