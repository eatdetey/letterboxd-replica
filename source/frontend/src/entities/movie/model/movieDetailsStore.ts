import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { moviesApi } from '../api/moviesApi'
import { moviesMockApi } from '../api/moviesMockApi'
import type { MovieDetailsDto } from './types'

export const useMovieDetailsStore = defineStore('movie-details', () => {
  const item = ref<MovieDetailsDto | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const playlists = computed(() => item.value?.playlists ?? [])

  async function loadMovie(movieId: string) {
    isLoading.value = true
    error.value = null

    try {
      const api = env.useMocks ? moviesMockApi : moviesApi
      item.value = await api.getMovie(movieId)
    } catch (loadError) {
      item.value = null
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoading.value = false
    }
  }

  return {
    item,
    isLoading,
    error,
    playlists,
    loadMovie,
  }
})
