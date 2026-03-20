import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { moviesApi } from '../api/moviesApi'
import { moviesMockApi } from '../api/moviesMockApi'
import type { MovieDto } from './types'

const PLAYLIST_MOVIES_PAGE_SIZE = 20

export const usePlaylistMoviesStore = defineStore('playlist-movies', () => {
  const items = ref<MovieDto[]>([])
  const isLoading = ref(false)
  const isLoadingMore = ref(false)
  const totalAvailable = ref(0)
  const currentPlaylistId = ref<string | null>(null)
  const error = ref<string | null>(null)

  const total = computed(() => items.value.length)
  const hasMore = computed(() => items.value.length < totalAvailable.value)

  async function loadPlaylistMovies(playlistId: string) {
    if (!playlistId || isLoading.value) {
      return
    }

    currentPlaylistId.value = playlistId
    isLoading.value = true
    isLoadingMore.value = false
    error.value = null

    try {
      const api = env.useMovieMocks ? moviesMockApi : moviesApi
      const response = await api.getMovies({
        playlistId,
        limit: PLAYLIST_MOVIES_PAGE_SIZE,
        offset: 0,
      })
      items.value = response.items
      totalAvailable.value = response.total
    } catch (loadError) {
      items.value = []
      totalAvailable.value = 0
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoading.value = false
    }
  }

  async function loadMorePlaylistMovies() {
    if (
      isLoading.value ||
      isLoadingMore.value ||
      !hasMore.value ||
      !currentPlaylistId.value
    ) {
      return
    }

    isLoadingMore.value = true
    error.value = null

    try {
      const api = env.useMovieMocks ? moviesMockApi : moviesApi
      const previousCount = items.value.length
      const response = await api.getMovies({
        playlistId: currentPlaylistId.value,
        limit: PLAYLIST_MOVIES_PAGE_SIZE,
        offset: items.value.length,
      })

      const seenIds = new Set(items.value.map((item) => item.id))
      const uniqueItems = response.items.filter((item) => !seenIds.has(item.id))
      items.value = [...items.value, ...uniqueItems]
      totalAvailable.value = response.total
      if (items.value.length === previousCount) {
        totalAvailable.value = items.value.length
      }
    } catch (loadError) {
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoadingMore.value = false
    }
  }

  function reset() {
    items.value = []
    isLoading.value = false
    isLoadingMore.value = false
    totalAvailable.value = 0
    currentPlaylistId.value = null
    error.value = null
  }

  return {
    items,
    isLoading,
    isLoadingMore,
    error,
    total,
    hasMore,
    loadPlaylistMovies,
    loadMorePlaylistMovies,
    reset,
  }
})
