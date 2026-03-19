import { getMock } from '@shared/api/mockClient'
import type { MovieDetailsDto, MoviesListResponseDto } from '../model/types'

type GetMoviesParams = {
  playlistId?: string
}

export const moviesMockApi = {
  async getMovies(params: GetMoviesParams = {}): Promise<MoviesListResponseDto> {
    if (params.playlistId) {
      try {
        return await getMock<MoviesListResponseDto>(`v1/playlists/${params.playlistId}/movies/response.json`)
      } catch {
        return { items: [], total: 0 }
      }
    }

    return getMock<MoviesListResponseDto>('v1/movies/response.json')
  },
  async getMovie(id: string): Promise<MovieDetailsDto> {
    try {
      return await getMock<MovieDetailsDto>(`v1/movies/${id}/response.json`)
    } catch {
      const response = await getMock<MoviesListResponseDto>('v1/movies/response.json')
      const movie = response.items.find((item) => item.id === id)

      if (!movie) {
        throw new Error(`Movie with id "${id}" not found in mocks`)
      }

      return {
        ...movie,
        playlists: [],
      }
    }
  },
}
