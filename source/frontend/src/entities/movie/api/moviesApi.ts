import { httpClient } from '@shared/api/httpClient'
import type { MovieDetailsDto, MoviesListResponseDto } from '../model/types'

type GetMoviesParams = {
  playlistId?: string
}

export const moviesApi = {
  getMovies(params: GetMoviesParams = {}): Promise<MoviesListResponseDto> {
    const searchParams = new URLSearchParams()

    if (params.playlistId) {
      searchParams.set('playlist_id', params.playlistId)
    }

    const query = searchParams.toString()
    return httpClient.get<MoviesListResponseDto>(`/v1/movies${query ? `?${query}` : ''}`)
  },
  getMovie(id: string): Promise<MovieDetailsDto> {
    return httpClient.get<MovieDetailsDto>(`/v1/movies/${id}`)
  },
}
