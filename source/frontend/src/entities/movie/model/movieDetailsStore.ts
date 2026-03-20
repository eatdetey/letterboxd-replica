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
      if (env.useMovieMocks) {
        item.value = await moviesMockApi.getMovie(movieId)
        return
      }

      try {
        item.value = await moviesApi.getMovie(movieId)
      } catch (movieError) {
        const message = movieError instanceof Error ? movieError.message : ''

        if (!message.includes('HTTP 404')) {
          throw movieError
        }

        const response = await moviesApi.getMovies({
          ids: [movieId],
          enrichPlaylists: true,
        })
        const fallbackMovie = response.items.find((movie) => movie.id === movieId)

        if (!fallbackMovie) {
          throw movieError
        }

        item.value = {
          ...fallbackMovie,
          playlists: fallbackMovie.playlists ?? [],
        }
      }
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
