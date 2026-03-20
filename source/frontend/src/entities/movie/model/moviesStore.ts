import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { moviesApi } from '../api/moviesApi'
import { moviesMockApi } from '../api/moviesMockApi'
import type { MovieDto } from './types'

export const useMoviesStore = defineStore('movies', () => {
  const items = ref<MovieDto[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const featured = computed<MovieDto | null>(() => items.value[0] || null)

  const all = computed<MovieDto[]>(() => items.value)

  async function loadMovies() {
    isLoading.value = true
    error.value = null

    try {
      const api = env.useMovieMocks ? moviesMockApi : moviesApi
      const response = await api.getMovies()
      items.value = response.items
    } catch (loadError) {
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoading.value = false
    }
  }

  return {
    items,
    isLoading,
    error,
    featured,
    all,
    loadMovies,
  }
})
