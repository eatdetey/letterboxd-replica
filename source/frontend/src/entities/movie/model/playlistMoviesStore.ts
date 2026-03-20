import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { moviesApi } from '../api/moviesApi'
import { moviesMockApi } from '../api/moviesMockApi'
import type { MovieDto } from './types'

export const usePlaylistMoviesStore = defineStore('playlist-movies', () => {
  const items = ref<MovieDto[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const total = computed(() => items.value.length)

  async function loadPlaylistMovies(playlistId: string) {
    isLoading.value = true
    error.value = null

    try {
      const api = env.useMovieMocks ? moviesMockApi : moviesApi
      const response = await api.getMovies({ playlistId })
      items.value = response.items
    } catch (loadError) {
      items.value = []
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoading.value = false
    }
  }

  function reset() {
    items.value = []
    isLoading.value = false
    error.value = null
  }

  return {
    items,
    isLoading,
    error,
    total,
    loadPlaylistMovies,
    reset,
  }
})
