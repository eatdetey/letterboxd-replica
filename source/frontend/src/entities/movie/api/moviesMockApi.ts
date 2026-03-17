import { getMock } from '@shared/api/mockClient'
import type { MovieDetailsDto, MoviesListResponseDto } from '../model/types'

export const moviesMockApi = {
  getMovies(): Promise<MoviesListResponseDto> {
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
