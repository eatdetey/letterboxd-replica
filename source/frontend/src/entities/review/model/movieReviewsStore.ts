import { ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { reviewsApi } from '../api/reviewsApi'
import { reviewsMockApi } from '../api/reviewsMockApi'
import type { CreateMovieReviewPayloadDto, MovieReviewDto } from './types'

export const useMovieReviewsStore = defineStore('movie-reviews', () => {
  const items = ref<MovieReviewDto[]>([])
  const isLoading = ref(false)
  const isSubmitting = ref(false)
  const error = ref<string | null>(null)

  async function loadReviews(movieId: string) {
    isLoading.value = true
    error.value = null

    try {
      const api = env.useReviewMocks ? reviewsMockApi : reviewsApi
      const response = await api.getMovieReviews(movieId)
      items.value = response.items
    } catch (loadError) {
      items.value = []
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
    } finally {
      isLoading.value = false
    }
  }

  async function createReview(movieId: string, payload: CreateMovieReviewPayloadDto) {
    isSubmitting.value = true
    error.value = null

    try {
      const api = env.useReviewMocks ? reviewsMockApi : reviewsApi
      const response = await api.createMovieReview(movieId, payload)
      items.value = [response.review, ...items.value]
      return response.review
    } catch (submitError) {
      error.value = submitError instanceof Error ? submitError.message : 'Unknown error'
      throw submitError
    } finally {
      isSubmitting.value = false
    }
  }

  return {
    items,
    isLoading,
    isSubmitting,
    error,
    loadReviews,
    createReview,
  }
})
