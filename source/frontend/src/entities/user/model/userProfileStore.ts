import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { env } from '@shared/config/env'
import { userProfileApi } from '../api/userProfileApi'
import { userProfileMockApi } from '../api/userProfileMockApi'

export interface MoviePreviewDto {
  id: string
  title: string
  release_year: number
  poster_url?: string | null
}

export interface UserProfileStatsDto {
  filmsWatched: number
  reviewsCount: number
  listsCount: number
  averageRating: number | null
  followersCount: number
  followingCount: number
  likesCount: number
}

export interface UserProfileDto {
  id: string
  username: string
  email?: string
  role?: 'user' | 'admin'
  avatar_url?: string | null
  cover_url?: string | null
  bio?: string | null
  tags?: string[]
  stats?: Partial<UserProfileStatsDto>
}

export interface UserListDto {
  id: string
  name: string
  description?: string | null
  is_public?: boolean
  movies_count?: number
  movies?: MoviePreviewDto[]
  preview_movies?: MoviePreviewDto[]
}

export interface UserReviewDto {
  id: string
  movie_id: string
  username: string
  text: string
  created_at: string
  movie?: MoviePreviewDto
}

export interface UserProfileResponse {
  user: UserProfileDto
  lists: UserListDto[]
  reviews: UserReviewDto[]
  recentRatings: MoviePreviewDto[]
  watchlist: MoviePreviewDto[]
}

export const useUserProfileStore = defineStore('userProfile', () => {
  const item = ref<UserProfileDto | null>(null)
  const lists = ref<UserListDto[]>([])
  const reviews = ref<UserReviewDto[]>([])
  const recentRatings = ref<MoviePreviewDto[]>([])
  const watchlist = ref<MoviePreviewDto[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const username = computed(() => item.value?.username ?? '')
  const avatarUrl = computed(() => item.value?.avatar_url ?? null)
  const coverUrl = computed(() => item.value?.cover_url ?? null)
  const stats = computed(() => item.value?.stats ?? null)

  const listsCount = computed(() => lists.value.length)
  const reviewsCount = computed(() => reviews.value.length)
  const watchlistCount = computed(() => watchlist.value.length)
  const recentRatingsCount = computed(() => recentRatings.value.length)

  async function loadProfile(idOrUsername: string) {
    if (!idOrUsername) {
      item.value = null
      lists.value = []
      reviews.value = []
      recentRatings.value = []
      watchlist.value = []
      return
    }

    isLoading.value = true
    error.value = null

    try {
      const api = env.useMocks ? userProfileMockApi : userProfileApi
      const response: UserProfileResponse = await api.getProfile(idOrUsername)

      item.value = response.user
      lists.value = response.lists ?? []
      reviews.value = response.reviews ?? []
      recentRatings.value = response.recentRatings ?? []
      watchlist.value = response.watchlist ?? []
    } catch (loadError) {
      error.value = loadError instanceof Error ? loadError.message : 'Unknown error'
      item.value = null
      lists.value = []
      reviews.value = []
      recentRatings.value = []
      watchlist.value = []
    } finally {
      isLoading.value = false
    }
  }

  function reset() {
    item.value = null
    lists.value = []
    reviews.value = []
    recentRatings.value = []
    watchlist.value = []
    isLoading.value = false
    error.value = null
  }

  return {
    item,
    lists,
    reviews,
    recentRatings,
    watchlist,
    isLoading,
    error,
    username,
    avatarUrl,
    coverUrl,
    stats,
    listsCount,
    reviewsCount,
    watchlistCount,
    recentRatingsCount,
    loadProfile,
    reset,
  }
})