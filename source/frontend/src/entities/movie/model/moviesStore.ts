import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { moviesApi } from '../api/moviesApi'
import { moviesMockApi } from '../api/moviesMockApi'
import type { MovieDto } from './types'

const MOVIES_PAGE_SIZE = 20

export const useMoviesStore = defineStore('movies', () => {
  const items = ref<MovieDto[]>([])
  const isLoading = ref(false)
  const isLoadingMore = ref(false)
  const total = ref(0)
  const error = ref<string | null>(null)

  const featured = computed<MovieDto | null>(() => items.value[0] || null)
  const all = computed<MovieDto[]>(() => items.value)
  const hasMore = computed(() => items.value.length < total.value)

  function mergeMovies(nextItems: MovieDto[]) {
    const seenIds = new Set(items.value.map((item) => item.id))
    const uniqueItems = nextItems.filter((item) => !seenIds.has(item.id))
    items.value = [...items.value, ...uniqueItems]
  }

  async function loadMovies() {
    if (isLoading.value) {
      return
    }

    isLoading.value = true
    isLoadingMore.value = false
    error.value = null

    try {
      const api = env.useMovieMocks ? moviesMockApi : moviesApi
      const response = await api.getMovies({
        limit: MOVIES_PAGE_SIZE,
        offset: 0,
      })
      items.value = response.items
      total.value = response.total
    } catch (loadError) {
      items.value = []
      total.value = 0
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoading.value = false
    }
  }

  async function loadMoreMovies() {
    if (isLoading.value || isLoadingMore.value || !hasMore.value) {
      return
    }

    isLoadingMore.value = true
    error.value = null

    try {
      const api = env.useMovieMocks ? moviesMockApi : moviesApi
      const previousCount = items.value.length
      const response = await api.getMovies({
        limit: MOVIES_PAGE_SIZE,
        offset: items.value.length,
      })

      mergeMovies(response.items)
      total.value = response.total
      if (items.value.length === previousCount) {
        total.value = items.value.length
      }
    } catch (loadError) {
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoadingMore.value = false
    }
  }

  return {
    items,
    isLoading,
    isLoadingMore,
    total,
    error,
    featured,
    all,
    hasMore,
    loadMovies,
    loadMoreMovies,
  }
})
