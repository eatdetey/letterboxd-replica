import { httpClient } from '@shared/api/httpClient'
import type { MovieDetailsDto, MoviesListResponseDto } from '../model/types'

type GetMoviesParams = {
  playlistId?: string
  search?: string
  ids?: string[]
  enrichPlaylists?: boolean
  limit?: number
  offset?: number
}

type GetMoviesRequestDto = {
  playlist_id?: string
  search?: string
  ids?: string[]
  enrich_playlists?: boolean
  limit?: number
  offset?: number
}

type MovieResponseDto = {
  movie: MovieDetailsDto
}

function isNotFound(error: unknown) {
  return error instanceof Error && error.message.includes('HTTP 404')
}

export const moviesApi = {
  async getMovies(params: GetMoviesParams = {}): Promise<MoviesListResponseDto> {
    const payload: GetMoviesRequestDto = {}

    if (params.playlistId) {
      payload.playlist_id = params.playlistId
    }

    if (params.search?.trim()) {
      payload.search = params.search.trim()
    }

    if (params.ids?.length) {
      payload.ids = params.ids
    }

    if (params.enrichPlaylists) {
      payload.enrich_playlists = true
    }

    if (typeof params.limit === 'number') {
      payload.limit = params.limit
    }

    if (typeof params.offset === 'number') {
      payload.offset = params.offset
    }

    try {
      return await httpClient.post<MoviesListResponseDto>('/v1/movies/search', payload)
    } catch (error) {
      if (!isNotFound(error)) {
        throw error
      }

      const searchParams = new URLSearchParams()

      if (params.playlistId) {
        searchParams.set('playlist_id', params.playlistId)
      }

      if (params.search?.trim()) {
        searchParams.set('search', params.search.trim())
      }

      if (typeof params.limit === 'number') {
        searchParams.set('limit', String(params.limit))
      }

      if (typeof params.offset === 'number') {
        searchParams.set('offset', String(params.offset))
      }

      const query = searchParams.toString()
      return httpClient.get<MoviesListResponseDto>(`/v1/movies${query ? `?${query}` : ''}`)
    }
  },

  async getMovie(id: string): Promise<MovieDetailsDto> {
    const response = await httpClient.get<MovieResponseDto | MovieDetailsDto>(`/v1/movies/${id}`)

    if ('movie' in response) {
      return response.movie
    }

    return response
  },
}