import { ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { reviewsApi } from '../api/reviewsApi'
import { reviewsMockApi } from '../api/reviewsMockApi'
import type { MovieReviewDto } from './types'

export const useMovieReviewsStore = defineStore('movie-reviews', () => {
  const items = ref<MovieReviewDto[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function loadReviews(movieId: string) {
    isLoading.value = true
    error.value = null

    try {
      const api = env.useMocks ? reviewsMockApi : reviewsApi
      const response = await api.getMovieReviews(movieId)
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
    loadReviews,
  }
})
